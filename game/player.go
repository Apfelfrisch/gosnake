package game

type Snake struct {
	Perks     Perks `json:"pk"`
	Lives     uint8 `json:"li"`
	grows     uint8
	occupied  []Position
	direction direction
}

func newSnake(x uint16, y uint16, direction direction) Snake {
	return Snake{
		Lives:     10,
		Perks:     Perks{walkWall: {Usages: 3}, dash: {Usages: 3}},
		direction: direction,
		occupied:  []Position{{X: x, Y: y}},
	}
}

func (snake *Snake) reset(x uint16, y uint16, direction direction) {
	snake.occupied = []Position{{X: x, Y: y}}
	snake.direction = direction
	snake.Perks = Perks{walkWall: {Usages: 3}, dash: {Usages: 3}}
}

func (snake *Snake) ChangeDirection(direction direction) {
	switch direction {
	case North:
		if snake.direction != South {
			snake.direction = direction
		}
	case East:
		if snake.direction != West {
			snake.direction = direction
		}
	case West:
		if snake.direction != East {
			snake.direction = direction
		}
	case South:
		if snake.direction != North {
			snake.direction = direction
		}
	}
}

func (snake *Snake) head() Position {
	if len(snake.occupied) == 0 {
		panic("Snake sould always have at least one length")
	}

	return snake.occupied[len(snake.occupied)-1]
}

func (snake *Snake) body() []Position {
	if len(snake.occupied) == 0 {
		panic("Snake sould always have at least one length")
	}

	return snake.occupied[:len(snake.occupied)-1]
}

func (snake *Snake) move() {
	if len(snake.occupied) == 0 {
		return
	}

	head := snake.head()

	var newHead Position
	switch snake.direction {
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
		snake.occupied = append(snake.occupied[1:], newHead)
	} else {
		// Move and grow the Snake
		snake.grows--
		snake.occupied = append(snake.occupied[:], newHead)
	}
}

func (snake *Snake) walkWalls(game Game) {
	position := snake.head()

	if ok := snake.Perks.use(walkWall); !ok {
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
		snake.Perks.reload(walkWall, 1)
		return
	}

	snake.occupied = append(snake.occupied[:len(snake.occupied)-1], position)
}
