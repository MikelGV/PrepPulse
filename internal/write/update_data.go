package write

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rocketlaunchr/dataframe-go"
	"github.com/rocketlaunchr/dataframe-go/exports"
)

func Update_data(path string, df *dataframe.DataFrame) error {
	if df == nil {
		return nil
	}

	f, err := os.Create(path)

	if err != nil {
		return fmt.Errorf("there was an error trying to update the file: %w", err)
	}

	defer f.Close()
	ctx := context.Background()

	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		return exports.ExportToJSON(ctx, f, df)
	case ".tsv":
		return exports.ExportToCSV(ctx, f, df, exports.CSVExportOptions{Separator: '\t'})
	default:
		return exports.ExportToCSV(ctx, f, df)
	}

}
