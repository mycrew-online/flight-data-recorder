package logadapter

import (
	logz "github.com/mrlm-net/go-logz/pkg/logger"
)

// LogzWailsAdapter implements the Wails logger interface using go-logz
// and can be used as a drop-in logger for both Wails and your app.
type LogzWailsAdapter struct {
	logz logz.ILogger
}

func New(logzLogger *logz.Logger) *LogzWailsAdapter {
	return &LogzWailsAdapter{logz: logzLogger}
}

func (l *LogzWailsAdapter) Print(message string) {
	l.logz.Info(message)
}

func (l *LogzWailsAdapter) Printf(format string, args ...interface{}) {
	l.logz.Info(format, map[string]interface{}{"args": args})
}

func (l *LogzWailsAdapter) Trace(message string) {
	// Suppress TRACE logs for 'No listeners for event' messages
	if len(message) >= 22 && message[:22] == "No listeners for event" {
		return
	}
	l.logz.Debug("TRACE: " + message)
}

func (l *LogzWailsAdapter) Debug(message string) {
	l.logz.Debug(message)
}

func (l *LogzWailsAdapter) Info(message string) {
	l.logz.Info(message)
}

func (l *LogzWailsAdapter) Warning(message string) {
	l.logz.Warning(message)
}

func (l *LogzWailsAdapter) Error(message string) {
	l.logz.Error(message)
}

func (l *LogzWailsAdapter) Fatal(message string) {
	l.logz.Critical(message)
}

func (l *LogzWailsAdapter) Panic(message string) {
	l.logz.Alert(message)
}
