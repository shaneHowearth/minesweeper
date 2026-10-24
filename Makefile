# Makefile for Minesweeper Android Deployment

APP_NAME = "Minesweeper"
APP_ID = com.shanehowearth.minesweeper
PACKAGE_DIR = ./cmd

.PHONY: all clean deps android release

all: android

# Install the Fyne CLI tool required for packaging
deps:
	go install fyne.io/tools/cmd/fyne@latest

# Build a debug/testing APK for Android devices
android:
	fyne-cross android -app-id $(PACKAGE_ID) -icon cmd/fyne/Icon.png $(SOURCE_DIR)

# Build a release-ready package (App Bundle/APK) for Google Play distribution
release:
	fyne release -os android -app-id $(APP_ID) -name $(APP_NAME) $(PACKAGE_DIR)

# Clean up generated package artifacts
clean:
	rm -f *.apk *.aab
