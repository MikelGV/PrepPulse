package reader

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/rocketlaunchr/dataframe-go"
	"github.com/rocketlaunchr/dataframe-go/imports"
)

func LoadData(path string) (*dataframe.DataFrame, error) {
    file, err := os.Open(path)
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
        f, err:= imports.LoadFromJSON(context.Background(), file)
        if err != nil {
            return nil, fmt.Errorf("Error reading json: %s", err)
        }
        fmt.Println(f)
        return f, nil

    } else if isCSV(content) {
        f, err:= imports.LoadFromCSV(context.Background(), file)
        if  err != nil {
            return nil, fmt.Errorf("Error reading csv: %s", err)
        }
        fmt.Println(f)
        return f, nil
    } else if isTSV(content) {
        f, err := imports.LoadFromCSV(context.Background(), file, imports.CSVLoadOptions{Comma: '\t'})
        if err != nil {
            return nil, fmt.Errorf("Error reading tsv: %s", err)
        }
        fmt.Println(f)
        return f, nil
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

