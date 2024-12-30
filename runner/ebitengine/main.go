package main

import (
	"flag"
	"log"
	"os"

	// _ "net/http/pprof"

	"github.com/apfelfrisch/gosnake/game"
	gclient "github.com/apfelfrisch/gosnake/game/network/client"
	gserver "github.com/apfelfrisch/gosnake/game/network/server"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Engine struct {
	client *gclient.GameClient
}

func NewEninge(serverAddr string) *Engine {
	player := NewPlayer()
	client := gclient.Connect(serverAddr)
	client.EventBus.Add(gclient.PlayerHasEaten{}, func(event gclient.Event) {
		player.Play(Eat)
	})
	client.EventBus.Add(gclient.PlayerDashed{}, func(event gclient.Event) {
		player.Play(Dash)
	})
	client.EventBus.Add(gclient.PlayerWalkedWall{}, func(event gclient.Event) {
		player.Play(WalkWall)
	})
	client.EventBus.Add(gclient.PlayerCrashed{}, func(event gclient.Event) {
		player.Play(Crash)
	})
	client.EventBus.Add(gclient.GameHasStarted{}, func(event gclient.Event) {
		player.PlayMusic()
	})
	client.EventBus.Add(gclient.GameHasEnded{}, func(event gclient.Event) {
		player.PauseMusic()
	})

	return &Engine{
		client: client,
	}
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

	if e.client.Payload.GameState == game.Paused {
		drawPausedScreen(screen)
	} else if e.client.Payload.GameState == game.GameFinished {
		drawFinishScreen(screen, e.client.Payload.Player)
	} else {
		drawCandies(screen, e.client.Payload.Candies)
		drawSnakes(screen, e.client.Payload.Player, e.client.Payload.Opponents)
		drawGameField(screen, e.client.World())
	}
	drawPlayerInfo(screen, e.client.Payload)
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

	if err := ebiten.RunGame(NewEninge(*serverAddr)); err != nil {
		log.Fatal(err)
	}
}

func buildServer(playerCount int, addr string) *gserver.GameServer {
	return gserver.New(
		playerCount,
		addr,
		game.NewGame(playerCount, gameWidth/gridSize, gameHeight/gridSize),
	)
}
