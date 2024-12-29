package game

import (
	"math/rand/v2"
)

const growsSize = 5
const MapSwitch = 10

type Game struct {
	level   uint16
	gameMap *Map
	state   GameState
	players []Snake
	candies []Position
}

func NewGame(player, width, height int) *Game {
	var players []Snake

	for i := 1; i <= player; i++ {
		players = append(players, newSnake(uint16(i*5), uint16(i*5), East))
	}

	game := &Game{
		level:   1,
		gameMap: NewMap(1, uint16(width), uint16(height)),
		players: players,
	}

	game.candies = []Position{game.randomPosition()}

	return game
}

func (game *Game) State() GameState {
	return game.state
}

func (game *Game) Height() uint16 {
	return game.gameMap.Height()
}

func (game *Game) Width() uint16 {
	return game.gameMap.Width()
}

func (game *Game) TooglePaused() {
	if game.state == Ongoing {
		game.state = Paused
	} else if game.state == Paused {
		game.state = Ongoing
	}
}

func (game *Game) Reset() {
	if game.state == RoundFinished {
		for i := range game.players {
			game.players[i].reset((i+1)*5, (i+1)*5, East)
		}

		game.state = Ongoing
		game.candies[0] = game.randomPosition()
		game.gameMap = NewMap(game.level, uint16(game.Width()), uint16(game.Height()))
	} else {
		var players []Snake
		for i := 1; i <= len(game.players); i++ {
			players = append(players, newSnake(uint16(i*5), uint16(i*5), East))
		}

		game.level = 1
		game.state = Paused
		game.players = players
		game.gameMap = NewMap(game.level, uint16(game.Width()), uint16(game.Height()))
		game.candies[0] = game.randomPosition()
	}
}

func (game *Game) Field(playerIndex int, position Position) Field {
	if game.gameMap.IsWall(position) {
		return Wall
	}

	for _, candyPos := range game.candies {
		if position.Y == candyPos.Y && position.X == candyPos.X {
			return Candy
		}
	}

	for _, snakePos := range game.players[playerIndex].Occupied {
		if position.Y == snakePos.Y && position.X == snakePos.X {
			return SnakePlayer
		}
	}

	for index, player := range game.players {
		if index == playerIndex {
			continue
		}

		for _, snakePos := range player.Occupied {
			if position.Y == snakePos.Y && position.X == snakePos.X {
				return SnakeOpponent
			}
		}
	}

	return Empty
}

func (game *Game) Tick() {
	if game.state != Ongoing {
		return
	}

	for index := range game.players {
		player := &game.players[index]

		player.move()
		player.walkWalls(game)
	}

	candyCount := 0
	for index := range game.players {
		game.handelCollision(index)
		candyCount += (len(game.players[index].Occupied) + int(game.players[index].grows)) / growsSize
	}

	if candyCount >= MapSwitch {
		game.level += 1

		if game.level > 10 {
			game.state = GameFinished
		} else {
			game.state = RoundFinished
		}
	}
}

func (game *Game) handelCollision(playerIndex int) {
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
		if collision := player.head().getCollision(collisionPlayer.Occupied); collision != nil {
			handleCollision()
			return
		}
	}

	// Snake gets Candy
	if candyIndex := player.head().getCollision(game.candies); candyIndex != nil {
		player.eat(growsSize)
		game.candies[*candyIndex] = game.randomPosition()
	}
}

func (game *Game) ChangeDirection(playerIndex int, direction direction) {
	if playerIndex >= 0 && playerIndex < len(game.players) {
		game.players[playerIndex].ChangeDirection(direction)
	}
}

func (game *Game) Dash(playerIndex int) {
	if game.state != Ongoing {
		return
	}

	if playerIndex >= 0 && playerIndex < len(game.players) {
		if ok := game.players[playerIndex].Perks.use(Dash); !ok {
			return
		}

		for i := 0; i < 5; i++ {
			game.players[playerIndex].move()
			game.players[playerIndex].walkWalls(game)
			game.handelCollision(playerIndex)
		}
	}
}

func (game *Game) Players() []Snake {
	return game.players
}

func (game *Game) randomPosition() Position {
	pos := Position{
		Y: uint16(rand.N(game.gameMap.Height()-2) + 1),
		X: uint16(rand.N(game.gameMap.Width()-2) + 1),
	}

	if game.gameMap.IsWall(pos) {
		return game.randomPosition()
	}

	for _, player := range game.players {
		if collision := pos.getCollision(player.Occupied); collision != nil {
			return game.randomPosition()
		}
	}

	return pos
}

func (self Position) getCollision(others []Position) *int {
	for index, other := range others {
		if self.X == other.X && self.Y == other.Y {
			return &index
		}
	}

	return nil
}
