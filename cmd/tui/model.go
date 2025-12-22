package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	dataframe "github.com/rocketlaunchr/dataframe-go"
)

type AppState int

const (
	ListState AppState = iota
	FileState
)

type UndoCommand struct {
	Row      int
	Col      int
	OldValue interface{}
	NewValue interface{}
}

const MaxUndo int = 50

type ColumnsStat struct {
	Name         string
	Type         string
	Min          float64
	Max          float64
	Mean         float64
	Median       float64
	Unique       int
	Missing      int
	Buckets      []int
	BucketLabels []string
}

type Model struct {
	df    *dataframe.DataFrame
	state AppState

	files       []string
	selected    int
	filePath    string
	fileContent string

	//Stats
	statsMode    bool
	statsColumns int
	chacheStats  map[int]ColumnsStat
	histCache    map[int]string

	scrollRow int
	scrollCol int
	cursorRow int
	cursorCol int

	filterMode   bool
	filterInput  string
	filterActive bool
	filteredRows []int

	editMode     bool
	editBuffer   string
	editOriginal string

	redoStack []UndoCommand
	undoStack []UndoCommand
	maxUndo   int

	width  int
	height int
	ready  bool
	err    error
}

func NewModel() Model {
	return Model{
		state:    ListState,
		selected: 0,
		maxUndo:  MaxUndo,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(scanCmd(), tea.EnterAltScreen)
}
