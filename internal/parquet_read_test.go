package internal

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/reader"
)

func TestReadProcessedParquet(t *testing.T) {
	// Read all parquet files under data/processed recursively
	parquetPath := "data/processed"
	files, err := filepath.Glob(filepath.Join(parquetPath, "**/*.parquet"))
	if err != nil {
		t.Fatalf("Failed to read Parquet files: %v", err)
	}

	for _, file := range files {
		fr, err := local.NewLocalFileReader(file)
		if err != nil {
			t.Fatalf("Failed to open Parquet file: %v", err)
		}

		defer fr.Close()

		pr, err := reader.NewParquetReader(fr, new(ProcessedDataRecord), 4)
		if err != nil {
			t.Fatalf("Failed to create Parquet reader: %v", err)
		}
		defer pr.ReadStop()

		num := int(pr.GetNumRows())
		records := make([]ProcessedDataRecord, num)
		if err := pr.Read(&records); err != nil {
			t.Fatalf("Failed to read records: %v", err)
		}

		for i, rec := range records {
			fmt.Printf("Record %d: %+v\n", i, rec)
		}
	}
}
