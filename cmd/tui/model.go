package tui

import (
    tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
}

func (m *Model) NewModel() {
}

func (m *Model) Init() tea.Cmd{
    return tea.Batch(tea.EnterAltScreen)
}
