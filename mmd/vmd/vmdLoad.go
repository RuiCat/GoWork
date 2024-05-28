package vmd

import (
	"encoding/binary"
	"errors"
	"io/ioutil"
	"mmd/basebyte"
)

func ifnil(b []byte) []byte {
	for i, n := 0, len(b); i < n; i++ {
		if b[i] == 0x00 {
			return b[:i]
		}
	}
	return b
}
func getF3(r *basebyte.Read) [3]float32 {
	return [3]float32{r.Float32(), r.Float32(), r.Float32()}
}
func getF4(r *basebyte.Read) [4]float32 {
	return [4]float32{r.Float32(), r.Float32(), r.Float32(), r.Float32()}
}
func setF3(v [3]float32, w *basebyte.Write) {
	w.Float32(v[0])
	w.Float32(v[1])
	w.Float32(v[2])
}
func setF4(v [4]float32, w *basebyte.Write) {
	w.Float32(v[0])
	w.Float32(v[1])
	w.Float32(v[2])
	w.Float32(v[3])
}

func (vmd *BoneFrame) Read(r *basebyte.Read) {
	vmd.Name = string(r.GetByte(15))
	vmd.Frame = int(r.Int32())
	vmd.Position = getF3(r)
	vmd.Orientation = getF4(r)
	n := 0
	data := r.GetByte(64)
	for i1 := 0; i1 < 4; i1++ {
		for i2 := 0; i2 < 4; i2++ {
			for i3 := 0; i3 < 4; i3++ {
				vmd.Interpolation[i1][i2][i3] = data[n]
				n++
			}
		}
	}
}

func (vmd *BoneFrame) Write(w *basebyte.Write) {
	name := make([]byte, 15)
	copy(name, []byte([]byte(vmd.Name)))
	w.SetByte(name)
	w.Int32(int32(vmd.Frame))
	setF3(vmd.Position, w)
	setF4(vmd.Orientation, w)

	n := 0
	var data [64]byte
	for i1 := 0; i1 < 4; i1++ {
		for i2 := 0; i2 < 4; i2++ {
			for i3 := 0; i3 < 4; i3++ {
				data[n] = vmd.Interpolation[i1][i2][i3]
				n++
			}
		}
	}
	w.SetByte(data[:])
}

func (vmd *FaceFrame) Read(r *basebyte.Read) {
	vmd.FaceName = string(r.GetByte(15))
	vmd.Frame = r.Uint32()
	vmd.Weight = r.Float32()
}

func (vmd *FaceFrame) Write(w *basebyte.Write) {
	name := make([]byte, 15)
	copy(name, []byte([]byte(vmd.FaceName)))
	w.SetByte(name)
	w.Uint32(vmd.Frame)
	w.Float32(vmd.Weight)
}

func (vmd *CameraFrame) Read(r *basebyte.Read) {
	vmd.Frame = int(r.Int32())
	vmd.Distance = r.Float32()
	vmd.Position = getF3(r)
	vmd.Orientation = getF3(r)
	n := 0
	data := r.GetByte(24)
	for i1 := 0; i1 < 6; i1++ {
		for i2 := 0; i2 < 4; i2++ {
			vmd.Interpolation[i1][i2] = data[n]
			n++
		}
	}
	vmd.Angle = r.Float32()
	copy(vmd.Unknown[0:], r.GetByte(3))
}

func (vmd *CameraFrame) Write(w *basebyte.Write) {
	w.Int32(int32(vmd.Frame))
	w.Float32(vmd.Distance)
	setF3(vmd.Position, w)
	setF3(vmd.Orientation, w)
	n := 0
	var data [24]byte
	for i1 := 0; i1 < 6; i1++ {
		for i2 := 0; i2 < 4; i2++ {
			data[n] = vmd.Interpolation[i1][i2]
			n++
		}
	}
	w.SetByte(data[0:])
	w.Float32(vmd.Angle)
	w.SetByte(vmd.Unknown[0:])
}

func (vmd *LightFrame) Read(r *basebyte.Read) {
	vmd.Frame = int(r.Int32())
	vmd.Color = getF3(r)
	vmd.Position = getF3(r)
}

func (vmd *LightFrame) Write(w *basebyte.Write) {
	w.Int32(int32(vmd.Frame))
	setF3(vmd.Color, w)
	setF3(vmd.Position, w)
}

