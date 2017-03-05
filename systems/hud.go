package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"

	"fmt"
	"image/color"
	"strconv"
)

type HUDSystem struct {
	World              *ecs.World
	CurrentActiveLabel *LabelGroup
	CurrentLabelIndex  int

	SelectionRect *SHAPE
	ActionRects   []*SHAPE

	BottomHUDWidth  int
	BottomHUDHeight int
	TopHUDWidth     int
	TopHUDHeight    int

	FoodValLabelPos engo.Point
	WoodValLabelPos engo.Point
	PrevFoodVal     int
	PrevWoodVal     int
}

type HUD struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type SHAPE struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type Label struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

type LabelGroup struct {
	Name string

	DescriptionLabel Label
	ActionLabels     [][]Label
	DynamicLabels    []Label
}

var (
	TownCenterLabels, MilitaryBlockLabels, ResouceBuildingLabels,
	HouseLabels, VillagerLabels LabelGroup

	LabelGroupMap map[string]LabelGroup
)

func (hs *HUDSystem) New(w *ecs.World) {
	hs.World = w
	hs.CurrentActiveLabel = nil
	hs.CurrentLabelIndex = 0
	hs.SelectionRect = nil

	HUDColor := color.RGBA{222, 184, 135, 250}

	//Render Top and Bottom HUD Backgrounds
	var (
		TopHud    HUD
		BottomHud HUD
	)
	func() {
		//Bottom Hud Definition
		hs.BottomHUDWidth = int(engo.WindowWidth())
		hs.BottomHUDHeight = 160

		BottomHud = HUD{
			BasicEntity: ecs.NewBasic(),
			RenderComponent: common.RenderComponent{
				Drawable: common.Rectangle{},
				Color:    HUDColor,
			},
			SpaceComponent: common.SpaceComponent{
				Position: engo.Point{0, engo.WindowHeight() - float32(hs.BottomHUDHeight)},
				Width:    float32(hs.BottomHUDWidth),
				Height:   float32(hs.BottomHUDHeight),
			},
		}

		BottomHud.RenderComponent.SetZIndex(100)
		BottomHud.RenderComponent.SetShader(common.HUDShader)
		ActiveSystems.RenderSys.Add(&BottomHud.BasicEntity, &BottomHud.RenderComponent, &BottomHud.SpaceComponent)

		//Top Hud Definition
		hs.TopHUDWidth = int(engo.WindowWidth())
		hs.TopHUDHeight = 64

		TopHud = HUD{
			BasicEntity: ecs.NewBasic(),
			RenderComponent: common.RenderComponent{
				Drawable: common.Rectangle{},
				Color:    HUDColor,
			},
			SpaceComponent: common.SpaceComponent{
				Position: engo.Point{0, 0},
				Width:    float32(hs.TopHUDWidth),
				Height:   float32(hs.TopHUDHeight),
			},
		}

		TopHud.RenderComponent.SetZIndex(100)
		TopHud.RenderComponent.SetShader(common.HUDShader)
		ActiveSystems.RenderSys.Add(&TopHud.BasicEntity, &TopHud.RenderComponent, &TopHud.SpaceComponent)
	}()

	// -----------------------------------------------------------------------------------------------------
	// -----------------------------------------------------------------------------------------------------
	// -----------------------------------------------------------------------------------------------------

	// Define all the Rectangles that the labels are displayed over
	var (
		DescriptionRect SHAPE
		Action1Rect     SHAPE
	)
	func() {
		//Bottom HUD Rectangles
		DescriptionRect = SHAPE{BasicEntity: ecs.NewBasic()} //First Big Rectangle
		DescriptionRect.SpaceComponent = common.SpaceComponent{Position: engo.Point{15, engo.WindowHeight() - float32(hs.BottomHUDHeight-15)}, Width: float32((hs.BottomHUDWidth / 3) - 80), Height: float32((hs.BottomHUDHeight) - 30)}
		DescriptionRect.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

		DescriptionRect.RenderComponent.SetZIndex(125)
		DescriptionRect.RenderComponent.SetShader(common.HUDShader)
		ActiveSystems.RenderSys.Add(&DescriptionRect.BasicEntity, &DescriptionRect.RenderComponent, &DescriptionRect.SpaceComponent)

		Action1Rect = SHAPE{BasicEntity: ecs.NewBasic()} //R2 , R3, R4 are small Rectangles
		Action1Rect.SpaceComponent = common.SpaceComponent{Position: engo.Point{15 + DescriptionRect.SpaceComponent.Width + 80, DescriptionRect.SpaceComponent.Position.Y}, Width: float32(DescriptionRect.SpaceComponent.Width / 3), Height: float32(DescriptionRect.SpaceComponent.Height/2) - 5}
		Action1Rect.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

		Action1Rect.RenderComponent.SetZIndex(125)
		Action1Rect.RenderComponent.SetShader(common.HUDShader)

		Action2Rect := SHAPE{BasicEntity: ecs.NewBasic()}
		Action2Rect.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action1Rect.SpaceComponent.Position.X + Action1Rect.SpaceComponent.Width + 20, Action1Rect.SpaceComponent.Position.Y}, Width: Action1Rect.SpaceComponent.Width, Height: Action1Rect.SpaceComponent.Height}
		Action2Rect.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

		Action2Rect.RenderComponent.SetZIndex(125)
		Action2Rect.RenderComponent.SetShader(common.HUDShader)

		Action3Rect := SHAPE{BasicEntity: ecs.NewBasic()}
		Action3Rect.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action2Rect.SpaceComponent.Position.X + Action1Rect.SpaceComponent.Width + 20, Action1Rect.SpaceComponent.Position.Y}, Width: Action1Rect.SpaceComponent.Width, Height: Action1Rect.SpaceComponent.Height}
		Action3Rect.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

		Action3Rect.RenderComponent.SetZIndex(125)
		Action3Rect.RenderComponent.SetShader(common.HUDShader)

		Action4Rect := SHAPE{BasicEntity: ecs.NewBasic()}
		Action4Rect.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action1Rect.SpaceComponent.Position.X, Action1Rect.SpaceComponent.Position.Y + Action1Rect.Height + 10}, Width: Action1Rect.SpaceComponent.Width, Height: Action1Rect.SpaceComponent.Height}
		Action4Rect.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

		Action4Rect.RenderComponent.SetZIndex(125)
		Action4Rect.RenderComponent.SetShader(common.HUDShader)

		Action5Rect := SHAPE{BasicEntity: ecs.NewBasic()}
		Action5Rect.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action4Rect.SpaceComponent.Position.X + Action4Rect.SpaceComponent.Width + 20, Action4Rect.SpaceComponent.Position.Y}, Width: Action4Rect.SpaceComponent.Width, Height: Action4Rect.SpaceComponent.Height}
		Action5Rect.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

		Action5Rect.RenderComponent.SetZIndex(125)
		Action5Rect.RenderComponent.SetShader(common.HUDShader)

		Action6Rect := SHAPE{BasicEntity: ecs.NewBasic()}
		Action6Rect.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action5Rect.SpaceComponent.Position.X + Action5Rect.SpaceComponent.Width + 20, Action5Rect.SpaceComponent.Position.Y}, Width: Action5Rect.SpaceComponent.Width, Height: Action5Rect.SpaceComponent.Height}
		Action6Rect.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

		Action6Rect.RenderComponent.SetZIndex(125)
		Action6Rect.RenderComponent.SetShader(common.HUDShader)

		ActiveSystems.RenderSys.Add(&Action1Rect.BasicEntity, &Action1Rect.RenderComponent, &Action1Rect.SpaceComponent)
		ActiveSystems.RenderSys.Add(&Action2Rect.BasicEntity, &Action2Rect.RenderComponent, &Action2Rect.SpaceComponent)
		ActiveSystems.RenderSys.Add(&Action3Rect.BasicEntity, &Action3Rect.RenderComponent, &Action3Rect.SpaceComponent)
		ActiveSystems.RenderSys.Add(&Action4Rect.BasicEntity, &Action4Rect.RenderComponent, &Action4Rect.SpaceComponent)
		ActiveSystems.RenderSys.Add(&Action5Rect.BasicEntity, &Action5Rect.RenderComponent, &Action5Rect.SpaceComponent)
		ActiveSystems.RenderSys.Add(&Action6Rect.BasicEntity, &Action6Rect.RenderComponent, &Action6Rect.SpaceComponent)
		hs.ActionRects = append(make([]*SHAPE, 0),
			&Action1Rect, &Action2Rect, &Action3Rect, &Action4Rect, &Action5Rect, &Action6Rect,
		)

		DeselectRect := SHAPE{BasicEntity: ecs.NewBasic()}
		wid := Action1Rect.SpaceComponent.Width - 30
		hig := float32(hs.BottomHUDHeight/2) - 15
		DeselectRect.SpaceComponent = common.SpaceComponent{Position: engo.Point{engo.WindowWidth() - wid - 20, engo.WindowHeight() - float32(hs.BottomHUDHeight) + 10}, Width: wid, Height: hig}
		DeselectRect.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

		DeselectRect.RenderComponent.SetZIndex(125)
		DeselectRect.RenderComponent.SetShader(common.HUDShader)

		HelpRect := SHAPE{BasicEntity: ecs.NewBasic()}
		HelpRect.SpaceComponent = common.SpaceComponent{Position: engo.Point{DeselectRect.SpaceComponent.Position.X, engo.WindowHeight() - float32(hs.BottomHUDHeight/2) + 5}, Width: wid, Height: hig}
		HelpRect.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

		HelpRect.RenderComponent.SetZIndex(125)
		HelpRect.RenderComponent.SetShader(common.HUDShader)

		ActiveSystems.RenderSys.Add(&DeselectRect.BasicEntity, &DeselectRect.RenderComponent, &DeselectRect.SpaceComponent)
		ActiveSystems.RenderSys.Add(&HelpRect.BasicEntity, &HelpRect.RenderComponent, &HelpRect.SpaceComponent)

		// Top HUD Rectangles
		FoodRect := SHAPE{BasicEntity: ecs.NewBasic()}
		FoodRect.SpaceComponent = common.SpaceComponent{Position: engo.Point{96, 16}, Width: 128, Height: TopHud.Height - 32}
		FoodRect.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

		FoodRect.RenderComponent.SetZIndex(125)
		FoodRect.RenderComponent.SetShader(common.HUDShader)

		WoodRect := SHAPE{BasicEntity: ecs.NewBasic()}
		WoodRect.SpaceComponent = common.SpaceComponent{Position: engo.Point{352, 16}, Width: 128, Height: TopHud.Height - 32}
		WoodRect.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

		WoodRect.RenderComponent.SetZIndex(125)
		WoodRect.RenderComponent.SetShader(common.HUDShader)

		ActiveSystems.RenderSys.Add(&FoodRect.BasicEntity, &FoodRect.RenderComponent, &FoodRect.SpaceComponent)
		ActiveSystems.RenderSys.Add(&WoodRect.BasicEntity, &WoodRect.RenderComponent, &WoodRect.SpaceComponent)
	}()

	// -----------------------------------------------------------------------------------------------------
	// -----------------------------------------------------------------------------------------------------
	// -----------------------------------------------------------------------------------------------------

	// Define all the Labels
	func() {
		fnt := &common.Font{
			URL:  "Roboto-Regular.ttf",
			FG:   color.Black,
			Size: 16,
		}

		err := fnt.CreatePreloaded()
		if err != nil {
			panic(err)
		}

		//Top HUD Labels
		FoodTitleLabel := Label{BasicEntity: ecs.NewBasic()}
		FoodTitleLabel.SpaceComponent = common.SpaceComponent{Position: engo.Point{32, 24}}
		FoodTitleLabel.RenderComponent.Drawable = common.Text{
			Font: fnt,
			Text: "FOOD :",
		}
		FoodTitleLabel.SetShader(common.TextHUDShader)
		FoodTitleLabel.SetZIndex(150)

		WoodTitleLabel := Label{BasicEntity: ecs.NewBasic()}
		WoodTitleLabel.SpaceComponent = common.SpaceComponent{Position: engo.Point{288, 24}}
		WoodTitleLabel.RenderComponent.Drawable = common.Text{
			Font: fnt,
			Text: "WOOD :",
		}
		WoodTitleLabel.SetShader(common.TextHUDShader)
		WoodTitleLabel.SetZIndex(150)

		hs.FoodValLabelPos = engo.Point{104, 24}
		FoodValLabel := Label{BasicEntity: ecs.NewBasic()}
		FoodValLabel.SpaceComponent = common.SpaceComponent{Position: hs.FoodValLabelPos}
		FoodValLabel.RenderComponent.Drawable = common.Text{
			Font: fnt,
			Text: strconv.Itoa(PlayerFood),
		}
		FoodValLabel.SetShader(common.TextHUDShader)
		FoodValLabel.SetZIndex(150)

		hs.WoodValLabelPos = engo.Point{360, 24}
		WoodValLabel := Label{BasicEntity: ecs.NewBasic()}
		WoodValLabel.SpaceComponent = common.SpaceComponent{Position: hs.WoodValLabelPos}
		WoodValLabel.RenderComponent.Drawable = common.Text{
			Font: fnt,
			Text: strconv.Itoa(PlayerWood),
		}
		WoodValLabel.SetShader(common.TextHUDShader)
		WoodValLabel.SetZIndex(150)

		ActiveSystems.RenderSys.Add(&FoodTitleLabel.BasicEntity, &FoodTitleLabel.RenderComponent, &FoodTitleLabel.SpaceComponent)
		ActiveSystems.RenderSys.Add(&FoodValLabel.BasicEntity, &FoodValLabel.RenderComponent, &FoodValLabel.SpaceComponent)
		ActiveSystems.RenderSys.Add(&WoodTitleLabel.BasicEntity, &WoodTitleLabel.RenderComponent, &WoodTitleLabel.SpaceComponent)
		ActiveSystems.RenderSys.Add(&WoodValLabel.BasicEntity, &WoodValLabel.RenderComponent, &WoodValLabel.SpaceComponent)

		//Bottom HUD Labels
		var temp1, temp2 Label

		// -----------------------------------------------------------------------------------------------------

		temp1 = Label{BasicEntity: ecs.NewBasic()}
		temp1.SpaceComponent = common.SpaceComponent{Position: engo.Point{DescriptionRect.SpaceComponent.Position.X + 48, DescriptionRect.SpaceComponent.Position.Y + 32}}
		temp1.RenderComponent.Drawable = common.Text{Font: fnt, Text: "TOWN CENTRE\n\n\nHealth : XX/YY"}
		temp1.SetShader(common.TextHUDShader)
		temp1.SetZIndex(250)

		temp2 = Label{BasicEntity: ecs.NewBasic()}
		temp2.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action1Rect.SpaceComponent.Position.X + 16, Action1Rect.SpaceComponent.Position.Y + 16}}
		temp2.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Create Villager"}
		temp2.SetShader(common.TextHUDShader)
		temp2.SetZIndex(250)

		TownCenterLabels = LabelGroup{Name: "Town Center"}
		TownCenterLabels.DescriptionLabel = temp1
		TownCenterLabels.ActionLabels = append(make([][]Label, 0), make([]Label, 0))
		TownCenterLabels.ActionLabels[0] = append(TownCenterLabels.ActionLabels[0], temp2)

		// -----------------------------------------------------------------------------------------------------

		temp1 = Label{BasicEntity: ecs.NewBasic()}
		temp1.SpaceComponent = common.SpaceComponent{Position: engo.Point{DescriptionRect.SpaceComponent.Position.X + 48, DescriptionRect.SpaceComponent.Position.Y + 32}}
		temp1.RenderComponent.Drawable = common.Text{Font: fnt, Text: "HOUSE\n\n\nHealth : XX/YY\n\n\nCapacity : xx/yy"}
		temp1.SetShader(common.TextHUDShader)
		temp1.SetZIndex(250)

		HouseLabels = LabelGroup{Name: "House"}
		HouseLabels.DescriptionLabel = temp1

		// -----------------------------------------------------------------------------------------------------

		temp1 = Label{BasicEntity: ecs.NewBasic()}
		temp1.SpaceComponent = common.SpaceComponent{Position: engo.Point{DescriptionRect.SpaceComponent.Position.X + 48, DescriptionRect.SpaceComponent.Position.Y + 32}}
		temp1.RenderComponent.Drawable = common.Text{Font: fnt, Text: "MILITARY\n\n\nHealth : XX/YY"}
		temp1.SetShader(common.TextHUDShader)
		temp1.SetZIndex(250)

		temp2 = Label{BasicEntity: ecs.NewBasic()}
		temp2.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action1Rect.SpaceComponent.Position.X + 16, Action1Rect.SpaceComponent.Position.Y + 16}}
		temp2.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Create Warrior"}
		temp2.SetShader(common.TextHUDShader)
		temp2.SetZIndex(250)

		MilitaryBlockLabels = LabelGroup{Name: "Military Block"}
		MilitaryBlockLabels.DescriptionLabel = temp1
		MilitaryBlockLabels.ActionLabels = append(make([][]Label, 0), make([]Label, 0))
		MilitaryBlockLabels.ActionLabels[0] = append(MilitaryBlockLabels.ActionLabels[0], temp2)

		// -----------------------------------------------------------------------------------------------------

		temp1 = Label{BasicEntity: ecs.NewBasic()}
		temp1.SpaceComponent = common.SpaceComponent{Position: engo.Point{DescriptionRect.SpaceComponent.Position.X + 48, DescriptionRect.SpaceComponent.Position.Y + 32}}
		temp1.RenderComponent.Drawable = common.Text{Font: fnt, Text: "RESOURCE BUILDING\n\n\nHealth : XX/YY"}
		temp1.SetShader(common.TextHUDShader)
		temp1.SetZIndex(250)

		ResouceBuildingLabels = LabelGroup{Name: "Resource Building"}
		ResouceBuildingLabels.DescriptionLabel = temp1

		// -----------------------------------------------------------------------------------------------------

		LabelGroupMap = make(map[string]LabelGroup)
		LabelGroupMap["Town Center"] = TownCenterLabels
		LabelGroupMap["Military Block"] = MilitaryBlockLabels
		LabelGroupMap["Resource Building"] = ResouceBuildingLabels
		LabelGroupMap["House"] = HouseLabels
	}()

	// -----------------------------------------------------------------------------------------------------
	// -----------------------------------------------------------------------------------------------------
	// -----------------------------------------------------------------------------------------------------

	// Uncomment when proper items were implemented
	for {
		// temp1 = Label{BasicEntity: ecs.NewBasic()}
		// temp1.SpaceComponent = common.SpaceComponent{Position: engo.Point{DescriptionRect.SpaceComponent.Position.X + 48, DescriptionRect.SpaceComponent.Position.Y + 32}}
		// temp1.RenderComponent.Drawable = common.Text{Font: fnt, Text: "VILLAGER\n\n\nHealth : XX/YY"}
		// temp1.SetShader(common.TextHUDShader)
		// temp1.SetZIndex(250)

		// temp2 = Label{BasicEntity: ecs.NewBasic()}
		// temp2.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action1Rect.SpaceComponent.Position.X + 16, Action1Rect.SpaceComponent.Position.Y + 16}}
		// temp2.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Build"}
		// temp2.SetShader(common.TextHUDShader)
		// temp2.SetZIndex(250)

		// temp3 = Label{BasicEntity: ecs.NewBasic()}
		// temp3.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action2Rect.SpaceComponent.Position.X + 16, Action2Rect.SpaceComponent.Position.Y + 16}}
		// temp3.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Repair"}
		// temp3.SetShader(common.TextHUDShader)
		// temp3.SetZIndex(250)

		// VillagerLabels = LabelGroup{Name: "Villager"}
		// VillagerLabels.DescriptionLabel = temp1
		// VillagerLabels.ActionLabels = append(VillagerLabels.ActionLabels, temp2)
		// VillagerLabels.ActionLabels = append(VillagerLabels.ActionLabels, temp3)

		// lab15 := Label{BasicEntity: ecs.NewBasic()}
		// lab15.SpaceComponent = common.SpaceComponent{Position: engo.Point{DescriptionRect.SpaceComponent.Position.X + 48, DescriptionRect.SpaceComponent.Position.Y + 32}}
		// lab15.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Warrior"}
		// lab15.SetShader(common.TextHUDShader)
		// lab15.SetZIndex(250)

		// //If clicked on Build Then the following options are displayed

		// lab11 := Label{BasicEntity: ecs.NewBasic()}
		// lab11.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action1Rect.SpaceComponent.Position.X + 16, Action1Rect.SpaceComponent.Position.Y + 16}}
		// lab11.RenderComponent.Drawable = common.Text{Font: fnt, Text: "House"}
		// lab11.SetShader(common.TextHUDShader)
		// lab11.SetZIndex(250)

		// lab12 := Label{BasicEntity: ecs.NewBasic()}
		// lab12.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action2Rect.SpaceComponent.Position.X + 16, Action2Rect.SpaceComponent.Position.Y + 16}}
		// lab12.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Military Camp"}
		// lab12.SetShader(common.TextHUDShader)
		// lab12.SetZIndex(250)

		// lab13 := Label{BasicEntity: ecs.NewBasic()}
		// lab13.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action3Rect.SpaceComponent.Position.X + 16, Action3Rect.SpaceComponent.Position.Y + 16}}
		// lab13.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Resource"}
		// lab13.SetShader(common.TextHUDShader)
		// lab13.SetZIndex(250)

		// lab14 := Label{BasicEntity: ecs.NewBasic()}
		// lab14.SpaceComponent = common.SpaceComponent{Position: engo.Point{Action4Rect.SpaceComponent.Position.X + 16, Action4Rect.SpaceComponent.Position.Y + 16}}
		// lab14.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Go Back"}
		// lab14.SetShader(common.TextHUDShader)
		// lab14.SetZIndex(250)
		break
	}

	engo.Mailbox.Listen("BuildingMessage", func(_msg engo.Message) {
		msg, ok := _msg.(BuildingMessage)
		if !ok {
			panic("HUD recieved non BuildingMessage Message")
		}

		if hs.CurrentActiveLabel == nil ||
			(hs.CurrentActiveLabel.Name != msg.Name || hs.CurrentLabelIndex != msg.Index) {
			hs.SetBottomHUD(msg.Name, msg.Index)
		}
	})

	fmt.Println("HUD System Initialized")
}

