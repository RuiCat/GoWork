package components

import (
	"image"
	"ui/gioui/io/event"
	"ui/gioui/io/pointer"
	"ui/gioui/mat/f32"
	"ui/gioui/op"
	"ui/gioui/op/clip"
	"ui/gioui/op/paint"
	"ui/gioui/widget"
	"ui/gioui/widget/layout"
	"ui/gioui/widget/material"
)

// ViewsNew 创建控件
func ViewsNew(th *material.Theme, X, Y, Width, Height int, Size image.Point) *Views {
	views := &Views{
		Theme: th,
		Size:  Size,
		Position: image.Point{
			X: X,
			Y: Y,
		},
		Dimensions: layout.Dimensions{
			Size: image.Point{
				X: Width,
				Y: Height,
			},
			Baseline: 2,
		},
		UseGrid:      true,
		GridDistance: 50,
	}
	views.affine.Scale = 1
	views.clipRect = clip.Rect{Max: views.Dimensions.Size}.Op()
	views.Scrollbar[0].Use = true
	views.Scrollbar[0].ScrollbarStyle = material.Scrollbar(views.Theme, &views.Scrollbar[0].Scrollbar)
	views.Scrollbar[1].Use = true
	views.Scrollbar[1].ScrollbarStyle = material.Scrollbar(views.Theme, &views.Scrollbar[1].Scrollbar)
	views.Background()
	return views
}

// Views 视图控件
type Views struct {
	*material.Theme                   // 主题
	Size            image.Point       // 组件内画板大小
	Position        image.Point       // 组件位置
	Dimensions      layout.Dimensions // 组件大小
	UseGrid         bool              // 启用网格
	GridDistance    int               // 网格大小
	ComponentList   []struct {
		Size          image.Point // 组件大小
		Position      image.Point // 组件位置
		layout.Widget             // 组件对象
	}
	componentListIndex  []int // 组件索引
	componentListUpdate bool  // 组件需要更新
	Scrollbar           [2]struct {
		Use                     bool    // 滚动条生效
		Distance                float32 // 滚动条滚动距离
		widget.Scrollbar                // 滚动条对象
		material.ScrollbarStyle         // 滚动条对象风格
	}
	stroke   [5]clip.Op // 组件背景网格与边框
	clipRect clip.Op    // 用于组件剪裁
	moves    struct {
		Use      bool      // 背景移动状态标记
		Position f32.Point // 背景移动位置记录
	}
	affine struct {
		Scale        float32 // 缩放倍数
		f32.Affine2D         // 缩放矩阵
	}
}

// Add 添加组件
func (v *Views) Add(point image.Point, size image.Point, widget layout.Widget) {
	v.ComponentList = append(v.ComponentList, struct {
		Size     image.Point
		Position image.Point
		layout.Widget
	}{
		Size:     size,
		Position: point,
		Widget:   widget,
	})
	v.Scale(0)
}

// Scale 组件缩放
func (v *Views) Scale(scale float32) {
	if scale > 0 {
		v.affine.Scale = scale
	}
	v.affine.Affine2D = f32.Affine2D{}.Scale(
		f32.Point{X: 1, Y: 1},
		f32.Point{X: v.affine.Scale, Y: v.affine.Scale})
}

