package pmx

import "mmd/basebyte"

// Setting インデックス設定
type Setting struct {
	// エンコード方式
	Encoding uint8
	// 追加UV数
	Uv uint8
	// 頂点インデックスサイズ
	VertexIndexSize uint8
	// テクスチャインデックスサイズ
	TextureIndexSize uint8
	// マテリアルインデックスサイズ
	MaterialIndexSize uint8
	// ボーンインデックスサイズ
	BoneIndexSize uint8
	// モーフインデックスサイズ
	MorphIndexSize uint8
	// 剛体インデックスサイズ
	RigidbodyIndexSize uint8
}

// VertexSkinningType 頂点スキニングタイプ
type VertexSkinningType uint8

//
const (
	VertexBDEF1 VertexSkinningType = 0
	VertexBDEF2                    = 1
	VertexBDEF4                    = 2
	VertexSDEF                     = 3
	VertexQDEF                     = 4
)

// VertexSkinning 頂点スキニング
type VertexSkinning interface {
	read(r *basebyte.Read, setting *Setting)
}

// VertexSkinningBDEF1 BDEF1
type VertexSkinningBDEF1 struct {
	BoneIndex int
}

// VertexSkinningBDEF2 BDEF2
type VertexSkinningBDEF2 struct {
	BoneIndex1 int
	BoneIndex2 int
	BoneWeight float32
}

// VertexSkinningBDEF4 BDEF4
type VertexSkinningBDEF4 struct {
	BoneIndex1  int
	BoneIndex2  int
	BoneIndex3  int
	BoneIndex4  int
	BoneWeight1 float32
	BoneWeight2 float32
	BoneWeight3 float32
	BoneWeight4 float32
}

// VertexSkinningSDEF SDEF
type VertexSkinningSDEF struct {
	BoneIndex1 int
	BoneIndex2 int
	BoneWeight float32
	SdefC      [3]float32
	SdefR0     [3]float32
	SdefR1     [3]float32
}

// VertexSkinningQDEF QDEF
type VertexSkinningQDEF struct {
	BoneIndex1  int
	BoneIndex2  int
	BoneIndex3  int
	BoneIndex4  int
	BoneWeight1 float32
	BoneWeight2 float32
	BoneWeight3 float32
	BoneWeight4 float32
}

// Vertex 頂点
type Vertex struct {
	// 位置
	Positon [3]float32
	// 法線
	Normal [3]float32
	// テクスチャ座標
	Uv [2]float32
	// 追加テクスチャ座標
	Uva [4][4]float32
	// スキニングタイプ
	SkinningType VertexSkinningType
	// スキニング
	Skinning VertexSkinning
	// エッジ倍率
	Edge float32
}

// Material マテリアル
type Material struct {
	// モデル名
	MaterialName string
	// モデル英名
	MaterialEnglishName string
	// 減衰色
	Diffuse [4]float32
	// 光沢色
	Specular [3]float32
	// 光沢度
	Specularlity float32
	// 環境色
	Ambient [3]float32
	// 描画フラグ
	Flag uint8
	// エッジ色
	EdgeColor [4]float32
	// エッジサイズ
	EdgeSize float32
	// アルベドテクスチャインデックス
	DiffuseTextureIndex int
	// スフィアテクスチャインデックス
	SphereTextureIndex int
	// スフィアテクスチャ演算モード
	SphereOpMode uint8
	// 共有トゥーンフラグ
	CommonToonFlag uint8
	// トゥーンテクスチャインデックス
	ToonTextureIndex int
	// メモ
	Memo string
	// 頂点インデックス数
	IndexCount int
}

// IkLink リンク
type IkLink struct {
	// リンクボーンインデックス
	LinkTarget int
	// 角度制限
	AngleLock uint8
	// 最大制限角度
	MaxRadian [3]float32
	// 最小制限角度
	MinRadian [3]float32
}

