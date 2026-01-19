package tui

import (
	"context"
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
		case HistogramInputState:
			return m.updateHistogramInput(msg)
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
		fileInfo, _ := os.Stat(m.filePath)
		if m.df.NRows() > 10000 || fileInfo.Size() > 1<<20 {
			return m, computeAllStatsCMD(m.df)
		}
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
	case PreComputeDone:
		m.chacheStats = msg.Cache

	case StatsComputed:
		m.chacheStats[msg.Col] = msg.Stats

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
		if m.df != nil && !m.statsMode {
			m.state = HistogramInputState
			m.histInputMode = true
			m.histInputBuffer = ""
			m.histSelectedSugg = 0
			m.histSuggestions = m.getColumnSuggestions("")
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

	case "[":
		if m.statsMode && m.histBucketCount > 3 {
			m.histBucketCount--
			return m, computeStatsCMDWithBuckets(m.df, m.statsColumns, m.histBucketCount)
		}

	case "]":
		if m.statsMode && m.histBucketCount < 30 {
			m.histBucketCount++
			return m, computeStatsCMDWithBuckets(m.df, m.statsColumns, m.histBucketCount)
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

func computeCollumnStats(df *dataframe.DataFrame, col int) ColumnsStat {
	stat := ColumnsStat{}

	series := df.Series[col]
	stat.Name = series.Name()
	stat.Type = series.Type()

	switch s := series.(type) {
	case *dataframe.SeriesFloat64:
		vals := s.Values
		nRows := s.NRows()

		missing, _ := s.NilCount()
		stat.Missing = missing

		uniqueMap := make(map[float64]bool)
		nonNilVals := make([]float64, 0, nRows)

		for i := 0; i < nRows; i++ {
			val := vals[i]
			uniqueMap[val] = true
			nonNilVals = append(nonNilVals, val)
		}
		stat.Unique = len(uniqueMap)

		if len(nonNilVals) > 0 {
			min, max := nonNilVals[0], nonNilVals[0]
			sum := 0.0
			for _, v := range nonNilVals {
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
				sum += v
			}
			stat.Min = min
			stat.Max = max

			mean, _ := s.Mean(nil)
			stat.Mean = mean

			if len(nonNilVals) > 0 {
				median := calculateMedian(nonNilVals)
				stat.Median = median
			}

			if min != max {
				buckets, bucketLabels := createHistogramBuckets(nonNilVals, 10)
				stat.Buckets = buckets
				stat.BucketLabels = bucketLabels
			}
		}

	case *dataframe.SeriesInt64:
		nRows := s.NRows()

		missing, _ := s.NilCount()
		stat.Missing = missing

		uniqueMap := make(map[int64]bool)
		nonNilVals := make([]float64, 0, nRows)

		for i := 0; i < nRows; i++ {
			val := s.Value(i)
			if val != nil {
				intVal := val.(int64)
				uniqueMap[intVal] = true
				nonNilVals = append(nonNilVals, float64(intVal))
			}
		}
		stat.Unique = len(uniqueMap)

		if len(nonNilVals) > 0 {
			min, max := nonNilVals[0], nonNilVals[0]
			for _, v := range nonNilVals {
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
			}
			stat.Min = min
			stat.Max = max

			sum := 0.0
			for _, v := range nonNilVals {
				sum += v
			}
			stat.Mean = sum / float64(len(nonNilVals))

			if len(nonNilVals) > 0 {
				median := calculateMedian(nonNilVals)
				stat.Median = median
			}

			if min != max {
				buckets, bucketLabels := createHistogramBuckets(nonNilVals, 10)
				stat.Buckets = buckets
				stat.BucketLabels = bucketLabels
			}
		}

	case *dataframe.SeriesString:
		nRows := s.NRows()
		missing, _ := s.NilCount()
		stat.Missing = missing

		uniqueMap := make(map[string]bool)
		for i := 0; i < nRows; i++ {
			val := s.ValueString(i)
			if val != "NaN" {
				uniqueMap[val] = true
			}
		}
		stat.Unique = len(uniqueMap)

	case *dataframe.SeriesMixed:
		nRows := s.NRows()
		missing, _ := s.NilCount()
		stat.Missing = missing

		uniqueMap := make(map[string]bool)
		nonNilVals := make([]float64, 0, nRows)

		for i := 0; i < nRows; i++ {
			val := s.Value(i)
			valStr := fmt.Sprintf("%v", val)
			if valStr != "NaN" && val != nil {
				uniqueMap[valStr] = true
				if f, ok := val.(float64); ok {
					nonNilVals = append(nonNilVals, f)
				} else if f, ok := val.(float32); ok {
					nonNilVals = append(nonNilVals, float64(f))
				} else if f, ok := val.(int); ok {
					nonNilVals = append(nonNilVals, float64(f))
				} else if f, ok := val.(int64); ok {
					nonNilVals = append(nonNilVals, float64(f))
				} else if f, ok := val.(int32); ok {
					nonNilVals = append(nonNilVals, float64(f))
				}
			}
		}
		stat.Unique = len(uniqueMap)

		if len(nonNilVals) > 0 {
			min, max := nonNilVals[0], nonNilVals[0]
			sum := 0.0
			for _, v := range nonNilVals {
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
				sum += v
			}
			stat.Min = min
			stat.Max = max
			stat.Mean = sum / float64(len(nonNilVals))

			if len(nonNilVals) > 0 {
				median := calculateMedian(nonNilVals)
				stat.Median = median
			}

			if min != max {
				buckets, bucketLabels := createHistogramBuckets(nonNilVals, 10)
				stat.Buckets = buckets
				stat.BucketLabels = bucketLabels
			}
		}
	}

	return stat
}

func computeCollumnStatsWithBuckets(df *dataframe.DataFrame, col int, bucketCount int) ColumnsStat {
	stat := ColumnsStat{}

	series := df.Series[col]
	stat.Name = series.Name()
	stat.Type = series.Type()

	switch s := series.(type) {
	case *dataframe.SeriesFloat64:
		vals := s.Values
		nRows := s.NRows()

		missing, _ := s.NilCount()
		stat.Missing = missing

		uniqueMap := make(map[float64]bool)
		nonNilVals := make([]float64, 0, nRows)

		for i := 0; i < nRows; i++ {
			val := vals[i]
			uniqueMap[val] = true
			nonNilVals = append(nonNilVals, val)
		}
		stat.Unique = len(uniqueMap)

		if len(nonNilVals) > 0 {
			min, max := nonNilVals[0], nonNilVals[0]
			sum := 0.0
			for _, v := range nonNilVals {
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
				sum += v
			}
			stat.Min = min
			stat.Max = max

			mean, _ := s.Mean(context.TODO())
			stat.Mean = mean

			if len(nonNilVals) > 0 {
				median := calculateMedian(nonNilVals)
				stat.Median = median
			}

			if min != max {
				buckets, bucketLabels := createHistogramBuckets(nonNilVals, bucketCount)
				stat.Buckets = buckets
				stat.BucketLabels = bucketLabels
			}
		}

	case *dataframe.SeriesInt64:
		nRows := s.NRows()

		missing, _ := s.NilCount()
		stat.Missing = missing

		uniqueMap := make(map[int64]bool)
		nonNilVals := make([]float64, 0, nRows)

		for i := 0; i < nRows; i++ {
			val := s.Value(i)
			if val != nil {
				intVal := val.(int64)
				uniqueMap[intVal] = true
				nonNilVals = append(nonNilVals, float64(intVal))
			}
		}
		stat.Unique = len(uniqueMap)

		if len(nonNilVals) > 0 {
			min, max := nonNilVals[0], nonNilVals[0]
			for _, v := range nonNilVals {
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
			}
			stat.Min = min
			stat.Max = max

			sum := 0.0
			for _, v := range nonNilVals {
				sum += v
			}
			stat.Mean = sum / float64(len(nonNilVals))

			if len(nonNilVals) > 0 {
				median := calculateMedian(nonNilVals)
				stat.Median = median
			}

			if min != max {
				buckets, bucketLabels := createHistogramBuckets(nonNilVals, bucketCount)
				stat.Buckets = buckets
				stat.BucketLabels = bucketLabels
			}
		}

	case *dataframe.SeriesString:
		nRows := s.NRows()
		missing, _ := s.NilCount()
		stat.Missing = missing

		uniqueMap := make(map[string]bool)
		for i := 0; i < nRows; i++ {
			val := s.ValueString(i)
			if val != "NaN" {
				uniqueMap[val] = true
			}
		}
		stat.Unique = len(uniqueMap)

	case *dataframe.SeriesMixed:
		nRows := s.NRows()
		missing, _ := s.NilCount()
		stat.Missing = missing

		uniqueMap := make(map[string]bool)
		nonNilVals := make([]float64, 0, nRows)

		for i := 0; i < nRows; i++ {
			val := s.Value(i)
			valStr := fmt.Sprintf("%v", val)
			if valStr != "NaN" && val != nil {
				uniqueMap[valStr] = true
				if f, ok := val.(float64); ok {
					nonNilVals = append(nonNilVals, f)
				} else if f, ok := val.(float32); ok {
					nonNilVals = append(nonNilVals, float64(f))
				} else if f, ok := val.(int); ok {
					nonNilVals = append(nonNilVals, float64(f))
				} else if f, ok := val.(int64); ok {
					nonNilVals = append(nonNilVals, float64(f))
				} else if f, ok := val.(int32); ok {
					nonNilVals = append(nonNilVals, float64(f))
				}
			}
		}
		stat.Unique = len(uniqueMap)

		if len(nonNilVals) > 0 {
			min, max := nonNilVals[0], nonNilVals[0]
			sum := 0.0
			for _, v := range nonNilVals {
				if v < min {
					min = v
				}
				if v > max {
					max = v
				}
				sum += v
			}
			stat.Min = min
			stat.Max = max
			stat.Mean = sum / float64(len(nonNilVals))

			if len(nonNilVals) > 0 {
				median := calculateMedian(nonNilVals)
				stat.Median = median
			}

			if min != max {
				buckets, bucketLabels := createHistogramBuckets(nonNilVals, bucketCount)
				stat.Buckets = buckets
				stat.BucketLabels = bucketLabels
			}
		}
	}

	return stat
}

func calculateMedian(vals []float64) float64 {
	n := len(vals)
	if n == 0 {
		return 0
	}

	sorted := make([]float64, n)
	copy(sorted, vals)

	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	if n%2 == 0 {
		return (sorted[n/2-1] + sorted[n/2]) / 2.0
	}
	return sorted[n/2]
}

func createHistogramBuckets(vals []float64, numBuckets int) ([]int, []string) {
	if len(vals) == 0 {
		return nil, nil
	}

	min, max := vals[0], vals[0]
	for _, v := range vals {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	if min == max {
		return []int{len(vals)}, []string{fmt.Sprintf("%.2f", min)}
	}

	buckets := make([]int, numBuckets)
	bucketLabels := make([]string, numBuckets)
	bucketWidth := (max - min) / float64(numBuckets)

	for i := 0; i < numBuckets; i++ {
		bucketMin := min + float64(i)*bucketWidth
		bucketMax := min + float64(i+1)*bucketWidth
		bucketLabels[i] = fmt.Sprintf("%.1f-%.1f", bucketMin, bucketMax)
	}

	for _, v := range vals {
		bucketIdx := int((v - min) / bucketWidth)
		if bucketIdx >= numBuckets {
			bucketIdx = numBuckets - 1
		}
		buckets[bucketIdx]++
	}

	return buckets, bucketLabels
}

func computeStatsCMDWithBuckets(df *dataframe.DataFrame, col int, bucketCount int) tea.Cmd {
	return func() tea.Msg {
		stats := computeCollumnStatsWithBuckets(df, col, bucketCount)
		return StatsComputed{Col: col, Stats: stats}
	}
}

func computeAllStatsCMD(df *dataframe.DataFrame) tea.Cmd {
	return func() tea.Msg {
		cache := make(map[int]ColumnsStat)

		for col := 0; col < len(df.Series); col++ {
			stats := computeCollumnStats(df, col)
			cache[col] = stats
		}

		return PreComputeDone{Cache: cache}

	}
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

func (m Model) getColumnSuggestions(input string) []string {
	if m.df == nil {
		return []string{}
	}

	var suggestions []string
	lowerInput := strings.ToLower(input)
	columns := m.df.Names()

	for _, col := range columns {
		if strings.Contains(strings.ToLower(col), lowerInput) {
			suggestions = append(suggestions, col)
		}
	}

	return suggestions
}

func (m Model) updateHistogramInput(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyRunes:
		for _, r := range msg.Runes {
			m.histInputBuffer += string(r)
		}
		m.histSuggestions = m.getColumnSuggestions(m.histInputBuffer)
		m.histSelectedSugg = 0
		if len(m.histSuggestions) > 0 {
			m.histSelectedSugg = 0
		}
		return m, nil

	case tea.KeyBackspace, tea.KeyDelete:
		if len(m.histInputBuffer) > 0 {
			m.histInputBuffer = m.histInputBuffer[:len(m.histInputBuffer)-1]
			m.histSuggestions = m.getColumnSuggestions(m.histInputBuffer)
			m.histSelectedSugg = 0
			if len(m.histSuggestions) > 0 {
				m.histSelectedSugg = 0
			}
		}
		return m, nil

	case tea.KeyUp, tea.KeyLeft:
		if m.histSelectedSugg > 0 {
			m.histSelectedSugg--
		}
		return m, nil

	case tea.KeyDown, tea.KeyRight:
		if m.histSelectedSugg < len(m.histSuggestions)-1 {
			m.histSelectedSugg++
		}
		return m, nil

	case tea.KeyEnter:
		if len(m.histSuggestions) > 0 {
			selectedCol := m.histSuggestions[m.histSelectedSugg]
			for i, col := range m.df.Names() {
				if col == selectedCol {
					m.statsColumns = i
					break
				}
			}

			m.histInputMode = false
			m.histInputBuffer = ""
			m.histSuggestions = nil
			m.histSelectedSugg = 0
			m.state = FileState
			m.statsMode = true
			return m, computeStatsCMDWithBuckets(m.df, m.statsColumns, m.histBucketCount)
		}
		return m, nil

	case tea.KeyEsc:
		m.histInputMode = false
		m.histInputBuffer = ""
		m.histSuggestions = nil
		m.histSelectedSugg = 0
		m.state = FileState
		return m, nil
	}

	return m, nil
}
