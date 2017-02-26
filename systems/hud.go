package systems

import (
	"engo.io/ecs"
	"engo.io/engo"
	"engo.io/engo/common"
	"image/color"
)

type HUDSystem struct {
	world *ecs.World
}

func (*HUDSystem) Update(dt float32) {}

func (*HUDSystem) Remove(ecs.BasicEntity) {}

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

func (*HUDSystem) New(w *ecs.World) {
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

	BottomHud.RenderComponent.SetZIndex(1000)
	BottomHud.RenderComponent.SetShader(common.HUDShader)

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

	TopHud.RenderComponent.SetZIndex(1000)
	TopHud.RenderComponent.SetShader(common.HUDShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&BottomHud.BasicEntity, &BottomHud.RenderComponent, &BottomHud.SpaceComponent)
			sys.Add(&TopHud.BasicEntity, &TopHud.RenderComponent, &TopHud.SpaceComponent)
		}
	}

	/*

	   The above part is for the plane HUD's and from here the indicators are placed on them


	*/

	Rect1 := SHAPE{BasicEntity: ecs.NewBasic()} //First Big Rectangle
	Rect1.SpaceComponent = common.SpaceComponent{Position: engo.Point{15, engo.WindowHeight() - float32(HudHeight-15)}, Width: float32((HudWidth / 3) - 80), Height: float32((HudHeight) - 30)}
	Rect1.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect1.RenderComponent.SetZIndex(1500)
	Rect1.RenderComponent.SetShader(common.HUDShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&Rect1.BasicEntity, &Rect1.RenderComponent, &Rect1.SpaceComponent)
		}
	}

	Rect2 := SHAPE{BasicEntity: ecs.NewBasic()} //R2 , R3, R4 are small Rectangles
	Rect2.SpaceComponent = common.SpaceComponent{Position: engo.Point{15 + Rect1.SpaceComponent.Width + 80, engo.WindowHeight() - float32(HudHeight-15) + float32(Rect1.SpaceComponent.Height/4)}, Width: float32(Rect1.SpaceComponent.Width / 3), Height: float32(Rect1.SpaceComponent.Height / 2)}
	Rect2.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect2.RenderComponent.SetZIndex(1500)
	Rect2.RenderComponent.SetShader(common.HUDShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&Rect2.BasicEntity, &Rect2.RenderComponent, &Rect2.SpaceComponent)
		}
	}

	Rect3 := SHAPE{BasicEntity: ecs.NewBasic()}
	Rect3.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect2.SpaceComponent.Position.X + Rect2.SpaceComponent.Width + 20, Rect2.SpaceComponent.Position.Y}, Width: Rect2.SpaceComponent.Width, Height: Rect2.SpaceComponent.Height}
	Rect3.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect3.RenderComponent.SetZIndex(1500)
	Rect3.RenderComponent.SetShader(common.HUDShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&Rect3.BasicEntity, &Rect3.RenderComponent, &Rect3.SpaceComponent)
		}
	}

	Rect4 := SHAPE{BasicEntity: ecs.NewBasic()}
	Rect4.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect3.SpaceComponent.Position.X + Rect2.SpaceComponent.Width + 20, Rect2.SpaceComponent.Position.Y}, Width: Rect2.SpaceComponent.Width, Height: Rect2.SpaceComponent.Height}
	Rect4.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect4.RenderComponent.SetZIndex(1500)
	Rect4.RenderComponent.SetShader(common.HUDShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&Rect4.BasicEntity, &Rect4.RenderComponent, &Rect4.SpaceComponent)
		}
	}

	/*


		Last two rectangles in the Bottom Hud

	*/

	Rect5 := SHAPE{BasicEntity: ecs.NewBasic()}
	Rect5.SpaceComponent = common.SpaceComponent{Position: engo.Point{engo.WindowWidth() - (Rect2.SpaceComponent.Width - 30) - 20, engo.WindowHeight() - float32(HudHeight) + 10}, Width: Rect2.SpaceComponent.Width - 30, Height: float32(HudHeight/2) - 15}
	Rect5.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect5.RenderComponent.SetZIndex(1500)
	Rect5.RenderComponent.SetShader(common.HUDShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&Rect5.BasicEntity, &Rect5.RenderComponent, &Rect5.SpaceComponent)
		}
	}

	Rect6 := SHAPE{BasicEntity: ecs.NewBasic()}
	Rect6.SpaceComponent = common.SpaceComponent{Position: engo.Point{Rect5.SpaceComponent.Position.X, engo.WindowHeight() - float32(HudHeight/2) + 5}, Width: Rect2.SpaceComponent.Width - 30, Height: float32(HudHeight/2) - 15}
	Rect6.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect6.RenderComponent.SetZIndex(1500)
	Rect6.RenderComponent.SetShader(common.HUDShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&Rect6.BasicEntity, &Rect6.RenderComponent, &Rect6.SpaceComponent)
		}
	}

	/*

		Top Hud Rectangles

	*/

	Rect7 := SHAPE{BasicEntity: ecs.NewBasic()} //First Big Rectangle
	Rect7.SpaceComponent = common.SpaceComponent{Position: engo.Point{128, 16}, Width: 128, Height: TopHud.Height - 32}
	Rect7.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect7.RenderComponent.SetZIndex(1500)
	Rect7.RenderComponent.SetShader(common.HUDShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&Rect7.BasicEntity, &Rect7.RenderComponent, &Rect7.SpaceComponent)
		}
	}

	Rect8 := SHAPE{BasicEntity: ecs.NewBasic()} //First Big Rectangle
	Rect8.SpaceComponent = common.SpaceComponent{Position: engo.Point{384, 16}, Width: 128, Height: TopHud.Height - 32}
	Rect8.RenderComponent = common.RenderComponent{Drawable: common.Rectangle{}, Color: color.RGBA{255, 255, 255, 255}}

	Rect8.RenderComponent.SetZIndex(1500)
	Rect8.RenderComponent.SetShader(common.HUDShader)

	for _, system := range w.Systems() {
		switch sys := system.(type) {
		case *common.RenderSystem:
			sys.Add(&Rect8.BasicEntity, &Rect8.RenderComponent, &Rect8.SpaceComponent)
		}
	}

}
