package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
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

	if result["word"] == app.gameState.targetWord {
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

func (app *application) eventsHandler(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// inputPath := "ui/static/img/temp.png" // Path to the input PNG file

	rc := http.NewResponseController(w)

	startTime := time.Now()
	duration := 31 * time.Second // Total duration
	interval := 1 * time.Second  // Time between updates

	// Simulate image changes and push them to the client
	for i := 30; i >= 0; i -= 1 {
		elapsed := time.Since(startTime)
		fmt.Printf("elapsed %s:\n", elapsed)
		if elapsed >= duration {
			break
		}
		select {
		case <-app.gameState.stopGame:
			app.gameState.score = i * 512
			fmt.Printf("Stopping game... Score: %d\n", app.gameState.score)
			return
		default:
			// Apply Gaussian blur and send image update every 5 seconds
			hash := sha256.New()
			hash.Write([]byte(app.gameState.targetWord + fmt.Sprint(i)))
			name := hash.Sum(nil)

			nameString := strings.ReplaceAll(fmt.Sprintf("%x", name), "/", "_")
			imgPath := fmt.Sprintf("/static/img/test/%s.png", nameString)

			_, err := fmt.Fprintf(w, "data: %d\ndata: %s\n\n", i, imgPath)

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

func (app *application) showKeyboard(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)
	data.Blanks = make([]string, len(app.gameState.targetWord))

	ts, err := template.ParseFiles("./ui/html/keyboard.html")
	if err != nil {
		fmt.Printf("Failed to parse... Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "keyboard", data)
	if err != nil {
		fmt.Printf("Failed to execute... Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"./ui/html/base.html",
		"./ui/html/keyboard.html",
	}

	copyFile("ui/static/img/wizard.png", "ui/static/img/temp.png")

	data := app.newTemplateData(r)
	data.Blanks = make([]string, len(app.gameState.targetWord))

	ts, err := template.ParseFiles(files...)
	if err != nil {
		fmt.Printf("Failed to parse... Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Printf("Failed to execute... Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
