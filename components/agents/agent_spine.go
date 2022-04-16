package agents

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/skycoin/cx-game/cxmath"
	"github.com/skycoin/cx-game/render"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/Texture"
	"github.com/skycoin/cx-game/cmd/dynamicBatchShaderExample/geometry"
	"github.com/skycoin/cx-game/engine/spine"
	"github.com/skycoin/cx-game/test/spine-animation/animation"

	_ "image/png"
)

//todo structure better
type SpineData struct {
	Time  float64
	Play  bool
	Speed float64

	// TODO: replace this with atlas
	ImagesPath string
	Images     map[string]render.SpriteID
	Images2    map[string]*Texture.Texture

	Skeleton  *spine.Skeleton
	Animation *spine.Animation

	SkinIndex      int
	AnimationIndex int

	DebugCenter bool
	DebugBones  bool

	IgnoringPlatformsFor_2 float32
}

func (a *SpineData) LoadCharacter(loc animation.Location) (*SpineData, error) {
	fmt.Println("test here: ", loc.JSON)
	rd, err := os.Open(loc.JSON)
	if err != nil {
		fmt.Println("hit Error: ")
		return nil, err
	}
	fmt.Println("test data: ")
	data, err := spine.ReadJSON(rd)
	if err != nil {
		return nil, err
	}

	data.Name = loc.Name

	fmt.Println("test location: ", loc.Images)
	a.ImagesPath = loc.Images

	a.Images2 = make(map[string]*Texture.Texture)

	a.Images = make(map[string]render.SpriteID)

	//char.Images = make(map[string]*pixel.PictureData)

	//fmt.Println("____________Debug point_________")

	//fmt.Printf("%v", char.Images)

	//	fmt.Println("________________________________")

	a.Play = false
	a.DebugBones = true
	a.DebugCenter = true

	a.Speed = 0.05
	a.Skeleton = spine.NewSkeleton(data)
	a.Skeleton.Skin = a.Skeleton.Data.DefaultSkin
	a.Animation = a.Skeleton.Data.Animations[0]

	a.AnimationIndex = 0
	a.SkinIndex = 0

	a.Skeleton.FlipY = false
	a.Skeleton.UpdateAttachments()
	a.Skeleton.Update()
	fmt.Printf("%v", a)
	return a, nil
}

func (char *SpineData) Description() string {
	return char.Skeleton.Data.Name + " > " + char.Skeleton.Skin.Name + " > " + char.Animation.Name
}

func (char *SpineData) NextAnimation(offset int) {
	char.AnimationIndex += offset
	for char.AnimationIndex < 0 {
		char.AnimationIndex += len(char.Skeleton.Data.Animations)
	}
	char.AnimationIndex = char.AnimationIndex % len(char.Skeleton.Data.Animations)
	char.Animation = char.Skeleton.Data.Animations[char.AnimationIndex]
	char.Skeleton.SetToSetupPose()
	char.Skeleton.Update()
}
func (char *SpineData) SetAnimation(index int) {

	char.Animation = char.Skeleton.Data.Animations[index]
	char.Skeleton.SetToSetupPose()
	char.Skeleton.Update()
}

func (char *SpineData) NextSkin(offset int) {
	char.SkinIndex += offset
	for char.SkinIndex < 0 {
		char.SkinIndex += len(char.Skeleton.Data.Skins)
	}
	char.SkinIndex = char.SkinIndex % len(char.Skeleton.Data.Skins)
	char.Skeleton.Skin = char.Skeleton.Data.Skins[char.SkinIndex]
	char.Skeleton.SetToSetupPose()
	char.Skeleton.Update()
	char.Skeleton.UpdateAttachments()
}
func (char *SpineData) Update(dt float64, translate, scale cxmath.Vec2) {
	if char.Play {
		char.Time += dt * char.Speed
	}

	char.Skeleton.Local.Translate.Set(float32(translate.X), float32(translate.Y))
	char.Skeleton.Local.Scale.Set(scale.X, scale.Y)
	char.Animation.Apply(char.Skeleton, float32(char.Time), true)
	char.Skeleton.Update()
}

