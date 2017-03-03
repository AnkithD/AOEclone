package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"

	"fmt"
	"image/color"
	"sync"
)

// Defining the Map system
type MapSystem struct {
	world      *ecs.World
	GridSize   int
	vert_lines []GridLineEntity
	hor_lines  []GridLineEntity

	PrevXOffset int
	PrevYOffset int
	XOffset     int
	YOffset     int
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
	ms.GridSize = 32
	ms.PrevXOffset = 0
	ms.PrevYOffset = 0

	//Calculates how many vertical and horizontal grid lines we need
	vert_num := int(engo.WindowWidth()) / ms.GridSize
	hor_num := int(engo.WindowHeight()) / ms.GridSize

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
				Position: engo.Point{float32(i * ms.GridSize), 0},
				Width:    2,
				Height:   engo.WindowHeight(),
			},
		}
		ms.vert_lines[i].RenderComponent.SetZIndex(80)
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
				Position: engo.Point{0, float32(i * ms.GridSize)},
				Width:    engo.WindowWidth(),
				Height:   2,
			},
		}
		// Make the grid HUD, at a depth between 0 and HUD's
		ms.hor_lines[i].RenderComponent.SetZIndex(80)
		ms.hor_lines[i].RenderComponent.SetShader(common.HUDShader)
	}

	// Add each grid line entity to the render system
	for i := 0; i < vert_num; i++ {
		ActiveSystems.RenderSys.Add(&ms.vert_lines[i].BasicEntity, &ms.vert_lines[i].RenderComponent, &ms.vert_lines[i].SpaceComponent)
		ms.vert_lines[i].RenderComponent.Hidden = true
	}
	for i := 0; i < hor_num; i++ {
		ActiveSystems.RenderSys.Add(&ms.hor_lines[i].BasicEntity, &ms.hor_lines[i].RenderComponent, &ms.hor_lines[i].SpaceComponent)
		ms.hor_lines[i].RenderComponent.Hidden = true
	}

	fmt.Println("Map System initialized")
}

func (ms *MapSystem) Update(dt float32) {
	// Toggle the hidden attribute of every grid line's render component
	if engo.Input.Button(GridToggle).JustPressed() {
		for i, _ := range ms.vert_lines {
			ms.vert_lines[i].RenderComponent.Hidden = !ms.vert_lines[i].RenderComponent.Hidden
		}
		for i, _ := range ms.hor_lines {
			ms.hor_lines[i].RenderComponent.Hidden = !ms.hor_lines[i].RenderComponent.Hidden
		}
	}

	if ms.vert_lines[0].RenderComponent.Hidden == false {
		CameraSystem := ActiveSystems.CameraSys

		ms.XOffset = int(CameraSystem.X()) % ms.GridSize
		ms.YOffset = int(CameraSystem.Y()) % ms.GridSize

		wg := sync.WaitGroup{}

		wg.Add(2)
		// Updating hor and vert lines in parallel for faster execution
		go func() {
			for i, _ := range ms.vert_lines {
				ms.vert_lines[i].Position.Add(engo.Point{float32(ms.PrevXOffset - ms.XOffset), 0})
			}
			wg.Done()
		}()

		go func() {
			for i, _ := range ms.hor_lines {
				ms.hor_lines[i].Position.Add(engo.Point{0, float32(ms.PrevYOffset - ms.YOffset)})
			}
			wg.Done()
		}()
		wg.Wait()

		ms.PrevXOffset = ms.XOffset
		ms.PrevYOffset = ms.YOffset
	}
}
