package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
	// "image/color"
)

var (
	BuildingDetailsMap map[string]BuildingDetails

	SetHUD = "setHUD"
)

type BuildingSystem struct {
	world     *ecs.World
	Buildings []*BuildingEntity
}

func (bs *BuildingSystem) Remove(ecs.BasicEntity) {}

func (bs *BuildingSystem) New(w *ecs.World) {
	bs.world = w

	//Building Definitions (For loop to be able to collapse it)
	BuildingDetailsMap = make(map[string]BuildingDetails)

	for {
		TownCenterTexture, err := common.LoadedSprite("Town_centre.png")
		if err != nil {
			fmt.Println(err.Error())
		}
		TownCenterDetails := BuildingDetails{
			Name: "Town Center", Health: 150, Texture: TownCenterTexture,
			HUDSelectionIndex: 0,
		}
		BuildingDetailsMap[TownCenterDetails.Name] = TownCenterDetails

		MilitaryBlockTexture, err := common.LoadedSprite("Military_block.png")
		if err != nil {
			fmt.Println(err.Error())
		}
		MilitaryBlockDetails := BuildingDetails{
			Name: "Military Block", Health: 120, Texture: MilitaryBlockTexture,
			HUDSelectionIndex: 3,
		}
		BuildingDetailsMap[MilitaryBlockDetails.Name] = MilitaryBlockDetails

		ResourceBuildingTexture, err := common.LoadedSprite("Resource_Building.png")
		if err != nil {
			fmt.Println(err.Error())
		}
		ResourceBuildingDetails := BuildingDetails{
			Name: "Resource Building", Health: 75, Texture: ResourceBuildingTexture,
			HUDSelectionIndex: 4,
		}
		BuildingDetailsMap[ResourceBuildingDetails.Name] = ResourceBuildingDetails

		HouseTexture, err := common.LoadedSprite("House.png")
		if err != nil {
			fmt.Println(err.Error())
		}
		HouseDetails := BuildingDetails{
			Name: "House", Health: 30, Texture: HouseTexture,
			HUDSelectionIndex: 2,
		}
		BuildingDetailsMap[HouseDetails.Name] = HouseDetails

		break
	}

	bs.AddBuilding("Town Center", engo.Point{96, 320})
	bs.AddBuilding("Military Block", engo.Point{320, 320})
	bs.AddBuilding("Resource Building", engo.Point{544, 320})
	bs.AddBuilding("House", engo.Point{768, 320})

	fmt.Println("Building System Initialized")
}

func (bs *BuildingSystem) Update(dt float32) {
	// Mouse Bug is here!
	for _, item := range bs.Buildings {
		if item.MouseComponent.Clicked {
			engo.Mailbox.Dispatch(BuildingMessage{Action: SetHUD, Name: item.GetDetails().Name, Index: 0})
		}
	}
}

func (bs *BuildingSystem) AddBuilding(_BuildingName string, Pos engo.Point) {
	tex := BuildingDetailsMap[_BuildingName].Texture

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
		MouseComponent: common.MouseComponent{},
		BuildingName:   _BuildingName,
	}
	bs.Buildings = append(bs.Buildings, new_building)

	ActiveSystems.RenderSys.Add(&new_building.BasicEntity, &new_building.RenderComponent, &new_building.SpaceComponent)
	ActiveSystems.MouseSys.Add(&new_building.BasicEntity, &new_building.MouseComponent, &new_building.SpaceComponent, &new_building.RenderComponent)
}

type BuildingEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.MouseComponent

	BuildingName string
}

func (be *BuildingEntity) GetDetails() BuildingDetails {
	return BuildingDetailsMap[be.BuildingName]
}

type BuildingDetails struct {
	Name              string
	Health            int
	Texture           *common.Texture
	HUDSelectionIndex int
}

type BuildingMessage struct {
	Action string
	Name   string
	Index  int
}

func (BuildingMessage) Type() string {
	return "BuildingMessage"
}
