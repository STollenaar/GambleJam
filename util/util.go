package util

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"strings"

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
		fmt.Println(err)
		return err
	}
	defer mp3File.Close()

	// Decode the MP3 file
	decoder, err := mp3.NewDecoder(mp3File)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Create an audio context with a sample rate (usually 44100 Hz for MP3)
	context, err := oto.NewContext(decoder.SampleRate(), 2, 2, 8192)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer context.Close()

	// Create a player for the audio context
	player := context.NewPlayer()
	defer player.Close()

	// Read the entire MP3 file into memory and play it
	data, err := io.ReadAll(decoder)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Play the audio
	_, err = player.Write(data)
	if err != nil {
		fmt.Println(err)
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

func DrawText(screen *ebiten.Image, x, y float64, textColor color.Color, message string, font *text.GoTextFace) *Button {
	if font == nil {
		font = DefaultFont
	}
	var textWidth, textHeight float32
	for _, line := range strings.Split(message, "\n") {
		op := &text.DrawOptions{}
		op.GeoM.Translate(x, y)
		op.ColorScale.ScaleWithColor(textColor)
		tw, th := text.Measure(line, DefaultFont, 0)
		text.Draw(screen, line, font, op)
		if float32(tw) > textWidth {
			textWidth = float32(tw)
		}
		if float32(th) > textHeight {
			textHeight = float32(th)
		}
		y += float64(textHeight)
	}
	return &Button{
		Name:   message,
		X:      float32(x),
		Y:      float32(y),
		Width:  textWidth,
		Height: textHeight,
	}
}

func SelectActiveColor(b bool) color.RGBA {
	if b {
		return activeColor
	}
	return inactiveColor
}

// measureText measures the width and height of a line of text with a specific font face
func measureText(line string, face *text.GoTextFace) (width, height float64) {
	return text.Measure(line, face, 0)
}

// wrapText splits the message into lines that fit within the specified width
func wrapText(message string, face *text.GoTextFace, maxWidth float64) []string {
	words := strings.Fields(message) // Split message into words
	var lines []string
	var line string

	for _, word := range words {
		testLine := line + word + " "
		lineWidth, _ := measureText(testLine, face)

		if lineWidth > maxWidth && line != "" {
			// Add the current line to the lines slice and start a new line
			lines = append(lines, strings.TrimSpace(line))
			line = word + " "
		} else {
			// Add the word to the current line
			line = testLine
		}
	}

	// Append the final line
	if line != "" {
		lines = append(lines, strings.TrimSpace(line))
	}

	return lines
}

// adjustFontSize dynamically adjusts the font size to fit within the rectangle
func adjustFontSize(message string, font *text.GoTextFace, maxWidth, maxHeight, minSize float64) *text.GoTextFace {
	for size := font.Size; size >= minSize; size -= 1 {
		font = &text.GoTextFace{
			Source: font.Source,
			Size:   size,
		}
		lines := wrapText(message, font, maxWidth)

		// Measure total height of the wrapped text
		_, lineHeight := measureText(message, font) // Get height of one line of text
		totalLineHeight := lineHeight * float64(len(lines))
		if totalLineHeight <= maxHeight {
			return font
		}
	}
	return DefaultFont // Default font if no suitable size is found
}

// DrawTextInRect wraps and draws text within a rectangle, scaling font size to fit
func DrawTextInRect(screen *ebiten.Image, message string, x, y, width, height float64, textColor color.Color, font *text.GoTextFace) {
	if font == nil {
		font = DefaultFont
	}

	// Dynamically adjust font size to fit within the rectangle
	font = adjustFontSize(message, font, width, height, 5)
	// Wrap the text into lines that fit the rectangle width
	lines := wrapText(message, font, width)
	// Draw each line inside the rectangle
	drawY := y
	_, lineHeight := measureText(strings.Join(lines, "\n"), font)
	for _, line := range lines {
		lineWidth, _ := measureText(line, font)
		// Center each line horizontally
		op := &text.DrawOptions{}
		textX := x + (width-float64(lineWidth))/2
		op.GeoM.Translate(textX, drawY)
		op.ColorScale.ScaleWithColor(textColor)
		// Draw the line
		text.Draw(screen, line, font, op)
		// Move the Y position down for the next line
		drawY += float64(lineHeight)
	}
}