// Bone ボーン
type Bone struct {
	// ボーン名
	BoneName string
	// ボーン英名
	BoneEnglishName string
	// 位置
	Position [3]float32
	// 親ボーンインデックス
	ParentIndex int
	// 階層
	Level int
	// ボーンフラグ
	BoneFlag uint16
	// 座標オフセット(has Target)
	Offset [3]float32
	// 接続先ボーンインデックス(not has Target)
	TargetIndex int
	// 付与親ボーンインデックス
	GrantParentIndex int
	// 付与率
	GrantWeight float32
	// 固定軸の方向
	LockAxisOrientation [3]float32
	// ローカル軸のX軸方向
	LocalAxisXOrientation [3]float32
	// ローカル軸のY軸方向
	LocalAxisYOrientation [3]float32
	// 外部親変形のkey値
	Key int
	// IKターゲットボーン
	IkTargetBoneIndex int
	// IKループ回数
	IkLoop int
	// IKループ計算時の角度制限(ラジアン)
	IkLoopAngleLimit float32
	// IKリンク数
	IkLinkCount int
	// IKリンク
	IkLinks []IkLink
}

// MorphType Morph类型
type MorphType uint8

//
const (
	MorphGroup         MorphType = 0
	MorphVertex                  = 1
	MorphBone                    = 2
	MorphUV                      = 3
	MorphAdditionalUV1           = 4
	MorphAdditionalUV2           = 5
	MorphAdditionalUV3           = 6
	MorphAdditionalUV4           = 7
	MorphMatrial                 = 8
	MorphFlip                    = 9
	MorphImpulse                 = 10
)

// MorphCategory Morph种类
type MorphCategory uint8

//
const (
	MorphReservedCategory MorphCategory = 0
	MorphEyebrow                        = 1
	MorphEye                            = 2
	MorphMouth                          = 3
	MorphOther                          = 4
)

// MorphOffset 变形偏移量
type MorphOffset interface {
	read(r *basebyte.Read, setting *Setting)
}

// MorphVertexOffset 变形顶点偏移量
type MorphVertexOffset struct {
	VertexIndex    int
	PositionOffset [3]float32
}

// MorphUVOffset 变形UV偏移量
type MorphUVOffset struct {
	VertexIndex int
	UvOffset    [4]float32
}

// MorphBoneOffset 变形骨偏移
type MorphBoneOffset struct {
	Boneindex   int
	Translation [3]float32
	Rotation    [4]float32
}

// MorphMaterialOffset 变形材料偏移量
type MorphMaterialOffset struct {
	MaterialIndex     int
	OffsetOperation   uint8
	Diffuse           [4]float32
	Specular          [3]float32
	Specularity       float32
	Ambient           [3]float32
	EdgeColor         [4]float32
	EdgeSize          float32
	TextureArgb       [4]float32
	SphereTextureArgb [4]float32
	ToonTextureArgb   [4]float32
}

// MorphGroupOffset 变形组偏移
type MorphGroupOffset struct {
	MorphIndex  int
	MorphWeight float32
}

// MorphFlipOffset 变形反转偏移量
type MorphFlipOffset struct {
	MorphIndex int
	MorphValue float32
}

// MorphImpulseOffset 变形Impulse偏移量
type MorphImpulseOffset struct {
	RigidBodyIndex int
	IsLocal        uint8
	Velocity       [3]float32
	AngularTorque  [3]float32
}

// Morph モーフ
type Morph struct {
	// モーフ名
	MorphName string
	// モーフ英名
	MorphEnglishName string
	// カテゴリ
	Category MorphCategory
	// モーフタイプ
	MorphType MorphType
	// オフセット数
	OffsetCount int
	// 頂点モーフ配列
	VertexOffsets []MorphVertexOffset
	// UVモーフ配列
	UvOffsets []MorphUVOffset
	// ボーンモーフ配列
	BoneOffsets []MorphBoneOffset
	// マテリアルモーフ配列
	MaterialOffsets []MorphMaterialOffset
	// グループモーフ配列
	GroupOffsets []MorphGroupOffset
	// フリップモーフ配列
	FlipOffsets []MorphFlipOffset
	// インパルスモーフ配列
	ImpulseOffsets []MorphImpulseOffset
}

// FrameElement 枠内要素
type FrameElement struct {
	// 要素対象
	ElementTarget uint8
	// 要素対象インデックス
	Index int
}

// Frame 表示枠
type Frame struct {
	// 枠名
	FrameName string
	// 枠英名
	FrameEnglishName string
	// 特殊枠フラグ
	FrameFlag uint8
	// 枠内要素数
	ElementCount int
	// 枠内要素配列
	Elements []FrameElement
}

