package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fahimanzamdip/go-invoice-api/router"
	"github.com/joho/godotenv"
)

func init() {
	log.Println("<========= Go-Invoice-Api Starting =========>")
	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found, reading configuration from ENV")
		return
	}
	log.Println("Reading .env file successful || Alhamdulillaah")

	directoryPaths := []string{"./public/uploads", "./public/pdfs"}

	for _, directoryPath := range directoryPaths {
		_, err := os.Stat(directoryPath)
		if os.IsNotExist(err) {
			err = os.MkdirAll(directoryPath, os.ModePerm)
			if err != nil {
				log.Println("Error creating directory:", err)
				return
			}
			log.Println("Directory created:", directoryPath)
		} else if err != nil {
			log.Println("Error checking directory:", err)
			return
		} else {
			log.Println("Directory exists:", directoryPath)
		}
	}
}

func main() {
	r := router.Configure()

	port := os.Getenv("server_port")
	if port == "" {
		port = "8000"
	}

	log.Println("API URL: http://localhost/" + os.Getenv("api_uri"))
	log.Fatal(http.ListenAndServe(":"+port, r))
}
