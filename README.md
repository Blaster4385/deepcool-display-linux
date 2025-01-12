# DeepCool Display Linux Controller

This application is a replacement of the original DeepCool Windows application for the LP360 AIO cooler. I may add support for the entire LP series and any other new devices that use a similar pixel display. This currently only supports drawing custom patterns on the display. Support for displaying CPU temprature and usage will be added in future releases.

Special thanks to @Nortank12 for his work on [deepcool-digital-linux](https://github.com/Nortank12/deepcool-digital-linux). I would recommend checking out his app for additional functionality and support for other devices.

## Installation

1. Locate your directory, it can be /lib/udev/rules.d or /etc/udev/rules.d and create a new file named 99-deepcool.rules.
 ```bash
 sudo nano /etc/udev/rules.d/99-deepcool.rules
 ```
2. Insert the following:
 ```bash
SUBSYSTEM=="usb", ATTRS{idVendor}=="3633", ATTRS{idProduct}=="000c", GROUP="plugdev"
 ```
3. Ensure that your user is part of the plugdev group:
 ```bash
 sudo usermod -aG plugdev $USER
 ```
4. Reboot your system.
5. Download the latest release and run it.

## Usage

You can run the applications with or without providing any options. Running it without any options will launch the GUI.

```bash
./deepcool-display-linux [OPTIONS] 
```

```
Options:
  -d, --daemon       Run the application in daemon mode
  -f, --file         Specify the CSV file containing the pattern data (This is required in daemon mode)

Commands:
  -h, --help         Print help
  -v, --version      Print version
```
### Daemon Mode
Run the application in daemon mode to display the pattern from a CSV file:

```bash
./deepcool-display-linux -d -f /path/to/pattern.csv
```

In daemon mode, the \`-f\` or \`--file\` flag is required to specify the CSV file containing the pattern.

### Exporting patterns to CSV

The GUI has an option to export the current pattern to a CSV file. This can be done by clicking the "Export Layout" button. The CSV files are stored in ~/.config/deepcool-display-linux.

## Development

This application is built using [Wails](https://wails.io).

You can build the application from source by following the steps below.

### Dependencies

You need to install go, npm, and wails.

On Arch Linux, you can use a AUR helper to install them:

```bash
yay -S go npm wails
```
### Building

1. Clone the repository:

```bash
git clone https://github.com/Blaster4385/deepcool-display-linux
```

2. Open the directory:

```bash
cd deepcool-display-linux
```
3. Run a development server:

```bash
make dev
```
4. Build a release:

```bash
make build
```
5. Clean up:

```bash
make clean
```
## License

This project is licensed under the GPLv2 License - see the [LICENSE](LICENSE) file for details.
