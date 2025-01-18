package scenes

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"log"
	"sort"

	"github.com/apfelfrisch/gosnake/engine"
	"github.com/apfelfrisch/gosnake/game"
	netClient "github.com/apfelfrisch/gosnake/game/network/client"
	"github.com/apfelfrisch/gosnake/game/network/payload"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/joelschutz/stagehand"
	"golang.org/x/image/font/gofont/goregular"
)

type BaseScene struct {
	bounds         image.Rectangle
	localPlayer    engine.ClientSnake
	localOpponents []engine.ClientSnake
	client         *netClient.GameClient
	sm             *stagehand.SceneManager[game.GameState]
}

func (e *BaseScene) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return engine.DisplayWidth, engine.DisplayHeight
}

func (s *BaseScene) Load(st game.GameState, sm stagehand.SceneController[game.GameState]) {
	s.sm = sm.(*stagehand.SceneManager[game.GameState])
}

func (s *BaseScene) Unload() game.GameState {
	return s.client.Payload.GameState
}

const playerInfoXOffset = engine.GameWidth + 10

func drawPlayerInfo(screen *ebiten.Image, payload *payload.Payload) {
	// Background for the stats panel
	statsBgColor := color.RGBA{50, 50, 50, 255}
	statsPanelWidth := engine.DisplayWidth - engine.GameWidth
	statsPanelHeight := engine.DisplayHeight

	vector.DrawFilledRect(
		screen,
		float32(engine.GameWidth),
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
