package reader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-gota/gota/dataframe"
)

func LoadData() (*dataframe.DataFrame, error) {
    file, err := os.Open("./testing_data/dog_breeds.csv")
    if err != nil {
        return nil, fmt.Errorf("Error file name not found or not provided: %v", err)
    }

    defer file.Close()

    reader := bufio.NewReader(file)

    content, err := reader.ReadString('\n')
    if err != nil {
        return nil, fmt.Errorf("Error reading file: %s", err)
    }

    _, err = file.Seek(0, 0)
    if err != nil {
        return nil, fmt.Errorf("Error reseting file pointer: %s", err)
    }

    if isJSON(content) {
        f := dataframe.ReadJSON(file)
        if err := f.Err; err != nil {
            return nil, fmt.Errorf("Error reading json: %s", err)
        }
        fmt.Println(f)
        return &f, nil
    } else if isCSV(content) {
        f := dataframe.ReadCSV(file)
        if err := f.Err; err != nil {
            return nil, fmt.Errorf("Error reading csv: %s", err)
        }
        fmt.Println(f)
        return &f, nil
    } else if isTSV(content) {
        f := dataframe.ReadCSV(file, dataframe.WithDelimiter('\t'))
        if err := f.Err; err != nil {
            return nil, fmt.Errorf("Error reading tsv: %s", err)
        }
        fmt.Println(f)
        return &f, nil
    } else {
        return nil, fmt.Errorf("file is in a wrong format")
    }
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

