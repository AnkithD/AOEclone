package main

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"github.com/Ankithd/AOEClone/systems"

	"image/color"
)

type myScene struct{}

func (*myScene) Type() string { return "myGame" }

func (*myScene) Preload() {

}

func (*myScene) Setup(world *ecs.World) {
	world.AddSystem(new(common.RenderSystem))
	world.AddSystem(new(systems.MapSystem))
	world.AddSystem(new(systems.HUDSystem))

	common.SetBackground(color.RGBA{120, 120, 120, 255})
}
func main() {
	opts := engo.RunOptions{
		Title:         "AOE Clone",
		Width:         1280,
		Height:        720,
		ScaleOnResize: true,
	}

	engo.Run(opts, new(myScene))
}
