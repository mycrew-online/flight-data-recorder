package logger

import (
	logz "github.com/mrlm-net/go-logz/pkg/logger"
	"github.com/mycrew-online/flight-data-recorder/internal/logadapter"
)

var AppLogger = logadapter.New(logz.NewLogger(logz.LogOptions{
	Level:   logz.Debug,
	Format:  logz.StringOutput,
	Prefix:  "App",
	Outputs: []logz.OutputFunc{logz.ConsoleOutput()},
}))
