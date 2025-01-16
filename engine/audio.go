package engine

import (
	"fmt"
	"log"

	soundAsset "github.com/apfelfrisch/gosnake/game/assets/sound"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
)

const sampleRate = 44100

var sounds = [4]sound{Eat, Dash, WalkWall, Crash}

type sound int

const (
	Eat sound = iota
	Dash
	Crash
	WalkWall
)

func (s sound) file() string {
	switch s {
	case Crash:
		return "crash.mp3"
	case Dash:
		return "dash.mp3"
	case Eat:
		return "eat.mp3"
	case WalkWall:
		return "walkwall.mp3"
	default:
		panic(fmt.Sprintf("unexpected main.sound: %#v", s))
	}
}

var audioContext *audio.Context

func init() {
	audioContext = audio.NewContext(sampleRate)
}

type player struct {
	sound map[sound]*audio.Player
	music *audio.Player
}

func NewPlayer() *player {
	player := player{
		sound: make(map[sound]*audio.Player),
	}

	for _, sound := range sounds {
		f, err := soundAsset.Files.Open(sound.file())
		if err != nil {
			log.Fatal(err)
		}
		d, err := mp3.DecodeF32(f)
		if err != nil {
			log.Fatal(err)
		}

		soundPlayer, err := audioContext.NewPlayerF32(d)

		if err != nil {
			log.Fatal(err)
		}

		player.sound[sound] = soundPlayer
	}

	f, err := soundAsset.Files.Open("theme-b.mp3")

	if err != nil {
		log.Fatal(err)
	}

	d, err := mp3.DecodeF32(f)
	if err != nil {
		log.Fatal(err)
	}

	loop := audio.NewInfiniteLoop(d, d.Length())

	musicPlayer, err := audioContext.NewPlayerF32(loop)

	if err != nil {
		log.Fatal(err)
	}

	player.music = musicPlayer

	return &player
}

func (p *player) Play(sound sound) {
	p.sound[sound].Rewind()
	p.sound[sound].Play()
}

func (p *player) PlayMusic() {
	p.music.Play()
}

func (p *player) PauseMusic() {
	p.music.Pause()
	p.music.Rewind()
}
