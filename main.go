package main

import (
	"os"

	"github.com/AnshVM/flashpoll-backend/db"
	"github.com/AnshVM/flashpoll-backend/router"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	err := db.CreateConnection()

	if err != nil {
		panic(err)
	}

	router := router.SetupRouter()

	router.Run(os.Getenv("PORT"))
}
