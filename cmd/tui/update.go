package tui

import (
	"time"

	"github.com/MikelGV/PrepPulse/cmd/tui/components"
	"github.com/MikelGV/PrepPulse/internal/reader"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-gota/gota/dataframe"
)

type LoadErrMsg error
type LoadSuccessMsg struct { df *dataframe.DataFrame }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
        case LoadSuccessMsg:
            m.df = msg.df
            m.err = nil
            m.table = components.CreateTable(m.df)
            return m, loadDataCmd(m) 

        case LoadErrMsg:
            m.err = msg
            return m, loadDataCmd(m) 

        case tea.KeyMsg:
            switch msg.String() {
                case "q", "ctrl+c":
                    return m, tea.Quit
            }
    }

    return m, nil
}

func loadDataCmd(m Model) tea.Cmd {
    return func() tea.Msg {
        df, err := reader.LoadData()
        time.Sleep(500 * time.Millisecond)

        if err != nil {
            return LoadErrMsg(err)
        }

        return LoadSuccessMsg{
            df: df,
        }
    }
}
