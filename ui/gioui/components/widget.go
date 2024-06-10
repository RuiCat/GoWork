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

// OpsFace 上下文底层实现
type OpsFace interface {
	GetOps() *op.Ops           // 底层上下文
	GetTheme() *material.Theme // 底层主题
}

// LayoutFace 布局接口用于绘制组件
type LayoutFace interface {
	OpsFace                                       // 继承上下文
	Update()                                      // 状态更新
	Layout(gtx layout.Context) *layout.Dimensions // 组件布局
	GetPoint() *image.Point                       // 组件绘制位置
	GetDimensions() *layout.Dimensions            // 组件大小
	SetDisable(disable bool)                      // 组件禁用
	GetDisable() bool                             // 组件禁用
}

// UILayout 组件实现
type UILayout struct {
	LayoutFace
	Disable bool // 是否禁用组件
	*op.Ops
	*material.Theme    // 主题
	*image.Point       // 绘制位置
	*layout.Dimensions // 大小
}

func (ui *UILayout) GetOps() *op.Ops                   { return ui.Ops }
func (ui *UILayout) GetTheme() *material.Theme         { return ui.Theme }
func (ui *UILayout) GetPoint() *image.Point            { return ui.Point }
func (ui *UILayout) GetDimensions() *layout.Dimensions { return ui.Dimensions }
func (ui *UILayout) SetDisable(disable bool)           { ui.Disable = disable }
func (ui *UILayout) GetDisable() bool                  { return ui.Disable }
func (ui *UILayout) Update() {
	if ui.LayoutFace != nil {
		ui.LayoutFace.Update()
	}
}
func (ui *UILayout) Layout(gtx layout.Context) *layout.Dimensions { return ui.Dimensions }

// UIPosition 组件偏移
type UIPosition struct {
	LayoutFace
}

// NewUIPosition 组件偏移
func NewUIPosition(face LayoutFace) *UIPosition {
	return &UIPosition{LayoutFace: face}
}

// Layout 绘制
func (ui *UIPosition) Layout(gtx layout.Context) *layout.Dimensions {
	defer op.Offset(*ui.GetPoint()).Push(ui.GetOps()).Pop()
	return ui.LayoutFace.Layout(gtx)
}

// UIScale 组件缩放实现
type UIScale struct {
	LayoutFace
	value        float32 // 记录缩放倍数
	f32.Affine2D         // 缩放矩阵
}

// NewUIScale 组件缩放实现
func NewUIScale(face LayoutFace, value float32) *UIScale {
	ui := &UIScale{LayoutFace: face}
	ui.SetScale(value)
	return ui
}

// SetScale 设置缩放值
func (ui *UIScale) SetScale(value float32) {
	ui.value = value
	ui.Update()
}

// Update 更新
func (ui *UIScale) Update() {
	ui.Affine2D = ui.Affine2D.Scale(f32.Point{X: 1, Y: 1}, f32.Point{X: ui.value, Y: ui.value})
	ui.LayoutFace.Update()
}

// GetScale 获取缩放值
func (ui *UIScale) GetScale() float32 {
	return ui.value
}

// Layout 绘制
func (ui *UIScale) Layout(gtx layout.Context) *layout.Dimensions {
	op.Affine(ui.Affine2D).Push(ui.GetOps())
	return ui.LayoutFace.Layout(gtx)
}

// UIGrid 背景网格
type UIGrid struct {
	LayoutFace
	Use       bool        // 启用网格
	Size      image.Point // 网格大小
	OffsetPop *f32.Point  // 网格偏移
	distance  int         // 网格大小
	stroke    [4]clip.Op  // 组件背景网格与边框
	clipRect  clip.Op     // 用于组件剪裁
}

// NewUIGrid 背景网格
func NewUIGrid(face LayoutFace, distanceValue int, Size image.Point) *UIGrid {
	ui := &UIGrid{LayoutFace: face, Use: true, Size: Size}
	ui.SetDistance(distanceValue)
	return ui
}

// SetDistance 设置背景网格大小
func (ui *UIGrid) SetDistance(distanceValue int) {
	ui.distance = distanceValue
	ui.Update()
}

// GetDistance 获取缩放值
func (ui *UIGrid) GetDistance() int {
	return ui.distance
}

