package pmd

import (
	"encoding/binary"
	"errors"
	"mmd/basebyte"
	"os"
)

func ifnil(b []byte) []byte {
	for i, n := 0, len(b); i < n; i++ {
		if b[i] == 0x00 {
			return b[:i]
		}
	}
	return b
}

func (pmd *Header) read(r *basebyte.Read) {
	pmd.Name = string(ifnil(r.GetByte(20)))
	pmd.Comment = string(ifnil(r.GetByte(256)))

}
func (pmd *Header) readExtension(r *basebyte.Read) {
	pmd.NameEnglish = string(ifnil(r.GetByte(20)))
	pmd.CommentEnglish = string(ifnil(r.GetByte(256)))
}

func (pmd *Vertex) read(r *basebyte.Read) {
	pmd.Position[0] = r.Float32()
	pmd.Position[1] = r.Float32()
	pmd.Position[2] = r.Float32()
	pmd.Normal[0] = r.Float32()
	pmd.Normal[1] = r.Float32()
	pmd.Normal[2] = r.Float32()
	pmd.Uv[0] = r.Float32()
	pmd.Uv[1] = r.Float32()
	pmd.BoneIndex[0] = r.Uint16()
	pmd.BoneIndex[1] = r.Uint16()
	pmd.BoneWeight = r.Uint8()
	pmd.EdgeInvisible = r.Bool()
}

func (pmd *Material) read(r *basebyte.Read) {
	pmd.Diffuse[0] = r.Float32()
	pmd.Diffuse[1] = r.Float32()
	pmd.Diffuse[2] = r.Float32()
	pmd.Diffuse[3] = r.Float32()
	pmd.Power = r.Float32()
	pmd.Specular[0] = r.Float32()
	pmd.Specular[1] = r.Float32()
	pmd.Specular[2] = r.Float32()
	pmd.Ambient[0] = r.Float32()
	pmd.Ambient[1] = r.Float32()
	pmd.Ambient[2] = r.Float32()
	pmd.ToonIndex = r.Uint8()
	pmd.EdgeFlag = r.Uint8()
	pmd.IndexCount = r.Uint32()
	buffer := r.GetByte(20)
	var l = 0
	for i, n := 0, len(buffer); i < n; i++ {
		if buffer[i] == '*' {
			l = i
			break
		}
	}
	if l == 0 {
		pmd.TextureFilename = string(ifnil(buffer))
		pmd.SphereFilename = ""
	} else {
		pmd.TextureFilename = string(ifnil(buffer[:l]))
		pmd.SphereFilename = string(ifnil(buffer[l:]))
	}
}

func (pmd *Bone) read(r *basebyte.Read) {
	pmd.Name = string((ifnil(r.GetByte(20))))
	pmd.ParentBoneindex = r.Uint16()
	pmd.TailPosboneindex = r.Uint16()
	pmd.BoneType = BoneType(r.Uint8())
	pmd.IkParentboneindex = r.Uint16()
	pmd.BoneHeadpos[0] = r.Float32()
	pmd.BoneHeadpos[1] = r.Float32()
	pmd.BoneHeadpos[2] = r.Float32()
}
func (pmd *Bone) readExtension(r *basebyte.Read) {
	pmd.NameEnglish = string(ifnil(r.GetByte(20)))
}

func (pmd *FaceVertex) read(r *basebyte.Read) {
	pmd.VertexIndex = int(r.Int32())
	pmd.Position[0] = r.Float32()
	pmd.Position[1] = r.Float32()
	pmd.Position[2] = r.Float32()
}

func (pmd *Face) read(r *basebyte.Read) {
	pmd.Name = string((ifnil(r.GetByte(20))))
	count := int(r.Int32())
	pmd.Type = FaceCategory(r.Uint8())
	pmd.Vertices = make([]FaceVertex, count)
	for i := 0; i < count; i++ {
		pmd.Vertices[i].read(r)
	}
}
func (pmd *Face) readExtension(r *basebyte.Read) {
	pmd.NameEnglish = string(ifnil(r.GetByte(20)))
}

