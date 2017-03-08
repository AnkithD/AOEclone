package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"

	"container/heap"
	"fmt"
	"image/color"
	"math"
	"sync"
)

// Defining the Map system
type MapSystem struct {
	world      *ecs.World
	vert_lines []GridEntity
	hor_lines  []GridEntity
	ChunkBoxes [][]GridEntity

	LinePrevXOffset int
	LinePrevYOffset int
	BoxPrevXOffset  int
	BoxPrevYOffset  int
}

//Place holders to satisfy Interface

func (*MapSystem) Remove(ecs.BasicEntity) {}

// Every object of this entity is one grid line
type GridEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

var PathBlocks []*GridEntity
var PathBlocksMutex sync.Mutex
var item_tobe_placed bool = true
var mouseheld bool = false

// When system is created this func is executed
// Initialze the world variable and assign tab to toggle the grid
func (ms *MapSystem) New(w *ecs.World) {
	ms.world = w
	GridSize = 32
	ms.LinePrevXOffset = 0
	ms.LinePrevYOffset = 0
	ms.BoxPrevXOffset = 0
	ms.BoxPrevYOffset = 0

	PathBlocks = make([]*GridEntity, 0)

	// Initializes the Grid lines
	func() {
		//Calculates how many vertical and horizontal grid lines we need
		vert_num := int(engo.WindowWidth()) / GridSize
		hor_num := int(engo.WindowHeight()) / GridSize

		//Each grid line is an Entity, so we are storing all vert and hor lines in two
		//Seperate slices
		ms.vert_lines = make([]GridEntity, vert_num)
		ms.hor_lines = make([]GridEntity, hor_num)

		//Generating the vert grid lines
		for i := 0; i < vert_num; i++ {
			ms.vert_lines[i] = GridEntity{
				BasicEntity: ecs.NewBasic(),
				RenderComponent: common.RenderComponent{
					Drawable: common.Rectangle{},
					Color:    color.RGBA{0, 0, 0, 125},
				},
				SpaceComponent: common.SpaceComponent{
					Position: engo.Point{float32(i * GridSize), 0},
					Width:    2,
					Height:   engo.WindowHeight(),
				},
			}
			ms.vert_lines[i].RenderComponent.SetZIndex(80)
			ms.vert_lines[i].RenderComponent.SetShader(common.HUDShader)
		}
		//Generating the hor grid lines
		for i := 0; i < hor_num; i++ {
			ms.hor_lines[i] = GridEntity{
				BasicEntity: ecs.NewBasic(),
				RenderComponent: common.RenderComponent{
					Drawable: common.Rectangle{},
					Color:    color.RGBA{0, 0, 0, 125},
				},
				SpaceComponent: common.SpaceComponent{
					Position: engo.Point{0, float32(i * GridSize)},
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
	}()

	// Initializes the Chunk Rectangles
	func() {
		per_row := int(math.Ceil(float64(engo.WindowWidth())/float64(GridSize*ChunkSize))) + 1
		per_col := int(engo.WindowHeight())/(GridSize*ChunkSize) + 1

		ms.ChunkBoxes = make([][]GridEntity, per_row)
		for i := 0; i < per_row; i++ {
			ms.ChunkBoxes[i] = make([]GridEntity, per_col)
		}

		for i := 0; i < per_row; i++ {
			for j := 0; j < per_col; j++ {
				ms.ChunkBoxes[i][j] = GridEntity{
					BasicEntity: ecs.NewBasic(),
					SpaceComponent: common.SpaceComponent{
						Position: engo.Point{float32(i * GridSize * ChunkSize), float32(j * GridSize * ChunkSize)},
						Width:    float32(GridSize * ChunkSize),
						Height:   float32(GridSize * ChunkSize),
					},
					RenderComponent: common.RenderComponent{
						Drawable: common.Rectangle{
							BorderWidth: 2,
							BorderColor: color.RGBA{255, 255, 255, 255},
						},
						Color: color.RGBA{0, 0, 0, 0},
					},
				}

				ms.ChunkBoxes[i][j].SetShader(common.HUDShader)
				ms.ChunkBoxes[i][j].SetZIndex(81)
				ms.ChunkBoxes[i][j].Hidden = true

				cb := &ms.ChunkBoxes[i][j]
				ActiveSystems.RenderSys.Add(&cb.BasicEntity, &cb.RenderComponent, &cb.SpaceComponent)
			}
		}
	}()

	fmt.Println("Map System initialized")
}

func (ms *MapSystem) Update(dt float32) {

	/*
		For Placing Bushes and trees
	*/
	mx, my := GetAdjustedMousePos(false)

	func() {
		if engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonRight {
			fmt.Println(item_tobe_placed)
			item_tobe_placed = !item_tobe_placed
		}
		if engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonLeft {
			mouseheld = true
		}
		if engo.Input.Mouse.Action == engo.Release && engo.Input.Mouse.Button == engo.MouseButtonLeft {
			mouseheld = false
		}
		if mouseheld {
			var BuildingName string

			if item_tobe_placed {
				BuildingName = "Tree"
			} else {
				BuildingName = "Bush"
			}
			if WithinGameWindow(mx, my) {
				pik := float32(math.Floor(float64(mx)/float64(GridSize)) * float64(GridSize))
				cik := float32(math.Floor(float64(my)/float64(GridSize)) * float64(GridSize))
				if !Grid[int(pik/32)][int(cik/32)] {
					engo.Mailbox.Dispatch(CreateBuildingMessage{Name: BuildingName, Position: engo.Point{X: pik, Y: cik}})
				}
			}
		}
		if engo.Input.Button(R_remove).JustPressed() {

		}
	}()

	//Rendering the Gridlines and Chunk Boxes
	func() {
		// Toggle the hidden attribute of every grid line's render component
		if engo.Input.Button(GridToggle).JustPressed() {
			for i, _ := range ms.vert_lines {
				ms.vert_lines[i].RenderComponent.Hidden = !ms.vert_lines[i].RenderComponent.Hidden
			}
			for i, _ := range ms.hor_lines {
				ms.hor_lines[i].RenderComponent.Hidden = !ms.hor_lines[i].RenderComponent.Hidden
			}
			for i, _ := range ms.ChunkBoxes {
				for j, _ := range ms.ChunkBoxes[i] {
					ms.ChunkBoxes[i][j].RenderComponent.Hidden = !ms.ChunkBoxes[i][j].RenderComponent.Hidden
				}
			}
		}

		if ms.vert_lines[0].RenderComponent.Hidden == false {
			CamSys := ActiveSystems.CameraSys

			LineXOffset := int(CamSys.X()) % GridSize
			LineYOffset := int(CamSys.Y()) % GridSize
			BoxXOffset := int(CamSys.X()-engo.WindowWidth()/2) % (GridSize * ChunkSize)
			BoxYOffset := int(CamSys.Y()-engo.WindowHeight()/2) % (GridSize * ChunkSize)

			wg := sync.WaitGroup{}

			wg.Add(3)
			// Updating hor and vert lines in parallel for faster execution
			go func() {
				defer wg.Done()
				for i, _ := range ms.vert_lines {
					ms.vert_lines[i].Position.Add(engo.Point{float32(ms.LinePrevXOffset-LineXOffset) * CamSys.Z() * (engo.GameWidth() / engo.CanvasWidth()), 0})
				}
			}()

			go func() {
				defer wg.Done()
				for i, _ := range ms.hor_lines {
					ms.hor_lines[i].Position.Add(engo.Point{0, float32(ms.LinePrevYOffset-LineYOffset) * CamSys.Z() * (engo.GameHeight() / engo.CanvasHeight())})
				}
			}()

			go func() {
				defer wg.Done()
				for i, _ := range ms.ChunkBoxes {
					for j, _ := range ms.ChunkBoxes[i] {
						ms.ChunkBoxes[i][j].Position.Add(engo.Point{float32(ms.BoxPrevXOffset-BoxXOffset) * CamSys.Z() * (engo.GameWidth() / engo.CanvasWidth()), float32(ms.BoxPrevYOffset-BoxYOffset) * CamSys.Z() * (engo.GameHeight() / engo.CanvasHeight())})
					}
				}
			}()
			wg.Wait()

			ms.LinePrevXOffset = LineXOffset
			ms.LinePrevYOffset = LineYOffset
			ms.BoxPrevXOffset = BoxXOffset
			ms.BoxPrevYOffset = BoxYOffset
		}
	}()

	// Handle Middle Mouse clicks for debugging
	func() {
		mx, my := GetAdjustedMousePos(false)

		if engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonMiddle {
			fmt.Println("---------------------------------------------")
			fmt.Println("Mouse Pos is", mx, ",", my)
			ChunkRef, ChunkIndex := GetChunkFromPos(mx, my)
			Chunk := *ChunkRef

			if len(Chunk) > 0 {
				fmt.Println("-------------------------")
				for _, item := range Chunk {
					fmt.Println(item.GetStaticComponent().Name, "present in chunk:", ChunkIndex)
				}
				fmt.Println("-------------------------")
			} else {
				fmt.Println("Chunk", ChunkIndex, "Empty")
			}

			if GetGridAtPos(mx, my) {
				fmt.Println("Grid at", int(mx)/GridSize, ",", int(my)/GridSize, "is occupied")
			} else {
				fmt.Println("Grid at", int(mx)/GridSize, ",", int(my)/GridSize, "is not occupied")
			}
			fmt.Println("---------------------------------------------")
		}
	}()

	// Make path blocks fade
	func() {
		ShouldDelete := make([]bool, len(PathBlocks))

		for i, _ := range PathBlocks {
			r, g, b, a := PathBlocks[i].RenderComponent.Color.RGBA()
			A := float32(a) / 255
			A -= (255 * dt) / 2
			if A > 0 {
				A = float32(math.Floor(float64(A)))
				PathBlocks[i].RenderComponent.Color = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(A)}
			} else {
				ShouldDelete[i] = true
			}
		}
		i := 0
		for len(PathBlocks) > 0 {
			if ShouldDelete[i] {
				ActiveSystems.RenderSys.Remove(PathBlocks[i].BasicEntity)
				//Fast delete from slice
				PathBlocks[i] = PathBlocks[len(PathBlocks)-1]
				PathBlocks = PathBlocks[:len(PathBlocks)-1]
			} else {
				if i < len(PathBlocks) {
					i++
				}
			}

			if i >= len(PathBlocks) {
				break
			}
		}
	}()

	// A* Visualization
	func() {
		if engo.Input.Button(ShiftKey).JustReleased() && ShowDebugPathfinding {

			s, e := grid{x: 19, y: 15}, grid{x: int(mx) / GridSize, y: int(my) / GridSize}
			if (e.x < GridMaxX) && (e.y < GridMaxY) && !Grid[e.x][e.y] {
				DrawPathBlock(s.x, s.y, color.RGBA{0, 0, 255, 255})
				go GetPath(s, e, PathChannel)
			} else {
				fmt.Println(e.x, e.y, GridMaxX, GridMaxY)
			}
		}

		select {
		case res := <-PathChannel:
			for i, item := range res {
				if i == len(res)-1 {
					DrawPathBlock(item.x, item.y, color.RGBA{0, 255, 0, 255})
				} else {
					DrawPathBlock(item.x, item.y, color.RGBA{255, 0, 0, 255})
				}
			}
		default:
		}
	}()
}

//---------------------------------------------PathFinding Algorithm------------------------------------------------
type grid struct {
	x   int
	y   int
	f   float32
	g   float32
	par *grid
}

type gridHeap []*grid

func (h gridHeap) Len() int           { return len(h) }
func (h gridHeap) Less(i, j int) bool { return (h)[i].f < (h)[j].f }
func (h gridHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *gridHeap) Push(x interface{}) {
	*h = append(*h, x.(*grid))
}

func (h *gridHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func hvalue(x, y int, endgrid grid) float32 {
	var diagCost, sideCost float32
	diagCost, sideCost = 1, 1

	var a, b int
	X := endgrid.x
	Y := endgrid.y
	a = int(math.Abs(float64(X - x)))
	b = int(math.Abs(float64(Y - y)))
	if a > b {
		a, b = b, a
	}
	return (float32(b) - float32(a)*(diagCost-sideCost))
	//return b
}

func eval(neighbor *grid, block *grid, endgrid *grid, h *gridHeap, list *[][]bool) bool {
	if neighbor.x == endgrid.x && neighbor.y == endgrid.y {
		endgrid.par = block
		(*list)[endgrid.x][endgrid.y] = true
		return true
	}
	if neighbor.x == block.x || neighbor.y == block.y {
		neighbor.g = block.g + 1
	} else {
		neighbor.g = block.g + 1.414
	}
	hval := hvalue(neighbor.x, neighbor.y, *endgrid)
	neighbor.f = hval + neighbor.g
	neighbor.par = block
	//fmt.Println("f =", neighbor.f, ", Set par to", neighbor.par.x, ",", neighbor.par.y)
	(*list)[neighbor.x][neighbor.y] = true
	heap.Push(h, neighbor)
	return false
}

func open(block *grid, h *gridHeap, _list *[][]bool, endgrid *grid) {
	list := *_list

	if func() bool {
		//fmt.Println("-----------------------------------")
		//defer fmt.Println("-----------------------------------")
		//fmt.Println("Analyzing the neighbors of", block.x, ",", block.y)
		if block.x-1 >= 0 && block.y-1 >= 0 {
			if !Grid[block.x-1][block.y-1] {
				if !list[block.x-1][block.y-1] {
					neighbor := grid{
						x: block.x - 1,
						y: block.y - 1,
					}
					//fmt.Println("Evaluating grid at", neighbor.x, ",", neighbor.y)
					//DrawPathBlock(neighbor.x, neighbor.y)
					foundend := eval(&neighbor, block, endgrid, h, _list)
					if foundend {
						return true
					}
				}
			}
		}
		if block.y-1 >= 0 {
			if !Grid[block.x][block.y-1] {
				if !list[block.x][block.y-1] {
					neighbor := grid{
						x: block.x,
						y: block.y - 1,
					}
					//fmt.Println("Evaluating grid at", neighbor.x, ",", neighbor.y)
					//DrawPathBlock(neighbor.x, neighbor.y)
					foundend := eval(&neighbor, block, endgrid, h, _list)
					if foundend {
						return true
					}
				}
			}
		}
		if block.x+1 < int(engo.GameWidth()*ScaleFactor)/GridSize && block.y-1 >= 0 {
			if !Grid[block.x+1][block.y-1] {
				if !list[block.x+1][block.y-1] {
					neighbor := grid{
						x: block.x + 1,
						y: block.y - 1,
					}
					//fmt.Println("Evaluating grid at", neighbor.x, ",", neighbor.y)
					//DrawPathBlock(neighbor.x, neighbor.y)
					foundend := eval(&neighbor, block, endgrid, h, _list)
					if foundend {
						return true
					}
				}
			}
		}
		if block.x-1 >= 0 {
			if !Grid[block.x-1][block.y] {
				if !list[block.x-1][block.y] {
					neighbor := grid{
						x: block.x - 1,
						y: block.y,
					}
					//fmt.Println("Evaluating grid at", neighbor.x, ",", neighbor.y)
					//DrawPathBlock(neighbor.x, neighbor.y)
					foundend := eval(&neighbor, block, endgrid, h, _list)
					if foundend {
						return true
					}
				}
			}
		}
		if block.x+1 < int(engo.GameWidth()*ScaleFactor)/GridSize {
			if !Grid[block.x+1][block.y] {
				if !list[block.x+1][block.y] {
					neighbor := grid{
						x: block.x + 1,
						y: block.y,
					}
					//fmt.Println("Evaluating grid at", neighbor.x, ",", neighbor.y)
					//DrawPathBlock(neighbor.x, neighbor.y)
					foundend := eval(&neighbor, block, endgrid, h, _list)
					if foundend {
						return true
					}
				}
			}
		}
		if block.x-1 >= 0 && block.y+1 < int(engo.GameHeight()*ScaleFactor)/GridSize {
			if !Grid[block.x-1][block.y+1] {
				if !list[block.x-1][block.y+1] {
					neighbor := grid{
						x: block.x - 1,
						y: block.y + 1,
					}
					//fmt.Println("Evaluating grid at", neighbor.x, ",", neighbor.y)
					//DrawPathBlock(neighbor.x, neighbor.y)
					foundend := eval(&neighbor, block, endgrid, h, _list)
					if foundend {
						return true
					}
				}
			}
		}
		if block.y+1 < int(engo.GameHeight()*ScaleFactor)/GridSize {
			if !Grid[block.x][block.y+1] {
				if !list[block.x][block.y+1] {
					neighbor := grid{
						x: block.x,
						y: block.y + 1,
					}
					//fmt.Println("Evaluating grid at", neighbor.x, ",", neighbor.y)
					//DrawPathBlock(neighbor.x, neighbor.y)
					foundend := eval(&neighbor, block, endgrid, h, _list)
					if foundend {
						return true
					}
				}
			}
		}
		if block.x+1 < int(engo.GameWidth()*ScaleFactor)/GridSize && block.y+1 < int(engo.GameHeight()*ScaleFactor)/GridSize {
			if !Grid[block.x+1][block.y+1] {
				if !list[block.x+1][block.y+1] {
					neighbor := grid{
						x: block.x + 1,
						y: block.y + 1,
					}
					//fmt.Println("Evaluating grid at", neighbor.x, ",", neighbor.y)
					//DrawPathBlock(neighbor.x, neighbor.y)
					foundend := eval(&neighbor, block, endgrid, h, _list)
					if foundend {
						return true
					}
				}
			}
		}
		return false
	}() {
		return
	}

	block = heap.Pop(h).(*grid)
	p := block.x
	q := block.y
	(*_list)[p][q] = true
	if len(*h) <= 0 {
		return
	}

	open(block, h, _list, endgrid) //function calling recursively
}

func ReversePath(slice []grid) []grid {
	n := len(slice) - 1
	res := make([]grid, n+1)

	for i, item := range slice {
		res[n-i] = item
	}

	return res
}

func GetPath(startgrid grid, endgrid grid, c chan []grid) {

	x1 := startgrid.x
	y1 := startgrid.y
	h := &gridHeap{}
	heap.Init(h)
	startgrid.par = &startgrid
	startgrid.g = 0

	list := make([][]bool, int(engo.WindowWidth()*ScaleFactor)/GridSize)
	for i, _ := range list {
		list[i] = make([]bool, int(engo.WindowHeight()*ScaleFactor)/GridSize)
	}
	list[x1][y1] = true

	open(&startgrid, h, &list, &endgrid)

	var path []grid

	temp := endgrid
	path = append(path, temp)
	//fmt.Println("Adding grid at", temp.x, ",", temp.y, "to the path")
	for temp.par != &startgrid {
		prevtemp := temp
		temp = *(temp.par)
		//fmt.Println("Adding grid at", temp.x, ",", temp.y, "to the path")
		path = append(path, temp)
		if prevtemp.x == temp.x && prevtemp.y == temp.y {
			//fmt.Println("No change from ", prevtemp.x, ",", prevtemp.y, "and", temp.x, ",", temp.y)
			c <- make([]grid, 0)
		}
	}
	path = append(path, temp)
	c <- ReversePath(path)

}

func DrawPathBlock(x, y int, col color.RGBA) {
	myblock := GridEntity{
		BasicEntity: ecs.NewBasic(),
		SpaceComponent: common.SpaceComponent{
			Position: engo.Point{float32(x * GridSize), float32(y * GridSize)},
			Width:    float32(GridSize - 10),
			Height:   float32(GridSize - 10),
		},
		RenderComponent: common.RenderComponent{
			Drawable: common.Rectangle{},
			Color:    col,
		},
	}
	myblock.SetZIndex(70)

	PathBlocks = append(PathBlocks, &myblock)

	ActiveSystems.RenderSys.Add(&myblock.BasicEntity, &myblock.RenderComponent, &myblock.SpaceComponent)
}
