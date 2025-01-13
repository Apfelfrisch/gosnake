package game

type Snake struct {
	Perks        Perks      `json:"pk"`
	Lives        uint8      `json:"li"`
	Occupied     []Position `json:"oc"`
	Direction    Direction  `json:"dr"`
	NewDirection Direction  `json:"nd"`
	Points       uint16     `json:"pt"`
	grows        uint8
}

func NewSnake(x uint16, y uint16, direction Direction) Snake {
	return Snake{
		Lives:        10,
		Points:       0,
		Perks:        Perks{WalkWall: {Usages: 3}, Dash: {Usages: 3}},
		Direction:    direction,
		NewDirection: direction,
		Occupied:     []Position{{X: x, Y: y}},
		grows:        0,
	}
}

func (snake *Snake) reset(x uint16, y uint16, direction Direction) {
	snake.Occupied = []Position{{X: x, Y: y}}
	snake.Direction = direction
	snake.NewDirection = direction
	snake.Perks = Perks{WalkWall: {Usages: 3}, Dash: {Usages: 3}}
	snake.grows = 0
}

func (snake *Snake) ChangeDirection(direction Direction) {
	switch direction {
	case North:
		if snake.NewDirection != South {
			snake.NewDirection = direction
		}
	case East:
		if snake.NewDirection != West {
			snake.NewDirection = direction
		}
	case West:
		if snake.NewDirection != East {
			snake.NewDirection = direction
		}
	case South:
		if snake.NewDirection != North {
			snake.NewDirection = direction
		}
	}
}

func (snake *Snake) eat(grows uint8) {
	snake.grows = grows
	snake.Points += 1
}

func (snake *Snake) Head() Position {
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

	newHead := snake.Head().Move(snake.Direction)

	if snake.grows == 0 {
		// Move the Snake
		snake.Occupied = append(snake.Occupied[1:], newHead)
	} else {
		// Move and grow the Snake
		snake.grows--
		snake.Occupied = append(snake.Occupied[:], newHead)
	}
	snake.Direction = snake.NewDirection
}

func (snake *Snake) walkWalls(game *Game) {
	position := snake.Head()

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
