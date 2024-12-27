package server

import (
	"encoding/json"
	"net"
	"strings"
	"time"

	"github.com/apfelfrisch/gosnake/game"
)

const gameSpeed = time.Second / 10

func New(player int, addr string, game game.Game) *GameServer {
	return &GameServer{
		tcp:  NewTcpSever(":1200", player),
		game: game,
	}
}

type GameServer struct {
	tcp        *Tcp
	game       game.Game
	lastUpdate time.Time
}

func (s *GameServer) Addr() *net.TCPAddr {
	return s.tcp.addr
}

func (s *GameServer) Ready() bool {
	return s.tcp.Ready()
}

func (s *GameServer) RunBackground() {
	go func() {
		s.Run()
	}()
}

func (s *GameServer) Run() {
	s.tcp.Listen()
	s.broadcastState()

	for s.Ready() {
		s.Update()
	}

	// Reslisten for new Connections
	s.game.Reset()
	s.tcp = NewTcpSever(s.tcp.addr.String(), len(s.tcp.inputChans))

	s.Run()
}

func (s *GameServer) Update() {
	if time.Since(s.lastUpdate) < gameSpeed {
		time.Sleep(time.Millisecond)
		return
	}

	for connIndex := range s.tcp.conns.get() {
		pressedKey := s.tcp.ReadConn(connIndex)

		if pressedKey == nil {
			continue
		}

		if s.game.State() != game.Ongoing {
			if *pressedKey == rune('↵') {
				if s.game.State() == game.Paused {
					s.game.TooglePaused()
				} else {
					s.game.Reset()
				}
				return
			}
			continue
		}

		if *pressedKey == rune('w') {
			s.game.ChangeDirection(connIndex, game.North)
		} else if *pressedKey == rune('s') {
			s.game.ChangeDirection(connIndex, game.South)
		} else if *pressedKey == rune('a') {
			s.game.ChangeDirection(connIndex, game.West)
		} else if *pressedKey == rune('d') {
			s.game.ChangeDirection(connIndex, game.East)
		} else if *pressedKey == rune(' ') {
			s.game.Dash(connIndex)
		}
	}

	s.game.Tick()
	s.lastUpdate = time.Now()

	s.broadcastState()
}

func (s *GameServer) broadcastState() {
	players := s.game.Players()

	for i := range s.tcp.conns.get() {
		opponents := make([]game.Snake, 0, len(players)-1)
		opponents = append(opponents, players[:i]...)
		opponents = append(opponents, players[i+1:]...)

		bytes, err := json.Marshal(Payload{
			World:     SerializeWorld(i, s.game),
			GameState: s.game.State(),
			Player:    players[i],
			Opponents: opponents,
		})
		if err != nil {
			continue
		}
		s.tcp.WriteConn(i, bytes)
	}
}

func SerializeWorld(playerIndex int, g game.Game) string {
	var sb strings.Builder

	var x, y uint16
	for y = 1; y <= g.Height(); y++ {
		for x = 1; x <= g.Width(); x++ {
			sb.WriteRune(rune(g.Field(playerIndex, game.Position{Y: uint16(y), X: uint16(x)})))
		}
		sb.WriteRune('|')
	}

	return sb.String()
}
