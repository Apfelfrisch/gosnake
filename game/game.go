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
	candies []Candy
}

func NewGame(player, width, height int) *Game {
	game := &Game{
		level:   1,
		gameMap: NewMap(1, uint16(width), uint16(height)),
	}

	var players []Snake
	for i := 1; i <= player; i++ {
		startPos := game.randomPosition()
		players = append(players, NewSnake(startPos.X, startPos.Y, game.gameMap.FarestWall(startPos)))
	}

	game.players = players
	game.candies = []Candy{
		NewCandyGrow(game.randomPosition()),
	}

	return game
}

func (game *Game) Level() uint16 {
	return game.level
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
		game.state = Ongoing
		game.candies = []Candy{NewCandyGrow(game.randomPosition())}
		game.gameMap = NewMap(game.level, uint16(game.Width()), uint16(game.Height()))

		for i := range game.players {
			startPos := game.randomPosition()
			game.players[i].reset(startPos.X, startPos.Y, game.gameMap.FarestWall(startPos))
		}
	} else {
		game.level = 1
		game.state = Paused
		game.gameMap = NewMap(game.level, uint16(game.Width()), uint16(game.Height()))
		game.candies = []Candy{NewCandyGrow(game.randomPosition())}

		for i := range game.players {
			startPos := game.randomPosition()
			game.players[i] = NewSnake(startPos.X, startPos.Y, game.gameMap.FarestWall(startPos))
		}
	}
}

func (game *Game) Field(playerIndex int, position Position) Field {
	if game.gameMap.IsWall(position) {
		return FieldWall
	}

	for _, candyPos := range game.candies {
		if position.Y == candyPos.Y && position.X == candyPos.X {
			return FieldCandy
		}
	}

	for _, snakePos := range game.players[playerIndex].Occupied {
		if position.Y == snakePos.Y && position.X == snakePos.X {
			return FieldSnakePlayer
		}
	}

	for index, player := range game.players {
		if index == playerIndex {
			continue
		}

		for _, snakePos := range player.Occupied {
			if position.Y == snakePos.Y && position.X == snakePos.X {
				return FieldSnakeOpponent
			}
		}
	}

	return FieldEmpty
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

	// Spawn WalkWall
	if rand.IntN(250) == 0 {
		game.candies = append(game.candies, Candy{
			CandyTpe: CandyWalkWall,
			Position: game.randomPosition(),
		})
	}

	// Spawn Dash
	if rand.IntN(250) == 0 {
		game.candies = append(game.candies, Candy{
			CandyTpe: CandyDash,
			Position: game.randomPosition(),
		})
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

	if game.gameMap.IsWall(player.Head()) {
		handleCollision()
		return
	}
	if collision := player.Head().getCollision(player.body()); collision != nil {
		handleCollision()
		return
	}

	// Snake Crushed to other Snake
	for collisionIndex, collisionPlayer := range game.players {
		if collisionIndex == playerIndex {
			continue
		}
		if collision := player.Head().getCollision(collisionPlayer.Occupied); collision != nil {
			handleCollision()
			return
		}
	}

	for i := len(game.candies) - 1; i >= 0; i-- {
		candy := game.candies[i]
		if candyIndex := player.Head().getCollision([]Position{candy.Position}); candyIndex != nil {
			switch candy.CandyTpe {
			case CandyGrow:
				player.eat(growsSize)
				game.candies[i] = NewCandyGrow(game.randomPosition())
			case CandyDash:
				player.Perks.add(PerkTypeDash, 1)
				game.candies = append(game.candies[:i], game.candies[i+1:]...)
				continue
			case CandyWalkWall:
				player.Perks.add(PerkTypeWalkWall, 1)
				game.candies = append(game.candies[:i], game.candies[i+1:]...)
				continue
			}
		}
	}
}

func (game *Game) ChangeDirection(playerIndex int, direction Direction) {
	if playerIndex >= 0 && playerIndex < len(game.players) {
		game.players[playerIndex].ChangeDirection(direction)
	}
}

func (game *Game) Dash(playerIndex int) {
	if game.state != Ongoing {
		return
	}

	if playerIndex >= 0 && playerIndex < len(game.players) {
		if ok := game.players[playerIndex].Perks.use(PerkTypeDash); !ok {
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

func (game *Game) Candies() []Candy {
	return game.candies
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
