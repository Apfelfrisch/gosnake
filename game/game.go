package game

import (
	"fmt"
	"math/rand/v2"
)

type battleSnake struct {
	level   int
	gameMap *Map
	state   GameState
	players []Snake
	candies []Position
}

func NewBattleSnake(player, width, height int) *battleSnake {
	var players []Snake

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
	if game.state == RoundFinished {
		fmt.Printf("Befor reset:\n")
		fmt.Printf("Pointer: %T, Cap: %v, Len: , Values: %v \n", game.players, game.players, game.players)
		for i := range game.players {
			game.players[i].reset(uint16((i+1)*5), uint16((i+1)*5), East)
		}
		fmt.Printf("After reset:\n")
		fmt.Printf("Pointer: %T, Cap: %v, Len: , Values: %v \n", game.players, game.players, game.players)
		fmt.Println("---")
		game.state = Ongoing
		game.candies[0] = game.randomPosition()
	} else {
		var players []Snake

		for i := 1; i <= len(game.players); i++ {
			players = append(players, newSnake(uint16(i*5), uint16(i*5), East))
		}

		game = NewBattleSnake(len(game.players), int(game.Width()), int(game.Height()))
	}
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
	if game.state != Ongoing {
		return
	}

	for index := range game.players {
		player := &game.players[index]

		player.move()
		player.walkWalls(game)
	}

	for index := range game.players {
		game.handelCollision(index)
	}
}

func (game *battleSnake) handelCollision(playerIndex int) {
	player := &game.players[playerIndex]

	handleCollision := func() {
		player.Lives -= 1
		if player.Lives == 0 {
			game.state = GameFinished
		} else {
			game.state = RoundFinished
		}
	}

	if game.gameMap.IsWall(player.head()) {
		handleCollision()
		return
	}
	if collision := player.head().getCollision(player.body()); collision != nil {
		handleCollision()
		return
	}

	// Snake Crushed to other Snake
	for collisionIndex, collisionPlayer := range game.players {
		if collisionIndex == playerIndex {
			continue
		}
		if collision := player.head().getCollision(collisionPlayer.occupied); collision != nil {
			handleCollision()
			return
		}
	}

	// Snake gets Candy
	if candyIndex := player.head().getCollision(game.candies); candyIndex != nil {
		player.grows += 5
		game.candies[*candyIndex] = game.randomPosition()
	}
}

func (game *battleSnake) ChangeDirection(playerIndex int, direction direction) {
	if playerIndex >= 0 && playerIndex < len(game.players) {
		game.players[playerIndex].ChangeDirection(direction)
	}
}

func (game *battleSnake) Dash(playerIndex int) {
	if game.state != Ongoing {
		return
	}

	if playerIndex >= 0 && playerIndex < len(game.players) {
		if ok := game.players[playerIndex].Perks.use(dash); !ok {
			return
		}

		for i := 0; i < 5; i++ {
			game.players[playerIndex].move()
			game.players[playerIndex].walkWalls(game)
			game.handelCollision(playerIndex)
		}
	}
}

func (game *battleSnake) Players() []Snake {
	return game.players
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
