package main

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"

	"github.com/kirsle/configdir"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/build
var assets embed.FS

//go:embed data.db
var b []byte

func main() {
	// create the config path
	configPath := configdir.LocalConfig("httpxui", "/pb_data")
	err1 := configdir.MakePath(configPath)
	if err1 != nil {
		panic(err1)
	}

	// write the db file to the local system's config path
	data_file := configdir.LocalConfig("httpxui") + "/pb_data/data.db"
	_, err3 := os.Stat(data_file)
	if errors.Is(err3, fs.ErrNotExist) {
		err2 := os.WriteFile(fmt.Sprintf("%s/data.db", configPath), b, 0644)
		if err2 != nil {
			panic(err2)
		}
	}

	// start pocketbase server
	go func() {
		data_dir := fmt.Sprintf("%s/pb_data/", configdir.LocalConfig("httpxui"))

		pocketBaseApp := pocketbase.NewWithConfig(pocketbase.Config{
			DefaultDataDir: data_dir,
		})

		pocketBaseApp.Bootstrap()

		apis.Serve(pocketBaseApp, apis.ServeConfig{
			HttpAddr:        "127.0.0.1:8081",
			ShowStartBanner: false,
		})
	}()

	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "httpxui",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
