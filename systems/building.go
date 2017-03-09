package systems

import (
	"bufio"
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

var (
	BuildingDetailsMap map[string]BuildingDetails

	SetHUD = "setHUD"

	PathChannel chan []grid
)

type BuildingSystem struct {
	world     *ecs.World
	Buildings []*BuildingEntity
}

func (bs *BuildingSystem) Remove(ecs.BasicEntity) {}

func (bs *BuildingSystem) New(w *ecs.World) {
	rand.Seed(16548161)
	bs.world = w
	PathChannel = make(chan []grid)

	//Building Definitions (For loop to be able to collapse it)
	BuildingDetailsMap = make(map[string]BuildingDetails)
	func() {
		TownCenterTexture, err := common.LoadedSprite(TownCenterSprite)
		if err != nil {
			fmt.Println(err.Error())
		}
		TownCenterDetails := BuildingDetails{
			Name: "Town Center", MaxHealth: 150, Texture: TownCenterTexture,
		}
		BuildingDetailsMap[TownCenterDetails.Name] = TownCenterDetails

		MilitaryBlockTexture, err := common.LoadedSprite(MilitaryBlockSprite)
		if err != nil {
			fmt.Println(err.Error())
		}
		MilitaryBlockDetails := BuildingDetails{
			Name: "Military Block", MaxHealth: 120, Texture: MilitaryBlockTexture,
		}
		BuildingDetailsMap[MilitaryBlockDetails.Name] = MilitaryBlockDetails

		ResourceBuildingTexture, err := common.LoadedSprite(ResourceBuildingSprite)
		if err != nil {
			fmt.Println(err.Error())
		}
		ResourceBuildingDetails := BuildingDetails{
			Name: "Resource Building", MaxHealth: 75, Texture: ResourceBuildingTexture,
		}
		BuildingDetailsMap[ResourceBuildingDetails.Name] = ResourceBuildingDetails

		HouseTexture, err := common.LoadedSprite(HouseSprite)
		if err != nil {
			fmt.Println(err.Error())
		}
		HouseDetails := BuildingDetails{
			Name: "House", MaxHealth: 30, Texture: HouseTexture,
		}
		BuildingDetailsMap[HouseDetails.Name] = HouseDetails

		BushTexture, err := common.LoadedSprite(BushSprite)
		if err != nil {
			fmt.Println(err.Error())
		}
		BushDetails := BuildingDetails{
			Name: "Bush", MaxHealth: 50, Texture: BushTexture,
		}
		BuildingDetailsMap[BushDetails.Name] = BushDetails

		TreeTexture, err := common.LoadedSprite(TreeSprite)
		if err != nil {
			fmt.Println(err.Error())
		}
		TreeDetails := BuildingDetails{
			Name: "Tree", MaxHealth: 30, Texture: TreeTexture,
		}
		BuildingDetailsMap[TreeDetails.Name] = TreeDetails

	}()

	engo.Mailbox.Listen("HealthEnquiryMessage", func(_msg engo.Message) {
		msg, ok := _msg.(HealthEnquiryMessage)
		if !ok {
			panic("Building System expected HealthEnquiryMessage, instead got unexpected")
		}
		for _, item := range bs.Buildings {
			if item.BasicEntity.ID() == msg.ID {
				HealthEnquiryResponse.HealthResult = item.Health
				switch item.Name {
				case "Bush":
					HealthEnquiryResponse.ResourceName = "Food"
				case "Tree":
					HealthEnquiryResponse.ResourceName = "Wood"
				}
				HealthEnquiryResponse.set = true
				return
			}
		}

		panic("Health Enquiry for unkown building")
	})

	engo.Mailbox.Listen("CreateBuildingMessage", func(_msg engo.Message) {
		msg, ok := _msg.(CreateBuildingMessage)
		if !ok {
			panic("Building System expected CreateBuildingMessage, instead got unexpected")
		}

		bs.AddBuilding(msg.Name, msg.Position)
	})

	engo.Mailbox.Listen("DestroyBuildingMessage", func(_msg engo.Message) {
		msg, ok := _msg.(DestroyBuildingMessage)
		if !ok {
			panic("Building System expected DestroyBuildingMessage, instead got unexpected")
		}

		bs.RemoveBuilding(msg.obj)
	})

	engo.Mailbox.Listen("SaveMapMessage", func(_msg engo.Message) {
		msg, ok := _msg.(SaveMapMessage)

		if !ok {
			panic("Building System expected SaveMapMessage, instead got unexpected")
		}

		bs.SaveMap(msg.Fname)
	})

	bs.LoadMap("World.mapfile")

	fmt.Println("Building System Initialized")
}

func (bs *BuildingSystem) Update(dt float32) {
	// Handling of clicking building
	mx, my := GetAdjustedMousePos(false)
	mp := engo.Point{mx, my}

	// Debug info with middle mouse click
	func() {
		if engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonLeft {
			ChunkRef, _ := GetChunkFromPos(mx, my)
			Chunk := *ChunkRef

			if len(Chunk) > 0 {
				//fmt.Println("-------------------------")
				for _, item := range Chunk {
					sc := item.GetStaticComponent()
					if sc.Contains(mp) {
						engo.Mailbox.Dispatch(BuildingMessage{ID: sc.BasicEntity.ID(), Name: sc.Name, Index: 0})
					}
					//fmt.Println(item.GetStaticComponent().Name, "present in chunk:", ChunkIndex)
				}
				//fmt.Println("-------------------------")
			} else {
				//fmt.Println("Chunk", ChunkIndex, "Empty")
			}
		}
	}()

	if engo.Input.Button(SpaceButton).JustReleased() {
		bs.Buildings[int(math.Floor(rand.Float64()*float64(len(bs.Buildings))))].Health -= 10
	}

}

func (bs *BuildingSystem) AddBuilding(_Name string, Pos engo.Point) {
	tex := BuildingDetailsMap[_Name].Texture

	// Using reference so that the newly created building
	// doesn't get garbage collected after func return
	new_building := &BuildingEntity{
		StaticComponent: StaticComponent{
			BasicEntity: ecs.NewBasic(),
			RenderComponent: common.RenderComponent{
				Drawable: tex,
			},
			SpaceComponent: common.SpaceComponent{
				Position: Pos,
				Width:    tex.Width(),
				Height:   tex.Height(),
			},
			Name:   _Name,
			Width:  tex.Width(),
			Height: tex.Height(),
		},
		Health: BuildingDetailsMap[_Name].MaxHealth,
	}

	bs.Buildings = append(bs.Buildings, new_building)
	CacheInChunks(new_building)
	FillGrid(new_building, true)

	ActiveSystems.RenderSys.Add(&new_building.BasicEntity, &new_building.RenderComponent, &new_building.SpaceComponent)
}

func (bs *BuildingSystem) RemoveBuilding(obj StaticEntity) {
	ActiveSystems.RenderSys.Remove(obj.GetStaticComponent().BasicEntity)
	engo.Mailbox.Dispatch(CheckAndRemoveHUDMessage{ID: obj.GetStaticComponent().ID()})

	for i, _ := range bs.Buildings {
		building := bs.Buildings[i]
		if building.ID() == obj.GetStaticComponent().ID() {
			bs.Buildings[i] = bs.Buildings[len(bs.Buildings)-1]
			bs.Buildings = bs.Buildings[:len(bs.Buildings)-1]
			break
		}
	}
	UnCacheInChunks(obj)
	FillGrid(obj, false)
}

func (bs *BuildingSystem) SaveMap(Fname string) {
	fmt.Print("Saving Map! ... ")
	file, err := os.OpenFile(Fname, os.O_WRONLY|os.O_CREATE, 0777)
	defer file.Close()

	if err != nil {
		panic(err)
	}

	Writer := bufio.NewWriter(file)
	for _, item := range bs.Buildings {
		_, e := Writer.WriteString(item.Name + "," + strconv.Itoa(int(item.Position.X)) + "," + strconv.Itoa(int(item.Position.Y)) + "\n")
		if e != nil {
			panic(e)
		}
	}
	err = Writer.Flush()
	if err != nil {
		panic(err)
	}

	fmt.Println("Saved!")
}

func (bs *BuildingSystem) LoadMap(Fname string) {
	fmt.Print("Loading Map! ... ")
	file, err := os.Open(Fname)
	defer file.Close()

	if err != nil {
		panic(err)
	}

	Scanner := bufio.NewScanner(file)

	for Scanner.Scan() {
		line := Scanner.Text()
		args := strings.Split(line, ",")

		x, errx := strconv.ParseFloat(args[1], 32)
		y, erry := strconv.ParseFloat(args[2], 32)
		if errx != nil || erry != nil {
			panic("Could not convert " + args[1] + "," + args[2])
		}

		bs.AddBuilding(args[0], engo.Point{float32(x), float32(y)})
	}

}

type StaticComponent struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	Name   string
	Width  float32
	Height float32
}

func (se *StaticComponent) GetPos() (float32, float32) {
	return se.Position.X, se.Position.Y
}

func (se *StaticComponent) GetSize() (float32, float32) {
	return se.Width, se.Height
}

func (se *StaticComponent) GetStaticComponent() *StaticComponent {
	return se
}

type BuildingEntity struct {
	StaticComponent
	Health int
}

func (be *BuildingEntity) GetDetails() BuildingDetails {
	return BuildingDetailsMap[be.Name]
}

type BuildingDetails struct {
	Name      string
	MaxHealth int
	Texture   common.Drawable
}