// Update 更新
func (ui *UIGrid) Update() {
	dim := ui.GetDimensions()
	x := dim.Size.X/ui.distance + 1
	y := dim.Size.Y/ui.distance + 1
	// 绘制背景网格
	var gridsX clip.Path
	gridsX.Begin(&op.Ops{})
	var gridsXZ clip.Path
	gridsXZ.Begin(&op.Ops{})
	for i := 0; i < x; i++ {
		ix := float32(i * ui.distance)
		gridsX.MoveTo(f32.Pt(ix, 0))
		gridsX.LineTo(f32.Pt(ix, float32(dim.Size.Y)))
		iz := float32((i + 1) * ui.distance)
		for ; ix < iz; ix += 10 {
			gridsXZ.MoveTo(f32.Pt(ix, 0))
			gridsXZ.LineTo(f32.Pt(ix, float32(dim.Size.Y)))
		}
	}
	var gridsY clip.Path
	gridsY.Begin(&op.Ops{})
	var gridsYZ clip.Path
	gridsYZ.Begin(&op.Ops{})
	for i := 0; i < y; i++ {
		iy := float32(i * ui.distance)
		gridsY.MoveTo(f32.Pt(0, iy))
		gridsY.LineTo(f32.Pt(float32(dim.Size.X), iy))
		iz := float32((i + 1) * ui.distance)
		for ; iy < iz; iy += 10 {
			gridsY.MoveTo(f32.Pt(0, iy))
			gridsY.LineTo(f32.Pt(float32(dim.Size.X), iy))
		}
	}
	ui.stroke[0] = clip.Stroke{Path: gridsX.End(), Width: 1.5}.Op()
	ui.stroke[1] = clip.Stroke{Path: gridsXZ.End(), Width: 0.5}.Op()
	ui.stroke[2] = clip.Stroke{Path: gridsY.End(), Width: 0.5}.Op()
	ui.stroke[3] = clip.Stroke{Path: gridsYZ.End(), Width: 1}.Op()
	ui.clipRect = clip.Rect{Max: ui.GetDimensions().Size}.Op()
	ui.LayoutFace.Update()
}

// Layout 绘制
func (ui *UIGrid) Layout(gtx layout.Context) *layout.Dimensions {
	defer ui.clipRect.Push(gtx.Ops).Pop()
	if ui.Use && !ui.GetDisable() {
		ops := ui.GetOps()
		theme := ui.GetTheme()
		// 绘制背景线
		if ui.OffsetPop != nil {
			off := image.Point{X: -int(ui.OffsetPop.X * float32(ui.Size.X)), Y: -int(ui.OffsetPop.Y * float32(ui.Size.Y))}
			offest := op.Offset(image.Point{X: off.X % ui.distance}).Push(gtx.Ops)
			paint.FillShape(ops, theme.Fg, ui.stroke[0])
			paint.FillShape(ops, theme.ContrastBg, ui.stroke[1])
			offest.Pop()
			offest = op.Offset(image.Point{Y: off.Y % ui.distance}).Push(gtx.Ops)
			paint.FillShape(ops, theme.Fg, ui.stroke[2])
			paint.FillShape(ops, theme.ContrastBg, ui.stroke[3])
			offest.Pop()
			defer op.Offset(off).Push(gtx.Ops).Pop()
		} else {
			paint.FillShape(ops, theme.Fg, ui.stroke[0])
			paint.FillShape(ops, theme.ContrastBg, ui.stroke[1])
			paint.FillShape(ops, theme.Fg, ui.stroke[2])
			paint.FillShape(ops, theme.ContrastBg, ui.stroke[3])
		}
	}
	return ui.LayoutFace.Layout(gtx)
}

// UIFrame 边框
type UIFrame struct {
	LayoutFace
	width float32
	op    clip.Op
}

// NewUIFrame 边框
func NewUIFrame(face LayoutFace, width float32) *UIFrame {
	ui := &UIFrame{LayoutFace: face}
	ui.SetWidth(width)
	return ui
}

// UIFrame 设置边框大小
func (ui *UIFrame) SetWidth(width float32) {
	ui.width = width
	ui.Update()
}