func (vmd *IkFrame) Read(r *basebyte.Read) {
	vmd.Frame = int(r.Int32())
	vmd.Display = r.Bool()
	ikCount := int(r.Int32())
	vmd.IkEnable = make([]IkEnable, ikCount)
	for i := 0; i < ikCount; i++ {
		vmd.IkEnable[i].IkName = string(r.GetByte(20))
		vmd.IkEnable[i].Enable = r.Bool()
	}
}

func (vmd *IkFrame) Write(w *basebyte.Write) {
	w.Int32(int32(vmd.Frame))
	w.Bool(vmd.Display)
	ikCount := len(vmd.IkEnable)
	w.Int32(int32(ikCount))
	var buffer [20]byte
	for i := 0; i < ikCount; i++ {
		copy(buffer[0:], []byte(vmd.IkEnable[i].IkName))
		w.SetByte(buffer[:])
		w.Bool(vmd.IkEnable[i].Enable)
	}
}

// LoadFromFile 从文件加载
func (vmd *Motion) LoadFromFile(filename string) error {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return vmd.LoadFromStream(f)
}

// LoadFromStream 从流加载
func (vmd *Motion) LoadFromStream(stream []byte) error {
	if len(stream) == 0 {
		return errors.New("invalid stream")
	}
	// 读取器
	r := &basebyte.Read{Byte: &stream, Order: binary.LittleEndian}

	// magic and version
	buffer := r.GetByte(30)
	if string(buffer[:20]) != "Vocaloid Motion Data" {
		return errors.New("invalid vmd file")
	}
	vmd.Version = string(ifnil(buffer[20:]))

	// name
	vmd.ModelName = string(ifnil((r.GetByte(20))))

	// bone frames
	num := int(r.Int32())
	vmd.BoneFrames = make([]BoneFrame, num)
	for i := 0; i < num; i++ {
		vmd.BoneFrames[i].Read(r)
	}

	// face frames
	num = int(r.Int32())
	vmd.FaceFrames = make([]FaceFrame, num)
	for i := 0; i < num; i++ {
		vmd.FaceFrames[i].Read(r)
	}

	// camera frames
	num = int(r.Int32())
	vmd.CameraFrames = make([]CameraFrame, num)
	for i := 0; i < num; i++ {
		vmd.CameraFrames[i].Read(r)
	}

	// light frames
	num = int(r.Int32())
	vmd.LightFrames = make([]LightFrame, num)
	for i := 0; i < num; i++ {
		vmd.LightFrames[i].Read(r)
	}

	// unknown2
	r.GetByte(4)

	// ik frames
	if r.Offset < len(*r.Byte) {
		num = int(r.Int32())
		vmd.IkFrames = make([]IkFrame, num)
		for i := 0; i < num; i++ {
			vmd.IkFrames[i].Read(r)
		}
	}
	if r.Offset < len(*r.Byte) {
		return errors.New("vmd stream has unknown data")
	}
	return nil
}

// SaveToFile 导出到文件
func (vmd *Motion) SaveToFile(filename string) error {
	return ioutil.WriteFile(filename, vmd.SaveToStream(), 0666)
}

// SaveToStream  导出到流
func (vmd *Motion) SaveToStream() []byte {
	stream := make([]byte, 0)
	w := &basebyte.Write{Byte: &stream, Order: binary.LittleEndian}
	// magic and version
	var magic [30]byte
	copy(magic[0:], []byte(`Vocaloid Motion Data 0002`))
	w.SetByte(magic[:])
	// name
	copy(magic[0:], []byte(vmd.ModelName))
	w.SetByte(magic[:20])
	// bone frames
	w.Int32(int32(len(vmd.BoneFrames)))
	for _, v := range vmd.BoneFrames {
		v.Write(w)
	}
	// face frames
	w.Int32(int32(len(vmd.FaceFrames)))
	for _, v := range vmd.FaceFrames {
		v.Write(w)
	}
	// camera frames
	w.Int32(int32(len(vmd.CameraFrames)))
	for _, v := range vmd.CameraFrames {
		v.Write(w)
	}
	// light frames
	w.Int32(int32(len(vmd.LightFrames)))
	for _, v := range vmd.LightFrames {
		v.Write(w)
	}
	// self shadow datas
	var SelfShadowNum int32 = 0
	w.Int32(SelfShadowNum)
	// ik frames
	w.Int32(int32(len(vmd.IkFrames)))
	for _, v := range vmd.IkFrames {
		v.Write(w)
	}
	return stream
}
