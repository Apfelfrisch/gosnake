package server

import (
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
	for {
		s.Update()
	}
}

func (s *GameServer) Update() {
	if time.Since(s.lastUpdate) < gameSpeed {
		time.Sleep(time.Millisecond)
		return
	}

	for connIndex := range s.tcp.conns {
		pressedKey := s.tcp.ReadConn(connIndex)

		if pressedKey == nil {
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
		} else if *pressedKey == rune('↵') {
			s.game.Reset()
		}
	}

	s.game.Tick()
	s.lastUpdate = time.Now()

	s.broadcastState()
}

func (s *GameServer) broadcastState() {
	s.tcp.Broadcast(SerializeState(s.game))
}

func SerializeState(g game.Game) string {
	var sb strings.Builder

	var x, y uint16
	for y = 1; y <= g.Height(); y++ {
		for x = 1; x <= g.Width(); x++ {
			sb.WriteString(string(g.Field(game.Position{Y: uint16(y), X: uint16(x)})))
		}
		sb.WriteRune('|')
	}

	return sb.String()
}
