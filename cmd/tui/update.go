package tui

import tea "github.com/charmbracelet/bubbletea"

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    /**
        I have to add a case for when a dataframe is sent, when and a few cases 
        for when something happens(e.g. querying/filtering something)
    **/
        case error:
            m.err = msg
            return m, nil

        case tea.KeyMsg:
            switch msg.String() {
                case "q", "ctrl+c":
                    return m, tea.Quit
            }
    }

    return m, nil
}

func waitForDataFrames(m Model) tea.Cmd {
    return func() tea.Msg {
    }
}