func (char *SpineData) GetImage(attachment string) render.SpriteID {
	// if path != "" {
	// 	attachment = path
	// }
	if pdi, ok := char.Images[attachment]; ok {
		return pdi
	}
	//	fmt.Println("Loading " + attachment)

	// fallback := func() *ebiten.Image {
	// 	fmt.Println("missing: ", attachment)

	// 	m := image.NewRGBA(image.Rect(0, 0, 10, 10))
	// 	for i := range m.Pix {
	// 		m.Pix[i] = 0x80
	// 	}

	// 	pd, _ := ebiten.NewImageFromImage(m, ebiten.FilterDefault)
	// 	char.Images[attachment] = pd
	// 	return pd
	// }

	//fullpath := filepath.Join(char.ImagesPath, attachment+".png")
	// file, err := os.Open(fullpath)
	// if err != nil {
	// 	panic("Error opening file")
	// 	//return fallback()
	// }

	// m, _, err := image.Decode(file)
	// if err != nil {
	// 	panic("Error Decoding file")
	// 	// return fallback()
	// }
	spritename := attachment + ":0"
	pd := render.GetSpriteIDByName(spritename)

	//fmt.Println("Loaded: ", pd)
	//char.Images[attachment] = pd

	return pd
}

func (char *SpineData) GetImage2(attachment, path string) *Texture.Texture {
	if path != "" {
		attachment = path
	}
	if pd, ok := char.Images2[attachment]; ok {
		return pd
	}
	//fmt.Println("Loading " + attachment)
	//fmt.Println("path URL: " + path)

	// fallback := func() *ebiten.Image {
	// 	fmt.Println("missing: ", attachment)

	// 	m := image.NewRGBA(image.Rect(0, 0, 10, 10))
	// 	for i := range m.Pix {
	// 		m.Pix[i] = 0x80
	// 	}

	// 	pd, _ := ebiten.NewImageFromImage(m, ebiten.FilterDefault)
	// 	char.Images[attachment] = pd
	// 	return pd
	// }

	fullpath := filepath.Join(char.ImagesPath, attachment+".png")
	// file, err := os.Open(fullpath)
	// if err != nil {
	// 	panic("Error opening file")
	// 	//return fallback()
	// }

	// m, _, err := image.Decode(file)
	// if err != nil {
	// 	panic("Error Decoding file")
	// 	// return fallback()
	// }
	pd := Texture.SetUpTexture(fullpath)

	char.Images2[attachment] = pd

	return pd
}

func (char *SpineData) GetImageSize(id render.SpriteID) [2]float32 {
	sprite := render.GetSpriteByID(id)

	var width_height [2]float32

	width_height[0] = sprite.Width
	width_height[1] = sprite.Height
	return width_height
}

// func (char *SpineData) Draw() []float32 {
// 	var pos []float32
// 	//var offset float32 = 0
// 	var count = 0
// 	for i, slot := range char.Skeleton.Order {
// 		bone := slot.Bone
// 		switch attachment := slot.Attachment.(type) {
// 		case nil:
// 		case *spine.RegionAttachment:
// 			local := attachment.Local.Affine()
// 			final := bone.World.Mul(local)

// 			m := char.GetImage(attachment.Name, attachment.Path)
// 			xform := geometry.Matrix(final.Col64())
// 			//box := m.Bounds()

// 			// set blending mode
// 			// BUG: incorrect, should use blending mode not compositing mode
// 			switch slot.Data.Blend {
// 			case spine.Normal:
// 				// MISSING
// 			case spine.Additive:
// 			//	draw.CompositeMode = ebiten.CompositeModeLighter
// 			case spine.Multiply:
// 				// MISSING
// 			case spine.Screen:
// 				// MISSING
// 				//	draw.CompositeMode = ebiten.CompositeModeLighter
// 			}
// 			//_, _, _, _, tx, ty := draw.GeoM.GetElements()
// 			// poject := Project(mgl32.Vec2{final.Translation().X, final.Translation().Y}, final)
// 			m.Bind(uint32(i))
// 			// if attachment.Name == "goggles" {
// 			// 	fmt.Println("goggles sizex:", m.M_width)
// 			// 	fmt.Println("goggles sizey:", m.M_height)
// 			// 	fmt.Println("goggles transform:", final.Translation().Y)

