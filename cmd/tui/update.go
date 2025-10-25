package tui

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/MikelGV/PrepPulse/internal/reader"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
        case tea.WindowSizeMsg:
            m.width = msg.Width
            m.height = msg.Height
            return m, nil

        case tea.KeyMsg:
            if m.err != nil && msg.String() == "q" || msg.String() == "ctrl+c" {
                return m, tea.Quit
            }
            switch m.state {
            case ListState:
                return m.updateList(msg)
            case FileState:
                return  m.updateFile(msg)
            }

        case ScanResultMsg:
            if msg.Err != nil {
                m.err = msg.Err
                return m, nil
            }
            
            m.files = msg.Files
            m.selected = clamp(m.selected, 0, len(m.files)-1)
            m.ready = true
            return m, nil

        case OpenFileMsg:
            m.state = FileState
            m.filePath = msg.Path
            m.scrollRow = 0
            m.scrollCol = 0
            ext := strings.ToLower(filepath.Ext(msg.Path))

            if ext == ".csv" || ext == ".tsv" || ext == ".json" {
                return m, loadDataCmd(msg.Path)
            }
            return m, nil 

        case DataFrameContentMsg:
            m.df = msg.DataFrame
            m.err = msg.Err
            return m, nil

        case BackToList:
            m.state = ListState
            m.filePath = ""
            m.fileContent = ""
            m.df = nil
            m.err = nil
            return m, scanCmd()
        
        case ErrorMsg:
            m.err = error(msg)
            return m, nil
    }

    return m, nil

}

func (m Model) updateList(msg tea.KeyMsg) (Model, tea.Cmd) {
    switch msg.String() {

    case "q", "ctrl+c":
        return m, tea.Quit

    case "up", "k":
        if m.selected > 0 {
            m.selected--
        }

    case "down", "j":
        if m.selected < len(m.files)-1 {
            m.selected++
        } 

    case "enter", "v":
        if len(m.files) > 0 {
            return m, func() tea.Msg {
                return OpenFileMsg{Path: m.files[m.selected]}
            }
        }

    case "d":
        if len(m.files) > 0 {
            path := m.files[m.selected]
            if err := os.Remove(path); err != nil {
                return m, func() tea.Msg { return ErrorMsg(err) }
            }
            return m, scanCmd()
        }
    }

    return m, nil
}


func (m Model) updateFile(msg tea.KeyMsg) (Model, tea.Cmd) {
    switch msg.String() {

    case "esc":
        return m, func() tea.Msg { return BackToList{} }

    case "q", "ctrl+c":
        return m, tea.Quit

    case "e":
        if m.df == nil && !m.editing {
            m.editing = true
            m.textLines = strings.Split(m.fileContent, "\n")
        }

    case "up", "k":
        if m.df != nil && m.scrollRow > 0 {
            m.scrollRow--
        }

    case "down", "j":
        if m.df != nil && m.scrollRow < m.df.Nrow()-1 {
            m.scrollRow++
        }
    case "left", "h":
        if m.df != nil && m.scrollCol > 0 {
            m.scrollCol--
        }

    case "right", "l":
        if m.df != nil && m.scrollCol < m.df.Ncol()-1 {
            m.scrollCol++
        }
    }

    return m, nil
}

func loadDataCmd(path string) tea.Cmd {
    return func() tea.Msg {

        df, err := reader.LoadData(path)
        return DataFrameContentMsg{DataFrame: df, Err: err} 
    }
}

/**func loadTextCmd() tea.Cmd {
}
**/

func scanCmd() tea.Cmd {
    return func() tea.Msg {

        files, err := reader.GetFilesFromCWD(".")
        return ScanResultMsg{Files: files, Err: err}
    }
}

func clamp(v, min, max int) int {

    if v < min {
        return min
    }

    if v > max {
        return max
    }

    return v
}
