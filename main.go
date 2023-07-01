package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fahimanzamdip/go-invoice-api/routes"
	"github.com/joho/godotenv"
)

func init() {
	log.Println("Go-Shop-Api Starting.....")
	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found, reading configuration from ENV")
		return
	}
	log.Println("Reading .env file successful || Alhamdulillaah")
}

func main() {
	r := routes.Configure()

	port := os.Getenv("server_port")
	if port == "" {
		port = "8000"
	}

	log.Println("API URL: http://localhost/" + os.Getenv("api_uri"))
	log.Fatal(http.ListenAndServe(":"+port, r))
}
