package main

import (
	"embed"

	//"github.com/mrlm-net/go-logz/pkg/logger"
	"github.com/mycrew-online/flight-data-recorder/internal"
	"github.com/mycrew-online/flight-data-recorder/internal/logger"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:website/build
var assets embed.FS

// added comment to trigger rebuild
func main() {
	// Create an instance of the app structure
	app := internal.NewApp()

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
		Logger:           logger.AppLogger,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
