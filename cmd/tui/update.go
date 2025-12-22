package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/MikelGV/PrepPulse/internal/reader"
	"github.com/MikelGV/PrepPulse/internal/write"
	tea "github.com/charmbracelet/bubbletea"
	dataframe "github.com/rocketlaunchr/dataframe-go"
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
			return m.updateFile(msg)
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
			for i := 0; i < msg.Dataframe.NRows(); i++ {
				row := m.df.Row(i, false, dataframe.SeriesIdx)
				for j := 0; j < len(msg.Dataframe.Series); j++ {
					val := row[j]
					cell := fmt.Sprintf("%v", val)

					if strings.Contains(strings.ToLower(cell), query) {
						matches = append(matches, i)
						break
					}

				}
			}
			m.filteredRows = matches
		}

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

	if m.editMode {
		switch msg.Type {

		case tea.KeyRunes:
			for _, r := range msg.Runes {
				m.editBuffer += string(r)
			}
			return m, nil

		case tea.KeyBackspace:
			if len(m.editBuffer) > 0 {
				m.editBuffer = m.editBuffer[:len(m.editBuffer)-1]
			}
			return m, nil

		case tea.KeyEnter:
			if m.editMode {
				m.editMode = false

				oldVal := m.df.Series[m.cursorCol].Value(m.cursorRow)

				if m.editOriginal != m.editBuffer {
					colName := m.df.Names()[m.cursorCol]

					m.df.Update(m.cursorRow, colName, m.editBuffer)

					cmd := UndoCommand{
						Row:      m.cursorRow,
						Col:      m.cursorCol,
						OldValue: oldVal,
						NewValue: m.editBuffer,
					}

					m.undoStack = append(m.undoStack, cmd)
					if len(m.undoStack) > m.maxUndo {
						m.undoStack = m.undoStack[1:]
					}

					m.redoStack = nil

					return m, updateDataCmd(m.filePath, m.df)
				}

				return m, nil
			}

		case tea.KeyEsc:
			m.editBuffer = m.editOriginal
			m.editMode = false
			return m, nil

		}
	}

	switch msg.String() {

	case "esc":
		return m, func() tea.Msg { return BackToList{} }

	case "q", "ctrl+c":
		return m, tea.Quit

	case "e":
		if m.df != nil && !m.editMode {
			m.editMode = true
			val := m.df.Row(m.cursorRow, false, dataframe.SeriesIdx)[m.cursorCol]
			m.editBuffer = fmt.Sprintf("%v", val)
			m.editOriginal = m.editBuffer
			return m, nil
		}

	case "up", "k":
		if m.df != nil && m.scrollRow > 0 && m.cursorRow > 0 {
			m.scrollRow--
			m.cursorRow--
		}

	case "down", "j":
		if m.df != nil && m.scrollRow < m.df.NRows()-1 && m.cursorRow < m.df.NRows()-1 {
			m.scrollRow++
			m.cursorRow++
		}
	case "left", "h":
		if m.df != nil && m.scrollCol > 0 && m.cursorCol > 0 {
			m.scrollCol--
			m.cursorCol--
		}

	case "right", "l":
		if m.df != nil && m.scrollCol < len(m.df.Series)-1 && m.cursorCol < len(m.df.Series)-1 {
			m.scrollCol++
			m.cursorCol++
		}

	case "/":
		if m.df != nil {
			m.filterMode = true
			m.filterInput = ""
			return m, nil
		}

	case "u", "ctrl+z":
		if len(m.undoStack) > 0 {
			cmd := m.undoStack[len(m.undoStack)-1]
			m.undoStack = m.undoStack[:len(m.undoStack)-1]

			colName := m.df.Names()[cmd.Col]
			m.df.Update(cmd.Row, colName, cmd.OldValue)

			m.redoStack = append(m.redoStack, cmd)

			return m, nil
		}

	case "ctrl+r":
		if len(m.redoStack) > 0 {
			cmd := m.redoStack[len(m.redoStack)-1]
			m.redoStack = m.redoStack[:len(m.redoStack)-1]

			colName := m.df.Names()[cmd.Col]
			m.df.Update(cmd.Row, colName, cmd.NewValue)

			m.undoStack = append(m.undoStack, cmd)

			return m, nil
		}

	case "s":
		if m.df != nil {
			m.statsMode = !m.statsMode
			if m.statsMode {
				m.statsColumns = m.cursorCol

				if _, ok := m.chacheStats[m.statsColumns]; !ok {
					return m, computeStatsCMD(m.df, m.statsColumns)
				}
			}
		}

		return m, nil

	case "escape":
		if m.filterMode {
			m.filterMode = false

			if m.filterInput == "" {
				m.filterActive = false
				m.filteredRows = nil
			}
			return m, nil
		}

		if m.editMode {
			m.editMode = false
			if m.editBuffer == "" {
				m.editBuffer = m.editOriginal
			}
			return m, nil
		}

		if m.statsMode {
			m.statsMode = false
			m.statsColumns = 0

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
func updateDataCmd(path string, df *dataframe.DataFrame) tea.Cmd {
	return func() tea.Msg {

		err := write.Update_data(path, df)

		if err != nil {
			return ErrorMsg(err)
		}

		return UpdateContentMsg{DataFrame: df, Err: err}
	}
}

func filterCmd(df *dataframe.DataFrame, query string) tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return FilterQuery{Query: query, Dataframe: df}
	})
}

func computeStatsCMD(df *dataframe.DataFrame, statsColumns int) tea.Cmd {
	/**
		Here i handle the compute of the stats
	**/
	return nil
}

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
