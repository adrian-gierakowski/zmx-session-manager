package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/mdsakalu/zmx-session-manager/internal/tui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("zsm %s (%s, %s)\n", version, commit, date)
		return
	}

	zmxPath, err := exec.LookPath("zmx")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: zmx not found in PATH")
		os.Exit(1)
	}

	for {
		p := tea.NewProgram(tui.NewModel())
		finalModel, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// If the user pressed a key to attach, run zmx attach
		if m, ok := finalModel.(tui.Model); ok {
			target, replace := m.AttachTarget()
			if target != "" {
				if replace {
					env := os.Environ()
					syscall.Exec(zmxPath, []string{"zmx", "attach", target}, env)
					// If syscall.Exec fails, we might still be here
					fmt.Fprintf(os.Stderr, "Error: failed to exec into session %q\n", target)
					os.Exit(1)
				}

				cmd := exec.Command(zmxPath, "attach", target)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr

				if err := cmd.Run(); err != nil {
					fmt.Fprintf(os.Stderr, "Error: failed to attach to session %q: %v\n", target, err)
					// Wait for user to read the error before returning to zsm
					fmt.Print("Press Enter to return to zsm...")
					var dummy string
					fmt.Scanln(&dummy)
				}
				continue
			}
		}

		// Otherwise, the user quit the TUI normally
		break
	}
}