func (pmd *BoneDispName) read(r *basebyte.Read) {
	pmd.BoneDispname = string((ifnil(r.GetByte(50))))
}
func (pmd *BoneDispName) readExtension(r *basebyte.Read) {
	pmd.BoneDispnameenglish = string((ifnil(r.GetByte(50))))
}

func (pmd *BoneDisp) read(r *basebyte.Read) {
	pmd.BoneIndex = r.Uint16()
	pmd.BoneDispindex = r.Uint8()
}

func (pmd *Ik) read(r *basebyte.Read) {
	pmd.IkBoneindex = r.Uint16()
	pmd.TargetBoneindex = r.Uint16()
	length := int(r.Uint8())
	pmd.Interations = r.Uint16()
	pmd.AngleLimit = r.Float32()
	pmd.IkChildboneindex = make([]uint16, length)
	for i := 0; i < length; i++ {
		pmd.IkChildboneindex[i] = r.Uint16()
	}
}

func (pmd *RigidBody) read(r *basebyte.Read) {
	pmd.Name = string((ifnil(r.GetByte(20))))
	pmd.RelatedBoneindex = r.Uint16()
	pmd.GroupIndex = r.Uint8()
	pmd.Mask = r.Uint16()
	pmd.Shape = RigidBodyShape(r.Uint8())
	pmd.Size[0] = r.Float32()
	pmd.Size[1] = r.Float32()
	pmd.Size[2] = r.Float32()
	pmd.Position[0] = r.Float32()
	pmd.Position[1] = r.Float32()
	pmd.Position[2] = r.Float32()
	pmd.Orientation[0] = r.Float32()
	pmd.Orientation[1] = r.Float32()
	pmd.Orientation[2] = r.Float32()
	pmd.Weight = r.Float32()
	pmd.LinearDamping = r.Float32()
	pmd.AnglarDamping = r.Float32()
	pmd.Restitution = r.Float32()
	pmd.Friction = r.Float32()
	pmd.RigidType = RigidBodyType(r.Uint8())
}

func (pmd *Constraint) read(r *basebyte.Read) {
	pmd.Name = string((ifnil(r.GetByte(20))))
	pmd.RigidBodyindexa = r.Uint32()
	pmd.RigidBodyindexb = r.Uint32()
	pmd.Position[0] = r.Float32()
	pmd.Position[1] = r.Float32()
	pmd.Position[2] = r.Float32()
	pmd.Orientation[0] = r.Float32()
	pmd.Orientation[1] = r.Float32()
	pmd.Orientation[2] = r.Float32()
	pmd.LinearLowerlimit[0] = r.Float32()
	pmd.LinearLowerlimit[1] = r.Float32()
	pmd.LinearLowerlimit[2] = r.Float32()
	pmd.LinearUpperlimit[0] = r.Float32()
	pmd.LinearUpperlimit[1] = r.Float32()
	pmd.LinearUpperlimit[2] = r.Float32()
	pmd.AngularLowerlimit[0] = r.Float32()
	pmd.AngularLowerlimit[1] = r.Float32()
	pmd.AngularLowerlimit[2] = r.Float32()
	pmd.AngularUpperlimit[0] = r.Float32()
	pmd.AngularUpperlimit[1] = r.Float32()
	pmd.AngularUpperlimit[2] = r.Float32()
	pmd.LinearStiffness[0] = r.Float32()
	pmd.LinearStiffness[1] = r.Float32()
	pmd.LinearStiffness[2] = r.Float32()
	pmd.AngularStiffness[0] = r.Float32()
	pmd.AngularStiffness[1] = r.Float32()
	pmd.AngularStiffness[2] = r.Float32()
}

// LoadFromFile 从文件加载
func (pmd *Model) LoadFromFile(filename string) error {
	f, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return pmd.LoadFromStream(f)
}

