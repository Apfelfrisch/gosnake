package engine

import (
	"os"

	"github.com/apfelfrisch/gosnake/game"
	netClient "github.com/apfelfrisch/gosnake/game/network/client"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Engine struct {
	localPlayer    ClientSnake
	localOpponents []ClientSnake
	client         *netClient.GameClient
}

func New(serverAddr string, playerCount int) *Engine {
	player := NewPlayer()
	client := netClient.Connect(serverAddr, GameWidth/GridSize, GameHeight/GridSize)
	client.EventBus.Add(netClient.PlayerHasEaten{}, func(event netClient.Event) {
		player.Play(Eat)
	})
	client.EventBus.Add(netClient.PlayerDashed{}, func(event netClient.Event) {
		player.Play(Dash)
	})
	client.EventBus.Add(netClient.PlayerWalkedWall{}, func(event netClient.Event) {
		player.Play(WalkWall)
	})
	client.EventBus.Add(netClient.PlayerCrashed{}, func(event netClient.Event) {
		player.Play(Crash)
	})
	client.EventBus.Add(netClient.GameHasStarted{}, func(event netClient.Event) {
		player.PlayMusic()
	})
	client.EventBus.Add(netClient.GameHasEnded{}, func(event netClient.Event) {
		player.PauseMusic()
	})

	localPlayer := ClientSnake{
		gridSize:   GridSize,
		interPixel: 0,
	}
	return &Engine{
		client:         client,
		localPlayer:    localPlayer,
		localOpponents: []ClientSnake{},
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

	if e.client.Payload.GameState == game.Paused || e.client.Payload.GameState == game.RoundFinished {
		e.localPlayer.sync(e.client.Payload.Player)
		drawPausedScreen(screen)
	} else if e.client.Payload.GameState == game.GameFinished {
		e.localPlayer.sync(e.client.Payload.Player)
		drawFinishScreen(screen, e.client.Payload.Player)
	} else {
		drawCandies(screen, e.client.Payload.Candies)
		drawSnakes(screen, e)
		drawGameField(screen, e.client.World())
	}
	drawPlayerInfo(screen, e.client.Payload)
}

func (e *Engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return DisplayWidth, DisplayHeight
}
