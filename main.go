package main

import (
	"image"
	"log"
	"log/slog"
	"os"
)

func main() {
	// Load image
	file, err := os.Open("")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	im, format, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
		return
	}
	slog.Info("Image loaded successfully", "format", format)

}
