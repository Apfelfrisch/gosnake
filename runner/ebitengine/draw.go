package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"sort"

	"github.com/apfelfrisch/gosnake/game"
	"github.com/apfelfrisch/gosnake/game/client"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	displayWidth  = 1500
	displayHeight = 1000
	gameWidth     = 1000
	gameHeight    = 1000
	gridSize      = 20
)

const playerInfoXOffset = gameWidth + 10

func drawPlayerInfo(screen *ebiten.Image, payload *client.Payload) {
	// Background for the stats panel
	statsBgColor := color.RGBA{50, 50, 50, 255}
	statsPanelWidth := displayWidth - gameWidth
	statsPanelHeight := displayHeight

	vector.DrawFilledRect(
		screen,
		float32(gameWidth),
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
			float32(fieldPos.X*gridSize-gridSize),
			float32(fieldPos.Y*gridSize-gridSize),
			float32(gridSize),
			float32(gridSize),
			c,
			false,
		)
	}

	drawCircle := func(fieldPos game.FieldPos, c color.Color) {
		vector.DrawFilledCircle(
			screen,
			float32(fieldPos.X*gridSize-gridSize/2),
			float32(fieldPos.Y*gridSize-gridSize/2),
			float32(gridSize)/2,
			c,
			true,
		)
	}

	for _, fieldPos := range world {
		switch fieldPos.Field {
		case game.Wall:
			drawRect(fieldPos, color.Gray{150})
		case game.Empty:
			drawRect(fieldPos, color.RGBA{0, 0, 0, 0})
		case game.Candy:
			drawCircle(fieldPos, color.RGBA{255, 215, 0, 255})
		}
	}
}

var snakecolors = []color.Color{
	color.RGBA{204, 0, 0, 255},
	color.RGBA{204, 102, 0, 255},
	color.RGBA{102, 204, 0, 255},
	color.RGBA{204, 0, 204, 255},
}

func drawSnakes(screen *ebiten.Image, player game.Snake, oppenents []game.Snake) {
	drawRect := func(pos game.Position, c color.Color) {
		vector.DrawFilledRect(
			screen,
			float32(pos.X*gridSize-gridSize),
			float32(pos.Y*gridSize-gridSize),
			float32(gridSize),
			float32(gridSize),
			c,
			false,
		)
	}

	for i, opp := range oppenents {
		for _, pos := range opp.Occupied {
			if i >= 0 && i < len(snakecolors) {
				drawRect(pos, snakecolors[i])
			} else {
				drawRect(pos, color.RGBA{30, 144, 255, 255})
			}
		}
	}

	for _, pos := range player.Occupied {
		drawRect(pos, color.RGBA{30, 144, 255, 255})
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

	op.GeoM.Translate(gameWidth/2-300, displayHeight/2-50)
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

	op.GeoM.Translate(gameWidth/2-150, displayHeight/2-50)
	text.Draw(screen, message, face, op)
}
