package particles

import (
	"math/rand"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/skycoin/cx-game/render"
	"github.com/skycoin/cx-game/utility"
	"github.com/skycoin/cx-game/spriteloader"
)

type Particle struct {
	ID               int32
	ParticleMetaType int32
	Type             int32
	Sprite           uint32
	Size             int32
	// should position instead be stored as matrix
	// such that particles can spin??
	PosX             float32
	PosY             float32
	VelocityX        float32
	VelocityY        float32
	TimeToLive       float32
}
var particles = []Particle{}
var particleShader *utility.Shader
const initialVelocityScale = 3
const tileChunkLifetime = 1

func InitParticles() {
	particleShader = utility.NewShader(
		"./assets/shader/simple.vert", "./assets/shader/particle.frag" )
	InitLasers()
}

func TickParticles(dt float32) {
	TickLasers(dt)
	// age and kill particles
	newParticles := []Particle{}
	for _,laser := range particles {
		laser.TimeToLive -= dt
		if laser.TimeToLive > 0 {
			newParticles = append(newParticles,laser)
		}
	}
	particles = newParticles
}

func configureGlForParticles() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.TEXTURE_2D)
	// bits of tile are still pixel art; use nearest interpolation
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.ActiveTexture(gl.TEXTURE0)

	gl.BindVertexArray(spriteloader.QuadVao);
}

func DrawChunkParticles(ctx render.Context) {
	configureGlForParticles()
	particleShader.Use()
	// particles share both shader and projection matrix - set it only once
	particleShader.SetMat4("projection", &ctx.Projection)
	for _,particle := range particles {
		DrawChunkParticle(particle,ctx)
	}
}

// draw particles between mid and top layers
func DrawMidTopParticles(ctx render.Context) {
	DrawLasers(ctx)
}

// draw particles over top layer
func DrawTopParticles(ctx render.Context) {
	DrawChunkParticles(ctx)
}

func (p Particle) GetTransform() mgl32.Mat4 {
	size := float32(p.Size)
	return mgl32.Ident4().
		Mul4(mgl32.Translate3D(p.PosX,p.PosY,0)).
		Mul4(mgl32.Scale3D(size,size,1))
}

func DrawChunkParticle(particle Particle, ctx render.Context) {
	metadata := spriteloader.GetSpriteMetadata(particle.Sprite)
	particleShader.SetUint("tex", metadata.GpuTex)
	gl.BindTexture(gl.TEXTURE_2D, metadata.GpuTex)

	alpha := particle.TimeToLive / laserDuration
	alpha = 1
	particleShader.SetVec4F("color", 1,1,1, alpha)

	world := ctx.World.Mul4(particle.GetTransform())
	particleShader.SetMat4("world", &world)

	// TODO apply offset and scale to achieve a view 
	// of only the 2x2 chunk of the tile we are interested in
	particleShader.SetVec2F("offset", 0,0)
	particleShader.SetVec2F("scale", 1,1)

	gl.DrawArrays(gl.TRIANGLES,0,6) // draw quad
}

func CreateTileChunk(x,y float32, TileSpriteID uint32) {
	particle := Particle {
		ID: rand.Int31(),
		PosX: x, PosY: y,
		Size: 1,
		VelocityX: (rand.Float32()-0.5)*initialVelocityScale,
		VelocityY: (rand.Float32()-0.5)*initialVelocityScale,
		Sprite: TileSpriteID,
		TimeToLive: tileChunkLifetime,
	}
	particles = append(particles, particle)
}
