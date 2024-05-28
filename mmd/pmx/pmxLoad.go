package pmx

import (
	"encoding/binary"
	"errors"
	"mmd/basebyte"
	"os"
	"unicode/utf16"
)

func getF3(r *basebyte.Read) [3]float32 {
	return [3]float32{r.Float32(), r.Float32(), r.Float32()}
}
func getF4(r *basebyte.Read) [4]float32 {
	return [4]float32{r.Float32(), r.Float32(), r.Float32(), r.Float32()}
}

// ReadIndex インデックス値を読み込む
func ReadIndex(r *basebyte.Read, size uint8) int {
	switch size {
	case 1:
		tmp8 := r.Uint8()
		if 255 == tmp8 {
			return -1
		}
		return int(tmp8)
	case 2:
		tmp16 := r.Uint16()
		if 65535 == tmp16 {
			return -1
		}
		return int(tmp16)
	case 4:
		return int(r.Int32())
	default:
		return -1
	}
}

// ReadString 文字列を読み込む
func ReadString(r *basebyte.Read, encoding uint8) string {
	size := int(r.Int32())
	if size == 0 {
		return ""
	}
	if encoding == 0 {
		// 得到uint16字符
		size = size / 2
		buf := make([]uint16, size)
		for i := 0; i < size; i++ {
			buf[i] = r.Uint16()
		}
		return string(utf16.Decode(buf))
	}
	return string(r.GetByte(size))
}

func (pmx *Setting) read(r *basebyte.Read) {
	count := uint(r.Uint8())
	if count < 8 {
		return
	}
	pmx.Encoding = r.Uint8()
	pmx.Uv = r.Uint8()
	pmx.VertexIndexSize = r.Uint8()
	pmx.TextureIndexSize = r.Uint8()
	pmx.MaterialIndexSize = r.Uint8()
	pmx.BoneIndexSize = r.Uint8()
	pmx.MorphIndexSize = r.Uint8()
	pmx.RigidbodyIndexSize = r.Uint8()
	for i := uint(8); i < count; i++ {
		r.Uint8()
	}
}

func (pmx *VertexSkinningBDEF1) read(r *basebyte.Read, setting *Setting) {
	pmx.BoneIndex = ReadIndex(r, setting.BoneIndexSize)
}
func (pmx *VertexSkinningBDEF2) read(r *basebyte.Read, setting *Setting) {
	pmx.BoneIndex1 = ReadIndex(r, setting.BoneIndexSize)
	pmx.BoneIndex2 = ReadIndex(r, setting.BoneIndexSize)
	pmx.BoneWeight = r.Float32()
}
func (pmx *VertexSkinningBDEF4) read(r *basebyte.Read, setting *Setting) {
	pmx.BoneIndex1 = ReadIndex(r, setting.BoneIndexSize)
	pmx.BoneIndex2 = ReadIndex(r, setting.BoneIndexSize)
	pmx.BoneIndex3 = ReadIndex(r, setting.BoneIndexSize)
	pmx.BoneIndex4 = ReadIndex(r, setting.BoneIndexSize)
	pmx.BoneWeight1 = r.Float32()
	pmx.BoneWeight2 = r.Float32()
	pmx.BoneWeight3 = r.Float32()
	pmx.BoneWeight4 = r.Float32()
}
func (pmx *VertexSkinningSDEF) read(r *basebyte.Read, setting *Setting) {
	pmx.BoneIndex1 = ReadIndex(r, setting.BoneIndexSize)
	pmx.BoneIndex2 = ReadIndex(r, setting.BoneIndexSize)
	pmx.BoneWeight = r.Float32()
	pmx.SdefC = getF3(r)
	pmx.SdefR0 = getF3(r)
	pmx.SdefR1 = getF3(r)
}
func (pmx *VertexSkinningQDEF) read(r *basebyte.Read, setting *Setting) {
	pmx.BoneIndex1 = ReadIndex(r, setting.BoneIndexSize)
	pmx.BoneIndex2 = ReadIndex(r, setting.BoneIndexSize)
	pmx.BoneIndex3 = ReadIndex(r, setting.BoneIndexSize)
	pmx.BoneIndex4 = ReadIndex(r, setting.BoneIndexSize)
	pmx.BoneWeight1 = r.Float32()
	pmx.BoneWeight2 = r.Float32()
	pmx.BoneWeight3 = r.Float32()
	pmx.BoneWeight4 = r.Float32()
}

