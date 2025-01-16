package engine

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"sort"

	"github.com/apfelfrisch/gosnake/game"
	"github.com/apfelfrisch/gosnake/game/network/payload"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	DisplayWidth  = 1500
	DisplayHeight = 1000
	GameWidth     = 1000
	GameHeight    = 1000
	GridSize      = 20
)

const playerInfoXOffset = GameWidth + 10

func drawPlayerInfo(screen *ebiten.Image, payload *payload.Payload) {
	// Background for the stats panel
	statsBgColor := color.RGBA{50, 50, 50, 255}
	statsPanelWidth := DisplayWidth - GameWidth
	statsPanelHeight := DisplayHeight

	vector.DrawFilledRect(
		screen,
		float32(GameWidth),
		0,
		float32(statsPanelWidth),
		float32(statsPanelHeight),
		statsBgColor,
		false,
	)

	menuFont, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	face := &text.GoTextFace{
		Source: menuFont,
		Size:   20.0,
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)

	sortPerkNames := func(perks game.Perks) []game.PerkType {
		pNames := make([]game.PerkType, 0, len(payload.Player.Perks))
		for key := range payload.Player.Perks {
			pNames = append(pNames, key)
		}
		sort.Slice(pNames, func(i, j int) bool {
			return pNames[i] < pNames[j]
		})

		return pNames
	}

	op.GeoM.Translate(playerInfoXOffset, 50)
	text.Draw(screen, "Lives:", face, op)
	op.GeoM.Translate(70, 0)
	text.Draw(screen, fmt.Sprintf("%d", payload.Player.Lives), face, op)
	op.GeoM.Translate(-70, 30)
	text.Draw(screen, "Perks:", face, op)
	op.GeoM.Translate(70, 0)
	for _, pName := range sortPerkNames(payload.Player.Perks) {
		text.Draw(screen, fmt.Sprintf("%v (%v)", pName, payload.Player.Perks[pName].Usages), face, op)
		op.GeoM.Translate(0, 30)
	}

	for _, oppenent := range payload.Opponents {
		op.GeoM.Translate(-70, 0)
		text.Draw(screen, "---", face, op)
		op.GeoM.Translate(0, 30)
		text.Draw(screen, "Lives:", face, op)
		op.GeoM.Translate(70, 0)
		text.Draw(screen, fmt.Sprintf("%d", oppenent.Lives), face, op)
		op.GeoM.Translate(-70, 30)
		text.Draw(screen, "Perks:", face, op)
		op.GeoM.Translate(70, 0)
		for _, pName := range sortPerkNames(oppenent.Perks) {
			text.Draw(screen, fmt.Sprintf("%v (%v)", pName, oppenent.Perks[pName].Usages), face, op)
			op.GeoM.Translate(0, 30)
		}
	}
}

func drawGameField(screen *ebiten.Image, world []game.FieldPos) {
	drawRect := func(fieldPos game.FieldPos, c color.Color) {
		vector.DrawFilledRect(
			screen,
			float32(fieldPos.X*GridSize-GridSize),
			float32(fieldPos.Y*GridSize-GridSize),
			float32(GridSize),
			float32(GridSize),
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

var snakecolors = []color.Color{
	color.RGBA{204, 0, 0, 255},
	color.RGBA{204, 102, 0, 255},
	color.RGBA{102, 204, 0, 255},
	color.RGBA{204, 0, 204, 255},
}

func drawSnakes(screen *ebiten.Image, engine *Engine) {
	intermidiatPixel := 3

	player := engine.client.Payload.Player
	if engine.localPlayer.outOfsync(player) {
		engine.localPlayer.sync(player)
	}

	for _, body := range engine.localPlayer.positions(player.Direction, intermidiatPixel) {
		vector.DrawFilledRect(
			screen,
			body.x,
			body.y,
			body.width,
			body.height,
			color.RGBA{30, 144, 255, 255},
			false,
		)
	}

	var c color.Color
	for i, opp := range engine.client.Payload.Opponents {
		if len(engine.localOpponents) <= i {
			engine.localOpponents = append(engine.localOpponents, ClientSnake{
				gridSize:    GridSize,
				serverSnake: opp,
				interPixel:  0,
			})
		} else if engine.localOpponents[i].outOfsync(opp) {
			engine.localOpponents[i].sync(opp)
		}

		for _, body := range engine.localOpponents[i].positions(opp.Direction, intermidiatPixel) {
			if i >= 0 && i < len(snakecolors) {
				c = snakecolors[i]
			} else {
				c = color.RGBA{30, 144, 255, 255}
			}

			vector.DrawFilledRect(
				screen,
				body.x,
				body.y,
				body.width,
				body.height,
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
			float32(pos.X*GridSize-GridSize/2),
			float32(pos.Y*GridSize-GridSize/2),
			float32(GridSize)/2,
			c,
			true,
		)
	}

	for _, pos := range candies {
		drawCircle(pos, color.RGBA{255, 215, 0, 255})
	}
}

func drawPausedScreen(screen *ebiten.Image) {
	message := "Ready, Press 'Enter' to start"

	menuFont, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	face := &text.GoTextFace{
		Source: menuFont,
		Size:   50.0,
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)

	op.GeoM.Translate(GameWidth/2-300, DisplayHeight/2-50)
	text.Draw(screen, message, face, op)
}

func drawFinishScreen(screen *ebiten.Image, player game.Snake) {
	message := ""
	if player.Lives == 0 {
		message = "You Lost :("
	} else {
		message = "You Won :)"
	}

	menuFont, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}
	face := &text.GoTextFace{
		Source: menuFont,
		Size:   50.0,
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)

	op.GeoM.Translate(GameWidth/2-150, DisplayHeight/2-50)
	text.Draw(screen, message, face, op)
}
