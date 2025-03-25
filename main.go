package main

import (
	"context"
	"embed"
	fs "hyperion/backend/filesystem"
	ob "hyperion/backend/obfuscator"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()

	config := ob.NewConfig()

	file := fs.NewFileManager()
	dialog := fs.NewDialog()
	obfuscator := ob.NewObfuscator(config)

	menu := NewMenu(config)

	hyperion := application.NewWithOptions(&options.App{
		Title:     "hyperion",
		Width:     1024,
		Height:    768,
		MinWidth:  600,
		MinHeight: 400,
		Frameless: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			menu.startup(ctx)
			file.SetContext(ctx)
			dialog.SetContext(ctx)
			obfuscator.SetContext(ctx)
		},
		Bind: []interface{}{
			app, dialog,
			file, obfuscator,
		},
		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			BackdropType:         windows.Mica,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            true,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		Linux: &linux.Options{
			WindowIsTranslucent: true,
		},
	})
	if runtime.GOOS == "darwin" {
		hyperion.SetApplicationMenu(menu.CreateMenu(hyperion))
	}

	if err := hyperion.Run(); err != nil {
		log.Fatal(err)
	}
}
