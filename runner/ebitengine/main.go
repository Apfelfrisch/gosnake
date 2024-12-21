package main

import (
	"image/color"
	"log"
	"os"
	"time"

	"github.com/apfelfrisch/gosnake/game"
	gclient "github.com/apfelfrisch/gosnake/game/client"
	gserver "github.com/apfelfrisch/gosnake/game/server"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	displayWidth  = 1000
	displayHeight = 1000
	gridSize      = 20
)

type Engine struct {
	client *gclient.Tcp
}

func (e *Engine) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		e.client.Write('w')
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		e.client.Write('s')
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		e.client.Write('a')
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		e.client.Write('d')
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
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
		false,
	)
}

func main() {
	ebiten.SetWindowSize(displayWidth, displayHeight)
	ebiten.SetWindowTitle("Snake")

	playerCount := 2

	var client *gclient.Tcp
	if len(os.Args) <= 1 {
		startServer(playerCount)
		client = connectClient(":1200")
	} else {
		client = connectClient(os.Args[1])
	}

	if err := ebiten.RunGame(&Engine{client}); err != nil {
		log.Fatal(err)
	}
}

func startServer(playerCount int) *gserver.GameServer {
	server := gserver.New(
		playerCount,
		":1200",
		game.NewBattleSnake(playerCount, displayWidth/gridSize, displayHeight/gridSize),
	)
	server.RunBackground()

	return server
}

func connectClient(addr string) *gclient.Tcp {
	client := gclient.NewTcpClient(addr)

	for i := 0; i < 10; i++ {
		if err := client.Connect(); err == nil {
			break
		}
		time.Sleep(time.Second / 5)
	}

	for {
		if client.Read() != "" {
			break
		}
		time.Sleep(time.Second / 10)
	}

	return client
}
