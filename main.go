package main

import (
	"fmt"
	"os"

	"github.com/go-gota/gota/dataframe"
)

func main() {

    file, err := os.Open("./testing_data/dog_breeds.csv")
    if err != nil {
        fmt.Printf("Error file name not found or not provided: %v", err)
    }

    df := dataframe.ReadCSV(file)


    fmt.Println(df)
}
