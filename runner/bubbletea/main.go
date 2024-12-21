package main

import (
	"time"

	"github.com/apfelfrisch/gosnake/game"
	tea "github.com/charmbracelet/bubbletea"
)

// Key contains information about a keypress.
type Key struct {
	Runes []rune
	Alt   bool
	Paste bool
}

func tick(duration time.Duration) tea.Cmd {
	return tea.Tick(duration, func(t time.Time) tea.Msg {
		return t
	})
}

func (model gameModel) Init() tea.Cmd {
	return tea.Batch(tea.ClearScreen, tick(time.Millisecond*100))
}

func (m gameModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.game.State() == game.Finished {
				m.game.Reset()
			}
		case "up":
			m.game.ChangeDirection(0, game.North)
		case "down":
			m.game.ChangeDirection(0, game.South)
		case "left":
			m.game.ChangeDirection(0, game.West)
		case "right":
			m.game.ChangeDirection(0, game.East)
		case "w":
			m.game.ChangeDirection(1, game.North)
		case "a":
			m.game.ChangeDirection(1, game.West)
		case "s":
			m.game.ChangeDirection(1, game.South)
		case "d":
			m.game.ChangeDirection(1, game.East)
		}
	case time.Time:
		m.game.Tick()
		return m, tick(time.Millisecond * 100)
	}

	return m, nil
}

func (m gameModel) View() string {
	view := ""

	var x, y uint16

	for y = 1; y <= m.game.Height(); y++ {
		for x = 1; x <= m.game.Width(); x++ {
			view += " " + string(m.game.Field(
				game.Position{Y: uint16(y), X: uint16(x)},
			))
		}
		view += "\n"
	}

	return view
}

type gameModel struct {
	game game.Game
}

func main() {
	g := game.NewSingle(50, 30)

	tui := tea.NewProgram(gameModel{game: g})
	tea.ClearScreen()

	tui.Run()
}
