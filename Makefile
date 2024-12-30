.PHONY: run build clean wire

APP_NAME = root
BUILD_DIR = $(PWD)/build
MAIN_FILE = application/cmd/main.go
APP_FILE = application/app.go

setup:
	go mod tidy

clean:
	rm -rf $(BUILD_DIR)

wire:
	wire $(APP_FILE)

build:
	go build -ldflags="-s -w" -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

alpha:
	make wire
	make build
	APP_PHASE=alpha $(BUILD_DIR)/$(APP_NAME)

prod:
	make wire
	make build
	APP_PHASE=prod $(BUILD_DIR)/$(APP_NAME)
