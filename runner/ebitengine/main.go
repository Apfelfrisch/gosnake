package main

import (
	"flag"
	"log"
	"os"
	"time"

	// _ "net/http/pprof"

	"github.com/apfelfrisch/gosnake/game"
	gclient "github.com/apfelfrisch/gosnake/game/client"
	gserver "github.com/apfelfrisch/gosnake/game/server"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	displayWidth  = 1500
	displayHeight = 1000
	gameWidth     = 1000
	gameHeight    = 1000
	gridSize      = 20
)

type Engine struct {
	client *gclient.GameClient
}

func (e *Engine) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		e.client.PressKey('w')
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		e.client.PressKey('s')
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		e.client.PressKey('a')
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		e.client.PressKey('d')
	} else if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		e.client.PressKey(' ')
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		e.client.PressKey('↵')
	}

	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyC) {
		os.Exit(0)
	}

	return nil
}

func (e *Engine) Draw(screen *ebiten.Image) {
	e.client.UpdatePayload()

	if e.client.Payload.GameState != game.GameFinished {
		drawGameField(screen, e.client.World())
		drawPlayerInfo(screen, e.client.Payload)
	} else {
		drawFinishScreen(screen, e.client.Payload.Player)
	}
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return displayWidth, displayHeight
}

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	playerCount := flag.Int("player", 1, "Set Player count")
	serverAddr := flag.String("server-addr", ":1200", "Set Sever Address")
	onlyServer := flag.Bool("only-server", false, "Run only the server")
	onlyClient := flag.Bool("only-client", false, "Run only the server")

	flag.Parse()

	ebiten.SetWindowSize(displayWidth, displayHeight)
	ebiten.SetWindowTitle("Snake")

	if *onlyServer == true {
		buildServer(*playerCount, *serverAddr).Run()
	} else if *onlyClient == false {
		buildServer(*playerCount, *serverAddr).RunBackground()
	}

	client := gclient.Connect(*serverAddr)

	if err := ebiten.RunGame(&Engine{client}); err != nil {
		log.Fatal(err)
	}
}

func buildServer(playerCount int, addr string) *gserver.GameServer {
	return gserver.New(
		playerCount,
		addr,
		game.NewBattleSnake(playerCount, gameWidth/gridSize, gameHeight/gridSize),
	)
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