func (pmx *Vertex) read(r *basebyte.Read, setting *Setting) {
	pmx.Positon = getF3(r)
	pmx.Normal = getF3(r)
	pmx.Uv[0] = r.Float32()
	pmx.Uv[1] = r.Float32()
	for i := uint8(0); i < setting.Uv; i++ {
		pmx.Uva[i] = getF4(r)
	}
	pmx.SkinningType = VertexSkinningType(r.Uint8())
	switch pmx.SkinningType {
	case VertexBDEF1:
		pmx.Skinning = &VertexSkinningBDEF1{}
	case VertexBDEF2:
		pmx.Skinning = &VertexSkinningBDEF2{}
	case VertexBDEF4:
		pmx.Skinning = &VertexSkinningBDEF4{}
	case VertexSDEF:
		pmx.Skinning = &VertexSkinningSDEF{}
	case VertexQDEF:
		pmx.Skinning = &VertexSkinningQDEF{}
	default:
		return
	}
	pmx.Skinning.read(r, setting)
	pmx.Edge = r.Float32()
}

func (pmx *Material) read(r *basebyte.Read, setting *Setting) {
	pmx.MaterialName = ReadString(r, setting.Encoding)
	pmx.MaterialEnglishName = ReadString(r, setting.Encoding)
	pmx.Diffuse = getF4(r)
	pmx.Specular = getF3(r)
	pmx.Specularlity = r.Float32()
	pmx.Ambient = getF3(r)
	pmx.Flag = r.Uint8()
	pmx.EdgeColor = getF4(r)
	pmx.EdgeSize = r.Float32()
	pmx.DiffuseTextureIndex = ReadIndex(r, setting.TextureIndexSize)
	pmx.SphereTextureIndex = ReadIndex(r, setting.TextureIndexSize)
	pmx.SphereOpMode = r.Uint8()
	pmx.CommonToonFlag = r.Uint8()
	if pmx.CommonToonFlag != 0 {
		pmx.ToonTextureIndex = int(r.Uint8())
	} else {
		pmx.ToonTextureIndex = ReadIndex(r, setting.TextureIndexSize)
	}
	pmx.Memo = ReadString(r, setting.Encoding)
	pmx.IndexCount = int(r.Int32())
}

func (pmx *IkLink) read(r *basebyte.Read, setting *Setting) {
	pmx.LinkTarget = ReadIndex(r, setting.BoneIndexSize)
	pmx.AngleLock = r.Uint8()
	if pmx.AngleLock == 1 {
		pmx.MaxRadian = getF3(r)
		pmx.MinRadian = getF3(r)
	}
}

func (pmx *Bone) read(r *basebyte.Read, setting *Setting) {
	pmx.BoneName = ReadString(r, setting.Encoding)
	pmx.BoneEnglishName = ReadString(r, setting.Encoding)
	pmx.Position = getF3(r)
	pmx.ParentIndex = ReadIndex(r, setting.BoneIndexSize)
	pmx.Level = int(r.Int32())
	pmx.BoneFlag = r.Uint16()
	if (pmx.BoneFlag & 0x0001) != 0 {
		pmx.TargetIndex = ReadIndex(r, setting.BoneIndexSize)
	} else {
		pmx.Offset = getF3(r)
	}
	if (pmx.BoneFlag & (0x0100 | 0x0200)) != 0 {
		pmx.GrantParentIndex = ReadIndex(r, setting.BoneIndexSize)
		pmx.GrantWeight = r.Float32()
	}
	if (pmx.BoneFlag & 0x0400) != 0 {
		pmx.LockAxisOrientation = getF3(r)
	}
	if (pmx.BoneFlag & 0x0800) != 0 {
		pmx.LocalAxisXOrientation = getF3(r)
		pmx.LocalAxisYOrientation = getF3(r)
	}
	if (pmx.BoneFlag & 0x2000) != 0 {
		pmx.Key = int(r.Int32())
	}
	if (pmx.BoneFlag & 0x0020) > 1 {
		pmx.IkTargetBoneIndex = ReadIndex(r, setting.BoneIndexSize)
		pmx.IkLoop = int(r.Int32())
		pmx.IkLoopAngleLimit = r.Float32()
		pmx.IkLinkCount = int(r.Int32())
		pmx.IkLinks = make([]IkLink, pmx.IkLinkCount)
		for i := 0; i < pmx.IkLinkCount; i++ {
			pmx.IkLinks[i].read(r, setting)
		}
	}
}

