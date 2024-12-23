package game

type Game interface {
	Tick()
	Reset()
	Height() uint16
	Width() uint16
	State() GameState
	Field(position Position) Field
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
	Empty     Field = ' '
	Wall            = 'X'
	Candy           = '☀'
	SnakeBody       = '#'
)

type GameState int

const (
	Ongoing GameState = iota
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
