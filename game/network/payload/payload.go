package payload

import "github.com/apfelfrisch/gosnake/game"

type Payload struct {
	MapLevel  uint16         `json:"w"`
	GameState game.GameState `json:"gs"`
	Candies   []game.Candy   `json:"ca"`
	Player    game.Snake     `json:"pl"`
	Opponents []game.Snake   `json:"op"`
}

func PayloadFromProto(protoPayload *ProtoPayload) Payload {
	candies := make([]game.Candy, len(protoPayload.Candies))
	for i, protoCandy := range protoPayload.Candies {
		candies[i] = candyFromProto(protoCandy)
	}

	opponents := make([]game.Snake, len(protoPayload.Opponents))
	for i, protoOpponent := range protoPayload.Opponents {
		opponents[i] = snakeFromProto(protoOpponent)
	}

	return Payload{
		MapLevel:  uint16(protoPayload.MapLevel),
		GameState: game.GameState(protoPayload.GameState),
		Candies:   candies,
		Player:    snakeFromProto(protoPayload.Player),
		Opponents: opponents,
	}
}

func (payload Payload) ToProto() *ProtoPayload {
	candies := make([]*ProtoCandy, len(payload.Candies))
	for i, candy := range payload.Candies {
		candies[i] = candyToProto(candy)
	}

	opponents := make([]*ProtoSnake, len(payload.Opponents))
	for i, opponent := range payload.Opponents {
		opponents[i] = snakeToProto(opponent)
	}

	return &ProtoPayload{
		MapLevel:  uint32(payload.MapLevel),
		GameState: ProtoGameState(payload.GameState),
		Candies:   candies,
		Player:    snakeToProto(payload.Player),
		Opponents: opponents,
	}
}

// Convert Go Position to Protobuf Position
func candyToProto(candy game.Candy) *ProtoCandy {
	return &ProtoCandy{
		Type:     ProtoCandyType(candy.CandyTpe),
		Position: positionToProto(candy.Position),
	}
}

// Convert Go Position to Protobuf Position
func positionToProto(pos game.Position) *ProtoPosition {
	return &ProtoPosition{
		Y: uint32(pos.Y),
		X: uint32(pos.X),
	}
}

// Convert Protobuf Position to Go Position
func candyFromProto(protoCandy *ProtoCandy) game.Candy {
	return game.Candy{
		CandyTpe: game.CandyTpe(protoCandy.Type),
		Position: game.Position{
			Y: uint16(protoCandy.Position.Y),
			X: uint16(protoCandy.Position.X),
		},
	}

}

// Convert Protobuf Position to Go Position
func positionFromProto(protoPos *ProtoPosition) game.Position {
	return game.Position{
		Y: uint16(protoPos.Y),
		X: uint16(protoPos.X),
	}
}

// Convert Go Snake to Protobuf Snake
func snakeToProto(snake game.Snake) *ProtoSnake {
	perks := make(map[int32]*ProtoPerk)
	for perkType, perk := range snake.Perks {
		perks[int32(perkType)] = &ProtoPerk{
			Type:   ProtoPerkType(perkType),
			Usages: uint32(perk.Usages),
		}
	}

	occupied := make([]*ProtoPosition, len(snake.Occupied))
	for i, pos := range snake.Occupied {
		occupied[i] = positionToProto(pos)
	}

	return &ProtoSnake{
		Perks:     perks,
		Lives:     uint32(snake.Lives),
		Occupied:  occupied,
		Direction: ProtoDirection(snake.Direction),
		Points:    uint32(snake.Points),
		// Grows:     uint32(snake.grows),
	}
}

// Convert Protobuf Snake to Go Snake
func snakeFromProto(protoSnake *ProtoSnake) game.Snake {
	perks := make(game.Perks)
	for perkType, perk := range protoSnake.Perks {
		perks[game.PerkType(perkType)] = game.Perk{Usages: uint16(perk.Usages)}
	}

	occupied := make([]game.Position, len(protoSnake.Occupied))
	for i, protoPos := range protoSnake.Occupied {
		occupied[i] = positionFromProto(protoPos)
	}

	return game.Snake{
		Perks:     perks,
		Lives:     uint8(protoSnake.Lives),
		Occupied:  occupied,
		Direction: game.Direction(protoSnake.Direction),
		Points:    uint16(protoSnake.Points),
	}
}
