package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var wg sync.WaitGroup

func env() {
	err := godotenv.Overload(".env", ".env.local")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	env()

	jobs := run()
	enc := json.NewEncoder(os.Stdout)

	enc.Encode(jobs)
}
