package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/Shu-AFK/bm/internal/theme"
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

var usageErrorPrinted bool

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if shouldPrintUsageError(err) && !usageErrorPrinted {
			pterm.Error.Println(err)
		}
		usageErrorPrinted = false
		os.Exit(1)
	}
	usageErrorPrinted = false
}

func shouldPrintUsageError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	usageIndicators := []string{
		"unknown command",
		"unknown flag",
		"unknown shorthand flag",
		"flag provided but not defined",
	}

	for _, indicator := range usageIndicators {
		if strings.Contains(msg, indicator) {
			return true
		}
	}

	return false
}

func init() {
	theme.ConfigurePrinters()

	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		if err == nil {
			return nil
		}

		usageErrorPrinted = true
		pterm.Error.Println(err)

		return err
	})

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

		renderSectionTitle("Overview")
		description := cmd.Long
		if description == "" {
			description = cmd.Short
		}
		pterm.Println(description)
		pterm.Println()

		renderSectionTitle("Usage")
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
			renderSectionTitle("Commands")
			renderCommandsTable(commands)
			pterm.Println("Tip: append --help to a command to see its dedicated flags.")
			pterm.Println()
		}

		if len(flags) > 0 {
			renderSectionTitle("Flags")
			renderFlagList(flags)
			pterm.Println()
		}
	})
}

func printBanner() {
	ascii := `
██████╗ ███╗   ███╗
██╔══██╗████╗ ████║
██████╔╝██╔████╔██║
██╔══██╗██║╚██╔╝██║
██████╔╝██║ ╚═╝ ██║
╚═════╝ ╚═╝     ╚═╝
`

	lines := strings.Split(ascii, "\n")
	count := len(lines)

	for i, line := range lines {
		r := uint8(int(theme.AccentStart.R) + (int(theme.AccentEnd.R)-int(theme.AccentStart.R))*i/count)
		g := uint8(int(theme.AccentStart.G) + (int(theme.AccentEnd.G)-int(theme.AccentStart.G))*i/count)
		b := uint8(int(theme.AccentStart.B) + (int(theme.AccentEnd.B)-int(theme.AccentStart.B))*i/count)
		color := pterm.NewRGB(r, g, b)

		pterm.DefaultCenter.Print(color.Sprint(line))
	}
}

func renderSectionTitle(title string) {
	pterm.Println(theme.GradientText(strings.ToUpper(title)))
	pterm.Println(theme.GradientText(strings.Repeat("─", len(title)+2)))
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
		data = append(data, []string{fmt.Sprintf("%s %s", theme.SecondaryText.Sprint("➜"), c.Name), c.Desc})
	}
	_ = pterm.DefaultTable.
		WithHasHeader().
		WithHeaderStyle(theme.PrimaryStyle).
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
			Text: fmt.Sprintf("%s  %s", theme.SecondaryText.Sprint(f.Name), f.Usage),
		})
	}
	_ = pterm.DefaultBulletList.WithItems(items).Render()
}
