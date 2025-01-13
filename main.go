package main

import (
	"context"
	"deepcool-display-linux/modules"
	"embed"
	"flag"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"os"
	"time"
)

//go:embed all:frontend/dist
var assets embed.FS

var (
	daemonFlag  bool
	err         error
	filename    string
	tempFlag    bool
	celsiusFlag bool
	usageFlag   bool
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
    -t, --temp     Show CPU temperature
    -c, --celsius  Show CPU temperature in Celsius
    -u, --usage    Show CPU usage
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

func main() {
	flag.BoolVar(&daemonFlag, "daemon", false, "Run as daemon")
	flag.BoolVar(&daemonFlag, "d", false, "Run as daemon")
	flag.StringVar(&filename, "file", "", "CSV file")
	flag.StringVar(&filename, "f", "", "CSV file")
	flag.BoolVar(&tempFlag, "temp", false, "Show CPU temperature")
	flag.BoolVar(&tempFlag, "t", false, "Show CPU temperature")
	flag.BoolVar(&celsiusFlag, "celsius", false, "Show CPU temperature in Celsius")
	flag.BoolVar(&celsiusFlag, "c", false, "Show CPU temperature in Celsius")
	flag.BoolVar(&usageFlag, "usage", false, "Show CPU usage")
	flag.BoolVar(&usageFlag, "u", false, "Show CPU usage")
	flag.BoolVar(&helpFlag, "help", false, "Show help message")
	flag.BoolVar(&helpFlag, "h", false, "Show help message")
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
		ctx := context.Background()
		app.startup(ctx)
		if filename != "" {
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
		} else if tempFlag {
			ticker := time.NewTicker(3 * time.Second)
			defer ticker.Stop()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						var temp float64
						if celsiusFlag {
							temp, err = modules.GetCPUTemperature(false)
							grid, err := modules.CreateNumberGrid(int(temp), "celsius", 5)
							if err != nil {
								fmt.Printf("Error creating number grid: %v\n", err)
								os.Exit(1)
							}
							err = app.SendPattern(grid)
							if err != nil {
								fmt.Printf("Error sending pattern: %v\n", err)
								os.Exit(1)
							}
						} else {
							temp, err = modules.GetCPUTemperature(true)
							grid, err := modules.CreateNumberGrid(int(temp), "fahrenheit", 5)
							if err != nil {
								fmt.Printf("Error creating number grid: %v\n", err)
								os.Exit(1)
							}
							err = app.SendPattern(grid)
						}
					}
				}
			}()
		} else if usageFlag {
			ticker := time.NewTicker(3 * time.Second)
			defer ticker.Stop()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						usage, err := modules.GetCPUUsage()
						if err != nil {
							fmt.Printf("Error getting CPU usage: %v\n", err)
							os.Exit(1)
						}
						grid, err := modules.CreateNumberGrid(int(usage), "percent", 5)
						if err != nil {
							fmt.Printf("Error creating number grid: %v\n", err)
							os.Exit(1)
						}
						err = app.SendPattern(grid)
						if err != nil {
							fmt.Printf("Error sending pattern: %v\n", err)
							os.Exit(1)
						}
					}
				}
			}()
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
