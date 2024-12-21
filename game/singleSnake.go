package game

import "math/rand/v2"

type singleSnake struct {
	level   int
	gameMap *Map
	state   GameState
	score   int
	player  snake
	candies []Position
}

func NewSingle(width, height int) *singleSnake {
	return &singleSnake{
		level:   1,
		gameMap: NewMap(1, uint16(width), uint16(height)),
		player:  newSnake(5, 5, East),
		score:   0,
		candies: []Position{{
			Y: uint16(rand.N(height-2) + 1),
			X: uint16(rand.N(width-2) + 1),
		}},
	}
}

func (game *singleSnake) State() GameState {
	return game.state
}

func (game *singleSnake) Height() uint16 {
	return game.gameMap.Height()
}

func (game *singleSnake) Width() uint16 {
	return game.gameMap.Width()
}

func (game *singleSnake) Reset() {
	game.state = Ongoing
	game.player = newSnake(5, 5, East)
	game.candies[0] = game.randomPosition()
}

func (game *singleSnake) Field(position Position) Field {
	if game.gameMap.IsWall(position) {
		return Wall
	}

	for _, candyPos := range game.candies {
		if position.Y == candyPos.Y && position.X == candyPos.X {
			return Candy
		}
	}

	for _, snakePos := range game.player.occupied {
		if position.Y == snakePos.Y && position.X == snakePos.X {
			return SnakeBody
		}
	}

	return Empty
}

func (game *singleSnake) Tick() {
	if game.score >= 2 {
		game.level++
		game.score = 0
		game.gameMap = NewMap(uint16(game.level), game.Width(), game.Height())
		game.state = Finished
		game.player = newSnake(5, 5, East)
	}

	if game.state == Finished {
		return
	}

	game.player.move()

	game.player.beamThroughWall(game)

	// Snakes collides with Wall
	if collision := game.player.head().getCollision(game.gameMap.Walls()); collision != nil {
		game.state = Finished
		return
	}

	// Snakes collide with itself
	if collision := game.player.head().getCollision(game.player.body()); collision != nil {
		game.state = Finished
		return
	}

	// Snake gets Candy
	if candyIndex := game.player.head().getCollision(game.candies); candyIndex != nil {
		game.player.grows += 5
		game.candies[*candyIndex] = game.randomPosition()
		game.score++
	}
}

func (game *singleSnake) ChangeDirection(playerIndex int, direction direction) {
	game.player.ChangeDirection(direction)
}

func (game *singleSnake) randomPosition() Position {
	for {
		pos := Position{
			Y: uint16(rand.N(game.gameMap.Height()-2) + 1),
			X: uint16(rand.N(game.gameMap.Width()-2) + 1),
		}

		if !game.gameMap.IsWall(pos) {
			return pos
		}
	}
}