func (hs *HUDSystem) Update(dt float32) {
	//Rendering Selection Rect
	func() {
		CamSys := ActiveSystems.CameraSys

		// Converting Mouse Coordinates to be Independent of Camera ZooM
		mx := engo.Input.Mouse.X * CamSys.Z() * (engo.GameWidth() / engo.CanvasWidth())
		my := engo.Input.Mouse.Y * CamSys.Z() * (engo.GameHeight() / engo.CanvasHeight())

		//If left Mouse button is pressed within Active Game Area
		if engo.Input.Mouse.Action == engo.Press && engo.Input.Mouse.Button == engo.MouseButtonLeft &&
			my > float32(hs.TopHUDHeight) && my < engo.WindowHeight()-float32(hs.BottomHUDHeight) {

			hs.SelectionRect = &SHAPE{
				BasicEntity: ecs.NewBasic(),
				SpaceComponent: common.SpaceComponent{
					Position: engo.Point{mx, my},
					Width:    0,
					Height:   0,
				},
				RenderComponent: common.RenderComponent{
					Drawable: common.Rectangle{
						BorderColor: color.RGBA{255, 255, 255, 255},
						BorderWidth: 2,
					},
					Color: color.RGBA{0, 0, 0, 0},
				},
			}

			hs.SelectionRect.RenderComponent.SetShader(common.HUDShader)
			hs.SelectionRect.RenderComponent.SetZIndex(200)
			ActiveSystems.RenderSys.Add(
				&hs.SelectionRect.BasicEntity, &hs.SelectionRect.RenderComponent,
				&hs.SelectionRect.SpaceComponent,
			)
		}

		// If Left Mouse Button is released
		if engo.Input.Mouse.Action == engo.Release && engo.Input.Mouse.Button == engo.MouseButtonLeft {
			if hs.SelectionRect != nil {
				ActiveSystems.RenderSys.Remove(hs.SelectionRect.BasicEntity)
				hs.SelectionRect = nil
			}
		}

		// While Left Mouse Button is held down
		if hs.SelectionRect != nil {
			// Clamp the mouse cooridnates to be within Active Game Area
			if my < float32(hs.TopHUDHeight) {
				my = float32(hs.TopHUDHeight)
			}
			if my > (engo.WindowHeight() - float32(hs.BottomHUDHeight)) {
				my = (engo.WindowHeight() - float32(hs.BottomHUDHeight))
			}

			if mx < 0 {
				mx = 0
			}
			if mx > engo.WindowWidth() {
				mx = engo.WindowWidth()
			}

			SpaceCompRef := &hs.SelectionRect.SpaceComponent

			// Since mx, my represent the opposite cornor to the Position
			// Widht and Height is the difference between (mx, my) and Position
			SpaceCompRef.Width = mx - SpaceCompRef.Position.X
			SpaceCompRef.Height = my - SpaceCompRef.Position.Y
		}
	}()

	//Updating HUD with Food and Wood values
	func() {
		if hs.PrevFoodVal != PlayerFood {

		}
	}()
}

