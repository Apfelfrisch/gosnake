package game

type snake struct {
	grows     uint8
	state     PlayerState
	perk      perk
	occupied  []Position
	direction direction
}

func newSnake(x uint16, y uint16, direction direction) snake {
	return snake{
		grows:     0,
		perk:      perk{3},
		direction: direction,
		occupied:  []Position{{X: x, Y: y}},
	}
}

func (snake *snake) ChangeDirection(direction direction) {
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

func (snake *snake) State() PlayerState {
	return snake.state
}

func (snake *snake) head() Position {
	if len(snake.occupied) == 0 {
		panic("Snake sould always have at least one length")
	}

	return snake.occupied[len(snake.occupied)-1]
}

func (snake *snake) body() []Position {
	if len(snake.occupied) == 0 {
		panic("Snake sould always have at least one length")
	}

	return snake.occupied[:len(snake.occupied)-1]
}

func (snake *snake) move() {
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

func (snake *snake) walkWalls(game Game) {
	position := snake.head()

	if ok := snake.perk.use(); !ok {
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
		snake.perk.reload(1)
		return
	}

	snake.occupied = append(snake.occupied[:len(snake.occupied)-1], position)
}

type perk struct {
	usages uint16
}

func (p *perk) reload(usages uint16) {
	p.usages += usages
}

func (p *perk) use() bool {
	if p.usages == 0 {
		return false
	}

	p.usages -= 1
	return true
}

type PlayerState int

const (
	Undetermined PlayerState = iota
	Won
	Lost
)
