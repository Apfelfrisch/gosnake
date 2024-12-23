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

	op.GeoM.Translate(playerInfoXOffset, 50)
	text.Draw(screen, "Lives:", face, op)
	op.GeoM.Translate(0, 30)
	text.Draw(screen, "Perks:", face, op)

	op.GeoM.Reset()
	op.GeoM.Translate(playerInfoXOffset+70, 50)
	text.Draw(screen, fmt.Sprintf("%d", payload.Player.Lives), face, op)

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

	for _, pName := range sortPerkNames(payload.Player.Perks) {
		op.GeoM.Translate(0, 30)
		text.Draw(screen, fmt.Sprintf("%v (%v)", pName, payload.Player.Perks[pName].Usages), face, op)
	}

	for _, oppenent := range payload.Opponents {
		op.GeoM.Reset()
		op.GeoM.Translate(playerInfoXOffset, 140)
		text.Draw(screen, "---", face, op)
		op.GeoM.Translate(0, 30)
		text.Draw(screen, "Lives:", face, op)
		op.GeoM.Translate(0, 30)
		text.Draw(screen, "Perks:", face, op)

		op.GeoM.Reset()
		op.GeoM.Translate(playerInfoXOffset+70, 140)
		text.Draw(screen, fmt.Sprintf("%d", oppenent.Lives), face, op)
		for _, pName := range sortPerkNames(oppenent.Perks) {
			op.GeoM.Translate(0, 30)
			text.Draw(screen, fmt.Sprintf("%v (%v)", pName, oppenent.Perks[pName].Usages), face, op)
		}
	}
}

func drawGameField(screen *ebiten.Image, world []game.FieldPos) {
	var c color.Color

	for _, fieldPos := range world {
		switch fieldPos.Field {
		case game.Wall:
			c = color.Gray{150}
		case game.Empty:
			c = color.Black
		default:
			c = color.White
		}

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

	op.GeoM.Translate(displayWidth/2-100, displayHeight/2-50)
	text.Draw(screen, message, face, op)
}
