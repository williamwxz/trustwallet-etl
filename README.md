# TrustWallet ETL Service

A robust ETL (Extract, Transform, Load) service that fetches random user data, processes it, and stores it in both PostgreSQL and JSON files.

## Tech Stack

- Go 1.20
- PostgreSQL 15
- Prometheus metrics
- Docker & Docker Compose
- Alpine Linux

## Quickstart

1. Clone the repository:
```bash
git clone https://github.com/wzhang/trustwallet-etl.git
cd trustwallet-etl
```

2. Install Go dependencies:
```bash
go mod tidy
```

3. Start the services:
```bash
docker-compose up --build
```

4. Test parquet
```bash
go test -v ./internal -run TestReadProcessedParquet
```

The ETL service will start processing data every 30 seconds.

## Usage

### Logs
Logs are written to `logs/etl.log` and include:
- INFO: Successful operations
- ERROR: Failed operations
- FATAL: Critical errors

### Health & Metrics
- Health check: http://localhost:2112/health
- Prometheus metrics: http://localhost:2112/metrics

### Data Storage
- Raw data: `data/raw/raw_data.parquet` (Parquet format with Snappy compression)
  - Contains complete JSON response from randomuser.me API
  - Stored as UTF8 string with timestamp
- Processed data: `data/processed/processed_data.parquet` (Parquet format with Snappy compression)
  - Contains transformed user data with columns:
    - raw_id: Reference to raw data
    - full_name: Combined name (title + first + last)
    - email: User's email
    - gender: User's gender
    - registered_date: User registration timestamp
    - processed_at: Processing timestamp
    - created_at: Record creation timestamp
- PostgreSQL tables: `raw_data` and `processed_data`

## Project Structure

```
.
├── cmd/
│   └── etl/
│       └── main.go           # Application entry point
├── internal/
│   ├── extractor.go         # Data extraction from randomuser.me
│   ├── transformer.go       # Data transformation logic
│   ├── storage.go          # PostgreSQL and Parquet storage
│   ├── logger.go           # Logging utilities
│   ├── parquet_read_test.go           # parquet file test
│   └── metrics.go          # Prometheus metrics
├── data/
│   ├── raw/                # Raw Parquet data
│   └── processed/          # Processed Parquet data
├── logs/                   # Application logs
├── Dockerfile             # Multi-stage build
├── docker-compose.yml     # Service orchestration
└── README.md             # This file
```

## Productionization Tips

### Secrets Management
- Use environment variables or secrets management service
- Never commit sensitive data to version control
- Consider using AWS Secrets Manager or HashiCorp Vault

### Kubernetes/EKS Deployment
1. Create Kubernetes manifests:
   - Deployment
   - Service
   - ConfigMap
   - Secret
2. Use Helm charts for easier deployment
3. Configure resource limits and requests
4. Set up horizontal pod autoscaling

### Monitoring
1. Prometheus metrics are exposed at `/metrics`
2. Set up Grafana dashboards
3. Configure alerts for:
   - High failure rates
   - Processing delays
   - Storage capacity

### Storage

1. Consider using S3 for raw data:
   - Implement S3 upload in storage package
   - Use lifecycle policies for data retention

2. Use managed RDS for PostgreSQL:
   - Automated backups
   - Point-in-time recovery
   - Read replicas

3. Adopt Apache Iceberg for your data lake table format:
   - Store Parquet files in S3 (or other object storage) as Iceberg tables for scalable, cost-effective storage.
   - Partition data by date or other relevant columns (e.g., `date=YYYY-MM-DD/`) for efficient querying.
   - Use Iceberg-compatible engines (Spark, Trino, Flink, Presto, Athena) for analytics and batch/stream processing.
   - Manage table metadata with Iceberg's built-in catalog or integrate with AWS Glue/Hive Metastore.
   - Iceberg enables ACID transactions, schema evolution, time travel, and easy rollback to previous table versions.
   - Plan for data compaction, retention, and governance policies.

**References:**
- [Apache Iceberg](https://iceberg.apache.org/)
- [Iceberg on AWS](https://docs.aws.amazon.com/athena/latest/ug/querying-iceberg.html)
- [Iceberg with Spark](https://iceberg.apache.org/docs/latest/spark-quickstart/)

## Scaling & Reliability

### Parallelization
- Implement worker pools for parallel processing
- Use goroutines for concurrent operations
- Consider partitioning data for parallel processing

### Batching
- Implement batch processing for better performance
- Use bulk inserts for PostgreSQL
- Buffer writes to reduce I/O operations

### Retry & Backoff
- Implement exponential backoff for API calls
- Add retry logic for database operations
- Use circuit breakers for external services

### High Availability
1. Run multiple ETL instances
2. Use leader election for coordination
3. Implement idempotent operations
4. Use database transactions for consistency

### Streaming & Messaging

- Consider using Apache Kafka for scalable, decoupled data ingestion and processing.
- Use Kafka topics to buffer and distribute raw user data between extractors and ETL consumers.
- Benefits:
  - Handles spikes in data volume
  - Enables parallel and real-time processing
  - Supports replay and backfill scenarios
- Integrate with data lake tools (e.g., Spark Structured Streaming, Flink, Iceberg) for end-to-end streaming analytics. 