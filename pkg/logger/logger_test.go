package logger

import (
	"context"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		env  string
	}{
		{
			name: "development environment",
			env:  "development",
		},
		{
			name: "production environment",
			env:  "production",
		},
		{
			name: "empty environment",
			env:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.env)
			if logger == nil {
				t.Error("New() returned nil logger")
			}
		})
	}
}

func TestNewWithConfig(t *testing.T) {
	cfg := LogConfig{
		Env:       "development",
		LogToFile: false,
		LogDir:    "logs",
	}

	logger := NewWithConfig(cfg)
	if logger == nil {
		t.Error("NewWithConfig() returned nil logger")
	}
}

func TestWithTraceID(t *testing.T) {
	ctx := context.Background()
	traceID := "test-trace-id"
	ctxWithTrace := WithTraceID(ctx, traceID)

	retrievedTraceID := GetTraceID(ctxWithTrace)
	if retrievedTraceID != traceID {
		t.Errorf("GetTraceID() = %v, want %v", retrievedTraceID, traceID)
	}
}

func TestGetTraceID(t *testing.T) {
	ctx := context.Background()
	traceID := GetTraceID(ctx)
	// TraceID may be empty if not set, which is valid
	if traceID != "" {
		t.Errorf("GetTraceID() = %v, want empty string for context without trace ID", traceID)
	}

	ctxWithTrace := WithTraceID(ctx, "test-id")
	traceID = GetTraceID(ctxWithTrace)
	if traceID != "test-id" {
		t.Errorf("GetTraceID() = %v, want test-id", traceID)
	}
}

func TestNewWithConfig_FileLogging(t *testing.T) {
	// Test file logging in development mode
	cfg := LogConfig{
		Env:       "development",
		LogToFile: true,
		LogDir:    "test-logs",
	}

	logger := NewWithConfig(cfg)
	if logger == nil {
		t.Error("NewWithConfig() returned nil logger")
	}

	// Cleanup
	os.RemoveAll("test-logs")
}

func TestNewWithConfig_Production(t *testing.T) {
	cfg := LogConfig{
		Env:       "production",
		LogToFile: false,
	}

	logger := NewWithConfig(cfg)
	if logger == nil {
		t.Error("NewWithConfig() returned nil logger for production")
	}
}

func TestNewWithConfig_ProdEnv(t *testing.T) {
	cfg := LogConfig{
		Env:       "prod",
		LogToFile: false,
	}

	logger := NewWithConfig(cfg)
	if logger == nil {
		t.Error("NewWithConfig() returned nil logger for prod")
	}
}

func TestNewWithConfig_FileLogging_DefaultDir(t *testing.T) {
	cfg := LogConfig{
		Env:       "development",
		LogToFile: true,
		LogDir:    "", // Should use default
	}

	logger := NewWithConfig(cfg)
	if logger == nil {
		t.Error("NewWithConfig() returned nil logger")
	}

	// Cleanup
	os.RemoveAll(DefaultLogDir)
}

func TestServiceHandler_Handle(t *testing.T) {
	// Create a logger with service handler
	logger := New("development")
	if logger == nil {
		t.Fatal("New() returned nil logger")
	}

	// Test logging with trace ID
	ctx := WithTraceID(context.Background(), "test-trace-123")
	logger.InfoContext(ctx, "Test message", "key", "value")

	// Test logging without trace ID
	ctxNoTrace := context.Background()
	logger.InfoContext(ctxNoTrace, "Test message without trace", "key", "value")
}

func TestServiceHandler_Handle_EmptyTraceID(t *testing.T) {
	logger := New("development")
	ctx := context.WithValue(context.Background(), TraceIDKey, "")
	logger.InfoContext(ctx, "Test message")
}

func TestServiceHandler_Handle_NonStringTraceID(t *testing.T) {
	logger := New("development")
	ctx := context.WithValue(context.Background(), TraceIDKey, 123)
	logger.InfoContext(ctx, "Test message")
}