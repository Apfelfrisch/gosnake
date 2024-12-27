package game

type Game interface {
	Tick()
	Reset()
	TooglePaused()
	Height() uint16
	Width() uint16
	State() GameState
	Field(playerIndex int, position Position) Field
	ChangeDirection(playerIndex int, direction direction)
	Dash(playerIndex int)
	Players() []Snake
}

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
