package game

type Snake struct {
	Perks     Perks      `json:"pk"`
	Lives     uint8      `json:"li"`
	Occupied  []Position `json:"oc"`
	Direction direction  `json:"dr"`
	Points    uint16     `json:"pt"`
	grows     uint8
}

func newSnake(x uint16, y uint16, direction direction) Snake {
	return Snake{
		Lives:     10,
		Points:    0,
		Perks:     Perks{WalkWall: {Usages: 3}, Dash: {Usages: 3}},
		Direction: direction,
		Occupied:  []Position{{X: x, Y: y}},
		grows:     0,
	}
}

func (snake *Snake) reset(x int, y int, direction direction) {
	snake.Occupied = []Position{{X: uint16(x), Y: uint16(y)}}
	snake.Direction = direction
	snake.Perks = Perks{WalkWall: {Usages: 3}, Dash: {Usages: 3}}
	snake.grows = 0
}

func (snake *Snake) ChangeDirection(direction direction) {
	switch direction {
	case North:
		if snake.Direction != South {
			snake.Direction = direction
		}
	case East:
		if snake.Direction != West {
			snake.Direction = direction
		}
	case West:
		if snake.Direction != East {
			snake.Direction = direction
		}
	case South:
		if snake.Direction != North {
			snake.Direction = direction
		}
	}
}

func (snake *Snake) eat(grows uint8) {
	snake.grows = grows
	snake.Points += 1
}

func (snake *Snake) head() Position {
	if len(snake.Occupied) == 0 {
		panic("Snake sould always have at least one length")
	}

	return snake.Occupied[len(snake.Occupied)-1]
}

func (snake *Snake) body() []Position {
	if len(snake.Occupied) == 0 {
		panic("Snake sould always have at least one length")
	}

	return snake.Occupied[:len(snake.Occupied)-1]
}

func (snake *Snake) move() {
	if len(snake.Occupied) == 0 {
		return
	}

	head := snake.head()

	var newHead Position
	switch snake.Direction {
	case North:
		newHead = Position{Y: head.Y - 1, X: head.X}
	case East:
		newHead = Position{Y: head.Y, X: head.X + 1}
	case West:
		newHead = Position{Y: head.Y, X: head.X - 1}
	case South:
		newHead = Position{Y: head.Y + 1, X: head.X}
	default:
		panic("Unkow direction")
	}

	if snake.grows == 0 {
		// Move the Snake
		snake.Occupied = append(snake.Occupied[1:], newHead)
	} else {
		// Move and grow the Snake
		snake.grows--
		snake.Occupied = append(snake.Occupied[:], newHead)
	}
}

func (snake *Snake) walkWalls(game *Game) {
	position := snake.head()

	if ok := snake.Perks.use(WalkWall); !ok {
		return
	}

	// Walk through Walls
	if position.X > game.Width()-1 {
		position.X = 2
	} else if position.X == 1 {
		position.X = game.Width() - 1
	} else if position.Y > game.Height()-1 {
		position.Y = 2
	} else if position.Y == 1 {
		position.Y = game.Height() - 1
	} else {
		// Perk was not needed
		snake.Perks.reload(WalkWall, 1)
		return
	}

	snake.Occupied = append(snake.Occupied[:len(snake.Occupied)-1], position)
}
