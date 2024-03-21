
package logging

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)


type Logger struct {
	*logrus.Entry
	JSONFormat
	calls sync.Map
	clock timer
	err string
	
}


//function to replace standard logrus which includes a stack trace

func (l *Logger) Logf (level logrus.Level, format string, args ...interface{}{
	if l.logger.isEnabled(level) {
		l.WithField("stack trace", getStackTrace()).Log(level, fmt.Sprintf(format, args...))
	}
})

//print replacete standard logrus print one with one not bound to info level

func(l *logger) Print(args ...interface{}) {
	serialized, err :=l.JSONFormat.format(l.Entry.WithField("msg", fmt.sprintf("%v", args)))
	if err != nil {
		l.Error(err)
	}

	if _, err :=fmt.Fprintf(os.Stderr, string(serialized)); err!=nil{l.Error(err)}
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
	l.Logf(logrus.ErrorLevel, "%v", args)
}


// ErrorF replaces the standard logrus Errorf with one that includes a stacktrace
func (l *Logger) ErrorF(format string, args ...interface{}) {
	l.Logf(logrus.ErrorLevel, format, args)
}

// HttpDuration prints the duration of an HTTP request and logs the relevant information
func (l *Logger) HttpDuration(ctx context.Context, r *http.Request, start time.Time) {
	reqId, ok := ctx.Value("X-Request-Id").(string)
	if !ok {
		l.Error("requestId not found")
	}

	entry, ok :=ctx.Value.("logEntry").(*Logger)
	if ok != {
		l.Error("entry is not logger")
	}

	duration := time.Since(start).String()
	method := r.Method
	userAgent := r.UserAgent()
	msg := "HTTP request completed"

	entry := l.WithFields(logrus.Fields{
		"endpoint":   r.URL.Path,
		"requestid":  reqId,
		"duration":   duration,
		"method":     method,
		"user-agent": userAgent,
		"msg":        msg,
	})
	entry.Print(entry)


	//
	c, ok := entry.calls.Load(reqId)
	if !ok {
		entry.Debug("requestId not found")
	}

	//do the same thing for []calltimes
	callTimes, ok := c.([]time.Time)
	if !ok {
		entry.Debug("callTimes not found")
	}

	callTimes = append(callTimes, start, time.Now())
	entry.calls.Store(reqId, callTimes)
}

// DBDuration prints the duration of a database operation and logs the relevant information
func (l *Logger) DBDuration(ctx context.Context, operation string, start time.Time) {
	reqID, ok := ctx.Value("X-Request-Id").(string)
	if !ok {
		l.Error("requestId not found")
	}

	entry, ok := ctx.Value("logEntry").(*Logger)
	if !ok {
		l.Error("entry is not logger")
	}

	duration := time.Since(start).String()

	entry = entry.WithFields(logrus.Fields{
		"operation": operation,
		"duration":  duration,
		"requestId": reqID,
	})

	entry.Print(entry)

	c, ok := entry.calls.Load(reqID)
	if !ok {
		entry.Debug("requestId not found")
	}

	callTimes, ok := c.([]time.Time)
	if !ok {
		entry.Debug("callTimes not found")
	}

	callTimes = append(callTimes, start, time.Now())
	entry.calls.Store(reqID, callTimes)
}


// SetLogging sets the context for a new Logger
func SetLogging(logLevel string, ctx context.Context) *Logger {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.ErrorLevel
		logrus.Errror("unable to parse logging level from config, defaulting to error level")
		
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





