package util

import (
	"image/color"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
)

func DrawClearRectangle(screen *ebiten.Image, color color.RGBA, x, y, width, height, borderWidth float32) {
	// Draw top border
	vector.DrawFilledRect(screen, x, y, width, borderWidth, color, false)

	// Draw bottom border
	vector.DrawFilledRect(screen, x, y+height-borderWidth, width, borderWidth, color, false)

	// Draw left border
	vector.DrawFilledRect(screen, x, y, borderWidth, height, color, false)

	// Draw right border
	vector.DrawFilledRect(screen, x+width-borderWidth, y, borderWidth, height, color, false)
}

// Function to play an MP3 file using oto and go-mp3
func PlayMP3(path string) error {
	// Open the MP3 file
	mp3File, err := os.Open(path)
	if err != nil {
		return err
	}
	defer mp3File.Close()

	// Decode the MP3 file
	decoder, err := mp3.NewDecoder(mp3File)
	if err != nil {
		return err
	}

	// Create an audio context with a sample rate (usually 44100 Hz for MP3)
	context, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer context.Close()

	// Create a player for the audio context
	player := context.NewPlayer()
	defer player.Close()

	// Read the entire MP3 file into memory and play it
	data, err := io.ReadAll(decoder)
	if err != nil {
		return err
	}

	// Play the audio
	_, err = player.Write(data)
	if err != nil {
		return err
	}

	return nil
}