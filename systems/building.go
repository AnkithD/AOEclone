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

	BuildingDetailsMap map[string]BuildingDetails
	Buildings          []*BuildingEntity
}

func (bs *BuildingSystem) Update(dt float32)      {}
func (bs *BuildingSystem) Remove(ecs.BasicEntity) {}

func (bs *BuildingSystem) New(w *ecs.World) {
	bs.world = w

	bs.BuildingDetailsMap = make(map[string]BuildingDetails)
	//Building Definitions (For loop to be able to collapse it)
	for {
		TownCenterTexture, err := common.LoadedSprite("Town_centre.png")
		if err != nil {
			fmt.Println(err.Error())
		}
		TownCenterDetails := BuildingDetails{
			Name: "Town Center", Health: 150, Texture: TownCenterTexture,
			HUDSelectionIndex: 0,
		}
		bs.BuildingDetailsMap[TownCenterDetails.Name] = TownCenterDetails

		MilitaryBlockTexture, err := common.LoadedSprite("Military_block.png")
		if err != nil {
			fmt.Println(err.Error())
		}
		MilitaryBlockDetails := BuildingDetails{
			Name: "Military Block", Health: 120, Texture: MilitaryBlockTexture,
			HUDSelectionIndex: 3,
		}
		bs.BuildingDetailsMap[MilitaryBlockDetails.Name] = MilitaryBlockDetails

		ResourceBuildingTexture, err := common.LoadedSprite("Resource_Building.png")
		if err != nil {
			fmt.Println(err.Error())
		}
		ResourceBuildingDetails := BuildingDetails{
			Name: "Resource Building", Health: 75, Texture: ResourceBuildingTexture,
			HUDSelectionIndex: 4,
		}
		bs.BuildingDetailsMap[ResourceBuildingDetails.Name] = ResourceBuildingDetails

		HouseTexture, err := common.LoadedSprite("House.png")
		if err != nil {
			fmt.Println(err.Error())
		}
		HouseDetails := BuildingDetails{
			Name: "House", Health: 30, Texture: HouseTexture,
			HUDSelectionIndex: 2,
		}
		bs.BuildingDetailsMap[HouseDetails.Name] = HouseDetails

		break
	}

	bs.AddBuilding("Town Center", engo.Point{96, 320}, false)
	bs.AddBuilding("Military Block", engo.Point{320, 320}, false)
	bs.AddBuilding("Resource Building", engo.Point{544, 320}, false)
	bs.AddBuilding("Town Center", engo.Point{768, 320}, false)

	fmt.Println("Building System Initialized")
}

func (bs *BuildingSystem) AddBuilding(BuildingName string, Pos engo.Point, Render bool) {
	tex := bs.BuildingDetailsMap[BuildingName].Texture

	// Using reference so that the newly created building
	// doesn't get garbage collected after func return
	new_building := &BuildingEntity{
		BasicEntity: ecs.NewBasic(),
		RenderComponent: common.RenderComponent{
			Drawable: tex,
		},
		SpaceComponent: common.SpaceComponent{
			Position: Pos,
			Width:    tex.Width(),
			Height:   tex.Height(),
		},
	}
	bs.Buildings = append(bs.Buildings, new_building)

	if Render {
	}
}

type BuildingEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent

	BuildingName string
}

type BuildingDetails struct {
	Name              string
	Health            int
	Texture           *common.Texture
	HUDSelectionIndex int
}
