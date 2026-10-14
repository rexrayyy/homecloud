//go:build windows

// Package client contains the native Wails desktop shell and its local actions.
package client

import (
	"embed"
	"io/fs"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed frontend/index.html frontend/app.js frontend/styles.css frontend/vendor/xterm.js frontend/vendor/xterm.css frontend/vendor/xterm-addon-fit.js
var assets embed.FS

// Run starts the Wails desktop client.
func Run() error {
	root, err := fs.Sub(assets, "frontend")
	if err != nil {
		return err
	}
	// A modest initial window avoids the top title bar starting off-screen.
	a := &App{}
	return wails.Run(&options.App{Title: "HomeCloud", Width: 900, Height: 600, MinWidth: 640, MinHeight: 460, StartHidden: false, AssetServer: &assetserver.Options{Assets: root}, OnStartup: a.startup, OnShutdown: a.shutdown, Bind: []interface{}{a}})
}
