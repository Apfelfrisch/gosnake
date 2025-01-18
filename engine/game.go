package engine

import (
	"context"
	"math"

	"github.com/apfelfrisch/gosnake/game"
	netClient "github.com/apfelfrisch/gosnake/game/network/client"
	netServer "github.com/apfelfrisch/gosnake/game/network/server"
)

const (
	DisplayWidth  = 1500
	DisplayHeight = 1000
	GameWidth     = 1000
	GameHeight    = 1000
	GridSize      = 20
)

type interPosition struct {
	y int
	x int
}

type Rect struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

type ClientSnake struct {
	GridSize    uint16
	ServerSnake game.Snake
	InterPixel  int
	IsGrowing   bool
}

func (cs *ClientSnake) OutOfsync(serverSnake game.Snake) bool {
	return serverSnake.Head().X != cs.ServerSnake.Head().X || serverSnake.Head().Y != cs.ServerSnake.Head().Y
}

func (cs *ClientSnake) Sync(serverSnake game.Snake) {
	cs.InterPixel = 0
	if len(serverSnake.Occupied) != len(cs.ServerSnake.Occupied) {
		cs.IsGrowing = true
	} else {
		cs.IsGrowing = false
	}
	cs.ServerSnake = serverSnake
}

func (cs *ClientSnake) interPos(direction game.Direction) interPosition {
	switch direction {
	case game.North:
		return interPosition{y: -cs.InterPixel}
	case game.South:
		return interPosition{y: cs.InterPixel}
	case game.West:
		return interPosition{x: -cs.InterPixel}
	case game.East:
		return interPosition{x: cs.InterPixel}
	default:
		panic("Unkow direction")
	}
}

func (cs *ClientSnake) Positions(dir game.Direction, pixel int) []Rect {
	cs.InterPixel += pixel

	bodies := make([]Rect, 0, len(cs.ServerSnake.Occupied))

	for i, pos := range cs.ServerSnake.Occupied {
		body := Rect{
			X:      float32(pos.X*cs.GridSize - cs.GridSize),
			Y:      float32(pos.Y*cs.GridSize - cs.GridSize),
			Width:  float32(cs.GridSize),
			Height: float32(cs.GridSize),
		}

		// resize and replace head
		if i == len(cs.ServerSnake.Occupied)-1 {
			interPos := cs.interPos(dir)
			body.Width += float32(math.Abs(float64(interPos.x)))
			body.Height += float32(math.Abs(float64(interPos.y)))
			if interPos.x < 0 {
				body.X += float32(interPos.x)
			}
			if interPos.y < 0 {
				body.Y += float32(interPos.y)
			}
		}

		// resize and replace tail only if snake
		// is not growing, otherwise it gliches
		if !cs.IsGrowing && i == 0 {
			var interPos interPosition

			if len(cs.ServerSnake.Occupied) > i+1 {
				prevPos := cs.ServerSnake.Occupied[i+1]

				if prevPos.X < pos.X {
					interPos = cs.interPos(game.West)
				} else if prevPos.X > pos.X {
					interPos = cs.interPos(game.East)
				} else if prevPos.Y < pos.Y {
					interPos = cs.interPos(game.North)
				} else if prevPos.Y > pos.Y {
					interPos = cs.interPos(game.South)
				}
			} else {
				interPos = cs.interPos(dir)
			}

			body.Width -= float32(math.Abs(float64(interPos.x)))
			body.Height -= float32(math.Abs(float64(interPos.y)))
			if interPos.x > 0 {
				body.X += float32(interPos.x)
			}
			if interPos.y > 0 {
				body.Y += float32(interPos.y)
			}
		}

		bodies = append(bodies, body)
	}

	return bodies
}

func ConnectClient(ctx context.Context, serverAddr string) (*netClient.GameClient, error) {
	client, err := netClient.Connect(ctx, serverAddr, GameWidth/GridSize, GameHeight/GridSize)

	if err != nil {
		return nil, err
	}

	player := NewPlayer()

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

	return client, nil
}

func BuildServer(playerCount int, addr string) *netServer.GameServer {
	return netServer.New(
		playerCount,
		addr,
		game.NewGame(playerCount, GameWidth/GridSize, GameHeight/GridSize),
	)
}
