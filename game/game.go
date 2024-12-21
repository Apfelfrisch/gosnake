package game

type Game interface {
	Tick()
	Reset()
	Height() uint16
	Width() uint16
	State() GameState
	Field(position Position) Field
	ChangeDirection(playerIndex int, direction direction)
}

type direction int

const (
	East direction = iota
	North
	South
	West
)

type Field string

const (
	Empty     Field = " "
	Wall            = "X"
	Candy           = "☀"
	SnakeBody       = "#"
)

type GameState int

const (
	Ongoing GameState = iota
	Finished
)

type Position struct {
	Y uint16
	X uint16
}

type FieldPos struct {
	Field
	Position
}

func (self Position) getCollision(others []Position) *int {
	for index, other := range others {
		if self.X == other.X && self.Y == other.Y {
			return &index
		}
	}

	return nil
}
