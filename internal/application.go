package internal

import (
	"context"

	"github.com/mycrew-online/flight-data-recorder/internal/logger"
	simconnectmanager "github.com/mycrew-online/flight-data-recorder/pkg/simconnect-manager"
)

// App struct
type App struct {
	ctx        context.Context
	simconnect *simconnectmanager.SimConnectManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	mgr := simconnectmanager.NewSimConnectManager()
	mgr.SetLogger(logger.AppLogger)
	return &App{
		simconnect: mgr,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	a.simconnect.SetWailsContext(ctx)
	logger.AppLogger.Info("App has started")

	// Start SimConnect connection monitoring
	a.simconnect.StartConnection()

	// Listen for connection status changes
	go func() {
		for status := range a.simconnect.StatusChan() {
			if status {
				logger.AppLogger.Info("SimConnect connection established!")
			} else {
				logger.AppLogger.Warning("SimConnect disconnected.")
			}
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) {
	logger.AppLogger.Info("App is shutting down")
	a.simconnect.StopConnection()
}

// GetSimStatus returns the current SimConnect connection status
func (a *App) GetSimStatus() bool {
	return a.simconnect.Status()
}
