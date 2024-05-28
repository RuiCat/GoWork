package pmd

// Header ヘッダ
type Header struct {
	// モデル名
	Name string
	// モデル名(英語)
	NameEnglish string
	// コメント
	Comment string
	// コメント(英語)
	CommentEnglish string
}

// Vertex 頂点
type Vertex struct {
	// 位置
	Position [3]float32
	// 法線
	Normal [3]float32
	// UV座標
	Uv [2]float32
	// 関連ボーンインデックス
	BoneIndex [2]uint16
	// ボーンウェイト
	BoneWeight uint8
	// エッジ不可視
	EdgeInvisible bool
}

// Material 材質
type Material struct {
	// 減衰色
	Diffuse [4]float32
	// 光沢度
	Power float32
	// 光沢色
	Specular [3]float32
	// 環境色
	Ambient [3]float32
	// トーンインデックス
	ToonIndex uint8
	// エッジ
	EdgeFlag uint8
	// インデックス数
	IndexCount uint32
	// テクスチャファイル名
	TextureFilename string
	// スフィアファイル名
	SphereFilename string
}

// BoneType 骨骼类型
type BoneType uint8

// 骨骼类型
const (
	Rotation BoneType = iota
	RotationAndMove
	IkEffector
	Unknown
	IkEffectable
	RotationEffectable
	IkTarget
	Invisible
	Twist
	RotationMovement
)

// Bone ボーン
type Bone struct {
	// ボーン名
	Name string
	// ボーン名(英語)
	NameEnglish string
	// 親ボーン番号
	ParentBoneindex uint16
	// 末端ボーン番号
	TailPosboneindex uint16
	// ボーン種類
	BoneType BoneType
	// IKボーン番号
	IkParentboneindex uint16
	// ボーンのヘッドの位置
	BoneHeadpos [3]float32
}

// Ik IK
type Ik struct {
	// IKボーン番号
	IkBoneindex uint16
	// IKターゲットボーン番号
	TargetBoneindex uint16
	// 再帰回数
	Interations uint16
	// 角度制限
	AngleLimit float32
	// 影響下ボーン番号
	IkChildboneindex []uint16
}

// FaceVertex 面顶点
type FaceVertex struct {
	VertexIndex int
	Position    [3]float32
}

// FaceCategory 面类型
type FaceCategory uint8

// 面类型
const (
	Base FaceCategory = iota
	Eyebrow
	Eye
	Mouth
	Other
)

// Face 面
type Face struct {
	Name        string
	Type        FaceCategory
	Vertices    []FaceVertex
	NameEnglish string
}

// BoneDispName ボーン枠用の枠名
type BoneDispName struct {
	BoneDispname        string
	BoneDispnameenglish string
}

// BoneDisp 骨骼Disp
type BoneDisp struct {
	BoneIndex     uint16
	BoneDispindex uint8
}

// RigidBodyShape 衝突形状
type RigidBodyShape uint8

// 刚体形状
const (
	// 球
	Sphere RigidBodyShape = 0
	// 直方体
	Box RigidBodyShape = 1
	// カプセル
	Cpusel RigidBodyShape = 2
)

// RigidBodyType 剛体タイプ
type RigidBodyType uint8

// 刚体类型
const (
	// ボーン追従
	BoneConnected RigidBodyType = 0
	// 物理演算
	Physics RigidBodyType = 1
	// 物理演算(Bone位置合せ)
	ConnectedPhysics RigidBodyType = 2
)

// RigidBody 剛体
type RigidBody struct {
	// 名前
	Name string
	// 関連ボーン番号
	RelatedBoneindex uint16
	// グループ番号
	GroupIndex uint8
	// マスク
	Mask uint16
	// 形状
	Shape RigidBodyShape
	// 大きさ
	Size [3]float32
	// 位置
	Position [3]float32
	// 回転
	Orientation [3]float32
	// 質量
	Weight float32
	// 移動ダンピング
	LinearDamping float32
	// 回転ダンピング
	AnglarDamping float32
	// 反発係数
	Restitution float32
	// 摩擦係数
	Friction float32
	// 演算方法
	RigidType RigidBodyType
}

// Constraint 剛体の拘束
type Constraint struct {
	// 名前
	Name string
	// 剛体Aのインデックス
	RigidBodyindexa uint32
	// 剛体Bのインデックス
	RigidBodyindexb uint32
	// 位置
	Position [3]float32
	// 回転
	Orientation [3]float32
	// 最小移動制限
	LinearLowerlimit [3]float32
	// 最大移動制限
	LinearUpperlimit [3]float32
	// 最小回転制限
	AngularLowerlimit [3]float32
	// 最大回転制限
	AngularUpperlimit [3]float32
	// 移動に対する復元力
	LinearStiffness [3]float32
	// 回転に対する復元力
	AngularStiffness [3]float32
}

// Model モデル
type Model struct {
	Version       float32
	Header        Header
	Vertices      []Vertex
	Indices       []uint16
	Materials     []Material
	Bones         []Bone
	Iks           []Ik
	Faces         []Face
	Facesindices  []uint16
	Bonedispname  []BoneDispName
	Bonedisp      []BoneDisp
	Toonfilenames []string
	Rigidbodies   []RigidBody
	Constraints   []Constraint
}
