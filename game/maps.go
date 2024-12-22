package game

type Map struct {
	width  uint16
	height uint16
	walls  map[uint16]map[uint16]bool
}

func NewMap(level, gameWidth, gameHeight uint16) *Map {
	walls := outerWalls(gameWidth, gameHeight)

	switch level {
	case 2:
		for width := 10; width <= 40; width++ {
			walls[15][uint16(width)] = true
		}
	case 3:
		for height := 5; height <= 25; height++ {
			walls[uint16(height)][10] = true
			walls[uint16(height)][40] = true
		}
	case 4:
		for width := 2; width <= 25; width++ {
			walls[8][uint16(width)] = true
			walls[22][uint16(51-width)] = true
		}

		for height := 2; height <= 16; height++ {
			walls[uint16(height)][35] = true
			walls[uint16(32-height)][15] = true
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

func (self *Map) IsWall(pos Position) bool {
	_, exists := self.walls[pos.Y][pos.X]

	return exists
}

func outerWalls(gameWidth, gameHeight uint16) map[uint16]map[uint16]bool {
	walls := make(map[uint16]map[uint16]bool)
	for height := uint16(1); height <= gameHeight; height++ {
		walls[height] = make(map[uint16]bool)
		for width := uint16(1); width <= gameWidth; width++ {
			// Outer Wall
			if height == 1 || height == gameHeight {
				walls[height][width] = true
			}
			if width == 1 || width == gameWidth {
				walls[height][width] = true
			}
		}
	}
	return walls
}
