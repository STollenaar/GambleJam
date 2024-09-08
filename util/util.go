package util

import (
	"image/color"
	"io"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

// Draw a rectangle with centered text at (x, y), with the rectangle size dependent on the text size
func DrawCenteredTextInRect(screen *ebiten.Image, x, y float32, rectColor color.Color, textColor color.Color, message string) *Button {

	// Measure the text size
	textWidth, textHeight := text.Measure(message, DefaultFont, 0)

	// Calculate rectangle dimensions based on text size and margin
	rectWidth := float32(textWidth) + 2*ConfigFile.Margin
	rectHeight := float32(textHeight) + 2*ConfigFile.Margin

	// Draw the rectangle
	vector.DrawFilledRect(screen, x, y, rectWidth, rectHeight, rectColor, false)

	// Calculate the position to center the text within the rectangle
	textX := x + ConfigFile.Margin
	textY := y + ConfigFile.Margin 

	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(textX), float64(textY))
	op.ColorScale.ScaleWithColor(textColor)
	// Draw the text at the calculated position (centered)
	text.Draw(screen, message, DefaultFont, op)
	return &Button{
		Name:   message,
		X:      x,
		Y:      y,
		Width:  rectWidth,
		Height: rectHeight,
	}
}

func DrawText(screen *ebiten.Image, x, y float64, textColor color.Color, message string) *Button {
	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(textColor)
	textWidth, textHeight := text.Measure(message, DefaultFont, 0)
	text.Draw(screen, message, DefaultFont, op)
	return &Button{
		Name:   message,
		X:      float32(x),
		Y:      float32(y),
		Width:  float32(textWidth),
		Height: float32(textHeight),
	}
}
