package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
)

// Button mappings
var (
	GridToggle  = "gridtoggle"
	HorAxis     = "horAxis"
	VertAxis    = "vertAxis"
	SpaceButton = "SpaceButton"
	ShiftKey    = "shiftkey"
	R_remove    = "R_remove"
	SaveKey     = "savekey"
)

type ActiveSystemsStruct struct {
	RenderSys *common.RenderSystem
	MouseSys  *common.MouseSystem
	CameraSys *common.CameraSystem
}

// File Names

var (
	TownCenterSprite        = "towncentre1.png"
	ETownCenterSprite       = "towncentre2.png"
	MilitaryBlockSprite     = "militarybuilding1.png"
	EMilitaryBlockSprite    = "militarybuilding2.png"
	ResourceBuildingSprite  = "Resourcebuilding1.png"
	EResourceBuildingSprite = "Resourcebuilding2.png"
	HouseSprite             = "house1.png"
	EHouseSprite            = "house2.png"
	BushSprite              = "bush.png"
	TreeSprite              = "tree.png"
	WarriorSprite           = "warrior.png"
	EWarriorSprite          = "Ewarrior.png"
	BuildingSprites         = []string{TownCenterSprite, ETownCenterSprite, MilitaryBlockSprite, EMilitaryBlockSprite, ResourceBuildingSprite, EResourceBuildingSprite, HouseSprite, EHouseSprite, BushSprite}
)

// Other Variables
var (
	ActiveSystems ActiveSystemsStruct
	PlayerFood    int
	PlayerWood    int
	PlayerPop     int
	GridSize      int
	ScaleFactor   float32 // Ratio of Game World size with respect to Window Size
	Sectors       [][]*HumanEntity
	SectorSize    int
	Chunks        [][]StaticEntity
	ChunkSize     int
	Grid          [][]bool

	GridMaxX int
	GridMaxY int

	ShowDebugPathfinding bool
)

// Message Structs
type SetBottomHUDMessage struct {
	Action string
	Name   string
	Index  int
	ID     uint64
}

func (SetBottomHUDMessage) Type() string {
	return "SetBottomHUDMessage"
}

type CreateBuildingMessage struct {
	Name     string
	Position engo.Point
}

func (CreateBuildingMessage) Type() string {
	return "CreateBuildingMessage"
}

type DestroyBuildingMessage struct {
	obj StaticEntity
}

func (DestroyBuildingMessage) Type() string {
	return "DestroyBuildingMessage"
}

type BuildingHealthEnquiryMessage struct {
	ID uint64
}

func (BuildingHealthEnquiryMessage) Type() string {
	return "BuildingHealthEnquiryMessage"
}

type HumanHealthEnquiryMessage struct {
	ID uint64
}

func (HumanHealthEnquiryMessage) Type() string {
	return "HumanHealthEnquiryMessage"
}

type HealthEnquiryResponseStruct struct {
	HealthResult int
	ResourceName string
	set          bool
}

var HealthEnquiryResponse HealthEnquiryResponseStruct

type SaveMapMessage struct {
	Fname string
}

func (SaveMapMessage) Type() string {
	return "SaveMapMessage"
}

type CheckAndRemoveHUDMessage struct {
	ID uint64
}

func (CheckAndRemoveHUDMessage) Type() string {
	return "CheckAndRemoveHUDMessage"
}

type CreateHumanMessage struct {
	Name string
}

func (CreateHumanMessage) Type() string {
	return "CreateHumanMessage"
}

//Other types

type Fillable interface {
	GetPos() (float32, float32)
	GetSize() (float32, float32)
}

type StaticEntity interface {
	GetPos() (float32, float32)
	GetSize() (float32, float32)
	GetStaticComponent() *StaticComponent
}

// Functions

// Get the mouse position adjusted for zoom
func GetAdjustedMousePos(WRTWindow bool) (float32, float32) {
	CamSys := ActiveSystems.CameraSys
	x := engo.Input.Mouse.X * CamSys.Z() * (engo.GameWidth() / engo.CanvasWidth())
	y := engo.Input.Mouse.Y * CamSys.Z() * (engo.GameHeight() / engo.CanvasHeight())

	if !WRTWindow {
		x += CamSys.X() - (engo.GameWidth()/2)*CamSys.Z()
		y += CamSys.Y() - (engo.GameHeight()/2)*CamSys.Z()
	}

	return x, y

}

// Return chunk that the grid point belongs to
func GetChunkFromPos(x, y float32) (*[]StaticEntity, int) {
	rownum := int(engo.WindowWidth()*ScaleFactor) / (GridSize * ChunkSize)
	X := int(x) / (GridSize * ChunkSize)
	Y := int(y) / (GridSize * ChunkSize)

	//fmt.Println("Chunk of index", X, ",", Y)
	//fmt.Println("Row number is", rownum)

	return &Chunks[Y*rownum+X], (Y*rownum + X)
}

