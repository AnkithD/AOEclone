package systems

import (
	"engo.io/ecs"
	//"engo.io/engo"
	"engo.io/engo/common"
	//"image/color"
)

type HUDSystem struct {
	world ecs.World
}

func (*HUDSystem) Update(dt float32) {}

func (*HUDSystem) Remove(ecs.BasicEntity) {}

type BottomHUD struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}
