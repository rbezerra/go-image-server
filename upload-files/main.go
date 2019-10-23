package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"./db"
)

const uploadPath = "./temp-images"

func setupRoutes() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/upload", UploadFile)
	router.HandleFunc("/imagens", ListImages)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {
	fmt.Println("Hello World!!!")
	db.InitDB("postgres://user:pass@localhost/image_server?sslmode=disable")
	setupRoutes()
}
