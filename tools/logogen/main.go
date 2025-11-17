package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"

	"github.com/Shu-AFK/bm/internal/theme"
)

var logoLines = []string{
	"██████╗ ███╗   ███╗",
	"██╔══██╗████╗ ████║",
	"██████╔╝██╔████╔██║",
	"██╔══██╗██║╚██╔╝██║",
	"██████╔╝██║ ╚═╝ ██║",
	"╚═════╝ ╚═╝     ╚═╝",
}

func main() {
	const (
		fontPath = "tools/logogen/JetBrainsMono-Regular.ttf"
		fontSize = 20
		marginX  = 10
		marginY  = 10
		bgR      = 8
		bgG      = 8
		bgB      = 10
	)

	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		log.Fatalf("failed to read font file: %v", err)
	}

	ft, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}

	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("failed to create font face: %v", err)
	}

	maxLen := 0
	for _, ln := range logoLines {
		runes := []rune(ln)
		if len(runes) > maxLen {
			maxLen = len(runes)
		}
	}

	charWidth := int(float64(fontSize) * 0.6)
	lineHeight := int(float64(fontSize) * 1.3)

	imgWidth := marginX*2 + maxLen*charWidth
	imgHeight := marginY*2 + len(logoLines)*lineHeight

	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))
	bg := color.RGBA{R: bgR, G: bgG, B: bgB, A: 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{C: bg}, image.Point{}, draw.Src)

	for lineIndex, line := range logoLines {
		y := marginY + lineIndex*lineHeight + int(float64(fontSize))

		runes := []rune(line)
		count := len(runes)
		for i, r := range runes {
			c := gradientColorAtIndex(i, count)

			d := &font.Drawer{
				Dst:  img,
				Src:  image.NewUniform(c),
				Face: face,
				Dot:  fixed.P(marginX+i*charWidth, y),
			}
			d.DrawString(string(r))
		}
	}

	out, err := os.Create("logo.png")
	if err != nil {
		log.Fatalf("failed to create logo.png: %v", err)
	}
	defer out.Close()

	if err := png.Encode(out, img); err != nil {
		log.Fatalf("failed to encode PNG: %v", err)
	}

	log.Println("Wrote logo.png successfully.")
}

// gradientColorAtIndex uses theme.AccentStart / AccentEnd to compute the color
// for the i-th character out of count.
func gradientColorAtIndex(i, count int) color.RGBA {
	if count <= 1 {
		return color.RGBA{
			R: theme.AccentStart.R,
			G: theme.AccentStart.G,
			B: theme.AccentStart.B,
			A: 255,
		}
	}

	// If you want the last char to hit AccentEnd exactly, use count-1 in the denominator.
	mix := float64(i) / float64(count-1)

	rVal := uint8(float64(theme.AccentStart.R) + (float64(theme.AccentEnd.R)-float64(theme.AccentStart.R))*mix)
	gVal := uint8(float64(theme.AccentStart.G) + (float64(theme.AccentEnd.G)-float64(theme.AccentStart.G))*mix)
	bVal := uint8(float64(theme.AccentStart.B) + (float64(theme.AccentEnd.B)-float64(theme.AccentStart.B))*mix)

	return color.RGBA{R: rVal, G: gVal, B: bVal, A: 255}
}
