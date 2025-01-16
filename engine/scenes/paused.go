package scenes

import (
	"bytes"
	"image/color"
	"log"

	"github.com/apfelfrisch/gosnake/engine"
	"github.com/apfelfrisch/gosnake/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type MenuPaused struct {
	BaseScene
}

func (s *MenuPaused) Update() error {
	s.client.UpdatePayload()
	s.localPlayer.Sync(s.client.Payload.Player)

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s.client.PressKey('↵')
	}

	if s.client.Payload.GameState == game.Ongoing {
		s.sm.SwitchTo(&GameRunning{BaseScene: s.BaseScene})
	}

	return nil
}

func (s *MenuPaused) Draw(screen *ebiten.Image) {
	drawPausedScreen(screen)
	drawPlayerInfo(screen, s.client.Payload)
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

	op.GeoM.Translate(engine.GameWidth/2-300, engine.DisplayHeight/2-50)
	text.Draw(screen, message, face, op)
}
