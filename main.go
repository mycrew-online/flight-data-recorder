package main

import (
	"embed"

	//"github.com/mrlm-net/go-logz/pkg/logger"
	"github.com/mycrew-online/flight-data-recorder/internal"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:website/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := internal.NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "MyCrew.online FDR",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		OnShutdown:       app.Shutdown,
		// Logger:           logger.NewLogger(logger.LogOptions{}), // Removed due to interface incompatibility
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
