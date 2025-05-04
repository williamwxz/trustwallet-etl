package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/xitongsys/parquet-go-source/local"
	"github.com/xitongsys/parquet-go/parquet"
	"github.com/xitongsys/parquet-go/writer"
)

// ProcessedDataRecord represents a single record for Parquet storage
type ProcessedDataRecord struct {
	UUID            string    `parquet:"name=uuid, type=BYTE_ARRAY, convertedtype=UTF8"`
	FullName        string    `parquet:"name=full_name, type=BYTE_ARRAY, convertedtype=UTF8"`
	Email           string    `parquet:"name=email, type=BYTE_ARRAY, convertedtype=UTF8"`
	Gender          string    `parquet:"name=gender, type=BYTE_ARRAY, convertedtype=UTF8"`
	RegisteredDate  time.Time `parquet:"name=registered_date, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
	Nationality     string    `parquet:"name=nationality, type=BYTE_ARRAY, convertedtype=UTF8"`
	LocationCity    string    `parquet:"name=location_city, type=BYTE_ARRAY, convertedtype=UTF8"`
	LocationCountry string    `parquet:"name=location_country, type=BYTE_ARRAY, convertedtype=UTF8"`
	ProcessedAt     time.Time `parquet:"name=processed_at, type=INT64, convertedtype=TIMESTAMP_MILLIS"`
}

func NewPostgres(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS raw_data (
			id SERIAL PRIMARY KEY,
			data JSONB NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func StoreRaw(db *sql.DB, raw *RandomUserResponse) error {
	// Store in JSON file
	rawData, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal raw data: %w", err)
	}

	rawPath := filepath.Join("data", "raw", fmt.Sprintf("raw_data_%s.json", time.Now().Format("20060102-150405")))
	if err := os.MkdirAll(filepath.Dir(rawPath), 0755); err != nil {
		return fmt.Errorf("failed to create raw data directory: %w", err)
	}

	if err := os.WriteFile(rawPath, rawData, 0644); err != nil {
		return fmt.Errorf("failed to write raw data file: %w", err)
	}

	// Store in database
	var rawID int
	err = db.QueryRow(
		"INSERT INTO raw_data (data) VALUES ($1) RETURNING id",
		rawData,
	).Scan(&rawID)
	if err != nil {
		return fmt.Errorf("failed to insert raw data: %w", err)
	}

	return nil
}

func StoreProcessed(db *sql.DB, processed *ProcessedUser) error {
	// Create Hive-style date partition (date=YYYY-MM-DD)
	now := time.Now().UTC()
	datePartition := fmt.Sprintf("date=%s", now.Format("2006-01-02"))
	processedPath := filepath.Join("data", "processed", datePartition)
	if err := os.MkdirAll(processedPath, 0755); err != nil {
		return fmt.Errorf("failed to create date partition: %w", err)
	}

	// Store in Parquet file with timestamp and UUID
	uniqueName := fmt.Sprintf("processed_data_%s_%s.parquet", now.Format("150405"), uuid.New().String())
	processedPath = filepath.Join(processedPath, uniqueName)

	// Create or open Parquet file
	fw, err := local.NewLocalFileWriter(processedPath)
	if err != nil {
		return fmt.Errorf("failed to create parquet file: %w", err)
	}
	defer fw.Close()

	// Create Parquet writer
	pw, err := writer.NewParquetWriter(fw, new(ProcessedDataRecord), 4)
	if err != nil {
		return fmt.Errorf("failed to create parquet writer: %w", err)
	}
	defer pw.WriteStop()

	// Set compression
	pw.CompressionType = parquet.CompressionCodec_SNAPPY

	// Get the latest raw_id
	var rawID int64
	err = db.QueryRow("SELECT id FROM raw_data ORDER BY id DESC LIMIT 1").Scan(&rawID)
	if err != nil {
		return fmt.Errorf("failed to get latest raw_id: %w", err)
	}

	// Write record
	record := ProcessedDataRecord{
		UUID:            processed.UUID,
		FullName:        processed.FullName,
		Email:           processed.Email,
		Gender:          processed.Gender,
		RegisteredDate:  processed.RegisteredDate,
		Nationality:     processed.Nationality,
		LocationCity:    processed.Location.City,
		LocationCountry: processed.Location.Country,
		ProcessedAt:     processed.ProcessedAt,
	}
	if err := pw.Write(record); err != nil {
		return fmt.Errorf("failed to write parquet record: %w", err)
	}

	return nil
}
