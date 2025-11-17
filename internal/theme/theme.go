package theme

import (
	"strings"

	"github.com/pterm/pterm"
)

var (
	AccentStart = pterm.NewRGB(255, 160, 130)
	AccentEnd   = pterm.NewRGB(230, 120, 255)

	PrimaryText   = pterm.NewRGB(255, 160, 130)
	SecondaryText = pterm.NewRGB(230, 120, 255)
	EmphasisText  = pterm.NewRGB(242, 140, 220)

	PrimaryStyle   = pterm.NewStyle(pterm.FgLightYellow)
	SecondaryStyle = pterm.NewStyle(pterm.FgLightMagenta)
	EmphasisStyle  = pterm.NewStyle(pterm.FgLightMagenta)
)

func GradientText(text string) string {
	return GradientTextWith(text, AccentStart, AccentEnd)
}

func GradientTextWith(text string, start, end pterm.RGB) string {
	var b strings.Builder
	count := len([]rune(text))
	if count == 0 {
		return ""
	}

	for i, r := range text {
		mix := float64(i) / float64(count)
		rVal := uint8(float64(start.R) + (float64(end.R)-float64(start.R))*mix)
		gVal := uint8(float64(start.G) + (float64(end.G)-float64(start.G))*mix)
		bVal := uint8(float64(start.B) + (float64(end.B)-float64(start.B))*mix)
		b.WriteString(pterm.NewRGB(rVal, gVal, bVal).Sprint(string(r)))
	}

	return b.String()
}

func ConfigurePrinters() {
	pterm.Error.Prefix = pterm.Prefix{Text: "ERROR", Style: SecondaryStyle}
	pterm.Error.MessageStyle = SecondaryStyle

	pterm.Warning.Prefix = pterm.Prefix{Text: "WARN", Style: PrimaryStyle}
	pterm.Warning.MessageStyle = PrimaryStyle

	pterm.Success.Prefix = pterm.Prefix{Text: "SUCCESS", Style: EmphasisStyle}
	pterm.Success.MessageStyle = EmphasisStyle

	pterm.Info.Prefix = pterm.Prefix{Text: "INFO", Style: SecondaryStyle}
	pterm.Info.MessageStyle = SecondaryStyle
}