// Update 更新
func (ui *UIFrame) Update() {
	dim := ui.GetDimensions()
	ui.op = clip.Stroke{Path: clip.RRect{Rect: image.Rectangle{Max: dim.Size}}.Path(&op.Ops{}), Width: ui.width}.Op()
	ui.LayoutFace.Update()
}

// Layout 绘制
func (ui *UIFrame) Layout(gtx layout.Context) *layout.Dimensions {
	theme := ui.GetTheme()
	paint.FillShape(ui.GetOps(), theme.Fg, ui.op)
	return ui.LayoutFace.Layout(gtx)
}

// UIScroll 滚动条
type UIScroll struct {
	LayoutFace
	Use            bool                    // 滚动条生效
	Axis           layout.Axis             // 滚动条方向
	Distance       float32                 // 滚动条滚动距离
	Scrollbar      widget.Scrollbar        // 滚动条对象
	ScrollbarStyle material.ScrollbarStyle // 滚动条对象风格
	off            image.Point             // 绘制位置
	call           func(Distance float32)  // 位置回调
}

// NewUIScroll 创建滚动条
func NewUIScroll(face LayoutFace, Axis layout.Axis, call func(Distance float32)) *UIScroll {
	ui := &UIScroll{LayoutFace: face}
	ui.ScrollbarStyle = material.Scrollbar(ui.LayoutFace.GetTheme(), &ui.Scrollbar)
	ui.Use = true
	ui.call = call
	ui.SetAxis(Axis)
	return ui
}

// SetAxis 设置坐标轴
func (ui *UIScroll) SetAxis(Axis layout.Axis) {
	ui.Axis = Axis
	ui.Update()
}

// SetAxis 设置回调
func (ui *UIScroll) SetCall(call func(Distance float32)) {
	ui.call = call
}

// Update 更新
func (ui *UIScroll) Update() {
	dim := ui.GetDimensions()
	if ui.Axis == layout.Horizontal {
		ui.off = image.Point{Y: dim.Size.Y - 10}
	} else {
		ui.off = image.Point{X: dim.Size.X - 10}
	}
	ui.LayoutFace.Update()
}

// Layout 绘制
func (ui *UIScroll) Layout(gtx layout.Context) *layout.Dimensions {
	if ui.Use && !ui.GetDisable() {
		dim := ui.GetDimensions()
		gtx.Ops = ui.GetOps()
		gtx.Constraints.Max.X = dim.Size.X - 5
		gtx.Constraints.Max.Y = dim.Size.Y - 5
		// 计算滚动条位置
		if ui.Scrollbar.Dragging() {
			ui.Distance += ui.Scrollbar.ScrollDistance()
			if ui.call != nil {
				ui.call(ui.Distance)
			}
		}
		TransformStack := op.Offset(ui.off).Push(gtx.Ops)
		ui.ScrollbarStyle.Layout(
			gtx,
			ui.Axis,
			ui.Distance,
			ui.Distance,
		)
		TransformStack.Pop()
	}
	return ui.LayoutFace.Layout(gtx)
}

// UIDrag 拖放
type UIDrag struct {
	LayoutFace
	Use          bool // 启用拖放
	dragging     bool
	position     f32.Point
	IsAdjustThe  UIDragAdjustThe // 大小调整
	IsAdjustMask UIDragAdjustThe
}

// NewUIDrag 创建拖动效果
func NewUIDrag(face LayoutFace) *UIDrag {
	ui := &UIDrag{LayoutFace: face, IsAdjustMask: 255}
	return ui
}

type UIDragAdjustThe uint8

const (
	UIDragTheW        UIDragAdjustThe = 1 << iota // 右
	UIDragTheE                                    // 左
	UIDragTheN                                    // 顶
	UIDragTheR                                    // 底
	UIDragTheRE                                   // 左下角
	UIDragTheDragging                             // 移动

)

