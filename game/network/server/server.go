package server

import (
	"net"
	"time"

	"github.com/apfelfrisch/gosnake/game"
	"github.com/apfelfrisch/gosnake/game/network/payload"
	"google.golang.org/protobuf/proto"
)

const GameSpeed = time.Second / 10
const PackageIntervall = GameSpeed / 3

type byteBuffer [1][]byte
type byteBufferChan chan [1][]byte

func New(player int, addr string, game *game.Game) *GameServer {
	return &GameServer{
		udp:  NewUdpSever(":1200", player),
		game: game,
	}
}

type GameServer struct {
	udp             *UdpServer
	game            *game.Game
	lastUpdate      time.Time
	lastPackageSend time.Time
}

func (s *GameServer) Addr() *net.UDPAddr {
	return s.udp.addr
}

func (s *GameServer) Clients() []*net.UDPAddr {
	return s.udp.clients
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
		// Resend state to because of package lost
		if time.Since(s.lastPackageSend) > PackageIntervall {
			s.broadcastState()
			s.lastPackageSend = time.Now()
		}

		time.Sleep(time.Millisecond)
		return
	}

	for connIndex, conn := range s.udp.clients {
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
	s.broadcastState()

	s.lastUpdate = time.Now()
	s.lastPackageSend = time.Now()
}

func (s *GameServer) broadcastState() {
	players := s.game.Players()
	for i, conn := range s.udp.clients {
		opponents := make([]game.Snake, 0, len(players)-1)
		opponents = append(opponents, players[:i]...)
		opponents = append(opponents, players[i+1:]...)

		var bytes []byte
		var err error

		pl := payload.Payload{
			MapLevel:  s.game.Level(),
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
