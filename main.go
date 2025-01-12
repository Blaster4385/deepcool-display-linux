package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"os"
)

//go:embed all:frontend/dist
var assets embed.FS

var (
	daemonFlag  bool
	filename    string
	helpFlag    bool
	versionFlag bool
	version     = "0.1.0"
)

func printHelp() {
	fmt.Printf(`DeepCool Display Linux Controller

Usage:
    %s [options]

Options:
    -h, --help     Show this help message
    -d, --daemon   Run in daemon mode
    -f, --file     Specify CSV file path for pattern
    -v, --version  Show the version of the app

Modes:
    1. GUI Mode (default):
       Run without any flags to start the graphical interface
       Example: %s

    2. Daemon Mode:
       Run with -d flag and specify a CSV file to load a pattern
       Example: %s -d -f pattern.csv

For more information, visit: https://github.com/blaster4385/deepcool-display-linux
`, os.Args[0], os.Args[0], os.Args[0])
}

func main() {
	flag.BoolVar(&daemonFlag, "daemon", false, "Run as daemon")
	flag.StringVar(&filename, "file", "", "CSV file")
	flag.BoolVar(&helpFlag, "help", false, "Show help message")
	flag.BoolVar(&versionFlag, "version", false, "Show app version")
	flag.BoolVar(&versionFlag, "v", false, "Show app version")
	flag.Parse()

	if helpFlag {
		printHelp()
		os.Exit(0)
	}

	if versionFlag {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	app := NewApp()
	if daemonFlag {
		if filename == "" {
			fmt.Println("Error: CSV file path is required in daemon mode")
			fmt.Println("Use -h or --help for usage information")
			os.Exit(1)
		}

		ctx := context.Background()
		app.startup(ctx)
		grid, err := app.ParseCSV(filename)
		if err != nil {
			fmt.Printf("Error parsing CSV file: %v\n", err)
			os.Exit(1)
		}
		err = app.SendPattern(grid)
		if err != nil {
			fmt.Printf("Error sending pattern: %v\n", err)
			os.Exit(1)
		}
		select {}
	} else {
		err := wails.Run(&options.App{
			Title:  "deepcool-display-linux",
			Width:  1024,
			Height: 768,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
			BackgroundColour: &options.RGBA{R: 40, G: 40, B: 40, A: 1},
			OnStartup:        app.startup,
			Bind: []interface{}{
				app,
			},
		})
		if err != nil {
			fmt.Printf("Error starting application: %v\n", err)
			os.Exit(1)
		}
	}
}