// 			// }
// 			//if attachment.Name == "gun" || attachment.Name == "head" || attachment.Name == "goggles" {
// 			q := CreateQuad(float32(m.M_height), float32(m.M_width), float32(i), xform, m)
// 			pos = append(pos, q...)
// 			//}
// 			// fmt.Println(i)
// 			//offset += 50.0
// 			count += 1
// 			//target.DrawImage(m, &draw)
// 		default:
// 			panic(fmt.Sprintf("unknown attachment %v", attachment))
// 		}
// 	}

// 	// panic("stopped here")
// 	return pos
// }

type Vertex struct {
	Position  mgl32.Vec2
	TexCoords mgl32.Vec2
	TexID     float32
}

// var QuadPosition [4]mgl32.Vec4

func Project(u mgl32.Vec2, m spine.Affine) mgl32.Vec2 {
	return mgl32.Vec2{m.M00*u.X() + m.M01*u.Y() + m.M02, m.M10*u.X() + m.M11*u.Y() + m.M12}
}

func CreateQuad(h, w, textureID float32, matrix geometry.Matrix, m *Texture.Texture) []float32 {
	// QuadPosition[0] = mgl32.Vec4{-0.5, -0.5,0.0,1.0}
	// QuadPosition[1] = mgl32.Vec4{0.5, -0.5,0.0,1.0}
	// QuadPosition[2] = mgl32.Vec4{0.5, 0.5,0.0,1.0}
	// QuadPosition[3] = mgl32.Vec4{-0.5, 0.5,0.0,1.0}
	var (
		//center     = s.frame.Center()
		horizontal = geometry.V(float64(w/2), 0)
		vertical   = geometry.V(0, float64(h/2))
	)

	dirty := false
	if matrix != m.M_metrix {
		m.M_metrix = matrix
		dirty = true
	}
	if dirty {

	}
	//	var size float32 = 100.0
	var pos []float32
	v0 := Vertex{}
	// v0.Position = mgl32.Vec2{x, y}
	v0.TexCoords = mgl32.Vec2{0.0, 0.0}
	v0.TexID = textureID
	t0 := geometry.Vec{}.Sub(horizontal).Sub(vertical)

	xy0 := m.M_metrix.Project(t0)
	v0.Position = mgl32.Vec2{float32(xy0.X), float32(xy0.Y)}
	//	v0.Position = Project(v0.Position, final)
	pos = append(pos, v0.Position.X(), v0.Position.Y(), v0.TexCoords.X(), v0.TexCoords.Y(), v0.TexID)

	v1 := Vertex{}
	//v1.Position = mgl32.Vec2{x + w, y}
	v1.TexCoords = mgl32.Vec2{1.0, 0.0}
	v1.TexID = textureID
	t1 := geometry.Vec{}.Add(horizontal).Sub(vertical)
	xy1 := m.M_metrix.Project(t1)
	v1.Position = mgl32.Vec2{float32(xy1.X), float32(xy1.Y)}
	//v1.Position = Project(v1.Position, final)
	pos = append(pos, v1.Position.X(), v1.Position.Y(), v1.TexCoords.X(), v1.TexCoords.Y(), v1.TexID)

	v2 := Vertex{}
	//v2.Position = mgl32.Vec2{x + w, y + h}
	v2.TexCoords = mgl32.Vec2{1.0, 1.0}
	v2.TexID = textureID
	t2 := geometry.Vec{}.Add(horizontal).Add(vertical)
	xy2 := m.M_metrix.Project(t2)
	v2.Position = mgl32.Vec2{float32(xy2.X), float32(xy2.Y)}
	//	v2.Position = Project(v2.Position, final)
	pos = append(pos, v2.Position.X(), v2.Position.Y(), v2.TexCoords.X(), v2.TexCoords.Y(), v2.TexID)

	v3 := Vertex{}
	//	v3.Position = mgl32.Vec2{x, y + h}
	v3.TexCoords = mgl32.Vec2{0.0, 1.0}
	v3.TexID = textureID
	t3 := geometry.Vec{}.Sub(horizontal).Add(vertical)
	xy3 := m.M_metrix.Project(t3)
	v3.Position = mgl32.Vec2{float32(xy3.X), float32(xy3.Y)}
	//v3.Position = Project(v3.Position, final)
	pos = append(pos, v3.Position.X(), v3.Position.Y(), v3.TexCoords.X(), v3.TexCoords.Y(), v3.TexID)

	//fmt.Printf("numbers=%v\n", pos)
	return pos

}