func (pmx *MorphVertexOffset) read(r *basebyte.Read, setting *Setting) {
	pmx.VertexIndex = ReadIndex(r, setting.VertexIndexSize)
	pmx.PositionOffset = getF3(r)
}
func (pmx *MorphUVOffset) read(r *basebyte.Read, setting *Setting) {
	pmx.VertexIndex = ReadIndex(r, setting.VertexIndexSize)
	pmx.UvOffset = getF4(r)
}
func (pmx *MorphBoneOffset) read(r *basebyte.Read, setting *Setting) {
	pmx.Boneindex = ReadIndex(r, setting.BoneIndexSize)
	pmx.Translation = getF3(r)
	pmx.Rotation = getF4(r)
}
func (pmx *MorphMaterialOffset) read(r *basebyte.Read, setting *Setting) {
	pmx.MaterialIndex = ReadIndex(r, setting.MorphIndexSize)
	pmx.OffsetOperation = r.Uint8()
	pmx.Diffuse = getF4(r)
	pmx.Specular = getF3(r)
	pmx.Specularity = r.Float32()
	pmx.Ambient = getF3(r)
	pmx.EdgeColor = getF4(r)
	pmx.EdgeSize = r.Float32()
	pmx.TextureArgb = getF4(r)
	pmx.SphereTextureArgb = getF4(r)
	pmx.ToonTextureArgb = getF4(r)
}
func (pmx *MorphGroupOffset) read(r *basebyte.Read, setting *Setting) {
	pmx.MorphIndex = ReadIndex(r, setting.MorphIndexSize)
	pmx.MorphWeight = r.Float32()
}
func (pmx *MorphFlipOffset) read(r *basebyte.Read, setting *Setting) {
	pmx.MorphIndex = ReadIndex(r, setting.MorphIndexSize)
	pmx.MorphValue = r.Float32()
}
func (pmx *MorphImpulseOffset) read(r *basebyte.Read, setting *Setting) {
	pmx.RigidBodyIndex = ReadIndex(r, setting.MorphIndexSize)
	pmx.IsLocal = r.Uint8()
	pmx.Velocity = getF3(r)
	pmx.AngularTorque = getF3(r)
}

func (pmx *Morph) read(r *basebyte.Read, setting *Setting) {
	pmx.MorphName = ReadString(r, setting.Encoding)
	pmx.MorphEnglishName = ReadString(r, setting.Encoding)
	pmx.Category = MorphCategory(r.Uint8())
	pmx.MorphType = MorphType(r.Uint8())
	pmx.OffsetCount = int(r.Int32())
	switch pmx.MorphType {
	case MorphGroup:
		pmx.GroupOffsets = make([]MorphGroupOffset, pmx.OffsetCount)
		for i := 0; i < pmx.OffsetCount; i++ {
			pmx.GroupOffsets[i].read(r, setting)
		}
	case MorphVertex:
		pmx.VertexOffsets = make([]MorphVertexOffset, pmx.OffsetCount)
		for i := 0; i < pmx.OffsetCount; i++ {
			pmx.VertexOffsets[i].read(r, setting)
		}
	case MorphBone:
		pmx.BoneOffsets = make([]MorphBoneOffset, pmx.OffsetCount)
		for i := 0; i < pmx.OffsetCount; i++ {
			pmx.BoneOffsets[i].read(r, setting)
		}
	case MorphMatrial:
		pmx.MaterialOffsets = make([]MorphMaterialOffset, pmx.OffsetCount)
		for i := 0; i < pmx.OffsetCount; i++ {
			pmx.MaterialOffsets[i].read(r, setting)
		}
	case MorphUV,
		MorphAdditionalUV1,
		MorphAdditionalUV2,
		MorphAdditionalUV3,
		MorphAdditionalUV4:
		pmx.UvOffsets = make([]MorphUVOffset, pmx.OffsetCount)
		for i := 0; i < pmx.OffsetCount; i++ {
			pmx.UvOffsets[i].read(r, setting)
		}
	}
}

