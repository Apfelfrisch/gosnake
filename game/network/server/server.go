package server

import (
	"net"
	"strings"
	"time"

	"github.com/apfelfrisch/gosnake/game"
	"github.com/apfelfrisch/gosnake/game/network/payload"
	"google.golang.org/protobuf/proto"
)

const GameSpeed = time.Second / 10

type byteBuffer [1][]byte
type byteBufferChan chan [1][]byte

func New(player int, addr string, game *game.Game) *GameServer {
	return &GameServer{
		udp:  NewUdpSever(":1200", player),
		game: game,
	}
}

type GameServer struct {
	udp        *UdpServer
	game       *game.Game
	lastUpdate time.Time
}

func (s *GameServer) Addr() *net.UDPAddr {
	return s.udp.addr
}

func (s *GameServer) Ready() bool {
	return s.udp.Ready()
}

func (s *GameServer) RunBackground() {
	go func() {
		s.Run()
	}()
}

func (s *GameServer) Run() {
	s.udp.Listen()

	for s.Ready() {
		s.Update()
	}

	// Reslisten for new Connections
	s.game.Reset()
	s.udp = NewUdpSever(s.udp.addr.String(), len(s.udp.inputChans))

	s.Run()
}

func (s *GameServer) Update() {
	if time.Since(s.lastUpdate) < GameSpeed {
		time.Sleep(time.Millisecond)
		return
	}

	for connIndex, conn := range s.udp.conns {
		pressedKey := s.udp.ReadConn(conn)

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

	for i, conn := range s.udp.conns {
		opponents := make([]game.Snake, 0, len(players)-1)
		opponents = append(opponents, players[:i]...)
		opponents = append(opponents, players[i+1:]...)

		var bytes []byte
		var err error

		world := SerializeWorld(i, s.game)

		pl := payload.Payload{
			World:     world,
			GameState: s.game.State(),
			Candies:   s.game.Candies(),
			Player:    players[i],
			Opponents: opponents,
		}

		bytes, err = proto.Marshal(pl.ToProto())
		if err != nil {
			panic(err)
		}

		s.udp.WriteConn(conn, bytes)
	}
}

func SerializeWorld(playerIndex int, g *game.Game) string {
	var sb strings.Builder

	var x, y uint16
	for y = 1; y <= g.Height(); y++ {
		for x = 1; x <= g.Width(); x++ {
			field := g.Field(playerIndex, game.Position{Y: uint16(y), X: uint16(x)})
			if field == game.SnakePlayer || field == game.SnakeOpponent {
				sb.WriteRune(rune(game.Empty))
			} else {
				sb.WriteRune(rune(field))
			}
		}
		sb.WriteRune('|')
	}

	return sb.String()
}
