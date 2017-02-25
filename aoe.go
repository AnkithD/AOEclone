package main

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"

	"image/color"
)

type myScene struct{}

func (*myScene) Type() string { return "myGame" }

func (*myScene) Preload() {

}

func (*myScene) Setup(world *ecs.World) {
	world.AddSystem(new(common.RenderSystem))

	common.SetBackground(color.Gray)
}
func main() {
	opts := engo.RunOptions{
		Title:  "AOE Clone",
		Width:  1280,
		Height: 720,
	}

	engo.Run(opts, new(myScene))
}
