package systems

import (
	"container/heap"
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
	"image/color"
	"math"
	"math/rand"
)

var (
	HumanDetailsMap map[string]HumanDetails

	timer float32

	n int

	PathBlocks []*GridEntity
)

type AISystem struct {
	world *ecs.World

	Humans []HumanEntity
}

const (
	StateWaiting  = iota
	StateMoving   = iota
	StateFighting = iota
)

func (ais *AISystem) New(w *ecs.World) {
	ais.world = w
	PathBlocks = make([]*GridEntity, 0)
	HumanDetailsMap = make(map[string]HumanDetails)

	// Defining Human Entities
	func() {
		tex, err := common.LoadedSprite(WarriorSprite)
		if err != nil {
			fmt.Println("Could not load warrior")
		}

		WarriorDetails := HumanDetails{
			Name:      "Warrior",
			Texture:   tex,
			Width:     tex.Width(),
			Height:    tex.Height(),
			MaxHealth: 100,
			Attack:    10,
		}

		tex, err = common.LoadedSprite(EWarriorSprite)
		if err != nil {
			fmt.Println("Could not load Ewarrior")
		}

		EWarriorDetails := HumanDetails{
			Name:      "Enemy",
			Texture:   tex,
			Width:     tex.Width(),
			Height:    tex.Height(),
			MaxHealth: 100,
			Attack:    10,
		}

		HumanDetailsMap[WarriorDetails.Name] = WarriorDetails
		HumanDetailsMap[EWarriorDetails.Name] = EWarriorDetails
	}()

	fmt.Println("AI System Initialized")
}
func (ais *AISystem) Update(dt float32) {

	func() {
		timer = timer + dt
		if timer >= float32(50000) {
			n = n + 2
			fmt.Println("soldiers have started at the coordinates:\n")
			timer = 0
			var x []int
			var y []int
			var p, q int
			for i := 0; i < n; i++ {
				p = rand.Intn(7) + GridMaxX - 2*ChunkSize
				q = rand.Intn(7) + GridMaxY - 2*ChunkSize
				x = append(x, p)
				y = append(y, q)
			}
			for i := 0; i < n; i++ {
				fmt.Printf("x=%d,y=%d\n", x[i], y[i])
				ais.CreateHuman("Enemy", engo.Point{float32(x[i] * GridSize), float32(y[i] * GridSize)})
			}
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

	mx, my := GetAdjustedMousePos(false)
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

func (ais *AISystem) CreateHuman(_Name string, Pos engo.Point) {
	details := HumanDetailsMap[_Name]

	new_human := HumanEntity{
		BasicEntity: ecs.NewBasic(),
		SpaceComponent: common.SpaceComponent{
			Position: Pos,
		},
		RenderComponent: common.RenderComponent{
			Drawable: details.Texture,
		},
		AIComponent: AIComponent{
			State:       StateWaiting,
			LastGridPos: engo.Point{X: float32(math.Floor(float64(Pos.X / float32(GridSize)))), Y: float32(math.Floor(float64(Pos.Y / float32(GridSize))))},
		},
		Health: details.MaxHealth,
		Name:   _Name,
	}

	ais.Humans = append(ais.Humans, new_human)
	CacheInSectors(&new_human)
	Grid[int(new_human.LastGridPos.X)][int(new_human.LastGridPos.Y)] = true

	ActiveSystems.RenderSys.Add(&new_human.BasicEntity, &new_human.RenderComponent, &new_human.SpaceComponent)
}

func (*AISystem) Remove(ecs.BasicEntity) {}

type HumanEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	AIComponent

	Name   string
	Health int
}

type AIComponent struct {
	StartPoint  engo.Point
	EndPoint    engo.Point
	CurrentPath []grid
	State       int
	LastGridPos engo.Point
}

func (aic *AIComponent) MoveTo(To engo.Point) {}

func (aic *AIComponent) Update(dt float32) {
	if aic.State == "Waiting" {

	}
}

type HumanDetails struct {
	Name      string
	MaxHealth int
	Attack    int
	Texture   common.Drawable
	Width     float32
	Height    float32
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
