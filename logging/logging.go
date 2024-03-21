package logging

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type JSONFormat struct {
	// Add JSONFormat struct fields here
}

type Logger struct {
	*logrus.Entry
	JSONFormat
	calls sync.Map
	clock time.Timer
	err   string
}

// Function to replace standard logrus which includes a stack trace
func (l *Logger) Logf(level logrus.Level, format string, args ...interface{}) {
	if l.Level >= level { // Fixed: Changed from `l.logger.isEnabled(level)` to direct logrus Level comparison
		l.WithField("stack trace", getStackTrace()).Log(level, fmt.Sprintf(format, args...))
	}
}

func getStackTrace() string {
	stack := make([]byte, 4096)
	length := runtime.Stack(stack, false)
	return string(stack[:length])
}

// Print replaces the standard logrus print, with one not bound to info level
func (l *Logger) Print(args ...interface{}) {
	serialized, err := l.JSONFormat.format(l.Entry.WithField("msg", fmt.Sprint(args...))) // Changed `fmt.Sprintf("%v", args...)` to `fmt.Sprint(args...)`
	if err != nil {
		l.Error(err)
	}

	if _, err := fmt.Fprintf(os.Stderr, string(serialized)); err != nil {
		l.Error(err)
	}
}

func (jf *JSONFormat) format(entry *logrus.Entry) ([]byte, error) {
	// Implement the logic to format the entry as JSON
	return nil, nil
}

// PrintDuration prints the duration of an entry and performs a custom print operation
func (l *Logger) PrintDuration(entry *logrus.Entry, duration time.Duration) {
	serialized, err := l.JSONFormat.format(entry.WithField("duration", duration.String()))
	if err != nil {
		l.Error(err)
	}

	if _, err := fmt.Fprintf(os.Stderr, string(serialized)); err != nil {
		l.Error(err)
	}
}

// ErrorWithStackTrace replaces the standard logrus error with one that includes a stacktrace
func (l *Logger) ErrorWithStackTrace(err error) {
	l.WithField("stack trace", getStackTrace()).Error(err)
}

func (l *Logger) Error(err error) {
	l.Logf(logrus.ErrorLevel, "%v", err) // Fixed: Changed `args` to `err`
}

// ErrorF replaces the standard logrus Errorf with one that includes a stacktrace
func (l *Logger) ErrorF(format string, args ...interface{}) {
	l.Logf(logrus.ErrorLevel, format, args...)
}

// HttpDuration prints the duration of an HTTP request and logs the relevant information
func (l *Logger) HttpDuration(ctx context.Context, r *http.Request, start time.Time) {
	// Your existing HttpDuration code seems correct, assuming your context handling and Logger setup are as intended.
}

// DBDuration prints the duration of a database operation and logs the relevant information
func (l *Logger) DBDuration(ctx context.Context, operation string, start time.Time) {
	// Your existing DBDuration code seems correct, assuming your context handling and Logger setup are as intended.
}

// SetLogging sets the context for a new Logger
func SetLogging(logLevel string, ctx context.Context) *Logger {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.ErrorLevel
		logrus.Error("unable to parse logging level from config, defaulting to error level") // Fixed typo: Changed `logrus.Errror` to `logrus.Error`
	}

	ctx = context.WithValue(ctx, "logLevel", level)
	std := logrus.New()
	std.SetLevel(level)

	logger := &Logger{
		Entry: logrus.NewEntry(std),
	}

	if entry, ok := ctx.Value("logEntry").(*Logger); ok {
		logger.Entry = entry.Entry
	}

	return logger
}
