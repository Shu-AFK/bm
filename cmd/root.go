package cmd

import (
	"bufio"
	"bytes"
	"os"
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
		printBanner()
		_ = cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		pterm.Error.Println("Error:", err)
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

		scanner := bufio.NewScanner(strings.NewReader(buf.String()))

		cyan := pterm.NewRGB(90, 170, 200)
		blue := pterm.NewRGB(100, 130, 180)
		purple := pterm.NewRGB(130, 90, 160)

		for scanner.Scan() {
			line := scanner.Text()

			switch {
			case strings.HasPrefix(line, "Usage:"):
				cyan.Print("Usage:")
				pterm.Println(strings.TrimPrefix(line, "Usage:"))

			case strings.HasPrefix(line, "Available Commands:"):
				cyan.Print("Available Commands:")
				pterm.Println(strings.TrimPrefix(line, "Available Commands:"))

			case strings.HasPrefix(line, "Flags:"):
				cyan.Print("Flags:")
				pterm.Println(strings.TrimPrefix(line, "Flags:"))

			case isCommandLine(line):
				parts := strings.Fields(line)
				name := parts[0]
				desc := strings.Join(parts[1:], " ")
				purple.Print(name)
				pterm.Println(" " + desc)

				for _, sub := range cmd.Commands() {
					if sub.Name() == name && sub.LocalFlags().HasAvailableFlags() {
						sub.LocalFlags().VisitAll(func(f *pflag.Flag) {
							blue.Print("    --" + f.Name)
							pterm.Println("   " + f.Usage)
						})
					}
				}

			case strings.TrimSpace(line) != "" && strings.HasPrefix(strings.TrimSpace(line), "-"):
				fields := strings.Fields(line)
				flagName := fields[0]
				rest := strings.Join(fields[1:], " ")
				blue.Print(flagName)
				if rest != "" {
					pterm.Println(" " + rest)
				} else {
					pterm.Println()
				}

			default:
				pterm.Println(line)
			}
		}
	})
}

func isCommandLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	parts := strings.Fields(trimmed)
	if len(parts) < 2 {
		return false
	}
	return true
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

	start := pterm.NewRGB(50, 200, 255)
	end := pterm.NewRGB(150, 50, 255)

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
