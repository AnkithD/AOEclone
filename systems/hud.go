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

func (*HUDSystem) New(w *ecs.World) {
	BottomHud := HUD{BasicEntity: ecs.NewBasic()}

	HudWidth := int(engo.WindowWidth())
	HudHeight := 150

	colorA := color.RGBA{160, 82, 45, 250}

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
	TopHeight := 50

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

}
