package game

type CandyTpe int

const (
	CandyGrow     CandyTpe = 0
	CandyWalkWall CandyTpe = 1
	CandyDash     CandyTpe = 2
)

type Candy struct {
	CandyTpe
	Position
}

func NewCandyGrow(pos Position) Candy {
	return Candy{
		CandyTpe: CandyGrow,
		Position: pos,
	}
}