func (pmx *FrameElement) read(r *basebyte.Read, setting *Setting) {
	pmx.ElementTarget = r.Uint8()
	if pmx.ElementTarget == 0x00 {
		pmx.Index = ReadIndex(r, setting.BoneIndexSize)
	} else {
		pmx.Index = ReadIndex(r, setting.MorphIndexSize)
	}
}

func (pmx *Frame) read(r *basebyte.Read, setting *Setting) {
	pmx.FrameName = ReadString(r, setting.Encoding)
	pmx.FrameEnglishName = ReadString(r, setting.Encoding)
	pmx.FrameFlag = r.Uint8()
	pmx.ElementCount = int(r.Int32())
	pmx.Elements = make([]FrameElement, pmx.ElementCount)
	for i := 0; i < pmx.ElementCount; i++ {
		pmx.Elements[i].read(r, setting)
	}
}

func (pmx *RigidBody) read(r *basebyte.Read, setting *Setting) {
	pmx.RigidBodyName = ReadString(r, setting.Encoding)
	pmx.RigidBodyEnglishName = ReadString(r, setting.Encoding)
	pmx.TargetBone = ReadIndex(r, setting.BoneIndexSize)
	pmx.Group = r.Uint8()
	pmx.Mask = r.Uint16()
	pmx.Shape = r.Uint8()
	pmx.Size = getF3(r)
	pmx.Position = getF3(r)
	pmx.Orientation = getF3(r)
	pmx.Mass = r.Float32()
	pmx.MoveAttenuation = r.Float32()
	pmx.RotationAttenuation = r.Float32()
	pmx.Repulsion = r.Float32()
	pmx.Friction = r.Float32()
	pmx.PhysicsCalcType = r.Uint8()
}

func (pmx *JointParam) read(r *basebyte.Read, setting *Setting) {
	pmx.RigidBody1 = ReadIndex(r, setting.RigidbodyIndexSize)
	pmx.RigidBody2 = ReadIndex(r, setting.RigidbodyIndexSize)
	pmx.Position = getF3(r)
	pmx.Orientaiton = getF3(r)
	pmx.MoveLimitationMin = getF3(r)
	pmx.MoveLimitationMax = getF3(r)
	pmx.RotationLimitationMin = getF3(r)
	pmx.RotationLimitationMax = getF3(r)
	pmx.SpringMoveCoefficient = getF3(r)
	pmx.SpringRotationCoefficient = getF3(r)
}

func (pmx *Joint) read(r *basebyte.Read, setting *Setting) {
	pmx.JointName = ReadString(r, setting.Encoding)
	pmx.JointEnglishName = ReadString(r, setting.Encoding)
	pmx.JointType = JointType(r.Uint8())
	pmx.Param.read(r, setting)
}

func (pmx *AnchorRigidBody) read(r *basebyte.Read, setting *Setting) {
	pmx.RelatedRigidBody = ReadIndex(r, setting.RigidbodyIndexSize)
	pmx.RelatedVertex = ReadIndex(r, setting.VertexIndexSize)
	pmx.IsNear = r.Bool()
}

func (pmx *SoftBody) read(r *basebyte.Read, setting *Setting) {
	panic("Not Implemented Exception")
}

