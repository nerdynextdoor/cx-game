package particle_draw

import (
	"log"

	"github.com/skycoin/cx-game/components/types"
	"github.com/skycoin/cx-game/constants"
	"github.com/skycoin/cx-game/render"
)

type ParticleDrawHandler func()

var ParticleDrawHandlerList [constants.NUM_PARTICLE_DRAW_HANDLERS]ParticleDrawHandler

func Init() {

	particleShader := render.NewShader("./assets/")
	RegisterDrawHandler(constants.PARTICLE_DRAW_HANDLER_1, DrawSolid)
	RegisterDrawHandler(constants.PARTICLE_DRAW_HANDLER_2, DrawTransparent)

	AssertAllDrawHandlersRegistered()
}

func AssertAllDrawHandlersRegistered() {
	for _, handler := range ParticleDrawHandlerList {
		if handler == nil {
			log.Fatalln("Did not initialize all particle draw handlers")
		}
	}
}

func RegisterDrawHandler(id types.ParticleDrawHandlerId, handler ParticleDrawHandler) {
	ParticleDrawHandlerList[id] = handler
}

func GetDrawHandler(id types.ParticleDrawHandlerId) ParticleDrawHandler {
	return ParticleDrawHandlerList[id]
}
