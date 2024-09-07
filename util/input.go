package util

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input int

const (
	KeyNone Input = iota
	KeyLeft
	KeyRight
	KeyUp
	KeyDown
	KeyW
	KeyS
	KeyA
	KeyD
	KeyF
	KeyG
)

type InputHandler interface {
	Read() Input
	Cancel()
}

type KeyboardHandler struct {
	ch     chan Input
	cancel bool
}

func NewKBHandler() *KeyboardHandler {
	var ih KeyboardHandler
	ch := make(chan Input)

	go func(i KeyboardHandler) {
		for ih.Pooling() {
			if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
				ch <- KeyLeft
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
				ch <- KeyRight
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
				ch <- KeyUp
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
				ch <- KeyDown
			}

			if inpututil.IsKeyJustPressed(ebiten.KeyW) {
				ch <- KeyW
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyS) {
				ch <- KeyS
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyA) {
				ch <- KeyA
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyD) {
				ch <- KeyD
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyF) {
				ch <- KeyF
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyG) {
				ch <- KeyG
			}
			
		}
	}(ih)

	ih.ch = ch
	return &ih
}

func (ih *KeyboardHandler) Pooling() bool {
	return !ih.cancel
}

func (ih *KeyboardHandler) Cancel() {
	ih.cancel = true
}

func (ih *KeyboardHandler) Read() Input {
	var input Input
	select {
	case input = <-ih.ch:
	default:
		input = KeyNone
	}
	return input
}

type MockHandler struct {
	keys []Input
}

func NewMockHandler() *MockHandler {
	return &MockHandler{}
}

func (mh *MockHandler) AppendKeys(keys []Input) {
	mh.keys = append(mh.keys, keys...)
}

func (mh *MockHandler) Read() Input {
	if len(mh.keys) == 0 {
		return KeyNone
	}
	key := mh.keys[0]
	mh.keys = mh.keys[1:]
	return key
}

func (mh *MockHandler) Cancel() {

}
