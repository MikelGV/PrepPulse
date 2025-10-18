package reader

import (
	"io/fs"
	"path/filepath"
	"strings"
)

func GetFilesFromCWD(root string) ([]string, error) {
    var matches []string
    err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            return  nil
        }


        ext := strings.ToLower(filepath.Ext(path))

        if ext == ".csv" || ext == ".tsv" ||  ext == ".json" {
            relPath, _ := filepath.Rel(root, path)
            matches = append(matches, relPath)
        }

        return nil
    })

    return matches, err 
}
