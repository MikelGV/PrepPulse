package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/go-gota/gota/dataframe"
)

type Model struct {
    df *dataframe.DataFrame
    err error
    state string
    load bool
    table *table.Table
}

func NewModel() Model {
    return Model{
        load: true,
        state: "idle",
        table: table.New(),
    }
}

func (m Model) Init() tea.Cmd{
    return tea.Batch(loadDataCmd(m), tea.EnterAltScreen)
}
