package internal

import (
	"context"
	"fmt"

	simconnectmanager "github.com/mycrew-online/flight-data-recorder/pkg/simconnect-manager"
)

// App struct
type App struct {
	ctx        context.Context
	simconnect *simconnectmanager.SimConnectManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		simconnect: simconnectmanager.NewSimConnectManager(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	fmt.Println("App has started")

	// Start SimConnect connection monitoring
	a.simconnect.StartConnection()

	// Listen for connection status changes
	go func() {
		for status := range a.simconnect.StatusChan() {
			if status {
				fmt.Println("SimConnect connection established!")
			} else {
				fmt.Println("SimConnect disconnected.")
			}
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) {
	fmt.Println("App is shutting down")
	a.simconnect.StopConnection()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
