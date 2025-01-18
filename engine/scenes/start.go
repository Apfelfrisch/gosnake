package scenes

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"strconv"
	"time"

	"github.com/apfelfrisch/gosnake/engine"
	"github.com/apfelfrisch/gosnake/game"
	netServer "github.com/apfelfrisch/gosnake/game/network/server"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type gametype int

const (
	singleplayer gametype = 0
	client       gametype = 1
	server       gametype = 2
)

func (gt gametype) prev() gametype {
	index := int(gt) - 1
	if index < 0 {
		index = 2
	}
	return gametype(index)
}

func (gt gametype) next() gametype {
	index := int(gt) + 1
	if index > 2 {
		index = 0
	}
	return gametype(index)
}

type connState int

const (
	connClosed   connState = 0
	connPending  connState = 1
	connFinished connState = 2
)

type blink struct {
	visible   bool
	lastBlink time.Time
}

func (c *blink) blink() {
	if time.Since(c.lastBlink) > 500*time.Millisecond {
		c.visible = !c.visible
		c.lastBlink = time.Now()
	}
}

func (c *blink) Show(text string) string {
	if c.visible {
		return text
	}
	return ""
}

type MenuStart struct {
	BaseScene
	gametype    gametype
	connection  connState
	playerCount int
	blink       blink
	serverAddr  string
	server      *netServer.GameServer
	ctx         context.Context
	cancle      context.CancelFunc
}

func New() *MenuStart {
	ctx, cancel := context.WithCancel(context.Background())

	return &MenuStart{
		ctx:    ctx,
		cancle: cancel,
		BaseScene: BaseScene{
			bounds: image.Rectangle{},
			localPlayer: engine.ClientSnake{
				GridSize:   engine.GridSize,
				InterPixel: 0,
			},
			localOpponents: []engine.ClientSnake{},
		},
	}
}

func (s *MenuStart) Update() error {
	s.blink.blink()

	if s.client != nil {
		s.client.UpdatePayload()
		s.localPlayer.Sync(s.client.Payload.Player)

		if s.client.Payload.GameState == game.Ongoing {
			s.sm.SwitchTo(&GameRunning{BaseScene: s.BaseScene})
		}

		return nil
	}

	if s.connection == connPending {
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			s.cancle()
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		s.connect()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		s.gametype = s.gametype.prev()
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		s.gametype = s.gametype.next()
	}

	if s.gametype == client {
		for _, char := range ebiten.AppendInputChars(nil) {
			if char >= '0' && char <= '9' || char == '.' {
				s.serverAddr += string(char)
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(s.serverAddr) > 0 {
			s.serverAddr = s.serverAddr[:len(s.serverAddr)-1]
		}
	} else if s.gametype == server {
		if s.playerCount < 2 {
			s.playerCount = 2
		}

		for _, char := range ebiten.AppendInputChars(nil) {
			if char >= '2' && char <= '9' {
				s.playerCount, _ = strconv.Atoi(string(char))
			}
		}
	}

	return nil
}

func (s *MenuStart) Draw(screen *ebiten.Image) {
	menuFont, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal(err)
	}

	face := &text.GoTextFace{
		Source: menuFont,
		Size:   30.0,
	}

	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)
	op.ColorScale.ScaleAlpha(0.5)

	s.drawLeftMenu(screen, face, op)
	s.drawContextMenu(screen, face, op)
}

func (s *MenuStart) drawLeftMenu(screen *ebiten.Image, face *text.GoTextFace, op *text.DrawOptions) {
	op.GeoM.Translate(50, 50)
	text.Draw(screen, "Singleplayer", face, op)

	op.GeoM.Translate(0, 50)
	text.Draw(screen, "Client", face, op)

	op.GeoM.Translate(0, 50)
	text.Draw(screen, "Server", face, op)

	switch s.gametype {
	case singleplayer:
		op.GeoM.Translate(-40, -100)
		text.Draw(screen, "->", face, op)
	case client:
		op.GeoM.Translate(-40, -50)
		text.Draw(screen, "->", face, op)
		op.GeoM.Translate(350, -50)
		if s.connection == connPending {
			text.Draw(screen, "Server Adresse: "+s.serverAddr, face, op)
		} else {
			text.Draw(screen, "Server Adresse: "+s.serverAddr+s.blink.Show("|"), face, op)
		}
	case server:
		op.GeoM.Translate(-40, 0)
		text.Draw(screen, "->", face, op)
		op.GeoM.Translate(350, -100)
		if s.connection == connPending {
			text.Draw(screen, "Anzahl Spieler: "+strconv.Itoa(s.playerCount), face, op)
		} else {
			text.Draw(screen, "Anzahl Spieler: "+s.blink.Show(strconv.Itoa(s.playerCount)), face, op)
		}
	default:
		panic("unexpected scenes.gametype")
	}
}

func (s *MenuStart) drawContextMenu(screen *ebiten.Image, face *text.GoTextFace, op *text.DrawOptions) {
	if s.gametype == singleplayer {
		return
	}

	if s.gametype == client {
		op.GeoM.Reset()
		op.GeoM.Translate(360, 150)
		if s.connection == connPending {
			text.Draw(screen, "Verbinde"+s.blink.Show("..."), face, op)
		} else if s.connection == connFinished {
			text.Draw(screen, "Verbunden", face, op)
		}

		return
	}

	if s.server == nil || !s.server.IsListining() {
		return
	}

	op.GeoM.Reset()
	op.GeoM.Translate(360, 150)

	clients := s.server.Clients()
	for i := 1; i <= s.playerCount; i++ {
		if len(clients) >= i {
			text.Draw(screen, fmt.Sprintf("Spieler %v : verbunden", i), face, op)
		} else {
			text.Draw(screen, fmt.Sprintf("Spieler %v : ", i)+s.blink.Show("..."), face, op)
		}

		op.GeoM.Translate(0, 40)
	}
}

func (s *MenuStart) connect() {
	connClient := func() {
		var err error
		s.client, err = engine.ConnectClient(s.ctx, s.serverAddr+":1200")
		if err != nil {
			s.connection = connClosed
			s.ctx, s.cancle = context.WithCancel(context.Background())
			return
		}
		s.client.PressKey('↵')
	}

	switch s.gametype {
	case client:
		s.connection = connPending
		go connClient()
	case server:
		s.connection = connPending
		s.server = engine.BuildServer(s.playerCount, ":1200")
		s.server.RunBackground(s.ctx)
		go connClient()
	case singleplayer:
		s.server = engine.BuildServer(1, ":1200")
		s.server.RunBackground(s.ctx)
		connClient()
	default:
		panic(fmt.Sprintf("unexpected scenes.gametype: %#v", s.gametype))
	}
}
