package client

import (
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

	return &GameClient{tcp}
}

type GameClient struct {
	tcp *Tcp
}

func (gc GameClient) PressKey(char rune) {
	gc.tcp.Write(char)
}

func (gc GameClient) World() []game.FieldPos {
	worldString := gc.tcp.Read()

	fieldPos := make([]game.FieldPos, 0, len(worldString))

	var x, y uint16 = 1, 1
	for _, char := range []rune(worldString) {
		if char == '|' {
			x = 1
			y += 1
			continue
		}

		fieldPos = append(fieldPos, game.FieldPos{
			Field:    game.Field(string(char)),
			Position: game.Position{Y: uint16(y), X: x},
		})

		x += 1
	}

	return fieldPos
}

func DeserializeState(state string) []game.FieldPos {
	fieldPos := make([]game.FieldPos, 0, len(state))

	var x, y uint16 = 1, 1
	for _, char := range []rune(state) {
		if char == '|' {
			x = 1
			y += 1
			continue
		}

		fieldPos = append(fieldPos, game.FieldPos{
			Field:    game.Field(string(char)),
			Position: game.Position{Y: uint16(y), X: x},
		})

		x += 1
	}

	return fieldPos
}
