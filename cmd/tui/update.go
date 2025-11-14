package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MikelGV/PrepPulse/internal/reader"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-gota/gota/dataframe"
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

        case FilterQuery:
            query := strings.ToLower(msg.Query)
            var matches []int

            if m.filterInput == "" {
                m.filterActive = false
                m.filteredRows = nil
            } else {
                m.filterActive = true
                for i := 0; i < msg.Dataframe.Nrow(); i++ {
                    for j := 0; j < msg.Dataframe.Ncol(); j++ {
                        cell := fmt.Sprintf("%v", msg.Dataframe.Elem(i, j))
                        
                        if strings.Contains(strings.ToLower(cell), query) {
                            matches = append(matches, i)
                            break
                        }
                        
                    }
                }
                m.filteredRows = matches
            }

        case EditQuery:
            query := strings.ToLower(msg.Content)
        
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
    if m.filterMode {
        switch msg.Type {
        case tea.KeyRunes:
            for _, r := range msg.Runes {
                m.filterInput += string(r)
            }

            return m, filterCmd(m.df, m.filterInput)

        case tea.KeyBackspace, tea.KeyDelete:
            if len(m.filterInput) > 0 {
                m.filterInput = m.filterInput[:len(m.filterInput)-1]
                return m, filterCmd(m.df, m.filterInput)
            }

        case tea.KeyEnter:
            m.filterMode = false
            return m, nil

        case tea.KeyEsc:
            m.filterMode = false
            if m.filterInput == "" {
                m.filterActive = false
                m.filteredRows = nil
            }
            return m, nil
        }

    }

    if m.editMode {}

    switch msg.String() {

        case "esc":
            return m, func() tea.Msg { return BackToList{} }

        case "q", "ctrl+c":
            return m, tea.Quit

        case "e":
            if m.df == nil && !m.editMode {
                m.editMode = true
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

        case "/":
            if m.df != nil {
                m.filterMode = true
                m.filterInput = ""
                return m, nil
            }

        case "escape":
            if m.filterMode {
                m.filterMode = false

                if m.filterInput == "" {
                    m.filterActive = false
                    m.filteredRows = nil
                }
                return m, nil
            }

        return m, nil
    }

    return m, nil
}

func loadDataCmd(path string) tea.Cmd {
    return func() tea.Msg {

        df, err := reader.LoadData(path)
        return DataFrameContentMsg{DataFrame: df, Err: err} 
    }
}

func filterCmd(df *dataframe.DataFrame, query string) tea.Cmd {
    return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
        return FilterQuery{Query: query, Dataframe: df}
    })
}

func editCmd(df *dataframe.DataFrame, query string) tea.Cmd {
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
