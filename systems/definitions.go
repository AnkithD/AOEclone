package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
)

// Button mappings
var (
	gridToggle = "gridtoggle"
)

var ActiveSystems ActiveSystemsStruct

func RegisterButtons() {
	engo.Input.RegisterButton(gridToggle, engo.Tab)
}

func CacheActiveSystems(world *ecs.World) {
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			ActiveSystems.RenderSys = sys
		case *common.MouseSystem:
			ActiveSystems.MouseSys = sys
		}
	}
}

type ActiveSystemsStruct struct {
	RenderSys *common.RenderSystem
	MouseSys  *common.MouseSystem
}
