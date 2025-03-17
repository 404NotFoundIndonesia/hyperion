package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"hyperion/backend/obfuscator"
	"os/exec"
	"reflect"
	"regexp"
	rt "runtime"
)

type Menu struct {
	ctx    context.Context
	config *obfuscator.Config
}

func NewMenu(config *obfuscator.Config) *Menu {
	return &Menu{config: config}
}

func (m *Menu) startup(ctx context.Context) {
	m.ctx = ctx
}

func (m *Menu) CreateMenu(app *application.Application) *menu.Menu {
	appMenu := menu.NewMenu()

	// File Menu
	fileMenu := appMenu.AddSubmenu("File")
	fileMenu.AddText("Open", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(m.ctx, "menu:open")
	})
	fileMenu.AddText("Open Folder", keys.Combo("o", keys.ShiftKey, keys.CmdOrCtrlKey), func(_ *menu.CallbackData) {
		runtime.EventsEmit(m.ctx, "menu:open-folder")
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("Save", keys.CmdOrCtrl("s"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(m.ctx, "menu:save")
	})
	fileMenu.AddText("Save All", keys.Combo("s", keys.ShiftKey, keys.CmdOrCtrlKey), func(_ *menu.CallbackData) {
		runtime.EventsEmit(m.ctx, "menu:save-all")
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("Close", keys.CmdOrCtrl("w"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(m.ctx, "menu:close")
	})
	fileMenu.AddText("Close Folder", keys.Combo("w", keys.ShiftKey, keys.CmdOrCtrlKey), func(_ *menu.CallbackData) {
		runtime.EventsEmit(m.ctx, "menu:close-folder")
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		app.Quit()
	})

	// Run Menu
	runMenu := appMenu.AddSubmenu("Run")
	runMenu.AddText("Obfuscate", keys.CmdOrCtrl("r"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(m.ctx, "run:obfuscate")
	})
	runMenu.AddText("Obfuscate All", keys.Combo("r", keys.ShiftKey, keys.CmdOrCtrlKey), func(_ *menu.CallbackData) {
		runtime.EventsEmit(m.ctx, "run:obfuscate-all")
	})
	runMenu.AddSeparator()
	runMenu.AddText("Cancel Obfuscation", keys.CmdOrCtrl("e"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(m.ctx, "run:cancel")
	})
	runMenu.AddText("Cancel All Obfuscation", keys.Combo("e", keys.ShiftKey, keys.CmdOrCtrlKey), func(_ *menu.CallbackData) {
		runtime.EventsEmit(m.ctx, "run:cancel-all")
	})
	runMenu.AddSeparator()
	configMenu := runMenu.AddSubmenu("Configuration")
	configValue := reflect.ValueOf(m.config).Elem()
	configType := configValue.Type()
	for i := 0; i < configValue.NumField(); i++ {
		field := configType.Field(i)
		if field.Name == "path" {
			continue
		}

		configMenu.AddCheckbox(toCamelCaseLabel(field.Name), configValue.FieldByName(field.Name).Bool(), nil, func(_ *menu.CallbackData) {
			currentValue := configValue.FieldByName(field.Name).Bool()
			configValue.FieldByName(field.Name).SetBool(!currentValue)

			m.config.Save(*m.config)
		})
	}
	runMenu.AddText("Open Configuration Folder", nil, func(_ *menu.CallbackData) {
		switch rt.GOOS {
		case "windows":
			_ = exec.Command("explorer", m.config.GetPath()).Start()
			break
		case "darwin":
			_ = exec.Command("open", "-R", m.config.GetPath()).Start()
			break
		case "linux":
			_ = exec.Command("xdg-open", m.config.GetPath()).Start()
			break
		default:
			break
		}
	})

	// Help Menu
	helpMenu := appMenu.AddSubmenu("Help")
	helpMenu.AddText("How to Contribute", nil, func(_ *menu.CallbackData) {
		println("How to Contribute clicked")
	})
	helpMenu.AddText("Report Issue", nil, func(_ *menu.CallbackData) {
		runtime.BrowserOpenURL(m.ctx, "https://github.com/404NotFoundIndonesia/hyperion/issues")
	})

	return appMenu
}

func toCamelCaseLabel(input string) string {
	re := regexp.MustCompile("([a-z])([A-Z])")
	return re.ReplaceAllString(input, "$1 $2")
}
