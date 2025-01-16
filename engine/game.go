package engine

import (
	"math"

	"github.com/apfelfrisch/gosnake/game"
)

type interPosition struct {
	y int
	x int
}

type rect struct {
	x      float32
	y      float32
	width  float32
	height float32
}

type ClientSnake struct {
	gridSize    uint16
	serverSnake game.Snake
	interPixel  int
	isGrowing   bool
}

func (cs *ClientSnake) outOfsync(serverSnake game.Snake) bool {
	return serverSnake.Head().X != cs.serverSnake.Head().X || serverSnake.Head().Y != cs.serverSnake.Head().Y
}

func (cs *ClientSnake) sync(serverSnake game.Snake) {
	cs.interPixel = 0
	if len(serverSnake.Occupied) != len(cs.serverSnake.Occupied) {
		cs.isGrowing = true
	} else {
		cs.isGrowing = false
	}
	cs.serverSnake = serverSnake
}

func (cs *ClientSnake) interPos(direction game.Direction) interPosition {
	switch direction {
	case game.North:
		return interPosition{y: -cs.interPixel}
	case game.South:
		return interPosition{y: cs.interPixel}
	case game.West:
		return interPosition{x: -cs.interPixel}
	case game.East:
		return interPosition{x: cs.interPixel}
	default:
		panic("Unkow direction")
	}
}

func (cs *ClientSnake) positions(dir game.Direction, pixel int) []rect {
	cs.interPixel += pixel

	bodies := make([]rect, 0, len(cs.serverSnake.Occupied))

	for i, pos := range cs.serverSnake.Occupied {
		body := rect{
			x:      float32(pos.X*cs.gridSize - cs.gridSize),
			y:      float32(pos.Y*cs.gridSize - cs.gridSize),
			width:  float32(cs.gridSize),
			height: float32(cs.gridSize),
		}

		// resize and replace head
		if i == len(cs.serverSnake.Occupied)-1 {
			interPos := cs.interPos(dir)
			body.width += float32(math.Abs(float64(interPos.x)))
			body.height += float32(math.Abs(float64(interPos.y)))
			if interPos.x < 0 {
				body.x += float32(interPos.x)
			}
			if interPos.y < 0 {
				body.y += float32(interPos.y)
			}
		}

		// resize and replace tail only if snake
		// is not growing, otherwise it gliches
		if !cs.isGrowing && i == 0 {
			var interPos interPosition

			if len(cs.serverSnake.Occupied) > i+1 {
				prevPos := cs.serverSnake.Occupied[i+1]

				if prevPos.X < pos.X {
					interPos = cs.interPos(game.West)
				} else if prevPos.X > pos.X {
					interPos = cs.interPos(game.East)
				} else if prevPos.Y < pos.Y {
					interPos = cs.interPos(game.North)
				} else if prevPos.Y > pos.Y {
					interPos = cs.interPos(game.South)
				}
			} else {
				interPos = cs.interPos(dir)
			}

			body.width -= float32(math.Abs(float64(interPos.x)))
			body.height -= float32(math.Abs(float64(interPos.y)))
			if interPos.x > 0 {
				body.x += float32(interPos.x)
			}
			if interPos.y > 0 {
				body.y += float32(interPos.y)
			}
		}

		bodies = append(bodies, body)
	}

	return bodies
}
