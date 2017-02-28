package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
	// "image/color"
)

type BuildingSystem struct {
	world *ecs.World
}

func (bs *BuildingSystem) Update(dt float32)      {}
func (bs *BuildingSystem) Remove(ecs.BasicEntity) {}

func (bs *BuildingSystem) New(w *ecs.World) {
	bs.world = w

	TownCenterTexture, err := common.LoadedSprite("House.png")
	if err != nil {
		fmt.Println(err.Error())
	}
	testcenter := BuildingEntity{
		BasicEntity: ecs.NewBasic(),
		RenderComponent: common.RenderComponent{
			Drawable: TownCenterTexture,
		},
		SpaceComponent: common.SpaceComponent{
			Position: engo.Point{320, 320},
			Width:    160,
			Height:   160,
		},
		BuildingType: "towncenter",
	}

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&testcenter.BasicEntity, &testcenter.RenderComponent, &testcenter.SpaceComponent)
		}
	}

	fmt.Println("Building System Initialized")
}

type BuildingEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent

	BuildingType string
}
