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
)

type ActiveSystemsStruct struct {
	RenderSys *common.RenderSystem
	MouseSys  *common.MouseSystem
	CameraSys *common.CameraSystem
}

// Other Variables
var (
	ActiveSystems ActiveSystemsStruct
	PlayerFood    int
	PlayerWood    int
	PlayerPop     int
)

// Message Structs
type BuildingMessage struct {
	Action string
	Name   string
	Index  int
	ID     uint64
}

func (BuildingMessage) Type() string {
	return "BuildingMessage"
}

type HealthEnquiryMessage struct {
	ID uint64
}

func (HealthEnquiryMessage) Type() string {
	return "HealthEnquiryMessage"
}

type HealthEnquiryResponseStruct struct {
	HealthResult int
	set          bool
}

var HealthEnquiryResponse HealthEnquiryResponseStruct

func RegisterButtons() {
	engo.Input.RegisterButton(GridToggle, engo.Tab)
	engo.Input.RegisterButton(SpaceButton, engo.Space)
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
			fmt.Println("Found Camera System")
		}
	}

	fmt.Println("Cached Important System References")
}

func InitializeVariables() {
	PlayerFood = 100
	PlayerWood = 50
	PlayerPop = 0

	HealthEnquiryResponse = HealthEnquiryResponseStruct{set: false}
}