// LoadFromStream 从流加载
func (pmd *Model) LoadFromStream(stream []byte) error {
	if len(stream) == 0 {
		return errors.New("invalid stream")
	}
	// 读取器
	r := &basebyte.Read{Byte: &stream, Order: binary.LittleEndian}
	// magic
	if r.Uint8() != 'P' || r.Uint8() != 'm' || r.Uint8() != 'd' {
		return errors.New("invalid file")
	}
	// version
	pmd.Version = r.Float32()
	if pmd.Version != 1.0 {
		return errors.New("invalid version")
	}
	// header
	pmd.Header.read(r)
	// vertices
	vertexNum := int(r.Uint32())
	pmd.Vertices = make([]Vertex, vertexNum)
	for i := 0; i < vertexNum; i++ {
		pmd.Vertices[i].read(r)
	}
	// indices
	indexNum := int(r.Uint32())
	pmd.Indices = make([]uint16, indexNum)
	for i := 0; i < indexNum; i++ {
		pmd.Indices[i] = r.Uint16()
	}
	// materials
	materialNum := int(r.Uint32())
	pmd.Materials = make([]Material, materialNum)
	for i := 0; i < materialNum; i++ {
		pmd.Materials[i].read(r)
	}
	// bones
	boneNum := int(r.Uint16())
	pmd.Bones = make([]Bone, boneNum)
	for i := 0; i < boneNum; i++ {
		pmd.Bones[i].read(r)
	}
	// iks
	ikNum := int(r.Uint16())
	pmd.Iks = make([]Ik, ikNum)
	for i := 0; i < ikNum; i++ {
		pmd.Iks[i].read(r)
	}
	// faces
	faceNum := int(r.Uint16())
	pmd.Faces = make([]Face, faceNum)
	for i := 0; i < faceNum; i++ {
		pmd.Faces[i].read(r)
	}
	// face frames
	faceFrameNum := int(r.Uint8())
	pmd.Facesindices = make([]uint16, faceFrameNum)
	for i := 0; i < faceFrameNum; i++ {
		pmd.Facesindices[i] = r.Uint16()
	}
	// bone names
	boneDispNum := int(r.Uint8())
	pmd.Bonedispname = make([]BoneDispName, boneDispNum)
	for i := 0; i < boneDispNum; i++ {
		pmd.Bonedispname[i].read(r)
	}
	// bone frame
	boneFrameNum := int(r.Uint32())
	pmd.Bonedisp = make([]BoneDisp, boneFrameNum)
	for i := 0; i < boneFrameNum; i++ {
		pmd.Bonedisp[i].read(r)
	}
	// english name
	if r.Bool() {
		pmd.Header.readExtension(r)
		for i := 0; i < boneNum; i++ {
			pmd.Bones[i].readExtension(r)
		}
		for i := 0; i < faceNum; i++ {
			if pmd.Faces[i].Type == Base {
				continue
			}
			pmd.Faces[i].readExtension(r)
		}
		for i := 0; i < boneDispNum; i++ {
			pmd.Bonedispname[i].readExtension(r)
		}
	}
	// toon textures
	if r.Offset >= len(*r.Byte) {
		return nil
	}
	pmd.Toonfilenames = make([]string, 10)
	for i := 0; i < 10; i++ {
		pmd.Toonfilenames[i] = string((ifnil(r.GetByte(100))))
	}
	// physics
	if r.Offset >= len(*r.Byte) {
		return nil
	}
	rigidBodyNum := int(r.Uint32())
	pmd.Rigidbodies = make([]RigidBody, rigidBodyNum)
	for i := 0; i < rigidBodyNum; i++ {
		pmd.Rigidbodies[i].read(r)
	}
	constraintNum := int(r.Uint32())
	pmd.Constraints = make([]Constraint, constraintNum)
	for i := 0; i < constraintNum; i++ {
		pmd.Constraints[i].read(r)
	}
	if r.Offset != len(*r.Byte) {
		return errors.New("there is unknown data")
	}
	return nil
}
