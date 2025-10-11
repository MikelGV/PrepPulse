package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-gota/gota/dataframe"
)

/**
    I need to be able to automatically detect the csv's tables.
    Also I need to detect the file type.
type File_Handler struct {
FilenName string
}
**/


func OpenFile() ( *os.File, error){
    file, err := os.Open("../testing_data/dog_breeds.csv")
    if err != nil {
        return nil, fmt.Errorf("Error file name not found or not provided: %v", err)
    }

    return file, nil 
}

func DetectFileType() (*os.File, error) {
    file, err := OpenFile()
    if err != nil {
        return nil, fmt.Errorf("Error opening file: %s", err)
    }
    reader := bufio.NewReader(file)

    content, err := reader.ReadString('\n')
    if err != nil {
        return nil, fmt.Errorf("Error reading file: %s", err)
    }

    if isJSON(content) {
        f := dataframe.ReadJSON(file)
        fmt.Println(f)
    } else if isCSV(content) {
        f := dataframe.ReadCSV(file)
        fmt.Println(f)
    } else if isTSV(content) {
        f := dataframe.ReadCSV(file)
        fmt.Println(f)
    } else {
        return nil, fmt.Errorf("file is in a wrong format")
    }
    return file, nil
}

func isJSON(content string) bool {
    var js json.RawMessage
    return json.Unmarshal([]byte(content), &js) == nil
}

func isTSV(content string) bool {
    return countOcurrence(content, '\t') > 0
}

func isCSV(content string) bool {
    return countOcurrence(content, ',') > 0
}

func countOcurrence(content string, char rune) int {
    count := 0
    for _, c := range content {
        if c == char {
            count++
        }
    }
    return count
}

