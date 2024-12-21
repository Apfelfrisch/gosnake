package main

import (
	"image/color"
	"log"
	"time"

	"github.com/apfelfrisch/gosnake/game"
	gclient "github.com/apfelfrisch/gosnake/game/client"
	gserver "github.com/apfelfrisch/gosnake/game/server"
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
	client gclient.Tcp
}

func (e *Engine) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		e.client.Write('w')
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		e.client.Write('s')
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		e.client.Write('a')
	} else if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		e.client.Write('d')
	} else if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		e.client.Write('↵')
	}

	return nil
}

func (e *Engine) Draw(screen *ebiten.Image) {
	for _, fieldPos := range gclient.DeserializeState(e.client.Read()) {
		drawField(screen, fieldPos.Field, fieldPos.X, fieldPos.Y)
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
	ebiten.SetWindowTitle("Snake")

	server := gserver.New(
		1,
		":1200",
		game.NewSingle(displayWidth/gridSize, displayHeight/gridSize),
	)
	server.RunBackground()

	time.Sleep(time.Second)
	engine := Engine{
		client: *gclient.NewTcpClient(":1200"),
	}
	engine.client.Connect()

	for {
		if server.Ready() {
			break
		}
		time.Sleep(time.Second / 10)
	}

	if err := ebiten.RunGame(&engine); err != nil {
		log.Fatal(err)
	}
}
