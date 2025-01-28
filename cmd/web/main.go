package main

import (
	"log"
	"net/http"

	"pixblur.jkaisix/ui"
)

type templateData struct {
	Blanks   []string
	Keyboard [][]string
}
type application struct {
	gameState *GameState
}

type GameState struct {
	score      int
	stopGame   chan int
	targetWord string
}

func main() {
	inputPath := "ui/static/img/wizard.png" // Path to the input PNG file
	outputPath := "ui/static/img/base.png"  // Path to save the output PNG file

	app := &application{
		gameState: &GameState{
			stopGame:   make(chan int),
			targetWord: "WIZARD",
		},
	}

	// Forced false conditional for now to avoid redundant calls each time server starts
	// Turn to true when new blurred images are needed
	// Temporary until game admin functionality is setup
	if false {
		app.generateBlurredImages()
	}

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /keyboard", app.showKeyboard)
	mux.HandleFunc("/events", app.eventsHandler)
	mux.HandleFunc("POST /check_word", app.checkWord)

	log.Print("Serving on port 4000, http://localhost:4000")
	log.Printf("input: %s, output: %s", inputPath, outputPath)

	err := http.ListenAndServe("localhost:4000", mux)
	log.Fatal(err)
}
