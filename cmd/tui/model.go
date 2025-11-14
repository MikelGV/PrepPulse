package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-gota/gota/dataframe"
)

type AppState int

const (
    ListState AppState = iota
    FileState
)

type Model struct {
    df *dataframe.DataFrame
    state AppState 
    files []string
    selected int
    filePath string
    fileContent string
    scrollRow int
    scrollCol int
    cursorRow int
    cursorCol int
    filterMode bool
    filterInput string 
    filterActive bool 
    filteredRows []int
    editMode bool
    editBuffer string
    editContent string
    width int
    height int
    ready bool
    textLines []string
    err error
}

func NewModel() Model {
    return Model{
        state: ListState,
        selected: 0,
    }
}

func (m Model) Init() tea.Cmd{
    return tea.Batch( scanCmd(), tea.EnterAltScreen)
}
