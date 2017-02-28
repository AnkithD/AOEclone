package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"fmt"
	"image/color"
)

type HUDSystem struct {
	world       *ecs.World
	Bottomlabel [][]Details
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

type Details struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
}

func (rect *HUDSystem) New(w *ecs.World) {

	rect.world = w

	//Bottom Hud part

	BottomHud := HUD{BasicEntity: ecs.NewBasic()}

	HudWidth := int(engo.WindowWidth())
	HudHeight := 160

	colorA := color.RGBA{222, 184, 135, 250}

	BottomHud.RenderComponent = common.RenderComponent{
		Drawable: common.Rectangle{},
		Color:    colorA,
	}
	BottomHud.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, engo.WindowHeight() - float32(HudHeight)},
		Width:    float32(HudWidth),
		Height:   float32(HudHeight),
	}

	BottomHud.RenderComponent.SetZIndex(100)
	BottomHud.RenderComponent.SetShader(common.HUDShader)

	//Top Hud Part

	TopHud := HUD{BasicEntity: ecs.NewBasic()}

	TopWidth := int(engo.WindowWidth())
	TopHeight := 64

	TopHud.RenderComponent = common.RenderComponent{
		Drawable: common.Rectangle{},
		Color:    colorA,
	}
	TopHud.SpaceComponent = common.SpaceComponent{
		Position: engo.Point{0, 0},
		Width:    float32(TopWidth),
		Height:   float32(TopHeight),
	}

	TopHud.RenderComponent.SetZIndex(100)
	TopHud.RenderComponent.SetShader(common.HUDShader)

	ActiveSystems.RenderSys.Add(&BottomHud.BasicEntity, &BottomHud.RenderComponent, &BottomHud.SpaceComponent)
	ActiveSystems.RenderSys.Add(&TopHud.BasicEntity, &TopHud.RenderComponent, &TopHud.SpaceComponent)

	/*

	  Boxes on the HUD's to display texts

	*/

	Rect1 := SHAPE{BasicEntity: ecs.NewBasic()} //First Big Rectangle
	Rect1.SpaceComponent = common.SpaceComponent{Position: engo.Point{15, engo.WindowHeight() - float32(HudHeight-15)}, Width: float32((HudWidth / 3) - 80), Height: float32((HudHeight) - 30)}
	Rect1.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect1.RenderComponent.SetZIndex(125)
	Rect1.RenderComponent.SetShader(common.HUDShader)

	Rect2 := SHAPE{BasicEntity: ecs.NewBasic()} //R2 , R3, R4 are small Rectangles
	Rect2.SpaceComponent = common.SpaceComponent{Position: engo.Point{15 + Rect1.SpaceComponent.Width + 80, Rect1.SpaceComponent.Position.Y}, Width: float32(Rect1.SpaceComponent.Width / 3), Height: float32(Rect1.SpaceComponent.Height/2) - 5}
	Rect2.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect2.RenderComponent.SetZIndex(125)
	Rect2.RenderComponent.SetShader(common.HUDShader)

	Rect3 := SHAPE{BasicEntity: ecs.NewBasic()}
	Rect3.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect2.SpaceComponent.Position.X + Rect2.SpaceComponent.Width + 20, Rect2.SpaceComponent.Position.Y}, Width: Rect2.SpaceComponent.Width, Height: Rect2.SpaceComponent.Height}
	Rect3.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect3.RenderComponent.SetZIndex(125)
	Rect3.RenderComponent.SetShader(common.HUDShader)

	Rect4 := SHAPE{BasicEntity: ecs.NewBasic()}
	Rect4.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect3.SpaceComponent.Position.X + Rect2.SpaceComponent.Width + 20, Rect2.SpaceComponent.Position.Y}, Width: Rect2.SpaceComponent.Width, Height: Rect2.SpaceComponent.Height}
	Rect4.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect4.RenderComponent.SetZIndex(125)
	Rect4.RenderComponent.SetShader(common.HUDShader)

	ActiveSystems.RenderSys.Add(&Rect1.BasicEntity, &Rect1.RenderComponent, &Rect1.SpaceComponent)
	ActiveSystems.RenderSys.Add(&Rect2.BasicEntity, &Rect2.RenderComponent, &Rect2.SpaceComponent)
	ActiveSystems.RenderSys.Add(&Rect3.BasicEntity, &Rect3.RenderComponent, &Rect3.SpaceComponent)
	ActiveSystems.RenderSys.Add(&Rect4.BasicEntity, &Rect4.RenderComponent, &Rect4.SpaceComponent)

	Rect5 := SHAPE{BasicEntity: ecs.NewBasic()}
	Rect5.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect2.SpaceComponent.Position.X, Rect2.SpaceComponent.Position.Y + Rect2.Height + 10}, Width: Rect2.SpaceComponent.Width, Height: Rect2.SpaceComponent.Height}
	Rect5.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect5.RenderComponent.SetZIndex(125)
	Rect5.RenderComponent.SetShader(common.HUDShader)

	Rect6 := SHAPE{BasicEntity: ecs.NewBasic()}
	Rect6.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect5.SpaceComponent.Position.X + Rect5.SpaceComponent.Width + 20, Rect5.SpaceComponent.Position.Y}, Width: Rect5.SpaceComponent.Width, Height: Rect5.SpaceComponent.Height}
	Rect6.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect6.RenderComponent.SetZIndex(125)
	Rect6.RenderComponent.SetShader(common.HUDShader)

	Rect7 := SHAPE{BasicEntity: ecs.NewBasic()}
	Rect7.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect6.SpaceComponent.Position.X + Rect6.SpaceComponent.Width + 20, Rect6.SpaceComponent.Position.Y}, Width: Rect6.SpaceComponent.Width, Height: Rect6.SpaceComponent.Height}
	Rect7.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect7.RenderComponent.SetZIndex(125)
	Rect7.RenderComponent.SetShader(common.HUDShader)

	ActiveSystems.RenderSys.Add(&Rect5.BasicEntity, &Rect5.RenderComponent, &Rect5.SpaceComponent)
	ActiveSystems.RenderSys.Add(&Rect6.BasicEntity, &Rect6.RenderComponent, &Rect6.SpaceComponent)
	ActiveSystems.RenderSys.Add(&Rect7.BasicEntity, &Rect7.RenderComponent, &Rect7.SpaceComponent)

	/*


		Last two rectangles in the Bottom Hud

	*/

	Rect8 := SHAPE{BasicEntity: ecs.NewBasic()}
	wid := Rect2.SpaceComponent.Width - 30
	hig := float32(HudHeight/2) - 15
	Rect8.SpaceComponent = common.SpaceComponent{Position: engo.Point{engo.WindowWidth() - wid - 20, engo.WindowHeight() - float32(HudHeight) + 10}, Width: wid, Height: hig}
	Rect8.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect8.RenderComponent.SetZIndex(125)
	Rect8.RenderComponent.SetShader(common.HUDShader)

	Rect9 := SHAPE{BasicEntity: ecs.NewBasic()}
	Rect9.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect8.SpaceComponent.Position.X, engo.WindowHeight() - float32(HudHeight/2) + 5}, Width: wid, Height: hig}
	Rect9.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect9.RenderComponent.SetZIndex(125)
	Rect9.RenderComponent.SetShader(common.HUDShader)

	ActiveSystems.RenderSys.Add(&Rect8.BasicEntity, &Rect8.RenderComponent, &Rect8.SpaceComponent)
	ActiveSystems.RenderSys.Add(&Rect9.BasicEntity, &Rect9.RenderComponent, &Rect9.SpaceComponent)

	/*

		Top Hud Rectangles

	*/

	Rect10 := SHAPE{BasicEntity: ecs.NewBasic()}
	Rect10.SpaceComponent = common.SpaceComponent{Position: engo.Point{96, 16}, Width: 128, Height: TopHud.Height - 32}
	Rect10.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect10.RenderComponent.SetZIndex(125)
	Rect10.RenderComponent.SetShader(common.HUDShader)

	Rect11 := SHAPE{BasicEntity: ecs.NewBasic()}
	Rect11.SpaceComponent = common.SpaceComponent{Position: engo.Point{352, 16}, Width: 128, Height: TopHud.Height - 32}
	Rect11.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect11.RenderComponent.SetZIndex(125)
	Rect11.RenderComponent.SetShader(common.HUDShader)

	ActiveSystems.RenderSys.Add(&Rect10.BasicEntity, &Rect10.RenderComponent, &Rect10.SpaceComponent)
	ActiveSystems.RenderSys.Add(&Rect11.BasicEntity, &Rect11.RenderComponent, &Rect11.SpaceComponent)

	/*

		For the Text on The HUD's

	*/

	fnt := &common.Font{
		URL:  "Roboto-Regular.ttf",
		FG:   color.Black,
		Size: 16,
	}

	err := fnt.CreatePreloaded()
	if err != nil {
		panic(err)
	}

	/*
	   On the Top HUD----Food and Wood

	*/

	label1 := Details{BasicEntity: ecs.NewBasic()}
	label1.SpaceComponent = common.SpaceComponent{Position: engo.Point{32, 24}}
	label1.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: "FOOD :",
	}
	label1.SetShader(common.TextHUDShader)
	label1.SetZIndex(150)

	label2 := Details{BasicEntity: ecs.NewBasic()}
	label2.SpaceComponent = common.SpaceComponent{Position: engo.Point{288, 24}}
	label2.RenderComponent.Drawable = common.Text{
		Font: fnt,
		Text: "WOOD :",
	}
	label2.SetShader(common.TextHUDShader)
	label2.SetZIndex(150)

	ActiveSystems.RenderSys.Add(&label1.BasicEntity, &label1.RenderComponent, &label1.SpaceComponent)
	ActiveSystems.RenderSys.Add(&label2.BasicEntity, &label2.RenderComponent, &label2.SpaceComponent)

	/*

	   TEXT On the Bottom HUD----
	   for loop to be able to collapse code

	*/
	lab1 := Details{BasicEntity: ecs.NewBasic()}
	lab1.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect1.SpaceComponent.Position.X + 48, Rect1.SpaceComponent.Position.Y + 32}}
	lab1.RenderComponent.Drawable = common.Text{Font: fnt, Text: "TOWN CENTRE\n\n\nHealth : XX/YY"}
	lab1.SetShader(common.TextHUDShader)
	lab1.SetZIndex(250)

	lab2 := Details{BasicEntity: ecs.NewBasic()}
	lab2.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect1.SpaceComponent.Position.X + 48, Rect1.SpaceComponent.Position.Y + 32}}
	lab2.RenderComponent.Drawable = common.Text{Font: fnt, Text: "VILLAGER\n\n\nHealth : XX/YY"}
	lab2.SetShader(common.TextHUDShader)
	lab2.SetZIndex(250)

	lab3 := Details{BasicEntity: ecs.NewBasic()}
	lab3.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect1.SpaceComponent.Position.X + 48, Rect1.SpaceComponent.Position.Y + 32}}
	lab3.RenderComponent.Drawable = common.Text{Font: fnt, Text: "HOUSE\n\n\nHealth : XX/YY\n\n\nCapacity : xx/yy"}
	lab3.SetShader(common.TextHUDShader)
	lab3.SetZIndex(250)

	lab4 := Details{BasicEntity: ecs.NewBasic()}
	lab4.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect1.SpaceComponent.Position.X + 48, Rect1.SpaceComponent.Position.Y + 32}}
	lab4.RenderComponent.Drawable = common.Text{Font: fnt, Text: "MILITARY\n\n\nHealth : XX/YY"}
	lab4.SetShader(common.TextHUDShader)
	lab4.SetZIndex(250)

	lab5 := Details{BasicEntity: ecs.NewBasic()}
	lab5.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect1.SpaceComponent.Position.X + 48, Rect1.SpaceComponent.Position.Y + 32}}
	lab5.RenderComponent.Drawable = common.Text{Font: fnt, Text: "RESOURCE\n\n\nHealth : XX/YY"}
	lab5.SetShader(common.TextHUDShader)
	lab5.SetZIndex(250)

	lab6 := Details{BasicEntity: ecs.NewBasic()}
	lab6.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect2.SpaceComponent.Position.X + 16, Rect2.SpaceComponent.Position.Y + 16}}
	lab6.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Create Villager"}
	lab6.SetShader(common.TextHUDShader)
	lab6.SetZIndex(250)

	lab7 := Details{BasicEntity: ecs.NewBasic()}
	lab7.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect2.SpaceComponent.Position.X + 16, Rect2.SpaceComponent.Position.Y + 16}}
	lab7.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Build"}
	lab7.SetShader(common.TextHUDShader)
	lab7.SetZIndex(250)

	lab8 := Details{BasicEntity: ecs.NewBasic()}
	lab8.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect2.SpaceComponent.Position.X + 16, Rect2.SpaceComponent.Position.Y + 16}}
	lab8.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Create Soldier"}
	lab8.SetShader(common.TextHUDShader)
	lab8.SetZIndex(250)

	lab9 := Details{BasicEntity: ecs.NewBasic()}
	lab9.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect3.SpaceComponent.Position.X + 16, Rect3.SpaceComponent.Position.Y + 16}}
	lab9.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Repair"}
	lab9.SetShader(common.TextHUDShader)
	lab9.SetZIndex(250)

	//If clicked on Build Then the following options are displayed

	lab11 := Details{BasicEntity: ecs.NewBasic()}
	lab11.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect2.SpaceComponent.Position.X + 16, Rect2.SpaceComponent.Position.Y + 16}}
	lab11.RenderComponent.Drawable = common.Text{Font: fnt, Text: "House"}
	lab11.SetShader(common.TextHUDShader)
	lab11.SetZIndex(250)

	lab12 := Details{BasicEntity: ecs.NewBasic()}
	lab12.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect3.SpaceComponent.Position.X + 16, Rect3.SpaceComponent.Position.Y + 16}}
	lab12.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Military Camp"}
	lab12.SetShader(common.TextHUDShader)
	lab12.SetZIndex(250)

	lab13 := Details{BasicEntity: ecs.NewBasic()}
	lab13.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect4.SpaceComponent.Position.X + 16, Rect4.SpaceComponent.Position.Y + 16}}
	lab13.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Resource"}
	lab13.SetShader(common.TextHUDShader)
	lab13.SetZIndex(250)

	lab14 := Details{BasicEntity: ecs.NewBasic()}
	lab14.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect5.SpaceComponent.Position.X + 16, Rect5.SpaceComponent.Position.Y + 16}}
	lab14.RenderComponent.Drawable = common.Text{Font: fnt, Text: "Go Back"}
	lab14.SetShader(common.TextHUDShader)
	lab14.SetZIndex(250)

	/*

		Appending to the slices
		for loop to be able to collapse code

	*/
	rect.Bottomlabel = make([][]Details, 0)

	tempslice := make([]Details, 0)
	tempslice = append(tempslice, lab1, lab6)
	rect.Bottomlabel = append(rect.Bottomlabel, tempslice)

	tempslice = make([]Details, 0)
	tempslice = append(tempslice, lab2, lab7, lab9)
	rect.Bottomlabel = append(rect.Bottomlabel, tempslice)

	tempslice = make([]Details, 0)
	tempslice = append(tempslice, lab3)
	rect.Bottomlabel = append(rect.Bottomlabel, tempslice)

	tempslice = make([]Details, 0)
	tempslice = append(tempslice, lab4, lab8)
	rect.Bottomlabel = append(rect.Bottomlabel, tempslice)

	tempslice = make([]Details, 0)
	tempslice = append(tempslice, lab5)
	rect.Bottomlabel = append(rect.Bottomlabel, tempslice)

	lab := rect.Bottomlabel[1]
	fmt.Println(len(lab))
	for i, _ := range lab {
		ActiveSystems.RenderSys.Add(&lab[i].BasicEntity, &lab[i].RenderComponent, &lab[i].SpaceComponent)
	}

	fmt.Println("HUD System Initialized")

}

func (rect *HUDSystem) Update(dt float32) {}

func (*HUDSystem) Remove(ecs.BasicEntity) {}
