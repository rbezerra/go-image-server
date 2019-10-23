package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"./handlers"

	"./db"
)

func setupRoutes() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/upload", handlers.UploadFile)
	router.HandleFunc("/imagens", handlers.ListImages)
	router.HandleFunc("/imagem/{uuid}/{tamanho}", handlers.GetImage)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	fmt.Println("Hello World!!!")
	db.InitDB("postgres://user:pass@localhost/image_server?sslmode=disable")
	setupRoutes()
}
