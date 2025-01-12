.PHONY: dev build clean

build:
	wails build

dev:
	wails dev

clean:
	rm -rf build/bin
	rm -rf frontend/dist
