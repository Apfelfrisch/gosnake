package client

import "github.com/apfelfrisch/gosnake/game"

type Payload struct {
	World     string         `json:"w"`
	GameState game.GameState `json:"gs"`
	Player    game.Snake     `json:"pl"`
	Opponents []game.Snake   `json:"op"`
}
