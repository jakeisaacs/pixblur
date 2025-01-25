package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/disintegration/imaging"
)

type application struct {
	gameState *GameState
}

type GameState struct {
	score    int
	stopGame chan int
}

type Keyboard struct {
	Row1    []string
	Row2    []string
	Row3    []string
	Word    string
	IsWrong bool
	Message string
}

// Create a struct to pass to the template
var (
	keyboard = Keyboard{
		Row1: []string{"Q", "W", "E", "R", "T", "Y", "U", "I", "O"},
		Row2: []string{"Z", "A", "S", "D", "F", "G", "H", "J", "K"},
		Row3: []string{"Z", "X", "C", "V", "B", "N", "M"},
	}

	targetWord = "WIZARD"
)

func (app *application) checkWord(w http.ResponseWriter, r *http.Request) {
	var (
		status string
		result map[string]interface{}
	)
	err := json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println(result)

	if result["word"] == targetWord {
		app.gameState.stopGame <- 1
		status = "success"
	} else {
		status = "fail"
	}

	response := map[string]interface{}{
		"status": status,
	}

	w.Header().Set("Content-Type", "GameState/json")
	json.NewEncoder(w).Encode(response)
}

// Function to read a PNG image, apply Gaussian blur, and send it to WebSocket clients
func applyGaussianBlurAndSendImage(inputPath string, radius float64) error {
	// Open the PNG file
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("unable to open image file: %v", err)
	}
	defer file.Close()

	// Decode the PNG image
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("unable to decode image: %v", err)
	}

	// Apply Gaussian blur to the image
	blurredImg := imaging.Blur(img, radius)

	// Save the blurred image to a temporary file
	tmpFile, err := os.Create("ui/static/img/temp_blurred_image.png")
	if err != nil {
		return fmt.Errorf("unable to create temporary file: %v", err)
	}
	defer tmpFile.Close()

	// Encode and save the blurred image as PNG
	err = png.Encode(tmpFile, blurredImg)
	if err != nil {
		return fmt.Errorf("unable to encode and save image: %v", err)
	}

	return nil
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

// WebSocket handler
func (app *application) eventsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	inputPath := "ui/static/img/temp.png" // Path to the input PNG file

	rc := http.NewResponseController(w)

	startTime := time.Now()
	duration := 31 * time.Second // Total duration
	interval := 1 * time.Second  // Time between updates

	// Simulate image changes and push them to the client
	for blur := 30.0; blur >= 0.0; blur -= 1.0 {
		elapsed := time.Since(startTime)
		if elapsed >= duration {
			break
		}
		select {
		case <-app.gameState.stopGame:
			app.gameState.score = int(elapsed) * 512
			fmt.Println("Stopping game...")
			return
		default:
			// Apply Gaussian blur and send image update every 5 seconds
			err := applyGaussianBlurAndSendImage(inputPath, (blur * 1.5))
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			_, err = fmt.Fprintf(w, "data: %f\n\n", (blur))

			if err != nil {
				return
			}
			err = rc.Flush()
			if err != nil {
				return
			}

			// Sleep for 1 seconds before sending the next update
			sleepTime := interval - time.Since(startTime)%interval
			time.Sleep(sleepTime)
		}
	}

	ctx := r.Context()
	<-ctx.Done()

	fmt.Println("Finished sending image updates")
}

func showKeyboard(w http.ResponseWriter, r *http.Request) {

	ts, err := template.ParseFiles("./ui/html/keyboard.html")
	if err != nil {
		fmt.Printf("Failed to parse... Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "keyboard", keyboard)
	if err != nil {
		fmt.Printf("Failed to execute... Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/keyboard.html",
	}

	// Define the 18 characters (this is just an example, you can change them)
	row1 := []string{"Q", "W", "E", "R", "T", "Y", "U", "I", "O"}
	row2 := []string{"A", "S", "D", "F", "G", "H", "J", "K", "L"}
	row3 := []string{"Z", "X", "C", "V", "B", "N", "M"}

	// Create a struct to pass to the template
	keyboard := Keyboard{
		Row1: row1,
		Row2: row2,
		Row3: row3,
	}

	copyFile("ui/static/img/wizard.png", "ui/static/img/temp.png")

	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Printf("Failed to parse... Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", keyboard)
	if err != nil {
		fmt.Printf("Failed to execute... Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func main() {
	inputPath := "ui/static/img/wizard.png" // Path to the input PNG file
	outputPath := "ui/static/img/base.png"  // Path to save the output PNG file

	app := &application{
		gameState: &GameState{
			stopGame: make(chan int),
		},
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /keyboard", showKeyboard)
	mux.HandleFunc("/events", app.eventsHandler)
	mux.HandleFunc("POST /check_word", app.checkWord)

	log.Print("Serving on port 4000, http://localhost:4000")
	log.Printf("input: %s, output: %s", inputPath, outputPath)

	err := http.ListenAndServe("localhost:4000", mux)
	log.Fatal(err)
}
