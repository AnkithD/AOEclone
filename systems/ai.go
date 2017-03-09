package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
	"image/color"
	"math/rand"
)

var (
	HumanDetailsMap map[string]HumanDetails

	timer float32

	f bool

	n int
)

type AISystem struct {
	world *ecs.World
}

func (ais *AISystem) New(w *ecs.World) {
	ais.world = w

	// WarriorDetails := HumanDetails{
	// 	Name:      "Warrior",
	// 	Texture:   common.Circle{},
	// 	Color:     color.RGBA{0, 0, 255, 255},
	// 	Width:     float32(GridSize - 5),
	// 	Height:    float32(GridSize - 5),
	// 	MaxHealth: 100,
	// 	Attack:    10,
	// }

	fmt.Println("AI System Initialized")
}
func (*AISystem) Update(dt float32) {

	func() {
		timer = timer + dt
		if timer >= float32(2) {
			f = true
			n = n + 2
			fmt.Println("soldiers have started at the coordinates: ")
			timer = 0
		}
	}()
	func() {
		if f {
			var x []int
			var y []int
			var p, q int
			for i := 0; i < n; i++ {
				p = rand.Intn(7)
				q = rand.Intn(7)
				x = append(x, p)
				y = append(y, q)
			}
			for i := 0; i < n; i++ {
				fmt.Printf("%d,%d", x[i], y[i])
			}
		}
	}()
}
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
