package main

import (
	"crypto/sha256"
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/disintegration/imaging"
)

func (app *application) generateBlurredImage(inputFile, outputPath string, i int, radius float64, ch chan string) {

	//Generate blurred image
	img, err := applyGaussianBlur(inputFile, float64(i)*radius)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	hash := sha256.New()
	hash.Write([]byte(app.gameState.targetWord + fmt.Sprint(i)))
	name := hash.Sum(nil)

	nameString := strings.ReplaceAll(fmt.Sprintf("%x", name), "/", "_")
	outputFile := fmt.Sprintf("%s/%s.png", outputPath, nameString)

	// Save the blurred image to a temporary file
	tmpFile, err := os.Create(outputFile)
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
	ch <- outputFile
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

func (app *application) newTemplateData(r *http.Request) templateData {
	// Create a struct to pass to the template
	keyboard := [][]string{
		{"Q", "W", "E", "R", "T", "Y", "U", "I", "O"},
		{"A", "S", "D", "F", "G", "H", "J", "K", "L"},
		{"Z", "X", "C", "V", "B", "N", "M"},
	}

	return templateData{
		Keyboard: keyboard,
	}
}
