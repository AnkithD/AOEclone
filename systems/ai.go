package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
	"image/color"
)

var (
	HumanDetailsMap map[string]HumanDetails
)

type AISystem struct {
	world *ecs.World
}

func (ais *AISystem) New(w *ecs.World) {
	ais.world = w

	WarriorDetails := HumanDetails{
		Name:      "Warrior",
		Texture:   common.Circle{},
		Color:     color.RGBA{0, 0, 255, 255},
		Width:     float32(GridSize - 5),
		Height:    float32(GridSize - 5),
		MaxHealth: 100,
		Attack:    10,
	}

	fmt.Println("AI System Initialized")
}
func (*AISystem) Update(dt float32)      {}
func (*AISystem) Remove(ecs.BasicEntity) {}

type Human struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	AIComponent

	Name   string
	Health int
	Attack int
}

type AIComponent struct {
	StartPoint  engo.Point
	EndPoint    engo.Point
	CurrentPath []grid
	State       string
}

type HumanDetails struct {
	Name      string
	MaxHealth int
	Attack    int
	Texture   common.Drawable
	Color     color.RGBA
	Width     float32
	Height    float32
}

func CreateHuman(_Name string, Position engo.Point) {

}
