package scenes

import (
	"image/color"
	"os"

	"github.com/apfelfrisch/gosnake/engine"
	"github.com/apfelfrisch/gosnake/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type GameRunning struct {
	BaseScene
}

func (s *GameRunning) Update() error {
	s.client.UpdatePayload()

	if s.client.Payload.GameState == game.Paused || s.client.Payload.GameState == game.RoundFinished {
		s.sm.SwitchTo(&MenuPaused{BaseScene: s.BaseScene})
		return nil
	}
	if s.client.Payload.GameState == game.GameFinished {
		s.sm.SwitchTo(&MenuFinished{BaseScene: s.BaseScene})
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s.client.PressKey('↵')
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		s.client.PressKey('w')
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		s.client.PressKey('s')
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		s.client.PressKey('a')
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		s.client.PressKey('d')
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.client.PressKey(' ')
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s.client.PressKey('↵')
	}

	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyC) {
		os.Exit(0)
	}

	return nil
}

func (s *GameRunning) Draw(screen *ebiten.Image) {
	drawCandies(screen, s.client.Payload.Candies)
	drawSnakes(screen, &s.BaseScene)
	drawGameField(screen, s.client.World())
	drawPlayerInfo(screen, s.client.Payload)
}

func drawGameField(screen *ebiten.Image, world []game.FieldPos) {
	drawRect := func(fieldPos game.FieldPos, c color.Color) {
		vector.DrawFilledRect(
			screen,
			float32(fieldPos.X*engine.GridSize-engine.GridSize),
			float32(fieldPos.Y*engine.GridSize-engine.GridSize),
			float32(engine.GridSize),
			float32(engine.GridSize),
			c,
			false,
		)
	}

	for _, fieldPos := range world {
		switch fieldPos.Field {
		case game.Wall:
			drawRect(fieldPos, color.Gray{150})
		case game.Empty:
			if (fieldPos.X+fieldPos.Y)%2 == 0 {
				drawRect(fieldPos, color.RGBA{13, 13, 13, 0})
			} else {
				drawRect(fieldPos, color.RGBA{0, 0, 0, 0})
			}
		}
	}
}

var snakecolors = [4]color.Color{
	color.RGBA{204, 0, 0, 255},
	color.RGBA{204, 102, 0, 255},
	color.RGBA{102, 204, 0, 255},
	color.RGBA{204, 0, 204, 255},
}

func drawSnakes(screen *ebiten.Image, base *BaseScene) {
	intermidiatPixel := 3

	player := base.client.Payload.Player
	if base.localPlayer.OutOfsync(player) {
		base.localPlayer.Sync(player)
	}

	for _, body := range base.localPlayer.Positions(player.Direction, intermidiatPixel) {
		vector.DrawFilledRect(
			screen,
			body.X,
			body.Y,
			body.Width,
			body.Height,
			color.RGBA{30, 144, 255, 255},
			false,
		)
	}

	var c color.Color
	for i, opp := range base.client.Payload.Opponents {
		if len(base.localOpponents) <= i {
			base.localOpponents = append(base.localOpponents, engine.ClientSnake{
				GridSize:    engine.GridSize,
				ServerSnake: opp,
				InterPixel:  0,
			})
		} else if base.localOpponents[i].OutOfsync(opp) {
			base.localOpponents[i].Sync(opp)
		}

		for _, body := range base.localOpponents[i].Positions(opp.Direction, intermidiatPixel) {
			if i >= 0 && i < len(snakecolors) {
				c = snakecolors[i]
			} else {
				c = color.RGBA{30, 144, 255, 255}
			}

			vector.DrawFilledRect(
				screen,
				body.X,
				body.Y,
				body.Width,
				body.Height,
				c,
				false,
			)
		}
	}
}

func drawCandies(screen *ebiten.Image, candies []game.Position) {
	drawCircle := func(pos game.Position, c color.Color) {
		vector.DrawFilledCircle(
			screen,
			float32(pos.X*engine.GridSize-engine.GridSize/2),
			float32(pos.Y*engine.GridSize-engine.GridSize/2),
			float32(engine.GridSize)/2,
			c,
			true,
		)
	}

	for _, pos := range candies {
		drawCircle(pos, color.RGBA{255, 215, 0, 255})
	}
}
