package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var rootCmd = &cobra.Command{
	Use:           "bm",
	Short:         "CLI Bookmark Manager: save and open files, folders, and links",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	defaultHelp := rootCmd.HelpFunc()

	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		buf := new(bytes.Buffer)
		origOut := cmd.OutOrStdout()
		cmd.SetOut(buf)
		defaultHelp(cmd, args)
		cmd.SetOut(origOut)

		printBanner()
		pterm.Println()

		help := buf.String()
		usage := strings.TrimSpace(extractUsage(help))
		if usage == "" {
			usage = strings.TrimSpace(cmd.UseLine())
		}
		if usage == "" {
			usage = cmd.Use
		}
		flags := extractFlags(cmd)
		commands := extractCommands(cmd)

		accentStart := pterm.NewRGB(255, 160, 130)
		accentEnd := pterm.NewRGB(230, 120, 255)
		errorAccent := pterm.NewStyle(pterm.FgLightMagenta)
		pterm.Error.Prefix = pterm.Prefix{Text: "ERROR", Style: errorAccent}
		pterm.Error.MessageStyle = errorAccent

		renderSectionTitle("Overview", accentStart, accentEnd)
		description := cmd.Long
		if description == "" {
			description = cmd.Short
		}
		pterm.Println(description)
		pterm.Println()

		renderSectionTitle("Usage", accentStart, accentEnd)
		title := "How to run it"
		usageLen := len([]rune(usage))
		titleLen := len([]rune(title))
		if usageLen < titleLen {
			usage += strings.Repeat(" ", titleLen-usageLen)
		}
		pterm.DefaultBox.
			WithLeftPadding(6).
			WithRightPadding(6).
			WithTopPadding(0).
			WithBottomPadding(0).
			WithTitle(title).
			WithTitleBottomRight().
			Println(usage)
		pterm.Println()

		if len(commands) > 0 {
			renderSectionTitle("Commands", accentStart, accentEnd)
			renderCommandsTable(commands)
			pterm.Println("Tip: append --help to a command to see its dedicated flags.")
			pterm.Println()
		}

		if len(flags) > 0 {
			renderSectionTitle("Flags", accentStart, accentEnd)
			renderFlagList(flags)
			pterm.Println()
		}
	})
}

func printBanner() {
	ascii := `
 ____    __  __
| __ )  |  \/  |
|  _ \  | |\/| |
| |_) | | |  | |
|____/  |_|  |_|
bookmark manager
`

	start := pterm.NewRGB(255, 170, 140)
	end := pterm.NewRGB(230, 120, 255)

	lines := strings.Split(ascii, "\n")
	count := len(lines)

	for i, line := range lines {
		r := uint8(int(start.R) + (int(end.R)-int(start.R))*i/count)
		g := uint8(int(start.G) + (int(end.G)-int(start.G))*i/count)
		b := uint8(int(start.B) + (int(end.B)-int(start.B))*i/count)
		color := pterm.NewRGB(r, g, b)

		pterm.DefaultCenter.Print(color.Sprint(line))
	}
}

func renderSectionTitle(title string, start, end pterm.RGB) {
	pterm.Println(gradientText(strings.ToUpper(title), start, end))
	pterm.Println(gradientText(strings.Repeat("─", len(title)+2), start, end))
}

func gradientText(text string, start, end pterm.RGB) string {
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

func extractUsage(help string) string {
	scanner := bufio.NewScanner(strings.NewReader(help))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Usage:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "Usage:"))
		}
	}
	return ""
}

type commandInfo struct {
	Name string
	Desc string
}

func extractCommands(cmd *cobra.Command) []commandInfo {
	var commands []commandInfo
	for _, c := range cmd.Commands() {
		if c.Hidden {
			continue
		}
		commands = append(commands, commandInfo{
			Name: c.Name(),
			Desc: c.Short,
		})
	}

	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Name < commands[j].Name
	})

	return commands
}

func renderCommandsTable(commands []commandInfo) {
	data := pterm.TableData{{"Command", "What it does"}}
	for _, c := range commands {
		data = append(data, []string{fmt.Sprintf("%s %s", pterm.LightMagenta("➜"), c.Name), c.Desc})
	}
	_ = pterm.DefaultTable.
		WithHasHeader().
		WithSeparator("│").
		WithHeaderRowSeparator("─").
		WithData(data).
		WithBoxed(true).
		Render()
}

type flagInfo struct {
	Name  string
	Usage string
}

func extractFlags(cmd *cobra.Command) []flagInfo {
	var flags []flagInfo

	visitor := func(f *pflag.Flag) {
		if f.Hidden {
			return
		}
		label := "--" + f.Name
		if f.Shorthand != "" {
			label = fmt.Sprintf("-%s, %s", f.Shorthand, label)
		}
		flags = append(flags, flagInfo{Name: label, Usage: f.Usage})
	}

	cmd.InheritedFlags().VisitAll(visitor)
	cmd.NonInheritedFlags().VisitAll(visitor)

	sort.Slice(flags, func(i, j int) bool {
		return flags[i].Name < flags[j].Name
	})

	return flags
}

func renderFlagList(flags []flagInfo) {
	var items []pterm.BulletListItem
	for _, f := range flags {
		items = append(items, pterm.BulletListItem{
			Text: fmt.Sprintf("%s  %s", pterm.LightMagenta(f.Name), f.Usage),
		})
	}
	_ = pterm.DefaultBulletList.WithItems(items).Render()
}