// RigidBody 刚体
type RigidBody struct {
	// 剛体名
	RigidBodyName string
	// 剛体英名
	RigidBodyEnglishName string
	// 関連ボーンインデックス
	TargetBone int
	// グループ
	Group uint8
	// マスク
	Mask uint16
	// 形状
	Shape               uint8
	Size                [3]float32
	Position            [3]float32
	Orientation         [3]float32
	Mass                float32
	MoveAttenuation     float32
	RotationAttenuation float32
	Repulsion           float32
	Friction            float32
	PhysicsCalcType     uint8
}

// JointType 链接类型
type JointType uint8

//
const (
	JointGeneric6DofSpring JointType = 0
	JointGeneric6Dof                 = 1
	JointPoint2Point                 = 2
	JointConeTwist                   = 3
	JointSlider                      = 5
	JointHinge                       = 6
)

// JointParam 链接参数
type JointParam struct {
	RigidBody1                int
	RigidBody2                int
	Position                  [3]float32
	Orientaiton               [3]float32
	MoveLimitationMin         [3]float32
	MoveLimitationMax         [3]float32
	RotationLimitationMin     [3]float32
	RotationLimitationMax     [3]float32
	SpringMoveCoefficient     [3]float32
	SpringRotationCoefficient [3]float32
}

// Joint 链接
type Joint struct {
	JointName        string
	JointEnglishName string
	JointType        JointType
	Param            JointParam
}

// SoftBodyFlag 软体标记
type SoftBodyFlag uint8

//
const (
	SoftBodyBLink   SoftBodyFlag = 0x01
	SoftBodyCluster              = 0x02
	SoftBodyLink                 = 0x04
)

// AnchorRigidBody 锚固刚体
type AnchorRigidBody struct {
	RelatedRigidBody int
	RelatedVertex    int
	IsNear           bool
}

// SoftBody 软体
type SoftBody struct {
	SoftBodyName        string
	SoftBodyEnglishName string
	Shape               uint8
	TargetMaterial      int
	Group               uint8
	Mask                uint16
	Flag                SoftBodyFlag
	BlinkDistance       int
	ClusterCount        int
	Mass                float32
	CollisioniMargin    float32
	AeroModel           int
	VCF                 float32
	DP                  float32
	DG                  float32
	LF                  float32
	PR                  float32
	VC                  float32
	DF                  float32
	MT                  float32
	CHR                 float32
	KHR                 float32
	SHR                 float32
	AHR                 float32
	SRHRCL              float32
	SKHRCL              float32
	SSHRCL              float32
	SRSPLTCL            float32
	SKSPLTCL            float32
	SSSPLTCL            float32
	VIT                 int
	PIT                 int
	DIT                 int
	CIT                 int
	LST                 float32
	AST                 float32
	VST                 float32
	AnchorCount         int
	Anchors             []AnchorRigidBody
	PinVertexCount      int
	pinVertices         []int
}

// Model モデル
type Model struct {
	// バージョン
	Version float32
	// 設定
	Setting Setting
	// モデル名
	ModelName string
	// モデル英名
	ModelEnglishName string
	// コメント
	ModelComment string
	// 英語コメント
	ModelEnglishComment string
	// 頂点数
	VertexCount int
	// 頂点配列
	Vertices []Vertex
	// インデックス数
	IndexCount int
	// インデックス配列
	Indices []int
	// テクスチャ数
	TextureCount int
	// テクスチャ配列
	Textures []string
	// マテリアル数
	MaterialCount int
	// マテリアル
	Materials []Material
	// ボーン数
	BoneCount int
	// ボーン配列
	Bones []Bone
	// モーフ数
	MorphCount int
	// モーフ配列
	Morphs []Morph
	// 表示枠数
	FrameCount int
	// 表示枠配列
	Frames []Frame
	// 剛体数
	RigidBodyCount int
	// 剛体配列
	RigidBodies []RigidBody
	// ジョイント数
	JointCount int
	// ジョイント配列
	Joints []Joint
	// ソフトボディ数
	SoftBodyCount int
	// ソフトボディ配列
	SoftBodies []SoftBody
}
