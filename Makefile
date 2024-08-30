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

local:
	make wire
	make build
	APP_PHASE=local $(BUILD_DIR)/$(APP_NAME)
