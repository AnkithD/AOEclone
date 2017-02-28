package systems

import (
	"engo.io/ecs"
	"engo.io/engo/common"
)

// Button mappings
var (
	gridToggle = "gridToggle"
)

var ActiveSystems ActiveSystemsStruct

func CacheActiveSystems(world *ecs.World) {
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			ActiveSystems.RenderSys = sys
		}
	}
}

type ActiveSystemsStruct struct {
	RenderSys *common.RenderSystem
}
