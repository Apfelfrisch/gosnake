package client

import (
	"github.com/apfelfrisch/gosnake/game"
)

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
