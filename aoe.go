package main

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/Ankithd/AOEClone/systems"

	"image/color"
	"runtime"
)

type myScene struct{}

// Place holder methods to satisfy interface
func (*myScene) Type() string { return "myGame" }

func (*myScene) Preload() {
	err := engo.Files.Load(
		"Roboto-Regular.ttf", "Town_centre.png", "Military_block.png", "Resource_Building.png",
		"House.png",
	)
	if err != nil {
		panic(err)
	}
}

func (*myScene) Setup(world *ecs.World) {
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(&common.MouseSystem{})
	world.AddSystem(common.NewKeyboardScroller(640, systems.HorAxis, systems.VertAxis))
	systems.CacheActiveSystems(world)
	systems.RegisterButtons()
	systems.InitializeVariables()

	world.AddSystem(&systems.MapSystem{})
	world.AddSystem(&systems.HUDSystem{})
	world.AddSystem(&systems.BuildingSystem{})

	common.SetBackground(color.RGBA{182, 204, 104, 255})
}
func main() {
	runtime.GOMAXPROCS(8)
	opts := engo.RunOptions{
		Title:         "AOE Clone",
		Width:         1280,
		Height:        768,
		ScaleOnResize: true,
		MSAA:          2,
		VSync:         true,
	}

	engo.Run(opts, new(myScene))
}
