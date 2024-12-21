package game

type Map struct {
	width  uint16
	height uint16
	walls  []Position
}

func NewMap(level, gameWidth, gameHeight uint16) *Map {
	walls := outerWalls(gameWidth, gameHeight)

	switch level {
	case 2:
		for width := 10; width <= 40; width++ {
			walls = append(walls, Position{15, uint16(width)})
		}
	case 3:
		for height := 5; height <= 25; height++ {
			walls = append(walls, Position{uint16(height), 10}, Position{uint16(height), 40})
		}
	case 4:
		for width := 2; width <= 25; width++ {
			walls = append(walls, Position{8, uint16(width)}, Position{22, uint16(51 - width)})
		}

		for height := 2; height <= 16; height++ {
			walls = append(walls, Position{uint16(height), 35}, Position{uint16(32 - height), 15})
		}
	}

	return &Map{gameWidth, gameHeight, walls}
}

func (self *Map) Width() uint16 {
	return self.width
}

func (self *Map) Height() uint16 {
	return self.height
}

func (self *Map) Walls() []Position {
	return self.walls
}

func (self *Map) IsWall(pos Position) bool {
	for _, wall := range self.walls {
		if wall.X == pos.X && wall.Y == pos.Y {
			return true
		}
	}

	return false
}

func outerWalls(gameWidth, gameHeight uint16) []Position {
	var walls []Position
	for height := uint16(1); height <= gameHeight; height++ {
		for width := uint16(1); width <= gameWidth; width++ {
			// Outer Wall
			if height == 1 || height == gameHeight {
				walls = append(walls, Position{height, width})
			}
			if width == 1 || width == gameWidth {
				walls = append(walls, Position{height, width})
			}
		}
	}
	return walls
}
