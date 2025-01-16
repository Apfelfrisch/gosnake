package scenes

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"strconv"
	"time"

	netClient "github.com/apfelfrisch/gosnake/game/network/client"
	netServer "github.com/apfelfrisch/gosnake/game/network/server"

	"github.com/apfelfrisch/gosnake/engine"
	"github.com/apfelfrisch/gosnake/game"
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

type connection int

const (
	closed     connection = 0
	connecting connection = 1
	connected  connection = 2
)

type MenuStart struct {
	BaseScene
	gametype    gametype
	connection  connection
	playerCount int
	serverAddr  string
	blink       blink
	server      *netServer.GameServer
}

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

func (s *MenuStart) Update() error {
	s.blink.blink()

	if s.client != nil {
		s.client.UpdatePayload()
		s.localPlayer.Sync(s.client.Payload.Player)

		if s.client.Payload.GameState == game.Ongoing {
			s.sm.SwitchTo(&GameRunning{BaseScene: s.BaseScene})
		}
	} else if s.connection == connecting {
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		switch s.gametype {
		case client:
			s.connection = connecting
			go func() {
				s.client = buildClient(s.serverAddr + ":1200")
				s.client.PressKey('↵')
			}()
		case server:
			s.connection = connecting
			s.server = buildServer(s.playerCount, ":1200")
			s.server.RunBackground()
			go func() {
				s.client = buildClient(":1200")
				s.client.PressKey('↵')
			}()
		case singleplayer:
			s.server = buildServer(1, ":1200")
			s.server.RunBackground()
			s.client = buildClient(":1200")
			s.client.PressKey('↵')
		default:
			panic(fmt.Sprintf("unexpected scenes.gametype: %#v", s.gametype))
		}

	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		index := int(s.gametype) - 1
		if index < 0 {
			index = 2
		}
		s.gametype = gametype(index)
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		index := int(s.gametype) + 1
		if index > 2 {
			index = 0
		}
		s.gametype = gametype(index)
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
		if s.connection == connecting {
			text.Draw(screen, "Server Adresse: "+s.serverAddr, face, op)
		} else {
			text.Draw(screen, "Server Adresse: "+s.serverAddr+s.blink.Show("|"), face, op)
		}
	case server:
		op.GeoM.Translate(-40, 0)
		text.Draw(screen, "->", face, op)
		op.GeoM.Translate(350, -100)
		if s.connection == connecting {
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
		if s.connection == connecting {
			text.Draw(screen, "Verbinde"+s.blink.Show("..."), face, op)
		} else if s.connection == connected {
			text.Draw(screen, "Verbunden", face, op)
		}

		return
	}

	if s.server == nil {
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

func buildServer(playerCount int, addr string) *netServer.GameServer {
	return netServer.New(
		playerCount,
		addr,
		game.NewGame(playerCount, engine.GameWidth/engine.GridSize, engine.GameHeight/engine.GridSize),
	)
}

func buildClient(serverAddr string) *netClient.GameClient {
	client := netClient.Connect(serverAddr, engine.GameWidth/engine.GridSize, engine.GameHeight/engine.GridSize)
	player := engine.NewPlayer()

	client.EventBus.Add(netClient.PlayerHasEaten{}, func(event netClient.Event) {
		player.Play(engine.Eat)
	})
	client.EventBus.Add(netClient.PlayerDashed{}, func(event netClient.Event) {
		player.Play(engine.Dash)
	})
	client.EventBus.Add(netClient.PlayerWalkedWall{}, func(event netClient.Event) {
		player.Play(engine.WalkWall)
	})
	client.EventBus.Add(netClient.PlayerCrashed{}, func(event netClient.Event) {
		player.Play(engine.Crash)
	})
	client.EventBus.Add(netClient.GameHasStarted{}, func(event netClient.Event) {
		player.PlayMusic()
	})
	client.EventBus.Add(netClient.GameHasEnded{}, func(event netClient.Event) {
		player.PauseMusic()
	})

	return client
}
