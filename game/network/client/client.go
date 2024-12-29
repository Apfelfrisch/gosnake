package client

import (
	"encoding/json"
	"time"

	"github.com/apfelfrisch/gosnake/game"
	"github.com/apfelfrisch/gosnake/game/network/payload"
)

func Connect(serverAddr string) *GameClient {
	tcp := NewTcpClient(serverAddr)

	for i := 0; i < 10; i++ {
		if err := tcp.Connect(); err == nil {
			break
		}
		time.Sleep(time.Second / 5)
	}

	for {
		if tcp.Read() != "" {
			break
		}
		time.Sleep(time.Second / 10)
	}

	return &GameClient{tcp, &payload.Payload{}, NewEventBus(), ""}
}

type GameClient struct {
	tcp         *Tcp
	Payload     *payload.Payload
	EventBus    *EventBus
	cachedWorld string
}

func (gc *GameClient) PressKey(char rune) {
	gc.tcp.Write(char)
}

func (gc *GameClient) UpdatePayload() {
	stalePayload := gc.Payload
	gc.Payload = &payload.Payload{}

	json.Unmarshal([]byte(gc.tcp.Read()), gc.Payload)

	if gc.Payload.World != "" {
		gc.cachedWorld = gc.Payload.World
	} else {
		gc.Payload.World = gc.cachedWorld
	}

	if stalePayload.GameState != gc.Payload.GameState {
		if gc.Payload.GameState == game.Ongoing {
			gc.EventBus.Dispatch(GameHasStarted{})
		} else {
			gc.EventBus.Dispatch(GameHasEnded{})
		}
	}

	if stalePayload.GameState != game.Ongoing {
		return
	}

	if stalePayload.Player.Lives != gc.Payload.Player.Lives {
		gc.EventBus.Dispatch(PlayerCrashed{})
	} else {
		for i, opp := range gc.Payload.Opponents {
			if opp.Lives != stalePayload.Opponents[i].Lives {
				gc.EventBus.Dispatch(PlayerCrashed{})
				break
			}
		}
	}

	if gc.Payload.GameState != game.Ongoing {
		return
	}

	if stalePayload.Player.Points != gc.Payload.Player.Points {
		gc.EventBus.Dispatch(PlayerHasEaten{})
	} else {
		for i, opp := range gc.Payload.Opponents {
			if opp.Points != stalePayload.Opponents[i].Points {
				gc.EventBus.Dispatch(PlayerHasEaten{})
				break
			}
		}
	}

	if stalePayload.Player.Perks.Get(game.Dash).Usages != gc.Payload.Player.Perks.Get(game.Dash).Usages {
		gc.EventBus.Dispatch(PlayerDashed{})
	} else {
		for i, opp := range gc.Payload.Opponents {
			if opp.Perks.Get(game.Dash).Usages != stalePayload.Opponents[i].Perks.Get(game.Dash).Usages {
				gc.EventBus.Dispatch(PlayerDashed{})
				break
			}
		}
	}

	if stalePayload.Player.Perks.Get(game.WalkWall).Usages != gc.Payload.Player.Perks.Get(game.WalkWall).Usages {
		gc.EventBus.Dispatch(PlayerWalkedWall{})
	} else {
		for i, opp := range gc.Payload.Opponents {
			if opp.Perks.Get(game.WalkWall).Usages != stalePayload.Opponents[i].Perks.Get(game.WalkWall).Usages {
				gc.EventBus.Dispatch(PlayerWalkedWall{})
				break
			}
		}
	}
}

func (gc *GameClient) AddListener(e Event, l EventListener) {
	gc.EventBus.Add(e, l)
}

func (gc *GameClient) World() []game.FieldPos {
	fieldPos := make([]game.FieldPos, 0, len(gc.Payload.World))

	var x, y uint16 = 1, 1
	for _, char := range []rune(gc.Payload.World) {
		if char == '|' {
			x = 1
			y += 1
			continue
		}

		fieldPos = append(fieldPos, game.FieldPos{
			Field:    game.Field(char),
			Position: game.Position{Y: uint16(y), X: x},
		})

		x += 1
	}

	return fieldPos
}
