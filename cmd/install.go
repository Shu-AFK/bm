package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var InstallCmd = &cobra.Command{
	Use:   "install <shell> <show|apply>",
	Short: "prints or installs the shell wrapper in order for opening paths to work",
	RunE:  install,
}

func init() {
	rootCmd.AddCommand(InstallCmd)
}

const bmMarker = "# bm shell wrapper"

func install(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		pterm.Error.Println("usage: bm install <shell> <show|apply>")
		return errors.New("wrong usage")
	}

	shell := strings.ToLower(args[0])
	mode := strings.ToLower(args[1])

	if mode != "show" && mode != "apply" {
		pterm.Error.Println("second argument must be 'show' or 'apply'")
		return errors.New("second argument must be 'show' or 'apply'")
	}

	snippet, err := snippetForShell(shell)
	if err != nil {
		pterm.Error.Println(err.Error())
		return err
	}

	switch mode {
	case "show":
		fmt.Println(snippet)
		return nil

	case "apply":
		path, err := rcFileForShell(shell)
		if err != nil {
			pterm.Error.Println(err.Error())
			return err
		}

		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			pterm.Error.Printf("failed to create config directory %s: %v\n", filepath.Dir(path), err)
			return err
		}

		if alreadyInstalled(path) {
			pterm.Info.Printf("bm shell wrapper already installed for %s in %s\n", shell, path)
			return nil
		}

		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			pterm.Error.Printf("failed to open %s: %v\n", path, err)
			return err
		}
		defer f.Close()

		if _, err := fmt.Fprintln(f, "\n"+bmMarker); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(f, snippet); err != nil {
			return err
		}

		pterm.Success.Printf("installed bm shell wrapper for %s in %s\n", shell, path)
		pterm.Info.Println("restart your shell or source the file to start using `bm`")
		return nil

	default:
		return errors.New("unsupported mode")
	}
}

func snippetForShell(shell string) (string, error) {
	switch shell {
	case "bash", "zsh":
		return `bm() {
  local out status

  out=$(command bm "$@")
  status=$?

  printf '%s\n' "$out"

  if [[ $1 == open && $out == cd\ * ]]; then
    local path=${out#cd }
    cd "$path" || return
  fi

  return "$status"
}`, nil

	case "fish":
		return `function bm
    set out (command bm $argv)
    set status $status

    printf '%s\n' "$out"

    if test (count $argv) -ge 1; and test "$argv[1]" = "open"; and string match -q "cd *" "$out"
        set path (string sub -s 4 "$out")
        cd "$path"
    end

    return $status
end`, nil

	case "powershell", "pwsh":
		return `function bm {
    param(
        [Parameter(ValueFromRemainingArguments = $true)]
        [string[]]$Args
    )

    $bmCmd = Get-Command bm -CommandType Application,ExternalScript -ErrorAction SilentlyContinue
    if (-not $bmCmd) {
        Write-Error "bm executable not found in PATH"
        return
    }

    $out = & $bmCmd @Args
    $nativeExit = $LASTEXITCODE

    Write-Output $out

    if ($Args.Length -ge 1 -and $Args[0] -eq "open" -and $out -like "cd *") {
        $path = $out.Substring(3)
        Set-Location $path
    }

    if ($null -ne $nativeExit) {
        $global:LASTEXITCODE = $nativeExit
    }
}`, nil

	default:
		return "", fmt.Errorf("unsupported shell %q (supported: bash, zsh, fish, powershell)", shell)
	}
}

// rcFileForShell determines which file to modify for each shell.
// - bash: ~/.bashrc
// - zsh:  ~/.zshrc
// - fish: ~/.config/fish/functions/bm.fish (idiomatic function file)
// - pwsh: use $PROFILE.CurrentUserAllHosts (via pwsh/powershell)
func rcFileForShell(shell string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}

	switch shell {
	case "bash":
		return filepath.Join(home, ".bashrc"), nil
	case "zsh":
		return filepath.Join(home, ".zshrc"), nil
	case "fish":
		return filepath.Join(home, ".config", "fish", "functions", "bm.fish"), nil
	case "powershell", "pwsh":
		return detectPowerShellProfile()
	default:
		return "", fmt.Errorf("unsupported shell %q", shell)
	}
}

// alreadyInstalled checks if the marker is already present in the config file.
func alreadyInstalled(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	return strings.Contains(string(data), bmMarker)
}

// detectPowerShellProfile shells out to pwsh or powershell to ask for $PROFILE.CurrentUserAllHosts.
func detectPowerShellProfile() (string, error) {
	// Try pwsh first (modern PowerShell)
	if path, err := getProfileFromCommand("pwsh"); err == nil && path != "" {
		return path, nil
	}

	// Fallback to Windows PowerShell
	if path, err := getProfileFromCommand("powershell"); err == nil && path != "" {
		return path, nil
	}

	return "", errors.New("could not determine PowerShell profile path (tried pwsh and powershell)")
}

func getProfileFromCommand(cmdName string) (string, error) {
	cmd := exec.Command(cmdName, "-NoLogo", "-NoProfile", "-Command", "$PROFILE.CurrentUserAllHosts")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	profile := strings.TrimSpace(string(out))
	if profile == "" {
		return "", errors.New("empty profile path")
	}
	return profile, nil
}
