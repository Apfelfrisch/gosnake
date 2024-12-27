package client

import (
	"encoding/json"
	"time"

	"github.com/apfelfrisch/gosnake/game"
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

	return &GameClient{tcp: tcp, Payload: &Payload{}}
}

type GameClient struct {
	tcp     *Tcp
	Payload *Payload
}

func (gc *GameClient) PressKey(char rune) {
	gc.tcp.Write(char)
}

func (gc *GameClient) UpdatePayload() {
	json.Unmarshal([]byte(gc.tcp.Read()), gc.Payload)
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
