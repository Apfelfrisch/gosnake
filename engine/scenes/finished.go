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

type MenuFinished struct {
	BaseScene
}

func (s *MenuFinished) Update() error {
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

func (s *MenuFinished) Draw(screen *ebiten.Image) {
	drawFinishScreen(screen, s.client.Payload.Player)
	drawPlayerInfo(screen, s.client.Payload)
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

	op.GeoM.Translate(engine.GameWidth/2-150, engine.DisplayHeight/2-50)
	text.Draw(screen, message, face, op)
}
