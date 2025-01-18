package client

import (
	"context"
	"log"
	"time"

	"github.com/apfelfrisch/gosnake/game"
	"github.com/apfelfrisch/gosnake/game/network/payload"
	"google.golang.org/protobuf/proto"
)

func Connect(ctx context.Context, serverAddr string, width, height int) (*GameClient, error) {
	udp := NewUdpClient(serverAddr)

	for i := 0; i < 10; i++ {
		if err := udp.Connect(ctx); err == nil {
			break
		}
		time.Sleep(time.Second / 10)
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			time.Sleep(time.Second / 10)
		}

		if len(udp.Read()) != 0 {
			break
		}
	}

	return &GameClient{
		udp,
		game.NewMap(1, uint16(width), uint16(height)),
		&payload.Payload{},
		NewEventBus(),
	}, nil
}

type GameClient struct {
	udp      *UdpClient
	gameMap  *game.Map
	Payload  *payload.Payload
	EventBus *EventBus
}

func (gc *GameClient) PressKey(char rune) {
	gc.udp.Write(char)
}

func (gc *GameClient) UpdatePayload() {
	data := gc.udp.Read()

	if string(data) == string(HANDSHAKE_RESP) {
		return
	}

	stalePayload := *gc.Payload

	ppl := &payload.ProtoPayload{}
	err := proto.Unmarshal(gc.udp.Read(), ppl)
	if err != nil {
		log.Println(err)
		log.Println(gc.udp.Read())
		return
	}

	*gc.Payload = payload.PayloadFromProto(ppl)

	if stalePayload.MapLevel != gc.Payload.MapLevel {
		*gc.gameMap = *game.NewMap(gc.Payload.MapLevel, gc.gameMap.Width(), gc.gameMap.Height())
	}

	go func() {
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
	}()
}

func (gc *GameClient) AddListener(e Event, l EventListener) {
	gc.EventBus.Add(e, l)
}

func (gc *GameClient) World() []game.FieldPos {
	fieldPos := make([]game.FieldPos, 0, gc.gameMap.Width()*gc.gameMap.Height())

	var x, y uint16
	for y = 1; y <= gc.gameMap.Height(); y++ {
		for x = 1; x <= gc.gameMap.Width(); x++ {
			pos := game.Position{Y: uint16(y), X: x}
			if gc.gameMap.IsWall(pos) {
				fieldPos = append(fieldPos, game.FieldPos{
					Field:    game.Wall,
					Position: pos,
				})
			} else {
				fieldPos = append(fieldPos, game.FieldPos{
					Field:    game.Empty,
					Position: pos,
				})
			}
		}
	}

	return fieldPos
}
