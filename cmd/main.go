package main

import (
	"fmt"
	"forum/database"
	"forum/internal"
	"log"
	"net/http"
)

func main() {
	database.InitDB()
	_, err := internal.LoadConfig("config.json")
	if err != nil {
		log.Println(err)
	}
	h := routes()

	fmt.Println("Server is running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", h); err != nil {
		log.Println(err)
	}
}
