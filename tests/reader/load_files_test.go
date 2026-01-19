package tests

import (
	"os"
	"path/filepath"
	"sort"
	"testing"

	"github.com/MikelGV/PrepPulse/internal/reader"
	"github.com/stretchr/testify/assert"
)

func TestLoadFiles(t *testing.T) {

	createFile := func(t *testing.T, dir, filename, content string) {
		t.Helper()
		path := filepath.Join(dir, filename)
		err := os.WriteFile(path, []byte(content), 0644)
		assert.NoError(t, err)
	}

	createDir := func(t *testing.T, dir string) {
		t.Helper()
		err := os.MkdirAll(dir, 0755)
		assert.NoError(t, err)
	}

	t.Run("Empty directory", func(t *testing.T) {
		tempDirp := t.TempDir()
		files, err := reader.GetFilesFromCWD(tempDirp)
		assert.NoError(t, err)
		assert.Empty(t, files)
	})

	t.Run("directory with files found", func(t *testing.T) {
		tempDir := t.TempDir()

		createFile(t, tempDir, "data1.csv", "name,age\njohn,29")
		createFile(t, tempDir, "data2.tsv", "name\age\njohn\t29")
		createFile(t, tempDir, "data3.json", `{ "key": "value"}`)

		file, err := reader.GetFilesFromCWD(tempDir)
		assert.NoError(t, err)

		sort.Strings(file)
		expected := []string{"data1.csv", "data2.tsv", "data3.json"}
		assert.Equal(t, expected, file)
	})

	t.Run("directory with mixed files found", func(t *testing.T) {
		tempDir := t.TempDir()

		createFile(t, tempDir, "data1.csv", "name,age\njohn,29")
		createFile(t, tempDir, "data2.txt", "some text")
		createFile(t, tempDir, "data3.json", `{ "key": "value"}`)

		file, err := reader.GetFilesFromCWD(tempDir)
		assert.NoError(t, err)

		sort.Strings(file)
		expected := []string{"data1.csv", "data3.json"}
		assert.Equal(t, expected, file)
	})

	t.Run("directory with subdirectories", func(t *testing.T) {
		tempDir := t.TempDir()

		subdir := filepath.Join(tempDir, "subdir")
		createDir(t, subdir)

		createFile(t, tempDir, "data1.csv", "name,age\njohn,29")
		createFile(t, tempDir, "data2.tsv", "name\age\njohn\t29")

		file, err := reader.GetFilesFromCWD(tempDir)
		assert.NoError(t, err)

		sort.Strings(file)
		expected := []string{"data1.csv", "data2.tsv"}
		assert.Equal(t, expected, file)
	})

	t.Run("non existent directory", func(t *testing.T) {
		nonExistent := filepath.Join(t.TempDir(), "does-not-exists")

		_, err := reader.GetFilesFromCWD(nonExistent)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
	})

	t.Run("case intensive extensidons", func(t *testing.T) {
		tempDir := t.TempDir()

		createFile(t, tempDir, "data1.Csv", "name,age\njohn,29")
		createFile(t, tempDir, "data2.tSv", "name\age\njohn\t29")
		createFile(t, tempDir, "data3.jsOn", `{ "key": "value"}`)

		file, err := reader.GetFilesFromCWD(tempDir)
		assert.NoError(t, err)

		sort.Strings(file)
		expected := []string{"data1.Csv", "data2.tSv", "data3.jsOn"}
		assert.Equal(t, expected, file)
	})

	t.Run("premission denied", func(t *testing.T) {
		tempDir := t.TempDir()

		restrictedDir := filepath.Join(tempDir, "restricted")
		createDir(t, restrictedDir)

		err := os.Chmod(restrictedDir, 0000)
		assert.NoError(t, err)

		defer os.Chmod(restrictedDir, 0755)

		_, err = reader.GetFilesFromCWD(restrictedDir)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission denied")
	})
}
