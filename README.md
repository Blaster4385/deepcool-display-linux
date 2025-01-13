# DeepCool Display Linux

This application is a replacement of the original DeepCool Windows application for the LP360 AIO cooler. I may add support for the entire LP series and any other new devices that use a similar pixel display. This supports drawing custom patterns on the as well as displaying the CPU temperature and usage.

Special thanks to [@Nortank12](https://github.com/Nortank12) for his work on [deepcool-digital-linux](https://github.com/Nortank12/deepcool-digital-linux). I would recommend checking out his app for additional functionality and support for other devices. Additionally, thanks to [@rohan09-raj](https://github.com/rohan09-raj) for figuring out the logic of the commands for creating the patterns.

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
5. Download the latest [release](https://github.com/Blaster4385/deepcool-display-linux/releases/latest) and run it.

## Usage

You can run the applications with or without providing any options. Running it without any options will launch the GUI.

```bash
./deepcool-display-linux [OPTIONS] 
```

```
Options:
  -d, --daemon       Run the application in daemon mode
  -m, --mode         Specify the mode (temp, usage or custom)
  -f, --file         Specify the CSV file containing the pattern data (This is required in daemon mode)
  -c, --celcius      Display the CPU temperature in celcius

Commands:
  -h, --help         Print help
  -v, --version      Print version
```

### Daemon Mode
Run the application in daemon mode to display the pattern from a CSV file:

```bash
./deepcool-display-linux -d -m custom -f /path/to/pattern.csv
```

In custom mode, the \`-f\` or \`--file\` flag is required to specify the CSV file containing the pattern.

### Exporting patterns to CSV

The GUI has an option to export the current pattern to a CSV file. This can be done by clicking the "Export Layout" button. The CSV files are stored in ~/.config/deepcool-display-linux.

## Automatic Startup using systemd

1. Copy the binary to /usr/bin:

```bash
sudo cp deepcool-display-linux /usr/bin
```

2. Create a service file:

```bash
sudo nano /etc/systemd/system/deepcool-display-linux.service
```

3. Add the following to the file:

```bash
[Unit]
Description=DeepCool Display Linux

[Service]
ExecStart=/usr/bin/deepcool-display-linux -d -m temp -c
Restart=always

[Install]
WantedBy=multi-user.target
```

4. Enable the service:

```bash
sudo systemctl enable --now deepcool-display-linux.service
```

 **Note:** The application will automatically start when the system is booted.

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