// Store Static objects in respective Chunk(s)
func CacheInChunks(se StaticEntity) {
	x, y := se.GetPos()
	X, Y := se.GetSize()
	X = X + x
	Y = Y + y

	chunk1, _ := GetChunkFromPos(x, y)
	*chunk1 = append(*chunk1, se)

	chunk2, _ := GetChunkFromPos(X, y)
	if chunk1 != chunk2 {
		*chunk2 = append(*chunk2, se)
	}

	chunk3, _ := GetChunkFromPos(x, Y)
	if chunk1 != chunk3 {
		*chunk3 = append(*chunk3, se)
	}

	chunk4, _ := GetChunkFromPos(X, Y)
	if chunk4 != chunk2 && chunk4 != chunk3 {
		*chunk4 = append(*chunk4, se)
	}
}

func UnCacheInChunks(se StaticEntity) {
	x, y := se.GetPos()
	X, Y := se.GetSize()
	X = X + x
	Y = Y + y

	chunk1, _ := GetChunkFromPos(x, y)
	for i, _ := range *chunk1 {
		entity := (*chunk1)[i]
		if entity.GetStaticComponent().ID() == se.GetStaticComponent().ID() {
			(*chunk1)[i] = (*chunk1)[len(*chunk1)-1]
			*chunk1 = (*chunk1)[:len(*chunk1)-1]
			break
		}
	}

	chunk2, _ := GetChunkFromPos(X, y)
	if chunk1 != chunk2 {
		for i, _ := range *chunk2 {
			entity := (*chunk2)[i]
			if entity.GetStaticComponent().ID() == se.GetStaticComponent().ID() {
				(*chunk2)[i] = (*chunk2)[len(*chunk2)-1]
				*chunk2 = (*chunk2)[:len(*chunk2)-1]
				break
			}
		}
	}

	chunk3, _ := GetChunkFromPos(x, Y)
	if chunk1 != chunk3 {
		for i, _ := range *chunk3 {
			entity := (*chunk3)[i]
			if entity.GetStaticComponent().ID() == se.GetStaticComponent().ID() {
				(*chunk3)[i] = (*chunk3)[len(*chunk3)-1]
				*chunk3 = (*chunk3)[:len(*chunk3)-1]
				break
			}
		}
	}

	chunk4, _ := GetChunkFromPos(X, Y)
	if chunk4 != chunk2 && chunk4 != chunk3 {
		for i, _ := range *chunk4 {
			entity := (*chunk4)[i]
			if entity.GetStaticComponent().ID() == se.GetStaticComponent().ID() {
				(*chunk4)[i] = (*chunk4)[len(*chunk4)-1]
				*chunk4 = (*chunk4)[:len(*chunk4)-1]
				break
			}
		}
	}
}

func GetSectorFromPos(x, y float32) (*[]*HumanEntity, int) {
	rownum := int(engo.WindowWidth()*ScaleFactor) / (GridSize * ChunkSize)
	X := int(x) / (GridSize * SectorSize)
	Y := int(y) / (GridSize * SectorSize)

	//fmt.Println("Chunk of index", X, ",", Y)
	//fmt.Println("Row number is", rownum)

	return &Sectors[Y*rownum+X], (Y*rownum + X)
}

// Store Static objects in respective Chunk(s)
func CacheInSectors(he *HumanEntity) {
	x, y := he.Position.X, he.Position.Y

	sector, _ := GetSectorFromPos(x, y)
	//fmt.Println("Appending a", he.Name, "in Sec ", si)
	*sector = append(*sector, he)
}

func UnCacheInSectors(he *HumanEntity, pos engo.Point) {
	x, y := pos.X, pos.Y

	sector, _ := GetSectorFromPos(x, y)
	for i, _ := range *sector {
		entity := (*sector)[i]
		if entity.ID() == he.ID() {
			(*sector)[i] = (*sector)[len(*sector)-1]
			*sector = (*sector)[:len(*sector)-1]
			break
		}
	}
}

func GetGridAtPos(x, y float32) bool {
	return Grid[int(x)/GridSize][int(y)/GridSize]
}

func GetCenterOfGrid(x, y int) engo.Point {
	X, Y := x*GridSize, y*GridSize

	return engo.Point{float32(X + 16), float32(Y + 16)}
}

func WithinGameWindow(x, y float32) bool {
	CamSys := ActiveSystems.CameraSys
	cx, cy := CamSys.X()-engo.WindowWidth()/2, CamSys.Y()-engo.WindowHeight()/2
	ymin := cy + 64
	ymax := cy + engo.WindowHeight() - 160

	return (cx <= x && x <= cx+engo.WindowWidth() && ymin <= y && y <= ymax)
}

