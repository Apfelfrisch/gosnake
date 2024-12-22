package game

import "math/rand/v2"

type battleSnake struct {
	level   int
	gameMap *Map
	state   GameState
	players []snake
	candies []Position
}

func NewBattleSnake(player, width, height int) *battleSnake {
	var players []snake

	for i := 1; i <= player; i++ {
		players = append(players, newSnake(uint16(i*5), uint16(i*5), East))
	}

	return &battleSnake{
		level:   1,
		gameMap: NewMap(1, uint16(width), uint16(height)),
		players: players,
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
	return game.gameMap.Height()
}

func (game *battleSnake) Width() uint16 {
	return game.gameMap.Width()
}

func (game *battleSnake) Reset() {
	var players []snake

	for i := range game.players {
		players = append(players, newSnake(uint16(i*5), uint16(i*5), East))
	}

	game.state = Ongoing
	game.players = players
	game.candies[0] = game.randomPosition()
}

func (game *battleSnake) Field(position Position) Field {
	if game.gameMap.IsWall(position) {
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

		player.walkWalls(game)

		if game.gameMap.IsWall(player.head()) {
			// Snake Crushed to other Snake
			game.state = Finished
		}

		for collisionIndex, collisionPlayer := range game.players {
			// Snake Crushed to its own Body
			if collisionIndex == index {
				if collision := player.head().getCollision(collisionPlayer.body()); collision != nil {
					game.state = Finished
				}
				continue
			}

			// Snake Crushed to other Snake
			if collision := player.head().getCollision(collisionPlayer.occupied); collision != nil {
				game.state = Finished
			}
		}

		// Snake gets Candy
		if candyIndex := player.head().getCollision(game.candies); candyIndex != nil {
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

func (self Position) getCollision(others []Position) *int {
	for index, other := range others {
		if self.X == other.X && self.Y == other.Y {
			return &index
		}
	}

	return nil
}
