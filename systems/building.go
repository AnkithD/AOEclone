package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	// "fmt"
	// "image/color"
)

type BuildingSystem struct {
	world *ecs.World
}

func (bs *BuildingSystem) Update(dt float32)      {}
func (bs *BuildingSystem) Remove(ecs.BasicEntity) {}

func (bs *BuildingSystem) New(w *ecs.World) {
	bs.world = w

	testcenter := Building{
		BasicEntity: ecs.NewBasic(),
		RenderComponent: common.RenderComponent{
			Drawable: TownCenterTexture,
			Scale:    engo.Point{1, 1},
		},
		SpaceComponent: common.SpaceComponent{
			Position: engo.Point{320, 320},
			Width:    TownCenterTexture.Width(),
			Height:   TownCenterTexture.Height(),
		},
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&testcenter.BasicEntity, &testcenter.RenderComponent, &testcenter.SpaceComponent)
		}
	}
}

type Building struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent

	Type string
}
