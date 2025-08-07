# Technical Assessment: Event Processing Pipeline

**Duration:** 4 hours maximum  
**Target Role:** Software Engineer (Golang, Mid-level+)  
**Assessment Type:** Take-home coding challenge

---

## Challenge Overview

Create a **concurrent event processing system** in Go that demonstrates your expertise with:
- Go's concurrency patterns (goroutines, channels, context)
- System design and architecture principles
- Production-ready backend development practices

### Business Context
Build a system that processes real-time events (user actions, sensor data, system logs) with requirements for:
- High throughput processing (1000+ events/second)
- Concurrent event processing using goroutines and channels
- Real-time metrics and monitoring
- Graceful failure handling and recovery

---

## Technical Requirements

### Required API Endpoints
```
POST /events        - Accept single event
POST /events/batch  - Accept multiple events (up to 100)
GET  /metrics       - Show processing statistics
GET  /health        - Health check endpoint
```

### Event Data Model
```json
{
  "id": "uuid (auto-generated if not provided)",
  "type": "user_action|sensor_data|system_log",
  "source": "web|mobile|iot_device|system",
  "timestamp": "2025-01-15T10:30:00Z",
  "user_id": "string (optional)",
  "data": {
    "action": "login|purchase|click|error",
    "value": 100.50,
    "metadata": {}
  }
}
```

### Core Go Requirements
You **must demonstrate** these Go concepts:

#### Worker Pool Pattern
```go
type EventPipeline struct {
    ingestionChan   chan Event
    workerPool      []*Worker  
    storage         Storage
    metrics         *Metrics
    ctx             context.Context
}

type Worker struct {
    id       int
    jobChan  chan Event
    pipeline *EventPipeline
}
```

#### Required Interfaces
```go
type Processor interface {
    Process(ctx context.Context, event Event) (*ProcessedEvent, error)
}

type Storage interface {
    Store(ctx context.Context, events []ProcessedEvent) error
}

type Validator interface {
    Validate(ctx context.Context, event Event) error
}
```

#### Context Usage for Graceful Shutdown
```go
func (w *Worker) Start(ctx context.Context) {
    go func() {
        for {
            select {
            case job := <-w.jobChan:
                w.processJob(ctx, job)
            case <-ctx.Done():
                return
            }
        }
    }()
}
```

---

## Development Setup

### Prerequisites
```bash
# Go 1.21+
go version

# Docker & Docker Compose
docker --version
docker-compose --version
```

### Infrastructure Setup (Provided)

**Create this `docker-compose.yml` file in your project directory:**

```yaml
services:
  mysql:
    image: mysql:8.0
    container_name: event-pipeline-mysql
    environment:
      MYSQL_ROOT_PASSWORD: testpass
      MYSQL_DATABASE: eventdb
    ports:
      - "3306:3306"
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      timeout: 20s
      retries: 10

  mysql-init:
    image: mysql:8.0
    depends_on:
      mysql:
        condition: service_healthy
    command: |
      mysql -h mysql -uroot -ptestpass eventdb -e "
      CREATE TABLE IF NOT EXISTS processed_events (
        id VARCHAR(36) PRIMARY KEY,
        type ENUM('user_action', 'sensor_data', 'system_log') NOT NULL,
        source VARCHAR(50) NOT NULL,
        user_id VARCHAR(50),
        processed_data JSON,
        processing_time_ms INT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        
        INDEX idx_type_created (type, created_at),
        INDEX idx_user_created (user_id, created_at)
      );"
    restart: "no"
```

### Quick Setup (2 minutes)
```bash
# 1. Start infrastructure
docker-compose up -d

# 2. Verify setup (should show 'processed_events' table)
docker-compose exec mysql mysql -uroot -ptestpass eventdb -e "SHOW TABLES;"

# 3. Initialize Go project
go mod init event-pipeline
```

### Connection Details
```
MySQL Connection String: "root:testpass@tcp(localhost:3306)/eventdb"
```

---

## Required Metrics

Track and expose these metrics via `/metrics` endpoint:
```json
{
  "events_received": 15430,
  "events_processed": 15389,
  "events_failed": 41,
  "average_processing_latency_ms": 23.5,
  "current_queue_depth": 45,
  "active_workers": 10,
  "events_per_second": 157.3,
  "uptime_seconds": 3600
}
```

---

## Project Structure
```
event-pipeline/
├── docker-compose.yml           # Provided infrastructure
├── cmd/
│   └── main.go                  # Application entry point
├── internal/
│   ├── api/                     # HTTP handlers
│   ├── pipeline/                # Core processing pipeline
│   ├── storage/                 # Database operations
│   ├── metrics/                 # Metrics collection
│   └── config/                  # Configuration
├── pkg/
│   └── validator/               # Event validation
├── tests/
│   ├── integration/
│   └── unit/
├── Dockerfile                   # YOUR containerization
├── go.mod
├── go.sum
└── README.md
```

---

## Time Allocation Guidance

| Phase | Time | Focus Areas |
|-------|------|-------------|
| **Core Pipeline** | 90 min | Worker pool, goroutines, channels, processing logic |
| **HTTP API + Database** | 45 min | REST endpoints, MySQL integration |
| **Metrics & Error Handling** | 45 min | Real-time metrics, graceful failures |
| **Testing & Documentation** | 60 min | Unit tests, integration test, README |

