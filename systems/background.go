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
	world      *ecs.World
	grid_size  int
	vert_lines []GridLineEntity
	hor_lines  []GridLineEntity
}

//Place holders to satisfy Interface

func (*MapSystem) Remove(ecs.BasicEntity) {}

// Every object of this entity is one grid line
type GridLineEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

// When system is created this func is executed
// Initialze the world variable and assign tab to toggle the grid
func (ms *MapSystem) New(w *ecs.World) {
	ms.world = w

	engo.Input.RegisterButton(gridToggle, engo.Tab)

	ms.grid_size = 32
	//Calculates how many vertical and horizontal grid lines we need
	vert_num := int(engo.WindowWidth()) / ms.grid_size
	hor_num := int(engo.WindowHeight()) / ms.grid_size

	//Each grid line is an Entity, so we are storing all vert and hor lines in two
	//Seperate slices
	ms.vert_lines = make([]GridLineEntity, vert_num)
	ms.hor_lines = make([]GridLineEntity, hor_num)

	//Generating the vert grid lines
	for i := 0; i < vert_num; i++ {
		ms.vert_lines[i] = GridLineEntity{
			BasicEntity: ecs.NewBasic(),
			RenderComponent: common.RenderComponent{
				Drawable: common.Rectangle{},
				Color:    color.RGBA{0, 0, 0, 125},
			},
			SpaceComponent: common.SpaceComponent{
				Position: engo.Point{float32(i * ms.grid_size), 0},
				Width:    2,
				Height:   engo.WindowHeight(),
			},
		}
		ms.vert_lines[i].RenderComponent.SetZIndex(800)
		ms.vert_lines[i].RenderComponent.SetShader(common.HUDShader)
	}
	//Generating the hor grid lines
	for i := 0; i < hor_num; i++ {
		ms.hor_lines[i] = GridLineEntity{
			BasicEntity: ecs.NewBasic(),
			RenderComponent: common.RenderComponent{
				Drawable: common.Rectangle{},
				Color:    color.RGBA{0, 0, 0, 125},
			},
			SpaceComponent: common.SpaceComponent{
				Position: engo.Point{0, float32(i * ms.grid_size)},
				Width:    engo.WindowWidth(),
				Height:   2,
			},
		}
		// Make the grid HUD, at a depth between 0 and HUD's
		ms.hor_lines[i].RenderComponent.SetZIndex(800)
		ms.hor_lines[i].RenderComponent.SetShader(common.HUDShader)
	}

	// Add each grid line entity to the render system
	for _, system := range ms.world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			for i := 0; i < vert_num; i++ {
				sys.Add(&ms.vert_lines[i].BasicEntity, &ms.vert_lines[i].RenderComponent, &ms.vert_lines[i].SpaceComponent)
				ms.vert_lines[i].RenderComponent.Hidden = true
			}
			for i := 0; i < hor_num; i++ {
				sys.Add(&ms.hor_lines[i].BasicEntity, &ms.hor_lines[i].RenderComponent, &ms.hor_lines[i].SpaceComponent)
				ms.hor_lines[i].RenderComponent.Hidden = true
			}
		}
	}

	fmt.Println("Map System initialized")
}

func (ms *MapSystem) Update(dt float32) {
	// Toggle the hidden attribute of every grid line's render component
	if engo.Input.Button(gridToggle).JustPressed() {
		for i, _ := range ms.vert_lines {
			ms.vert_lines[i].RenderComponent.Hidden = !ms.vert_lines[i].RenderComponent.Hidden
		}
		for i, _ := range ms.hor_lines {
			ms.hor_lines[i].RenderComponent.Hidden = !ms.hor_lines[i].RenderComponent.Hidden
		}
	}
}
