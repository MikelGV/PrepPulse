package main

import (
	"fmt"
	"os"

	"github.com/MikelGV/PrepPulse/cmd/tui"
	tea "github.com/charmbracelet/bubbletea"
)


func main() {
    p := tea.NewProgram(tui.NewModel())

    if _, err := p.Run(); err != nil {
        fmt.Printf("Uh, oh something when wrong: %v\n", err)
        os.Exit(1)
    }
}
