package main

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/disintegration/imaging"
)

func (app *application) generateBlurredImages() {
	inputPath := "./ui/static/img/temp.png"

	for i := 30; i >= 0; i-- {
		img, err := applyGaussianBlur(inputPath, float64(i)*1.5)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		hash := sha256.New()
		hash.Write([]byte(app.gameState.targetWord + fmt.Sprint(i)))
		name := hash.Sum(nil)

		nameString := strings.ReplaceAll(fmt.Sprintf("%x", name), "/", "_")

		// Save the blurred image to a temporary file
		tmpFile, err := os.Create(fmt.Sprintf("./ui/static/img/test/%s.png", nameString))
		if err != nil {
			fmt.Printf("unable to create temporary file: %v\n", err)
			return
		}
		defer tmpFile.Close()

		// Encode and save the blurred image as PNG
		err = png.Encode(tmpFile, img)
		if err != nil {
			fmt.Printf("unable to encode and save image: %v\n", err)
			return
		}

	}
}

// Function to read a PNG image, apply Gaussian blur, and send it to WebSocket clients
func applyGaussianBlur(inputPath string, radius float64) (*image.NRGBA, error) {
	// Open the PNG file
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("unable to open image file: %v", err)
	}
	defer file.Close()

	// Decode the PNG image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("unable to decode image: %v", err)
	}

	// Apply Gaussian blur to the image
	blurredImg := imaging.Blur(img, radius)

	return blurredImg, nil
}

// Function to copy initial image to a temporary file
func copyFile(inputPath string, outputPath string) {
	// Read the input PNG file
	inputData, err := os.ReadFile(inputPath)
	if err != nil {
		log.Fatalf("unable to read input file: %v", err)
	}

	err = os.WriteFile(outputPath, inputData, 0644)
	if err != nil {
		log.Fatalf("unable to write output file: %v", err)
	}

	log.Printf("Copied %s to %s", inputPath, outputPath)
}