// LoadFromFile 从文件加载
func (pmx *Model) LoadFromFile(filename string) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return pmx.LoadFromStream(f)
}

// LoadFromStream 从流加载
func (pmx *Model) LoadFromStream(stream []byte) error {
	if len(stream) == 0 {
		return errors.New("invalid stream")
	}
	// 读取器
	r := &basebyte.Read{Byte: &stream, Order: binary.LittleEndian}

	// マジック
	magic := r.GetByte(4)
	if magic[0] != 0x50 || magic[1] != 0x4d || magic[2] != 0x58 || magic[3] != 0x20 {
		return errors.New("invalid magic number")
	}

	// バージョン
	pmx.Version = r.Float32()
	if pmx.Version != 2.0 && pmx.Version != 2.1 {
		return errors.New("this is not ver2.0 or ver2.1")
	}
	// ファイル設定
	pmx.Setting.read(r)

	// モデル情報
	pmx.ModelName = ReadString(r, pmx.Setting.Encoding)
	pmx.ModelEnglishName = ReadString(r, pmx.Setting.Encoding)
	pmx.ModelComment = ReadString(r, pmx.Setting.Encoding)
	pmx.ModelEnglishComment = ReadString(r, pmx.Setting.Encoding)

	// 頂点
	pmx.VertexCount = int(r.Int32())
	pmx.Vertices = make([]Vertex, pmx.VertexCount)
	for i := 0; i < pmx.VertexCount; i++ {
		pmx.Vertices[i].read(r, &pmx.Setting)
	}

	// 面
	pmx.IndexCount = int(r.Int32())
	pmx.Indices = make([]int, pmx.IndexCount)
	for i := 0; i < pmx.IndexCount; i++ {
		pmx.Indices[i] = ReadIndex(r, pmx.Setting.VertexIndexSize)
	}

	// テクスチャ
	pmx.TextureCount = int(r.Int32())
	pmx.Textures = make([]string, pmx.TextureCount)
	for i := 0; i < pmx.TextureCount; i++ {
		pmx.Textures[i] = ReadString(r, pmx.Setting.Encoding)
	}

	// マテリアル
	pmx.MaterialCount = int(r.Int32())
	pmx.Materials = make([]Material, pmx.MaterialCount)
	for i := 0; i < pmx.MaterialCount; i++ {
		pmx.Materials[i].read(r, &pmx.Setting)
	}

	// ボーン
	pmx.BoneCount = int(r.Int32())
	pmx.Bones = make([]Bone, pmx.BoneCount)
	for i := 0; i < pmx.BoneCount; i++ {
		pmx.Bones[i].read(r, &pmx.Setting)
	}

	// モーフ
	pmx.MorphCount = int(r.Int32())
	pmx.Morphs = make([]Morph, pmx.MorphCount)
	for i := 0; i < pmx.MorphCount; i++ {
		pmx.Morphs[i].read(r, &pmx.Setting)
	}

	// 表示枠
	pmx.FrameCount = int(r.Int32())
	pmx.Frames = make([]Frame, pmx.FrameCount)
	for i := 0; i < pmx.FrameCount; i++ {
		pmx.Frames[i].read(r, &pmx.Setting)
	}

	// 剛体
	pmx.RigidBodyCount = int(r.Int32())
	pmx.RigidBodies = make([]RigidBody, pmx.RigidBodyCount)
	for i := 0; i < pmx.RigidBodyCount; i++ {
		pmx.RigidBodies[i].read(r, &pmx.Setting)
	}

	// ジョイント
	pmx.JointCount = int(r.Int32())
	pmx.Joints = make([]Joint, pmx.JointCount)
	for i := 0; i < pmx.JointCount; i++ {
		pmx.Joints[i].read(r, &pmx.Setting)
	}

	// ソフトボディ
	if pmx.Version == 2.1 {
		pmx.SoftBodyCount = int(r.Int32())
		pmx.SoftBodies = make([]SoftBody, pmx.SoftBodyCount)
	}

	return nil
}
