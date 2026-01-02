package tui

import dataframe "github.com/rocketlaunchr/dataframe-go"

type ScanResultMsg struct {
	Files []string
	Err   error
}

type DataFrameContentMsg struct {
	DataFrame *dataframe.DataFrame
	Err       error
}

type PreComputeDone struct {
	Cache map[int]ColumnsStat
}

type StatsComputed struct {
	Col   int
	Stats ColumnsStat
}

type UpdateContentMsg struct {
	DataFrame *dataframe.DataFrame
	Err       error
}

type OpenFileMsg struct {
	Path string
}

type BackToList struct{}

type ErrorMsg error

type FilterQuery struct {
	Query     string
	Dataframe *dataframe.DataFrame
}
