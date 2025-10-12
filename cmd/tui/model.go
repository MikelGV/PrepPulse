package tui

import (

	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-gota/gota/dataframe"
)

type Model struct {
    df *dataframe.DataFrame
    err error
}

func NewModel(df *dataframe.DataFrame) Model {
    //ctx, cancel := context.WithCancel(context.Background())
    return Model{
        df: df,
    }
}

func (m Model) Init() tea.Cmd{
    return tea.Batch(tea.EnterAltScreen)
}