// Layout 布局
func (v *Views) Layout(gtx layout.Context) layout.Dimensions {
	// 偏移位置
	offest := f32.Point{
		X: v.Scrollbar[0].Distance * float32(v.Size.X),
		Y: v.Scrollbar[1].Distance * float32(v.Size.Y),
	}
	// 处理事件
	event.Op(gtx.Ops, v)
	if ev, ok := gtx.Source.Event(pointer.Filter{
		Target: v,
		Kinds:  pointer.Press | pointer.Release | pointer.Move | pointer.Drag,
	}); ok {
		event := ev.(pointer.Event)
		pos := event.Position.Sub(f32.Point{X: float32(v.Position.X), Y: float32(v.Position.Y)})
		if pos.X > 0 && pos.Y > 0 && pos.X < float32(v.Dimensions.Size.X)-10 && pos.Y < float32(v.Dimensions.Size.Y)-10 || v.moves.Use {
			// 背景移动实现
			if event.Buttons == pointer.ButtonTertiary {
				switch event.Kind {
				case pointer.Press:
					if !v.moves.Use {
						v.moves.Position = pos.Add(offest)
					}
					v.moves.Use = true
					v.componentListUpdate = true
				case pointer.Drag:
					if v.moves.Use {
						pos := pos.Sub(v.moves.Position)
						v.Scrollbar[0].Distance = -pos.X / float32(v.Size.X)
						v.Scrollbar[1].Distance = -pos.Y / float32(v.Size.Y)
					}
				}
			} else {
				v.moves.Use = false
			}
		}
	}
	// 更新显示组件
	if v.componentListUpdate {
		v.componentListIndex = v.componentListIndex[:0]
		// 计算内容
		for i, cl := range v.ComponentList {
			// 计算显示矩阵
			if cl.Position.X <= int(offest.X)+v.Dimensions.Size.X &&
				cl.Position.X+cl.Size.X >= int(offest.X) &&
				cl.Position.Y <= int(offest.Y)+v.Dimensions.Size.Y &&
				cl.Position.Y+cl.Size.Y >= int(offest.Y) {
				v.componentListIndex = append(v.componentListIndex, i)
			}
		}
	}
	// 计算滚动条位置
	if !v.moves.Use {
		v.Scrollbar[1].Distance += v.Scrollbar[1].Scrollbar.ScrollDistance()
		v.Scrollbar[0].Distance += v.Scrollbar[0].Scrollbar.ScrollDistance()
	} else {
		if v.Scrollbar[0].Distance < 0 {
			v.Scrollbar[0].Distance = 0
		}
		if v.Scrollbar[1].Distance < 0 {
			v.Scrollbar[1].Distance = 0
		}
		if v.Scrollbar[0].Distance >= 1 {
			v.Scrollbar[0].Distance = 1
		}
		if v.Scrollbar[1].Distance >= 1 {
			v.Scrollbar[1].Distance = 1
		}
	}
	// 绘制控件
	defer op.Offset(v.Position).Push(gtx.Ops).Pop()         // 移动画布
	defer v.clipRect.Push(gtx.Ops).Pop()                    // 剪裁组件位置
	AffinePop := op.Affine(v.affine.Affine2D).Push(gtx.Ops) // 缩放剩余对象
	// 绘制边框
	paint.FillShape(gtx.Ops, v.Fg, v.stroke[4])
	// 绘制背景网格
	if v.UseGrid {
		offest := op.Offset(image.Point{X: -int(v.Scrollbar[0].Distance*float32(v.Size.X)) % v.GridDistance}).Push(gtx.Ops)
		paint.FillShape(gtx.Ops, v.Fg, v.stroke[0])         // 背景线
		paint.FillShape(gtx.Ops, v.ContrastBg, v.stroke[1]) // 背景线
		offest.Pop()
		offest = op.Offset(image.Point{Y: -int(v.Scrollbar[1].Distance*float32(v.Size.Y)) % v.GridDistance}).Push(gtx.Ops)
		paint.FillShape(gtx.Ops, v.Fg, v.stroke[2])         // 背景线
		paint.FillShape(gtx.Ops, v.ContrastBg, v.stroke[3]) // 背景线
		offest.Pop()
	}
	// 更新并且绘制组件
	OffsetPop := op.Offset(image.Point{
		X: -int(offest.X),
		Y: -int(offest.Y)},
	).Push(gtx.Ops) // 移动剪裁位置
	for _, i := range v.componentListIndex {
		cl := v.ComponentList[i]
		offset := op.Offset(cl.Position).Push(gtx.Ops)
		gtx.Constraints.Max = cl.Size
		cl.Widget(gtx)
		offset.Pop()
	}
	OffsetPop.Pop() // 还原组件偏移移动
	AffinePop.Pop() // 还原组件缩放
	// 绘制滚动条
	gtx.Constraints.Max = v.Dimensions.Size.Sub(image.Point{X: 5, Y: 5}) // 定义长度并且减去右下角的空余位置偏移
	if v.Scrollbar[0].Use {
		TransformStack := op.Offset(image.Point{Y: v.Dimensions.Size.Y - 10}).Push(gtx.Ops)
		v.Scrollbar[0].ScrollbarStyle.Layout(gtx, layout.Horizontal, v.Scrollbar[0].Distance, v.Scrollbar[0].Distance)
		TransformStack.Pop()
	}
	if v.Scrollbar[1].Use {
		TransformStack := op.Offset(image.Point{X: v.Dimensions.Size.X - 10}).Push(gtx.Ops)
		v.Scrollbar[1].ScrollbarStyle.Layout(gtx, layout.Vertical, v.Scrollbar[1].Distance, v.Scrollbar[1].Distance)
		TransformStack.Pop()
	}
	return v.Dimensions
}

// Background 更新背景
func (v *Views) Background() {
	if v.UseGrid {
		x := v.Dimensions.Size.X/v.GridDistance + 1
		y := v.Dimensions.Size.Y/v.GridDistance + 1
		// 绘制背景网格
		var gridsX clip.Path
		gridsX.Begin(&op.Ops{})
		var gridsXZ clip.Path
		gridsXZ.Begin(&op.Ops{})
		for i := 0; i < x; i++ {
			ix := float32(i * v.GridDistance)
			gridsX.MoveTo(f32.Pt(ix, 0))
			gridsX.LineTo(f32.Pt(ix, float32(v.Dimensions.Size.Y)))
			iz := float32((i + 1) * v.GridDistance)
			for ; ix < iz; ix += 10 {
				gridsXZ.MoveTo(f32.Pt(ix, 0))
				gridsXZ.LineTo(f32.Pt(ix, float32(v.Dimensions.Size.Y)))
			}
		}
		var gridsY clip.Path
		gridsY.Begin(&op.Ops{})
		var gridsYZ clip.Path
		gridsYZ.Begin(&op.Ops{})
		for i := 0; i < y; i++ {
			iy := float32(i * v.GridDistance)
			gridsY.MoveTo(f32.Pt(0, iy))
			gridsY.LineTo(f32.Pt(float32(v.Dimensions.Size.X), iy))
			iz := float32((i + 1) * v.GridDistance)
			for ; iy < iz; iy += 10 {
				gridsY.MoveTo(f32.Pt(0, iy))
				gridsY.LineTo(f32.Pt(float32(v.Dimensions.Size.X), iy))
			}
		}
		v.stroke[0] = clip.Stroke{Path: gridsX.End(), Width: 1.5}.Op()
		v.stroke[1] = clip.Stroke{Path: gridsXZ.End(), Width: 0.5}.Op()
		v.stroke[2] = clip.Stroke{Path: gridsY.End(), Width: 0.5}.Op()
		v.stroke[3] = clip.Stroke{Path: gridsYZ.End(), Width: 1}.Op()
	}
	// 绘制组件边框
	v.stroke[4] = clip.Stroke{Path: clip.RRect{
		Rect: image.Rectangle{
			Max: v.Dimensions.Size,
		},
	}.Path(&op.Ops{}), Width: 1.5}.Op()
	v.componentListUpdate = true
}