**Total: 4 hours maximum**

---

## Required Testing

### Unit Tests (Minimum Required)
```go
func TestWorkerPoolProcessing(t *testing.T) {
    // Test concurrent processing with multiple workers
}

func TestEventValidation(t *testing.T) {
    // Test event validation logic
}

func TestMetricsAccuracy(t *testing.T) {
    // Test metrics collection and calculation
}

func TestGracefulShutdown(t *testing.T) {
    // Test context cancellation and cleanup
}
```

### Integration Test
```go
func TestEndToEndPipeline(t *testing.T) {
    // Send events through entire pipeline
    // Verify storage and metrics
}
```

---

## Deliverables

### 1. Working Application
- Complete Go application that compiles and runs
- All required API endpoints functional
- Concurrent processing pipeline working
- MySQL integration with proper connection pooling

### 2. Containerization
Create your own `Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o event-pipeline cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/event-pipeline .
EXPOSE 8080

CMD ["./event-pipeline"]
```

### 3. Documentation
**README.md must include:**
- Project overview and architecture
- Setup and installation instructions
- How to run the application and tests
- API endpoint documentation with examples
- Design decisions and trade-offs made

---

## Sample API Usage

### Single Event
```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
    "type": "user_action",
    "source": "web",
    "user_id": "user_123",
    "data": {
      "action": "purchase",
      "value": 99.99
    }
  }'
```

### Batch Events
```bash
curl -X POST http://localhost:8080/events/batch \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      {
        "type": "sensor_data",
        "source": "iot_device",
        "data": {"temperature": 23.5}
      },
      {
        "type": "system_log", 
        "source": "system",
        "data": {"level": "error", "message": "Connection timeout"}
      }
    ]
  }'
```

### Get Metrics
```bash
curl http://localhost:8080/metrics
```

### Health Check
```bash
curl http://localhost:8080/health
```

---

## Evaluation Criteria

### Go Expertise (50% weight)
- ✅ **Goroutines & Channels**: Proper usage of Go's concurrency primitives
- ✅ **Context Handling**: Appropriate use for cancellation and timeouts
- ✅ **Interface Design**: Clean abstractions and dependency injection
- ✅ **Error Handling**: Go idiomatic error patterns
- ✅ **Memory Safety**: Avoiding race conditions, proper channel closing

### System Design (30% weight)  
- ✅ **Architecture**: Clean separation of concerns, scalable design
- ✅ **Performance**: Efficient processing, proper buffering strategies
- ✅ **Reliability**: Error recovery, graceful degradation
- ✅ **Observability**: Meaningful metrics and logging

### Code Quality (20% weight)
- ✅ **Structure**: Well-organized, readable code
- ✅ **Testing**: Meaningful test coverage
- ✅ **Documentation**: Clear setup instructions and API docs
- ✅ **Production-Readiness**: Configuration, containerization

---

## Success Benchmarks

### Minimum Viable (Meets Expectations)
- ✅ Working pipeline processing events concurrently
- ✅ Basic worker pool with goroutines and channels  
- ✅ MySQL integration with event storage
- ✅ Core metrics collection and API
- ✅ Essential unit tests
- ✅ Clear documentation

### Exceeds Expectations (Senior-level)
- ✅ Sophisticated concurrency patterns
- ✅ Comprehensive error handling and recovery
- ✅ Performance optimizations (batching, connection pooling)
- ✅ Advanced testing (race condition detection, load testing)
- ✅ Production considerations (health checks, graceful shutdown)

---

## Submission Instructions

1. **Create a Git repository** with your complete solution
2. **Include all required files**: 
   - Source code with proper structure
   - Tests (unit and integration)
   - Dockerfile and docker-compose.yml
   - README.md with complete documentation
3. **Ensure everything works**:
   - Application starts successfully: `go run cmd/main.go`
   - All tests pass: `go test ./...`
   - Docker build succeeds: `docker build -t event-pipeline .`
4. **Submit the repository URL** when complete

---

## Final Checklist

- [ ] Application compiles and runs without errors
- [ ] All required API endpoints are implemented and working
- [ ] Database integration works (can store and retrieve events)
- [ ] Concurrent processing uses proper Go patterns (goroutines + channels)
- [ ] Metrics endpoint returns accurate real-time data
- [ ] Unit tests pass and cover key functionality  
- [ ] Integration test demonstrates end-to-end functionality
- [ ] README includes clear setup and usage instructions
- [ ] Dockerfile builds successfully

---

## Questions or Clarifications?

If you have questions about the requirements, feel free to ask for clarification. However, part of the assessment is your ability to make reasonable technical decisions when requirements are ambiguous.

**Focus your time on demonstrating Go concurrency expertise and clean system design** rather than trying to implement every possible feature.

---

**Good luck!** We're excited to see your Go skills in action. This assessment is designed to showcase your real-world Go development capabilities, especially your understanding of Go's unique strengths in building concurrent, scalable systems.

**Time limit: 4 hours maximum** ⏰