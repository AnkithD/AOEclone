package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
)

// Button mappings
var (
	GridToggle = "gridtoggle"
	HorAxis    = "horAxis"
	VertAxis   = "vertAxis"
)

var ActiveSystems ActiveSystemsStruct

func RegisterButtons() {
	engo.Input.RegisterButton(GridToggle, engo.Tab)
	engo.Input.RegisterAxis(HorAxis, engo.AxisKeyPair{engo.A, engo.D})
	engo.Input.RegisterAxis(VertAxis, engo.AxisKeyPair{engo.W, engo.S})

	fmt.Println("Registered Buttons")
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

	fmt.Println("Cached Important System References")
}

type ActiveSystemsStruct struct {
	RenderSys *common.RenderSystem
	MouseSys  *common.MouseSystem
}
