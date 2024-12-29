package game

import "encoding/json"

type direction int

const (
	East direction = iota
	North
	South
	West
)

type Field rune

const (
	Empty         Field = ' '
	Wall                = 'X'
	Candy               = '☀'
	SnakePlayer         = '0'
	SnakeOpponent       = '1'
)

type GameState int

const (
	Paused GameState = iota
	Ongoing
	RoundFinished
	GameFinished
)

type Position struct {
	Y uint16
	X uint16
}

func (p Position) MarshalJSON() ([]byte, error) {
	return json.Marshal([2]uint16{p.Y, p.X})
}

func (p *Position) UnmarshalJSON(data []byte) error {
	var arr [2]uint16
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	p.Y, p.X = arr[0], arr[1]
	return nil
}

type FieldPos struct {
	Field
	Position
}