// Layout 绘制
func (ui *UIDrag) Layout(gtx layout.Context) (dim *layout.Dimensions) {
	event.Op(gtx.Ops, ui)
	if ev, ok := gtx.Source.Event(pointer.Filter{
		Target: ui,
		Kinds:  pointer.Press | pointer.Release | pointer.Move | pointer.Drag | pointer.Cancel,
	}); ok {
		dim := ui.GetDimensions()
		point := ui.GetPoint()
		event := ev.(pointer.Event)
		pos := event.Position.Sub(f32.Point{X: float32(point.X), Y: float32(point.Y)})
		if event.Buttons == pointer.ButtonSecondary {
			ui.Use = true
			ui.SetDisable(true)
		} else if ((pos.X > 0 && pos.Y > 0 && pos.X < float32(dim.Size.X) && pos.Y < float32(dim.Size.Y)) && ui.Use) || ui.dragging {
			if !ui.dragging {
				//@ 处理鼠标所在点击的位置
				if pos.X+10 > float32(dim.Size.X) && pos.Y+10 > float32(dim.Size.Y) {
					ui.IsAdjustThe = UIDragTheRE & ui.IsAdjustMask
					pointer.CursorNorthWestResize.Add(gtx.Ops)
				} else if pos.X-10 < -5 {
					ui.IsAdjustThe = UIDragTheW & ui.IsAdjustMask
					pointer.CursorColResize.Add(gtx.Ops)
				} else if pos.X+10 > float32(dim.Size.X) {
					ui.IsAdjustThe = UIDragTheE & ui.IsAdjustMask
					pointer.CursorColResize.Add(gtx.Ops)
				} else if pos.Y+10 > float32(dim.Size.Y) {
					ui.IsAdjustThe = UIDragTheR & ui.IsAdjustMask
					pointer.CursorRowResize.Add(gtx.Ops)
				} else if pos.Y-10 < -5 {
					ui.IsAdjustThe = UIDragTheN & ui.IsAdjustMask
					pointer.CursorRowResize.Add(gtx.Ops)
				} else {
					ui.IsAdjustThe = UIDragTheDragging & ui.IsAdjustMask
					pointer.CursorGrabbing.Add(gtx.Ops)
				}
				// 处理事件
				if event.Buttons == pointer.ButtonPrimary && event.Kind == pointer.Press {
					switch ui.IsAdjustThe {
					case UIDragTheN:
						ui.dragging = true
						ui.position.Y = pos.Y
					case UIDragTheR:
						ui.dragging = true
						ui.position.Y = float32(dim.Size.Y) - pos.Y
					case UIDragTheE:
						ui.dragging = true
						ui.position.X = float32(dim.Size.X) - pos.X
					case UIDragTheW:
						ui.dragging = true
						ui.position.X = pos.X
					case UIDragTheRE:
						ui.dragging = true
						ui.position.Y = float32(dim.Size.Y) - pos.Y
						ui.position.X = float32(dim.Size.X) - pos.X
					case UIDragTheDragging:
						ui.dragging = true
						ui.position = pos
					}
				}
			} else {
				// 处理移动过程
				if event.Kind == pointer.Drag {
					r := pos.Sub(ui.position)
					switch ui.IsAdjustThe {
					case UIDragTheN:
						point.Y += int(r.Y)
						dim.Size.Y -= int(r.Y)
						pointer.CursorRowResize.Add(gtx.Ops)
					case UIDragTheR:
						dim.Size.Y = int(pos.Y - ui.position.Y)
						pointer.CursorRowResize.Add(gtx.Ops)
					case UIDragTheE:
						dim.Size.X = int(pos.X - ui.position.X)
						pointer.CursorColResize.Add(gtx.Ops)
					case UIDragTheW:
						point.X += int(r.X)
						dim.Size.X -= int(r.X)
						pointer.CursorColResize.Add(gtx.Ops)
					case UIDragTheRE:
						dim.Size.Y = int(pos.Y - ui.position.Y)
						dim.Size.X = int(pos.X - ui.position.X)
						pointer.CursorNorthWestResize.Add(gtx.Ops)
					case UIDragTheDragging:
						point.X += int(r.X)
						point.Y += int(r.Y)
						pointer.CursorGrabbing.Add(gtx.Ops)
					}
				} else {
					ui.Use = false
					ui.dragging = false
					ui.IsAdjustThe = 0
					ui.SetDisable(false)
					pointer.CursorDefault.Add(gtx.Ops)
				}
			}
		}
	}
	if ui.dragging {
		ui.LayoutFace.Update()

	}
	return ui.LayoutFace.Layout(gtx)
}
