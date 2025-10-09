package reader

import (
	"encoding/csv"
	"fmt"
	"os"
)

/**
    I need to be able to automatically detect the csv's tables.
    Also I need to detect the file type.
**/

type File_Handler struct {
    FilenName string
}

func (f *File_Handler) OpenFile() ( *os.File, error){
    /**
        I need to add a way to detect if the file is supported or not, and if it's
        not a csv, tsv, or json throw an error and close the connection
    **/
    file, err := os.Open(f.FilenName)
    if err != nil {
        return nil, fmt.Errorf("Error file name not found or not provided: %v", err)
    }

    if err := detectFileType(file); err != nil {
        file.Close()
        return nil, fmt.Errorf("Error file type doesn't match our valid file types: %s", err)
    }


    return file, nil 
}

func detectFileType(file *os.File) error {
}

func isJSON() {}

func isTSVorCSV() {}

