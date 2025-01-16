package main

import (
	"flag"
	"log"

	// _ "net/http/pprof"

	"github.com/apfelfrisch/gosnake/engine"
	"github.com/apfelfrisch/gosnake/game"
	netServer "github.com/apfelfrisch/gosnake/game/network/server"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	playerCount := flag.Int("player", 1, "Set Player count")
	serverAddr := flag.String("server-addr", ":1200", "Set Sever Address")
	onlyServer := flag.Bool("only-server", false, "Run only the server")
	onlyClient := flag.Bool("only-client", false, "Run only the server")

	flag.Parse()

	ebiten.SetWindowSize(engine.DisplayWidth, engine.DisplayHeight)
	ebiten.SetWindowTitle("Snake")

	if *onlyServer == true {
		buildServer(*playerCount, *serverAddr).Run()
	} else if *onlyClient == false {
		buildServer(*playerCount, *serverAddr).RunBackground()
	}

	if err := ebiten.RunGame(engine.New(*serverAddr, *playerCount)); err != nil {
		log.Fatal(err)
	}
}

func buildServer(playerCount int, addr string) *netServer.GameServer {
	return netServer.New(
		playerCount,
		addr,
		game.NewGame(playerCount, engine.GameWidth/engine.GridSize, engine.GameHeight/engine.GridSize),
	)
}
