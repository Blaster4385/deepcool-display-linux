package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"github.com/google/gousb"
	"os"
	// "strings"
	"sync"
	"time"
)

const (
	vendorID     = 0x3633
	productID    = 0x000c
	sendInterval = 750 * time.Millisecond
)

type App struct {
	ctx    context.Context
	usbCtx *gousb.Context
	device *gousb.Device
	iface  *gousb.Interface
	done   func()

	currentPattern []byte
	patternMutex   sync.RWMutex

	stopChan chan struct{}
	running  bool
	runMutex sync.Mutex
}

func NewApp() *App {
	return &App{
		stopChan: make(chan struct{}),
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) SendPattern(grid [][]bool) error {
	hexCommand, err := a.GenerateHexCommand(grid)
	if err != nil {
		return fmt.Errorf("failed to generate command: %v", err)
	}

	newPattern, err := hex.DecodeString(hexCommand)
	if err != nil {
		return fmt.Errorf("failed to decode hex command: %v", err)
	}

	a.patternMutex.Lock()
	patternChanged := !bytes.Equal(a.currentPattern, newPattern)
	a.patternMutex.Unlock()

	if !patternChanged && a.running {
		return nil
	}
	if !a.running {
		a.patternMutex.Lock()
		a.currentPattern = newPattern
		a.patternMutex.Unlock()

		return a.startContinuousSending()
	} else {
		a.patternMutex.Lock()
		a.currentPattern = newPattern
		a.patternMutex.Unlock()
	}

	return nil
}

func (a *App) initUSB() error {
	a.usbCtx = gousb.NewContext()

	var err error
	a.device, err = a.usbCtx.OpenDeviceWithVIDPID(vendorID, productID)
	if err != nil {
		return fmt.Errorf("could not open device: %v", err)
	}
	if a.device == nil {
		return fmt.Errorf("device not found")
	}

	a.device.SetAutoDetach(true)

	a.iface, a.done, err = a.device.DefaultInterface()
	if err != nil {
		return fmt.Errorf("failed to claim interface: %v", err)
	}

	return nil
}

func (a *App) cleanup() {
	a.runMutex.Lock()
	if a.running {
		close(a.stopChan)
		a.running = false
	}
	a.runMutex.Unlock()

	if a.done != nil {
		a.done()
	}
	if a.device != nil {
		a.device.Close()
	}
	if a.usbCtx != nil {
		a.usbCtx.Close()
	}
}

func (a *App) startContinuousSending() error {
	a.runMutex.Lock()
	defer a.runMutex.Unlock()

	if a.running {
		return nil
	}

	if a.device == nil {
		if err := a.initUSB(); err != nil {
			return fmt.Errorf("failed to initialize USB: %v", err)
		}
	}

	a.running = true
	a.stopChan = make(chan struct{})

	go func() {
		for {
			select {
			case <-a.stopChan:
				return
			default:
				a.patternMutex.RLock()
				currentData := make([]byte, len(a.currentPattern))
				copy(currentData, a.currentPattern)
				a.patternMutex.RUnlock()

				if err := a.sendCommand(currentData); err != nil {
					fmt.Printf("Error sending command: %v\n", err)
				}

				time.Sleep(sendInterval)
			}
		}
	}()

	return nil
}

func (a *App) sendCommand(data []byte) error {
	endpoint := 1
	outEndpoint, err := a.iface.OutEndpoint(endpoint)
	if err != nil {
		return fmt.Errorf("failed to get output endpoint: %v", err)
	}

	written, err := outEndpoint.Write(data)
	if err == nil {
		fmt.Printf("Success on endpoint 0x%02x! Sent %d bytes\n", endpoint, written)
		return nil
	}
	return fmt.Errorf("failed to send on any endpoint")
}

func (a *App) ExportGrid(filename string, grid [][]bool) error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}
	filePath := homedir + "/.config/deepcool-display-linux/" + filename + ".csv"
	print("creating file", filePath)
	file, err := os.Create(filePath)
	print("created file", filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	for _, row := range grid {
		csvRow := make([]string, len(row))
		for i, cell := range row {
			if cell {
				csvRow[i] = "1"
			} else {
				csvRow[i] = "0"
			}
		}
		if err := writer.Write(csvRow); err != nil {
			return fmt.Errorf("failed to write row to buffer: %w", err)
		}
	}
	writer.Flush()

	// Check for errors during CSV writing.
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error while flushing CSV writer: %w", err)
	}

	println(buffer.String())

	// Write the buffer content to the file instantly.
	// if err := os.WriteFile(filePath, buffer.Bytes(), 0644); err != nil {
	// 	return fmt.Errorf("failed to write buffer to file: %w", err)
	// }
	buffwriter := bufio.NewWriter(file)
	buffwriter.WriteString(buffer.String())
	buffwriter.Flush()

	file.Close()

	return nil
}

func (a *App) ParseCSV(filename string) ([][]bool, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %v", err)
	}

	grid := make([][]bool, len(records))
	for i, row := range records {
		grid[i] = make([]bool, len(row))
		for j, cell := range row {
			if cell == "1" {
				grid[i][j] = true
			} else {
				grid[i][j] = false
			}
		}
	}

	return grid, nil
}

func (a *App) GenerateHexCommand(grid [][]bool) (string, error) {
	header := []byte{0x10, 0x68, 0x01, 0x05, 0x1D, 0x01}

	rowValues := []byte{
		0x10, 0x10, // 1st and 2nd rows
		0x20, 0x20, // 3rd and 4th rows
		0x40, 0x40, // 5th and 6th rows
		0x80, 0x80, // 7th and 8th rows
		0x01, 0x01, // 9th and 10th rows
		0x02, 0x02, // 11th and 12th rows
		0x04, 0x04, // 13th and 14th rows
	}

	command := make([]byte, 68)

	copy(command[:len(header)], header)

	for col := 1; col <= 14; col++ {
		oddByte := 0
		evenByte := 0
		for row := 1; row <= 14; row++ {
			if grid[row-1][col-1] {
				if row%2 == 0 {
					evenByte += int(rowValues[row-1])
				} else {
					oddByte += int(rowValues[row-1])
				}
			}
		}
		command[len(header)+col-1] = byte(oddByte)
		command[len(header)+28-col] = byte(evenByte)
	}

	var checksum uint16
	for i := 1; i <= 33; i++ {
		checksum += uint16(command[i])
	}

	command[34] = byte(checksum % 256)
	command[35] = 22

	return hex.EncodeToString(command), nil
}
