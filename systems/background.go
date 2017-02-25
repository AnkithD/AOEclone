package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"

	"fmt"
	"image/color"
)

// Defining the Map system
type MapSystem struct {
	world     *ecs.World
	grid_size int
}

//Place holders to satisfy Interface
func (*MapSystem) Update(dt float32) {}

func (*MapSystem) Remove(ecs.BasicEntity) {}

type GridLineEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

//When system is created this func is executed
// Initialze the world variable and create an object of test entity
func (ms *MapSystem) New(w *ecs.World) {
	ms.world = w

	fmt.Println("Enter the grid size: ")
	fmt.Scanf("%d", &ms.grid_size)

	vert_num := int(engo.WindowHeight()) / ms.grid_size
	hor_num := int(engo.WindowWidth()) / ms.grid_size

	vert_lines := make([]GridLineEntity, vert_num)
	hor_lines := make([]GridLineEntity, hor_num)

	var render_sys *common.RenderSystem
	for _, system := range ms.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			render_sys = sys
		}
	}

	for i, line := range vert_lines {
		line = GridLineEntity{
			BasicEntity: ecs.NewBasic(),
			RenderComponent: common.RenderComponent{
				Drawable: common.Rectangle{},
				Color:    color.RGBA{255, 255, 255, 125},
			},
			SpaceComponent: common.SpaceComponent{
				Position: engo.Point{float32(i * ms.grid_size), 0},
				Width:    1,
				Height:   engo.WindowHeight(),
			},
		}
		render_sys.Add(&line.BasicEntity, &line.RenderComponent, &line.SpaceComponent)
	}
	for i, line := range hor_lines {
		line = GridLineEntity{
			BasicEntity: ecs.NewBasic(),
			RenderComponent: common.RenderComponent{
				Drawable: common.Rectangle{},
				Color:    color.RGBA{255, 255, 255, 125},
			},
			SpaceComponent: common.SpaceComponent{
				Position: engo.Point{0, float32(i * ms.grid_size)},
				Width:    engo.WindowWidth(),
				Height:   1,
			},
		}
		render_sys.Add(&line.BasicEntity, &line.RenderComponent, &line.SpaceComponent)
	}

}
