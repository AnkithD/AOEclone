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
		"Roboto-Regular.ttf", "Deselect_button.png", "bush.png",
	)
	if err != nil {
		panic(err)
	}

	for _, item := range systems.BuildingSprites {
		err := engo.Files.Load(item)
		if err != nil {
			panic(err)
		}
	}
}

func (*myScene) Setup(world *ecs.World) {
	world.AddSystem(&common.RenderSystem{})
	world.AddSystem(common.NewKeyboardScroller(640, systems.HorAxis, systems.VertAxis))
	systems.CacheActiveSystems(world)
	systems.RegisterButtons()
	systems.InitializeVariables()

	world.AddSystem(&systems.MapSystem{})
	world.AddSystem(&systems.BuildingSystem{})
	world.AddSystem(&systems.HUDSystem{})

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