func (hs *HUDSystem) SetBottomHUD(Name string, Index int) {
	LabelToSet := LabelGroupMap[Name]
	LastLabel := hs.CurrentActiveLabel
	LastIndex := hs.CurrentLabelIndex

	// fmt.Println("----------------------------")
	// fmt.Printf("Name: %v, Index: %v", Name, Index)

	// Remove the Labels currently being displayed
	if hs.CurrentActiveLabel != nil {

		// fmt.Printf(", LastLabel: %v, LastIndex: %v\n", LastLabel.Name, LastIndex)

		ActiveSystems.RenderSys.Remove(LastLabel.DescriptionLabel.BasicEntity)
		if len(LastLabel.ActionLabels) > 0 {
			for i, _ := range LastLabel.ActionLabels[LastIndex] {
				ActionLabel := &LastLabel.ActionLabels[LastIndex][i]

				// fmt.Printf("Action Label Remove: %v\n", ActionLabel.BasicEntity.ID())

				ActiveSystems.RenderSys.Remove(ActionLabel.BasicEntity)
			}
		}
	}

	// fmt.Println("\nRemoved Label")

	ActiveSystems.RenderSys.Add(
		&LabelToSet.DescriptionLabel.BasicEntity, &LabelToSet.DescriptionLabel.RenderComponent,
		&LabelToSet.DescriptionLabel.SpaceComponent,
	)
	if len(LabelToSet.ActionLabels) > 0 {
		for i, _ := range LabelToSet.ActionLabels[Index] {
			ActionLabel := &LabelToSet.ActionLabels[Index][i]

			// fmt.Printf("Action Label Add: %v\n", ActionLabel.BasicEntity.ID())

			ActiveSystems.RenderSys.Add(
				&ActionLabel.BasicEntity,
				&ActionLabel.RenderComponent,
				&ActionLabel.SpaceComponent,
			)
		}
	}

	hs.CurrentLabelIndex = Index
	hs.CurrentActiveLabel = &LabelToSet

	// fmt.Println("----------------------------")
}

func (*HUDSystem) Remove(ecs.BasicEntity) {}
