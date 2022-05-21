package logging

import (
	"io"
	"os"

	log "github.com/hashicorp/go-hclog"
)

var level string = "INFO"
var output io.Writer = os.Stdout

// Logger defines the interface that is used for logging operations.
type Logger struct {
	log.Logger
}

// SetLevel sets logging level.
func SetLevel(v string) {
	level = v
}

// Level returns logging level.
func Level() string {
	return level
}

// Output sets logging output.
func Output(v io.Writer) {
	output = v
}

func (l *Logger) Printf(msg string, args ...interface{}) {
	l.Logger.Info(msg, args...)
}

// NewLogger creates a new Logger instance with name n.
func NewLogger(n string) *Logger {
	return &Logger{
		Logger: log.New(&log.LoggerOptions{
			Level:  log.LevelFromString(level),
			Output: output,
			Color:  log.ForceColor,
		}),
	}
}
