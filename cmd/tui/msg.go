package tui

import "github.com/go-gota/gota/dataframe"

type ScanResultMsg struct {
    Files []string
    Err error
}

type DataFrameContentMsg struct {
    DataFrame *dataframe.DataFrame
    Err error
}

type OpenFileMsg struct {
    Path string
}

type BackToList struct {}

type ErrorMsg error 

type FilterQuery struct {
    Query string
    Dataframe *dataframe.DataFrame
}

