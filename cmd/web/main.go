package main

import (
	"log"
	"net/http"
)

type application struct {
	gameState *GameState
	keyboard  Keyboard
}

type GameState struct {
	score      int
	stopGame   chan int
	targetWord string
}

type Keyboard struct {
	Row1 []string
	Row2 []string
	Row3 []string
}

func main() {
	inputPath := "ui/static/img/wizard.png" // Path to the input PNG file
	outputPath := "ui/static/img/base.png"  // Path to save the output PNG file

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

	app := &application{
		gameState: &GameState{
			stopGame:   make(chan int),
			targetWord: "WIZARD",
		},
		keyboard: keyboard,
	}

	app.generateBlurredImages()

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /keyboard", app.showKeyboard)
	mux.HandleFunc("/events", app.eventsHandler)
	mux.HandleFunc("POST /check_word", app.checkWord)

	log.Print("Serving on port 4000, http://localhost:4000")
	log.Printf("input: %s, output: %s", inputPath, outputPath)

	err := http.ListenAndServe("localhost:4000", mux)
	log.Fatal(err)
}
