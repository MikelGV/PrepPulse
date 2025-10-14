package tui

import (

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-gota/gota/dataframe"
)

type Model struct {
    df *dataframe.DataFrame
    err error
    state string
    load bool
}

func NewModel() Model {
    return Model{
        load: true,
        state: "idle",
    }
}

func (m Model) Init() tea.Cmd{
    return tea.Batch(loadDataCmd, tea.EnterAltScreen)
}
