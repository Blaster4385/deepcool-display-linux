package modules

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	tempSensorCacheDuration = 24 * time.Hour
	cpuUsageSampleInterval  = 200 * time.Millisecond
	tempMilliCelsiusDivisor = 1000.0
	fahrenheitConversion    = 9.0 / 5.0
	fahrenheitBase          = 32.0
)

type CPUUsage struct {
	User    int64
	Nice    int64
	System  int64
	Idle    int64
	IOWait  int64
	IRQ     int64
	SoftIRQ int64
	Steal   int64
}

var (
	cachedTemp         float64
	lastTempUpdate     time.Time
	cachedTempSensor   string
	tempSensorCachedAt time.Time
)

func GetCPUTemperature(fahrenheit bool) (float64, error) {
	now := time.Now()
	if now.Sub(lastTempUpdate) < time.Second {
		return cachedTemp, nil
	}

	tempSensorPath, err := findTempSensor()
	if err != nil {
		return 0, fmt.Errorf("finding temp sensor: %w", err)
	}

	data, err := ioutil.ReadFile(tempSensorPath)
	if err != nil {
		return 0, fmt.Errorf("reading CPU temperature (%s): %w", tempSensorPath, err)
	}

	tempMilliCelsius, err := strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing CPU temperature: %w", err)
	}

	tempCelsius := float64(tempMilliCelsius) / tempMilliCelsiusDivisor
	if fahrenheit {
		tempCelsius = (tempCelsius * fahrenheitConversion) + fahrenheitBase
	}

	cachedTemp = tempCelsius
	lastTempUpdate = now
	return tempCelsius, nil
}

func GetCPUUsage() (float64, error) {
	prevUsage, err := readCPUUsage()
	if err != nil {
		return 0, fmt.Errorf("reading initial CPU usage: %w", err)
	}

	time.Sleep(cpuUsageSampleInterval)

	currUsage, err := readCPUUsage()
	if err != nil {
		return 0, fmt.Errorf("reading current CPU usage: %w", err)
	}

	usage := calculateCPUUsage(prevUsage, currUsage)
	return usage, nil
}

func findTempSensor() (string, error) {
	if cachedTempSensor != "" && time.Since(tempSensorCachedAt) < tempSensorCacheDuration {
		return cachedTempSensor, nil
	}

	hwmonPath := "/sys/class/hwmon"
	files, err := ioutil.ReadDir(hwmonPath)
	if err != nil {
		return "", fmt.Errorf("locating CPU temperature sensor directory: %w", err)
	}

	for _, file := range files {
		sensorPath := fmt.Sprintf("%s/%s", hwmonPath, file.Name())
		nameFilePath := fmt.Sprintf("%s/name", sensorPath)
		nameData, err := ioutil.ReadFile(nameFilePath)
		if err != nil {
			continue
		}
		name := strings.TrimSpace(string(nameData))
		if name == "coretemp" || name == "k10temp" || name == "zenpower" {
			cachedTempSensor = fmt.Sprintf("%s/temp1_input", sensorPath)
			tempSensorCachedAt = time.Now()
			return cachedTempSensor, nil
		}
	}

	return "", errors.New("appropriate CPU temperature sensor not found")
}

func readCPUUsage() (CPUUsage, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return CPUUsage{}, fmt.Errorf("reading /proc/stat: %w", err)
	}
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "cpu" {
			usage := CPUUsage{
				User:    parseInt64(fields[1]),
				Nice:    parseInt64(fields[2]),
				System:  parseInt64(fields[3]),
				Idle:    parseInt64(fields[4]),
				IOWait:  parseInt64(fields[5]),
				IRQ:     parseInt64(fields[6]),
				SoftIRQ: parseInt64(fields[7]),
			}
			return usage, nil
		}
	}

	return CPUUsage{}, errors.New("failed to parse CPU usage from /proc/stat")
}

func calculateCPUUsage(prev, curr CPUUsage) float64 {
	prevTotal := prev.User + prev.Nice + prev.System + prev.Idle + prev.IOWait + prev.IRQ + prev.SoftIRQ
	currTotal := curr.User + curr.Nice + curr.System + curr.Idle + curr.IOWait + curr.IRQ + curr.SoftIRQ

	totalDiff := float64(currTotal - prevTotal)
	idleDiff := float64(curr.Idle - prev.Idle)

	if totalDiff == 0 {
		return 0
	}

	return (totalDiff - idleDiff) / totalDiff * 100
}

func parseInt64(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}
