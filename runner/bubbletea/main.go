package main

import (
	"os"
	"time"

	gclient "github.com/apfelfrisch/gosnake/game/client"
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
			//
		case "up":
		case "w":
			m.client.Write('w')
		case "down":
		case "s":
			m.client.Write('s')
		case "left":
		case "a":
			m.client.Write('a')
		case "right":
		case "d":
			m.client.Write('d')
		}
	case time.Time:
		return m, tick(time.Millisecond * 100)
	}

	return m, nil
}

func (m gameModel) View() string {
	view := ""

	for _, fieldPos := range gclient.DeserializeState(m.client.Read()) {
		if fieldPos.X == 1 && view != "" {
			view += "\n"
		}

		view += " " + string(fieldPos.Field)
	}

	return view
}

type gameModel struct {
	client *gclient.Tcp
}

func main() {
	client := connectClient(os.Args[1])

	tui := tea.NewProgram(gameModel{client})
	tea.ClearScreen()

	tui.Run()
}

func connectClient(addr string) *gclient.Tcp {
	client := gclient.NewTcpClient(addr)

	for i := 0; i < 10; i++ {
		if err := client.Connect(); err == nil {
			break
		}
		time.Sleep(time.Second / 5)
	}

	for {
		if client.Read() != "" {
			break
		}
		time.Sleep(time.Second / 10)
	}

	return client
}
