package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"

	"image/color"
)

// Defining the Map system
type MapSystem struct {
	world *ecs.World
}

func (*MapSystem) Update(dt float32) {}

func (*MapSystem) Remove(ecs.BasicEntity) {}

type testEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func draw_grid() {
	test_entity := testEntity{BasicEntity}
}
