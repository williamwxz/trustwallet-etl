package main

import (
	"fmt"
	"os"
	"time"

	"trustwallet-etl/internal"
)

func main() {
	// Initialize logger
	logFile := fmt.Sprintf("logs/etl-%s.log", time.Now().Format("20060102-150405"))
	if err := internal.InitLogger(logFile); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	// Initialize metrics
	if err := internal.InitMetrics(); err != nil {
		internal.Fatal.Printf("Failed to initialize metrics: %v", err)
	}

	// Connect to Postgres
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := internal.NewPostgres(dsn)
	if err != nil {
		internal.Fatal.Printf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Main ETL loop
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Extract
			raw, err := internal.FetchRandomUser()
			if err != nil {
				internal.Error.Printf("Failed to fetch random user: %v", err)
				internal.IncrementFailures()
				continue
			}
			internal.IncrementExtractions()

			// Store raw data
			if err := internal.StoreRaw(db, raw); err != nil {
				internal.Error.Printf("Failed to store raw data: %v", err)
				internal.IncrementFailures()
				continue
			}

			// Transform
			processed, err := internal.Transform(raw)
			if err != nil {
				internal.Error.Printf("Failed to transform data: %v", err)
				internal.IncrementFailures()
				continue
			}
			internal.IncrementTransformations()

			// Store processed data
			if err := internal.StoreProcessed(db, processed); err != nil {
				internal.Error.Printf("Failed to store processed data: %v", err)
				internal.IncrementFailures()
				continue
			}
			internal.IncrementLoads()

			internal.Info.Printf("Successfully processed user: %s", processed.FullName)
		}
	}
}
