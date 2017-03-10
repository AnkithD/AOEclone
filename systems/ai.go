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
	"runtime"
	//"sync"
)

var (
	HumanDetailsMap map[string]HumanDetails

	timer        float32
	n            int
	PlacingHuman *HumanIcon

	PathBlocks []*GridEntity
	timeMax    float32
)

type AISystem struct {
	world *ecs.World

	Humans        []*HumanEntity
	HumanChannels chan HumanComStruct
	FChannelNum   bool
}

const (
	StateWaiting  = iota
	StateMoving   = iota
	StateFighting = iota
)

func (ais *AISystem) New(w *ecs.World) {
	ais.world = w
	ais.Humans = make([]*HumanEntity, 0)
	ais.HumanChannels = make(chan HumanComStruct, 100)

	timeMax = 30
	PlacingHuman = nil
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

	engo.Mailbox.Listen("HumanHealthEnquiryMessage", func(_msg engo.Message) {
		msg, ok := _msg.(HumanHealthEnquiryMessage)
		if !ok {
			panic("AI System expected HumanHealthEnquiryMessage, instead got unexpected")
		}
		for _, item := range ais.Humans {
			if item.BasicEntity.ID() == msg.ID {
				HealthEnquiryResponse.HealthResult = item.Health
				// switch item.Name {
				// case "Bush":
				// 	HealthEnquiryResponse.ResourceName = "Food"
				// case "Tree":
				// 	HealthEnquiryResponse.ResourceName = "Wood"
				// }
				HealthEnquiryResponse.set = true
				return
			}
		}

		panic("Health Enquiry for unkown building")
	})

	engo.Mailbox.Listen("CreateHumanMessage", func(_msg engo.Message) {
		fmt.Println("Got message")
		msg, ok := _msg.(CreateHumanMessage)
		if !ok {
			panic("AI System wants Create human message!")
		}

		mx, my := GetAdjustedMousePos(false)
		mp := engo.Point{mx, my}
		detail := HumanDetailsMap[msg.Name]
		PlacingHuman = &HumanIcon{
			BasicEntity: ecs.NewBasic(),
			SpaceComponent: common.SpaceComponent{
				Position: mp,
				Width:    detail.Width,
				Height:   detail.Height,
			},
			RenderComponent: common.RenderComponent{
				Drawable: detail.Texture,
				Color:    color.RGBA{255, 255, 255, 150},
			},
			Name: msg.Name,
		}

		ActiveSystems.RenderSys.Add(
			&PlacingHuman.BasicEntity,
			&PlacingHuman.RenderComponent,
			&PlacingHuman.SpaceComponent,
		)
	})
	for i, _ := range ais.Humans {
		go ais.Humans[i].Update(0, ais.HumanChannels)
	}

	fmt.Println("AI System Initialized: Note", runtime.GOMAXPROCS(0), "Physical Threads available")
}
func (ais *AISystem) Update(dt float32) {

	// Enemy Spawning
	func() {
		timer = timer + dt
		if timer >= timeMax {
			timeMax = 180
			n = n + 2
			//fmt.Println("soldiers have started at the coordinates:\n")
			timer = 0
			var x []int
			var y []int
			var p, q int
			for i := 0; i < n; i++ {
				p = rand.Intn(7) + GridMaxX - 2*ChunkSize
				q = rand.Intn(7) + GridMaxY - 2*ChunkSize
				if Grid[p][q] {
					i--
					continue
				}
				x = append(x, p)
				y = append(y, q)
			}
			for i := 0; i < n; i++ {
				//fmt.Printf("x=%d,y=%d\n", x[i], y[i])
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
			if ((255 * dt) / 2) < 1 {
				A -= 1
			} else {
				A -= ((255 * dt) / 2)
			}

			if A > 0 {
				A = float32(math.Floor(float64(A)))
				PathBlocks[i].RenderComponent.Color = color.RGBA{uint8(r), uint8(g), uint8(b), uint8(A)}
				ShouldDelete[i] = false
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
	mp := engo.Point{mx, my}

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

	// Comunications with Humans
	func() {
	ReadMessage:
		for {
			select {
			case msg := <-ais.HumanChannels:
				for i, _ := range ais.Humans {
					if ais.Humans[i].ID() == msg.ID {
						go ais.Humans[i].Update(dt, ais.HumanChannels)
					}
				}
			default:
				break ReadMessage
			}
		}

		// for i, _ := range ais.Humans {
		// 	ais.Humans[i].Update(dt, ais.HumanChannels)
		// }
	}()

	// Handling of Mouse click on human
	func() {
		if engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonLeft && PlacingHuman == nil {
			Sector, _ := GetSectorFromPos(mx, my)

			if len(*Sector) > 0 {
				//fmt.Println("-------------------------")
				for _, item := range *Sector {
					if item.SpaceComponent.Contains(mp) {
						engo.Mailbox.Dispatch(SetBottomHUDMessage{ID: item.ID(), Name: item.Name, Index: 0})
					}
					//fmt.Println(item.GetStaticComponent().Name, "present in chunk:", ChunkIndex)
				}
				//fmt.Println("-------------------------")
			} else {
				//fmt.Println("Chunk", ChunkIndex, "Empty")
			}
		}
	}()

	// Placing Human code
	func() {
		if PlacingHuman != nil {
			PlacingHuman.SpaceComponent.Position = mp
			if engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonLeft &&
				WithinGameWindow(mx, my) && !GetGridAtPos(mx, my) {
				pos := engo.Point{X: float32((int(mx) / GridSize) * GridSize), Y: float32((int(my) / GridSize) * GridSize)}
				ActiveSystems.RenderSys.Remove(PlacingHuman.BasicEntity)
				ais.CreateHuman(PlacingHuman.Name, pos)
				PlacingHuman = nil
			}
		}
	}()

}

func (ais *AISystem) CreateHuman(_Name string, Pos engo.Point) {
	details := HumanDetailsMap[_Name]

	new_human := HumanEntity{
		BasicEntity: ecs.NewBasic(),
		RenderComponent: common.RenderComponent{
			Drawable: details.Texture,
		},
		AIComponent: AIComponent{
			State:          StateWaiting,
			LastGridPos:    grid{x: int(Pos.X) / GridSize, y: int(Pos.Y) / GridSize},
			SpaceComponent: common.SpaceComponent{Position: Pos, Width: details.Width, Height: details.Height},
			Direction:      -1,
		},
		Health: details.MaxHealth,
		Name:   _Name,
	}

	ais.Humans = append(ais.Humans, &new_human)
	CacheInSectors(&new_human)
	Grid[new_human.LastGridPos.x][new_human.LastGridPos.y] = true

	ActiveSystems.RenderSys.Add(&new_human.BasicEntity, &new_human.RenderComponent, &new_human.SpaceComponent)
	new_human.Update(0, ais.HumanChannels)
}

func (*AISystem) Remove(ecs.BasicEntity) {}

type HumanIcon struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	Name string
}

type HumanEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	AIComponent

	Name   string
	Health int
}

type AIComponent struct {
	StartPoint  engo.Point
	EndPoint    engo.Point
	CurrentPath []grid
	State       int
	LastGridPos grid
	Direction   int

	common.SpaceComponent
}

func (he *HumanEntity) MoveTo(To engo.Point) {
	if WithinGameWindow(To.X, To.Y) {
		c := make(chan []grid)
		go GetPath(PointToGrid(he.Position), PointToGrid(To), c)
		he.CurrentPath = <-c
		he.Direction = -1
		he.State = StateMoving
		// for _, item := range he.CurrentPath {
		// 	DrawPathBlock(item.x, item.y, color.RGBA{0, 0, 255, 255})
		// }
		return
	}
}

func (he *HumanEntity) Update(dt float32, ComChannel chan HumanComStruct) {

	switch he.State {
	case StateWaiting:
		if he.Name == "Enemy" && dt == 0 {
			he.MoveTo(engo.Point{X: float32(32 * GridSize), Y: float32(21 * GridSize)})
		}

		if engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonRight &&
			ActiveHUDLabel != nil && ActiveHUDLabel.ID == he.ID() && he.Name == "Warrior" {
			mx, my := GetAdjustedMousePos(false)
			mp := engo.Point{mx, my}
			x, y := int(mx)/GridSize, int(my)/GridSize
			// fmt.Println(x != int(he.Position.X)/GridSize, y != int(he.Position.Y)/GridSize,
			// 	WithinGameWindow(mx, my), !Grid[x][y])
			if !(x == int(he.Position.X)/GridSize && y == int(he.Position.Y)/GridSize) &&
				WithinGameWindow(mx, my) && !Grid[x][y] {
				DrawPathBlock(x, y, color.RGBA{194, 24, 7, 150})
				he.MoveTo(mp)
			}
		}
	case StateMoving:
		speed := float32(3 * GridSize)
		TargetLocation := engo.Point{float32(he.CurrentPath[0].x * GridSize), float32(he.CurrentPath[0].y * GridSize)}
		if he.Direction == -1 {
			he.Direction = he.GetDirection(TargetLocation)
		}
		x := float32((he.Direction%10)-1) * speed * dt
		y := float32(((he.Direction/10)%10)-1) * speed * dt

		//fmt.Println(float32((he.Direction%10)-1), float32(((he.Direction/10)%10)-1))

		he.Position.Add(engo.Point{x, y})
		//fmt.Println(he.Position.PointDistance(TargetLocation))
		if he.Position.PointDistance(TargetLocation) < 2 {
			he.Position = TargetLocation
			i, j := int(he.Position.X)/GridSize, int(he.Position.Y)/GridSize
			I, J := he.LastGridPos.x, he.LastGridPos.y

			Grid[I][J] = false
			Grid[i][j] = true

			_, prevsec := GetSectorFromPos(float32(i*GridSize), float32(j*GridSize))
			_, thissec := GetSectorFromPos(float32(I*GridSize), float32(J*GridSize))
			if prevsec != thissec {
				UnCacheInSectors(he, engo.Point{float32(I * GridSize), float32(J * GridSize)})
				CacheInSectors(he)
			}

			he.Direction = -1
			if len(he.CurrentPath) > 1 {
				he.CurrentPath = he.CurrentPath[1:]
			} else {
				he.State = StateWaiting
			}
		}
	}
	ComChannel <- HumanComStruct{ID: he.BasicEntity.ID()}
}

func (he *HumanEntity) GetDirection(To engo.Point) int {
	x, y := int(he.Position.X)/GridSize, int(he.Position.Y)/GridSize
	X, Y := int(To.X)/GridSize, int(To.Y)/GridSize

	dir := 0
	if X > x {
		dir += 2
	} else if X == x {
		dir += 1
	}

	if Y > y {
		dir += 20
	} else if Y == y {
		dir += 10
	}

	if dir == 11 {
		panic("To path same as prev!")
	}
	//fmt.Println(x, ",", y, "|", X, ",", Y, "||", dir)
	return dir
}

type HumanDetails struct {
	Name      string
	MaxHealth int
	Attack    int
	Texture   common.Drawable
	Width     float32
	Height    float32
}

type HumanComStruct struct {
	RunUpdate bool
	ID        uint64
}

func PointToGrid(p engo.Point) grid {
	return grid{x: int(p.X) / GridSize, y: int(p.Y) / GridSize}
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
	// var diagCost, sideCost float32
	// diagCost, sideCost = 1, 1

	var a, b int
	X := endgrid.x
	Y := endgrid.y
	a = int(math.Abs(float64(X - x)))
	b = int(math.Abs(float64(Y - y)))
	if a > b {
		a, b = b, a
	}
	return float32(a + b)
	//return b
}

func eval(neighbor *grid, block *grid, endgrid *grid, h *gridHeap, list *[][]bool) bool {
	if neighbor.x == endgrid.x && neighbor.y == endgrid.y {
		endgrid.par = block
		(*list)[endgrid.x][endgrid.y] = true
		return true
	}
	if neighbor.x == block.x || neighbor.y == block.y {
		neighbor.g = block.g + 10
	} else {
		neighbor.g = block.g + 14
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

	list := make([][]bool, GridMaxX)
	for i, _ := range list {
		list[i] = make([]bool, GridMaxY)
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
	c <- ReversePath(path)

}

func DrawPathBlock(x, y int, col color.RGBA) {
	myblock := GridEntity{
		BasicEntity: ecs.NewBasic(),
		SpaceComponent: common.SpaceComponent{
			Position: engo.Point{float32(x * GridSize), float32(y * GridSize)},
			Width:    float32(GridSize),
			Height:   float32(GridSize),
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