// Returns mouse over object?, given button pressed?, given button released?
func StaticMouseCollision(obj StaticEntity, mb engo.MouseButton) (bool, bool, bool) {
	mx, my := GetAdjustedMousePos(false)
	mp := engo.Point{mx, my}
	if WithinGameWindow(mx, my) {
		if obj.GetStaticComponent().Contains(mp) {
			pressed := engo.Input.Mouse.Action == engo.Press
			button_matched := engo.Input.Mouse.Button == mb
			released := engo.Input.Mouse.Action == engo.Release
			return true, (pressed && button_matched), (released && button_matched)
		}
	}
	return false, false, false
}

func GetStaticClicked() StaticEntity {
	mx, my := GetAdjustedMousePos(false)
	mp := engo.Point{mx, my}
	if WithinGameWindow(mx, my) {
		Chunk, _ := GetChunkFromPos(mx, my)
		for i, _ := range *Chunk {
			if (*Chunk)[i].GetStaticComponent().Contains(mp) {
				pressed := engo.Input.Mouse.Action == engo.Press
				button_matched := engo.Input.Mouse.Button == engo.MouseButtonLeft
				if pressed && button_matched {
					return (*Chunk)[i]
				}
			}
		}
	}
	return nil
}

func GetStaticHover() StaticEntity {
	mx, my := GetAdjustedMousePos(false)
	mp := engo.Point{mx, my}
	if WithinGameWindow(mx, my) {
		Chunk, _ := GetChunkFromPos(mx, my)
		for i, _ := range *Chunk {
			if (*Chunk)[i].GetStaticComponent().Contains(mp) {
				return (*Chunk)[i]
			}
		}
	}
	return nil
}

// Mark the solids in the Grid
func FillGrid(f Fillable, val bool) {
	x, y := f.GetPos()
	w, h := f.GetSize()
	minx, miny := (int(x) / GridSize), (int(y) / GridSize)
	maxx, maxy := int(x+w)/GridSize, int(y+h)/GridSize

	Grid[minx][miny] = val

	for i := minx + 1; i < maxx; i += 1 {
		for j := miny + 1; j < maxy; j += 1 {
			Grid[i][j] = val
		}
	}
}

func RegisterButtons() {
	engo.Input.RegisterButton(GridToggle, engo.Tab)
	engo.Input.RegisterButton(SpaceButton, engo.Space)
	engo.Input.RegisterButton(ShiftKey, engo.LeftShift)
	engo.Input.RegisterButton(R_remove, engo.R)
	engo.Input.RegisterButton(SaveKey, engo.M)
	engo.Input.RegisterAxis(HorAxis, engo.AxisKeyPair{engo.A, engo.D})
	engo.Input.RegisterAxis(VertAxis, engo.AxisKeyPair{engo.W, engo.S})

	fmt.Println("Registered Buttons")
}

func CacheActiveSystems(world *ecs.World) {
	for _, system := range world.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			ActiveSystems.RenderSys = sys
		case *common.MouseSystem:
			ActiveSystems.MouseSys = sys
		case *common.CameraSystem:
			ActiveSystems.CameraSys = sys
		}
	}

	fmt.Println("Cached Important System References")
}

func InitializeVariables() {
	PlayerFood = 100
	PlayerWood = 50
	PlayerPop = 0

	ShowDebugPathfinding = false

	ScaleFactor = 2

	HealthEnquiryResponse = HealthEnquiryResponseStruct{set: false}

	GridSize = 32

	// Camera bounds is ScaleFactor times window size, also Go defaults to false
	GridMaxX = int(engo.WindowWidth()*ScaleFactor) / GridSize
	GridMaxY = int(engo.WindowHeight()*ScaleFactor) / GridSize
	Grid = make([][]bool, GridMaxX)
	for i, _ := range Grid {
		Grid[i] = make([]bool, GridMaxY)
	}

	// Chunks used to Cache Static Entities
	ChunkSize = 8
	ChunkNum := (int(engo.WindowHeight()*ScaleFactor) / (GridSize * ChunkSize)) * (int(engo.WindowWidth()*ScaleFactor) / (GridSize * ChunkSize))

	Chunks = make([][]StaticEntity, ChunkNum)
	for i, _ := range Chunks {
		Chunks[i] = make([]StaticEntity, 0)
	}

	// Sectors used to keep track of moving entities(i.e. humans)
	SectorSize = 8
	SectorNum := (int(engo.WindowHeight()*ScaleFactor) / (GridSize * SectorSize)) * (int(engo.WindowWidth()*ScaleFactor) / (GridSize * SectorSize))

	Sectors = make([][]*HumanEntity, SectorNum)
	for i, _ := range Sectors {
		Sectors[i] = make([]*HumanEntity, 0)
	}
}
