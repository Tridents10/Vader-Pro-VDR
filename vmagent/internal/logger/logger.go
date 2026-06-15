package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger writes to both the agent.log file and stdout.
type Logger struct {
	file *os.File
	std  *log.Logger
}

// New creates a logger writing to output/agent.log and stdout.
func New(outputDir string) (*Logger, error) {
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return nil, err
	}

	f, err := os.OpenFile(
		filepath.Join(outputDir, "agent.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0o644,
	)
	if err != nil {
		return nil, err
	}

	return &Logger{
		file: f,
		std:  log.New(os.Stdout, "", 0),
	}, nil
}

func (l *Logger) write(level, msg string) {
	line := fmt.Sprintf("%s %-5s %s",
		time.Now().UTC().Format(time.RFC3339), level, msg)
	l.std.Println(line)
	if l.file != nil {
		fmt.Fprintln(l.file, line)
	}
}

// Info logs an informational message.
func (l *Logger) Info(msg string) { l.write("INFO", msg) }

// Warn logs a warning message.
func (l *Logger) Warn(msg string) { l.write("WARN", msg) }

// Error logs an error message.
func (l *Logger) Error(msg string) { l.write("ERROR", msg) }

// Close closes the underlying log file.
func (l *Logger) Close() {
	if l.file != nil {
		_ = l.file.Close()
	}
}
