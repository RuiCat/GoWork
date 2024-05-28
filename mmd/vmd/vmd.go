package vmd

// BoneFrame ボーンフレーム
type BoneFrame struct {
	// ボーン名
	Name string
	// フレーム番号
	Frame int
	// 位置
	Position [3]float32
	// 回転
	Orientation [4]float32
	// 補間曲線
	Interpolation [4][4][4]byte
}

// FaceFrame 表情フレーム
type FaceFrame struct {
	// 表情名
	FaceName string
	// 表情の重み
	Weight float32
	// フレーム番号
	Frame uint32
}

// CameraFrame カメラフレーム
type CameraFrame struct {
	// フレーム番号
	Frame int
	// 距離
	Distance float32
	// 位置
	Position [3]float32
	// 回転
	Orientation [3]float32
	// 補間曲線
	Interpolation [6][4]byte
	// 視野角
	Angle float32
	// 不明データ
	Unknown [3]byte
}

// LightFrame ライトフレーム
type LightFrame struct {
	// フレーム番号
	Frame int
	// 色
	Color [3]float32
	// 位置
	Position [3]float32
}

// IkEnable IKの有効無効
type IkEnable struct {
	IkName string
	Enable bool
}

// IkFrame IKフレーム
type IkFrame struct {
	Frame    int
	Display  bool
	IkEnable []IkEnable
}

// Motion VMDモーション
type Motion struct {
	// モデル名
	ModelName string
	// バージョン
	Version string
	// ボーンフレーム
	BoneFrames []BoneFrame
	// 表情フレーム
	FaceFrames []FaceFrame
	// カメラフレーム
	CameraFrames []CameraFrame
	// ライトフレーム
	LightFrames []LightFrame
	// IKフレーム
	IkFrames []IkFrame
}
