module github.com/stollenaar/gamblingjam

go 1.23.0

require (
	github.com/hajimehoshi/ebiten/v2 v2.7.8
	github.com/hajimehoshi/go-mp3 v0.3.4
	github.com/hajimehoshi/oto v1.0.1
)

require (
	github.com/ebitengine/gomobile v0.0.0-20240518074828-e86332849895 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.7.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/exp v0.0.0-20190306152737-a1d7652674e8 // indirect
	golang.org/x/image v0.18.0 // indirect
	golang.org/x/mobile v0.0.0-20190415191353-3e0bab5405d6 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

replace (
	github.com/stollenaar/gamblingjam/internal => ./internal
	github.com/stollenaar/gamblingjam/util => ./util
)
