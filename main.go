package main

import (
	"embed"
	"fmt"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

func main() {
	application := NewApp()

	err := wails.Run(&options.App{
		Title:     appTitle,
		Width:     defaultWindowWidth,
		Height:    defaultWindowHeight,
		MinWidth:  minWindowWidth,
		MinHeight: minWindowHeight,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 18, G: 20, B: 24, A: 1},
		Linux:            &linux.Options{Icon: appIcon},
		OnStartup:        application.startup,
		Bind: []interface{}{
			application,
		},
	})
	if err != nil {
		fmt.Println("Error:", err.Error())
	}
}
