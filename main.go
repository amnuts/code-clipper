package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := CodeClipper()

	// Create application with options
	err := wails.Run(&options.App{
		Title:            "code-clipper",
		Width:            1024,
		Height:           768,
		Assets:           assets,
		BackgroundColour: &options.RGBA{R: 24, G: 24, B: 24, A: 1},
		Frameless:        false,
		OnStartup:        app.startup,
		OnBeforeClose:    app.beforeShutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
