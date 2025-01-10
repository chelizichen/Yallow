package main

import (
	"Yallow/backend/apis"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := apis.NewApp()
	AppMenu := menu.NewMenu()
	FileMenu := AppMenu.AddSubmenu("File")
	FileMenu.AddSeparator()
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.Quit(app.Ctx)
	})
	FileMenu.AddText("Copy", keys.CmdOrCtrl("c"), func(t *menu.CallbackData) {
		runtime.ClipboardSetText(app.Ctx, t.MenuItem.Label)
	})
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Yallow",
		Width:  1600,
		Height: 1200,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
		Menu: AppMenu,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}