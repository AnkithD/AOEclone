package main

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/Ankithd/AOEClone/systems"

	"image/color"
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
	world.AddSystem(&common.EdgeScroller{ScrollSpeed: 640, EdgeMargin: 32})
	world.AddSystem(&common.MouseSystem{})
	systems.CacheActiveSystems(world)
	systems.RegisterButtons()

	world.AddSystem(&systems.MapSystem{})
	world.AddSystem(&systems.HUDSystem{})
	world.AddSystem(&systems.BuildingSystem{})

	common.SetBackground(color.RGBA{120, 120, 120, 255})
}
func main() {
	opts := engo.RunOptions{
		Title:         "AOE Clone",
		Width:         1280,
		Height:        768,
		ScaleOnResize: true,
		MSAA:          2,
	}

	engo.Run(opts, new(myScene))
}
