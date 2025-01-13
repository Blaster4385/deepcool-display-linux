package main

import (
	"context"
	"deepcool-display-linux/modules"
	"embed"
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"os"
	"time"
)

//go:embed all:frontend/dist
var assets embed.FS

const (
	version = "0.1.0"
)

var (
	daemonFlag  = flag.BoolP("daemon", "d", false, "Run in daemon mode")
	modeFlag    = flag.StringP("mode", "m", "", "Specify the mode (temp, usage or custom)")
	celsiusFlag = flag.BoolP("celsius", "c", false, "Show CPU temperature in Celsius")
	helpFlag    = flag.BoolP("help", "h", false, "Show help message")
	versionFlag = flag.BoolP("version", "v", false, "Show app version")
	filename    = flag.StringP("file", "f", "", "Specify CSV file path for pattern")
)

func printHelp() {
	fmt.Printf(`DeepCool Display Linux Controller

Usage:
    %s [options]

Options:
    -h, --help     Show this help message
    -d, --daemon   Run in daemon mode
    -m, --mode     Specify the mode (temp, usage or custom)
    -f, --file     Specify CSV file path for pattern
    -c, --celsius  Show CPU temperature in Celsius
    -v, --version  Show the version of the app

Modes:
    1. GUI Mode (default):
       Run without any flags to start the graphical interface
       Example: %s

    2. Daemon Mode:
       Run with -d flag and specify an option.
       Example: %s -d -f pattern.csv

For more information, visit: https://github.com/blaster4385/deepcool-display-linux
`, os.Args[0], os.Args[0], os.Args[0])
}

func handleError(err error) {
	fmt.Printf("Error: %v\n", err)
	os.Exit(1)
}

func main() {
	flag.Parse()

	if *helpFlag {
		printHelp()
		return
	}
	if *versionFlag {
		fmt.Printf("Version: %s\n", version)
		return
	}

	app := NewApp()

	if *daemonFlag {
		daemonMode(app)
	} else {
		guiMode(app)
	}
}

func daemonMode(app *App) {
	contextBackground, contextCancel := context.WithCancel(context.Background())
	defer contextCancel()

	if *modeFlag == "custom" {
		if *filename == "" {
			handleError(fmt.Errorf("filename is required for custom mode"))
		}
		handleDaemonPattern(app, *filename, contextBackground)
		return
	}
	if *modeFlag == "temp" {
		daemonTemperatureDisplay(app, contextBackground)
		return
	}
	if *modeFlag == "usage" {
		daemonUsageDisplay(app, contextBackground)
		return
	}

	select {}
}

func handleDaemonPattern(app *App, filename string, ctx context.Context) {
	grid, err := app.ParseCSV(filename)
	if err != nil {
		handleError(fmt.Errorf("parsing CSV file: %w", err))
	}
	err = app.SendPattern(grid)
	if err != nil {
		handleError(fmt.Errorf("sending pattern: %w", err))
	}
	select {}
}

func daemonTemperatureDisplay(app *App, ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				temp, err := modules.GetCPUTemperature(!*celsiusFlag)
				if err != nil {
					handleError(fmt.Errorf("getting temperature: %w", err))
				}
				var tempSymbol string
				if *celsiusFlag {
					tempSymbol = "celsius"
				} else {
					tempSymbol = "fahrenheit"
				}
				grid, err := modules.CreateNumberGrid(int(temp), tempSymbol, 5)
				if err != nil {
					handleError(fmt.Errorf("creating temperature grid: %w", err))
				}
				err = app.SendPattern(grid)
				if err != nil {
					handleError(fmt.Errorf("sending pattern: %w", err))
				}
			}
		}
	}()
	select {}
}

func daemonUsageDisplay(app *App, ctx context.Context) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				usage, err := modules.GetCPUUsage()
				if err != nil {
					handleError(fmt.Errorf("getting CPU usage: %w", err))
				}
				grid, err := modules.CreateNumberGrid(int(usage), "percent", 5)
				if err != nil {
					handleError(fmt.Errorf("creating usage grid: %w", err))
				}
				err = app.SendPattern(grid)
				if err != nil {
					handleError(fmt.Errorf("sending pattern: %w", err))
				}
			}
		}
	}()
	select {}
}

func guiMode(app *App) {
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
		handleError(fmt.Errorf("starting application: %w", err))
	}
}
