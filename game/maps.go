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
		wallLen := gameWidth / 2
		for width := wallLen / 2; width <= gameWidth-wallLen/2; width++ {
			walls[gameHeight/2-1][uint16(width)] = true
		}
	case 3:
		for height := gameHeight / 4; height <= gameHeight-gameHeight/4; height++ {
			walls[uint16(height)][10] = true
			walls[uint16(height)][40] = true
		}
	case 4:
		wallLen := gameWidth / 2
		for width := uint16(1); width <= wallLen; width++ {
			walls[gameHeight/4][gameWidth-width] = true
			walls[gameHeight-gameHeight/4][width] = true
		}

		for height := uint16(1); height <= wallLen; height++ {
			walls[gameHeight-height][gameWidth-gameWidth/4] = true
			walls[height][gameWidth/4] = true
		}
	case 5:
		wallLen := gameWidth / 2
		for width := uint16(1); width <= wallLen-2; width++ {
			walls[gameHeight/4][gameWidth/4+width+1] = true
			walls[gameHeight-gameHeight/4][gameWidth/4+width+1] = true
		}
		for height := uint16(1); height <= wallLen-2; height++ {
			walls[gameHeight/4+height+1][gameWidth/4] = true
			walls[gameHeight/4+height+1][gameHeight-gameHeight/4] = true
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
