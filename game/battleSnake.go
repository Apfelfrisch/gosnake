package game

import "math/rand/v2"

type battleSnake struct {
	width   uint16
	height  uint16
	state   GameState
	scores  []int
	players []snake
	candies []Position
}

func NewBattleSnake(width int, height int) battleSnake {
	return battleSnake{
		width:   uint16(width),
		height:  uint16(height),
		players: []snake{newSnake(5, 5, East), newSnake(20, 20, West)},
		scores:  []int{0, 0},
		candies: []Position{{
			Y: uint16(rand.N(height-2) + 1),
			X: uint16(rand.N(width-2) + 1),
		}},
	}
}

func (game *battleSnake) State() GameState {
	return game.state
}

func (game *battleSnake) Height() uint16 {
	return game.height
}

func (game *battleSnake) Width() uint16 {
	return game.width
}

func (game *battleSnake) Reset() {
	game.state = Ongoing
	game.players = []snake{newSnake(5, 5, East), newSnake(20, 20, West)}
	game.candies[0] = game.randomPosition()
}

func (game *battleSnake) Field(position Position) Field {
	if position.X == 0 || position.X >= game.width-1 {
		return Wall
	}

	if position.Y == 0 || position.Y >= game.height-1 {
		return Wall
	}

	for _, candyPos := range game.candies {
		if position.Y == candyPos.Y && position.X == candyPos.X {
			return Candy
		}
	}

	for _, player := range game.players {
		for _, snakePos := range player.occupied {
			if position.Y == snakePos.Y && position.X == snakePos.X {
				return SnakeBody
			}
		}
	}

	return Empty
}

func (game *battleSnake) Tick() {
	if game.state == Finished {
		return
	}

	for index := range game.players {
		player := &game.players[index]

		player.move()

		player.beamThroughWall(game)

		for collisionIndex, collisionPlayer := range game.players {
			// Snake Crushed to its own Body
			if collisionIndex == index {
				if collision := player.head().getCollision(collisionPlayer.body()); collision != nil {
					game.state = Finished
					game.scores[index]--
				}
				continue
			}

			// Snake Crushed to other Snake
			if collision := player.head().getCollision(collisionPlayer.occupied); collision != nil {
				game.state = Finished
				game.scores[index]--
			}
		}

		// Snake gets Candy
		if candyIndex := player.head().getCollision(game.candies); candyIndex == nil {
			player.grows += 5
			game.candies[*candyIndex] = game.randomPosition()
		}
	}
}

func (game *battleSnake) ChangeDirection(playerIndex int, direction direction) {
	if playerIndex >= 0 && playerIndex < len(game.players) {
		game.players[playerIndex].ChangeDirection(direction)
	}
}

func (game *battleSnake) randomPosition() Position {
	return Position{
		Y: uint16(rand.N(game.height-2) + 1),
		X: uint16(rand.N(game.width-2) + 1),
	}
}
