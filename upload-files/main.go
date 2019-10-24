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
	router.HandleFunc("/upload", handlers.UploadFile).Methods("POST")
	router.HandleFunc("/imagens", handlers.ListImages).Methods("GET")
	router.HandleFunc("/imagem/{uuid}/{tamanho}", handlers.GetImage).Methods("GET")
	router.HandleFunc("/imagem/{uuid}/", handlers.GetImage).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	fmt.Println("Hello World!!!")
	db.InitDB("postgres://user:pass@localhost/image_server?sslmode=disable")
	setupRoutes()
}
