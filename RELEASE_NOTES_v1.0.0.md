# Email Worker v1.0.0

## Overview

The Email Worker is a production-ready, distributed email processing service that consumes email messages from RabbitMQ and sends them via SMTP. This is the initial stable release featuring a robust architecture with health checks, metrics, structured logging, and comprehensive error handling.

## Key Features

### Core Functionality
- **RabbitMQ Integration**: Consumes email messages from configurable RabbitMQ queues
- **SMTP Email Sending**: Sends emails via SMTP with support for TLS/STARTTLS
- **Multipart Email Support**: Handles both text and HTML email bodies
- **Dead Letter Queue (DLQ)**: Failed messages are routed to a dead letter exchange for analysis
- **Message Acknowledgment**: Manual acknowledgment ensures reliable message processing

### Production Features
- **Health Check Endpoint**: HTTP endpoint at `/healthz` (port 8080) for Kubernetes liveness/readiness probes
- **Prometheus Metrics**: Comprehensive metrics exposed at `/metrics` endpoint
  - Job processing counters and durations
  - Failed job tracking
  - Queue depth monitoring
  - Dead letter queue size tracking
- **Structured Logging**: JSON logging in production, text logging in development with trace ID support
- **Graceful Shutdown**: Handles SIGTERM/SIGINT signals for clean shutdown
- **Connection Retry Logic**: Automatic retry with exponential backoff for RabbitMQ connections

### Configuration
- **Environment-based Configuration**: All settings via environment variables
- **SMTP Configuration**: Supports TLS, authentication, and custom ports
- **RabbitMQ Configuration**: Flexible connection string or individual parameters
- **Queue Configuration**: Configurable queue names, exchanges, and routing keys
- **Prefetch Control**: Configurable message prefetch count for optimal throughput

### Code Quality
- **Comprehensive Testing**: Unit tests and integration tests included
- **Test Coverage**: Coverage reporting with threshold checks (70% minimum)
- **Clean Architecture**: Separation of concerns with internal and pkg packages
- **Go 1.23**: Built with latest Go version

## Architecture

```
┌─────────────┐
│  RabbitMQ   │
│   Queue     │
└──────┬──────┘
       │
       │ EmailEnvelope
       │
┌──────▼─────────────────┐
│   Email Worker         │
│                        │
│  ┌──────────────────┐  │
│  │  Queue Consumer  │  │
│  └────────┬─────────┘  │
│           │            │
│  ┌────────▼─────────┐  │
│  │  SMTP Sender     │  │
│  └────────┬─────────┘  │
└───────────┼────────────┘
            │
            │ SMTP Protocol
            │
    ┌───────▼────────┐
    │  SMTP Server   │
    └────────────────┘
```

## Dependencies

- **github.com/rabbitmq/amqp091-go**: RabbitMQ client library
- **github.com/prometheus/client_golang**: Prometheus metrics
- **github.com/google/uuid**: UUID generation for message IDs
- **github.com/stretchr/testify**: Testing utilities

## Configuration Variables

### SMTP Configuration
- `SMTP_HOST`: SMTP server hostname (required)
- `SMTP_PORT`: SMTP server port (default: 587)
- `SMTP_USERNAME`: SMTP authentication username
- `SMTP_PASSWORD`: SMTP authentication password
- `SMTP_FROM`: Sender email address (required)
- `SMTP_TLS`: Enable TLS (default: true)

### RabbitMQ Configuration
- `RABBITMQ_URL`: Full AMQP connection URL (overrides individual settings)
- `RABBITMQ_USER`: RabbitMQ username (default: "woragis")
- `RABBITMQ_PASSWORD`: RabbitMQ password (default: "woragis")
- `RABBITMQ_HOST`: RabbitMQ hostname (default: "rabbitmq")
- `RABBITMQ_PORT`: RabbitMQ port (default: "5672")
- `RABBITMQ_VHOST`: RabbitMQ virtual host (default: "woragis")

### Worker Configuration
- `EMAIL_QUEUE_NAME`: Queue name (default: "emails.queue")
- `EMAIL_EXCHANGE`: Exchange name (default: "woragis.notifications")
- `EMAIL_ROUTING_KEY`: Routing key (default: "emails.send")
- `EMAIL_PREFETCH_COUNT`: Prefetch count (default: 1)

### Environment
- `ENV`: Environment mode - "development" or "production" (affects logging)

## Message Format

The worker expects messages in the following JSON format:

```json
{
  "user_id": "uuid",
  "subject": "Email Subject",
  "text_message": "Plain text email body",
  "html_message": "HTML email body",
  "destination": "recipient@example.com"
}
```

## Health Check

The service exposes a health check endpoint at `GET /healthz`:

- **Healthy Response (200 OK)**: All checks pass
- **Unhealthy Response (503)**: RabbitMQ connection is down

The health check includes:
- RabbitMQ connection status check (2s timeout)
- Result caching (5s TTL) for performance

## Metrics

Prometheus metrics are exposed at `GET /metrics`:

- `worker_jobs_processed_total{worker, status}`: Total jobs processed
- `worker_job_duration_seconds{worker}`: Job processing duration histogram
- `worker_jobs_failed_total{worker, error_type}`: Failed job counter
- `worker_jobs_retried_total{worker}`: Retried job counter
- `queue_depth{queue_name}`: Current queue depth
- `queue_dlq_size{queue_name}`: Dead letter queue size

## Deployment

### Docker
The service can be containerized and deployed as a standalone worker or in Kubernetes.

### Kubernetes
- **Liveness Probe**: `GET /healthz` every 10s, timeout 5s
- **Readiness Probe**: `GET /healthz` every 10s, timeout 5s
- **Port**: 8080 (health and metrics)

### Scaling
Multiple worker instances can consume from the same queue for horizontal scaling. Prefetch count can be adjusted based on message processing time.

## Development

### Testing
```bash
# Run all tests
make test

# Run unit tests only
make test-unit

# Run with coverage
make test-cov

# Check coverage threshold (70%)
make test-cov-check
```

### Building
```bash
go build ./cmd/email-worker
```

## Documentation

- **HEALTH_CHECK.md**: Detailed health check documentation
- **LOGGING.md**: Structured logging guidelines and usage
- **tests/README.md**: Testing documentation
- **internal/integration/README.md**: Integration test documentation

## Breaking Changes

None - This is the initial release.

## Future Enhancements

Potential improvements for future versions:
- Email template support
- Retry policies with exponential backoff
- Rate limiting per recipient
- Email bounce handling
- Multiple SMTP provider support (round-robin, failover)
- Email delivery status tracking

## Contributors

Initial release by the Woragis team.

## License

Part of the Woragis backend infrastructure.

