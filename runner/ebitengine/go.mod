module github.com/apfelfrisch/gosnake/runner/ebitengine

go 1.23.2

replace github.com/apfelfrisch/gosnake/game => ../../game

require (
	github.com/apfelfrisch/gosnake/game v0.0.0-00010101000000-000000000000
	github.com/hajimehoshi/ebiten/v2 v2.8.5
)

require (
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)
