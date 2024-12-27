package game

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

type FieldPos struct {
	Field
	Position
}
