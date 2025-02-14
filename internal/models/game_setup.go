// CHANGE TO "models" ONCE BETTER CONFIGURED
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

type GameSetup struct {
	Name       string
	Data       string
	Date_Valid string
}

// CHANGE FUNCTION NAME ONCE BETTER CONFIGURED
func main() {
	name := "wizard"
	valid_date := time.Now().Format("01/02/2006")

	db, err := sql.Open("sqlite", "./sqlite.db")
	if err != nil {
		fmt.Errorf("Unable to open db: %v", err)
	}
	defer db.Close()

	insertImage := func() error {
		imagePath := "./ui/static/img/wizard.png"
		imageData, err := os.ReadFile(imagePath)

		// fmt.Printf("name: %s, valid_date: %s\nData: ", name, valid_date, imageData)
		if err != nil {
			return fmt.Errorf("Error reading image: %v", err)
		}

		stmt, err := db.Prepare(`INSERT INTO game_info
						(name, data, valid_date)
						VALUES (?, ?, ?)`)
		if err != nil {
			return fmt.Errorf("Error preparing statement: %v", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, imageData, valid_date)
		if err != nil {
			return fmt.Errorf("Error inserting image: %v", err)
		}

		return nil
	}

	err = insertImage()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully inserted image!")

}
