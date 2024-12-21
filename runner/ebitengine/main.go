package main

import (
	"image/color"
	"log"
	"time"

	"github.com/apfelfrisch/gosnake/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	gameSpeed     = time.Second / 10
	displayWidth  = 1000
	displayHeight = 1000
	gridSize      = 20
)

type Engine struct {
	game       game.Game
	lastUpdate time.Time
}

func (e *Engine) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		e.game.ChangeDirection(0, game.North)
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		e.game.ChangeDirection(0, game.South)
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		e.game.ChangeDirection(0, game.West)
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		e.game.ChangeDirection(0, game.East)
	} else if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if e.game.State() == game.Finished {
			e.game.Reset()
		}
	}

	if time.Since(e.lastUpdate) < gameSpeed {
		return nil
	}

	e.game.Tick()
	e.lastUpdate = time.Now()

	return nil
}

func (e *Engine) Draw(screen *ebiten.Image) {
	var x, y uint16
	for y = 1; y <= e.game.Height(); y++ {
		for x = 1; x <= e.game.Width(); x++ {
			drawField(screen, e.game.Field(game.Position{Y: uint16(y), X: uint16(x)}), x, y)
		}
	}
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return displayWidth, displayHeight
}

func drawField(screen *ebiten.Image, field game.Field, x uint16, y uint16) {
	var c color.Color
	switch field {
	case game.Wall:
		c = color.Gray{150}
	case game.Empty:
		c = color.Black
	default:
		c = color.White
	}

	vector.DrawFilledRect(
		screen,
		float32(x*gridSize-gridSize),
		float32(y*gridSize-gridSize),
		float32(gridSize),
		float32(gridSize),
		c,
		true,
	)
}

func main() {
	ebiten.SetWindowSize(displayWidth, displayHeight)
	ebiten.SetWindowTitle("Hello, World!")

	if err := ebiten.RunGame(&Engine{game: game.NewSingle(displayWidth/gridSize, displayHeight/gridSize)}); err != nil {
		log.Fatal(err)
	}
}
