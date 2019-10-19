package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
)

const uploadPath = "./temp-images"

func uploadFile() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("File Upload Endpoint Hit")

		r.ParseMultipartForm(10 << 20)

		file, handler, err := r.FormFile("image")

		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			fmt.Println("INVALID_FILE")
			fmt.Println(err)
			return
		}
		defer file.Close()
		fmt.Printf("Uploaded File: %+v\n", handler.Filename)
		fmt.Printf("File Size: %+v\n", handler.Size)
		fmt.Printf("MIME Header: %+v\n", handler.Header)

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusInternalServerError)
			fmt.Println("INVALID_FILE")
			fmt.Println(err)
		}

		filetype := http.DetectContentType(fileBytes)
		fileName := createUUIDFileName(handler.Filename)
		switch filetype {
		case "image/jpeg", "image/jpg", "image/gif", "image/png":
		default:
			renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			fmt.Println("INVALID_FILE_TYPE")
			fmt.Println(err)
			return
		}

		fileEndings, err := mime.ExtensionsByType(filetype)
		if err != nil {
			renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			fmt.Println("CANT_READ_FILE_TYPE")
			fmt.Println(err)
			return
		}

		newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
		fmt.Printf("FileType: %s, File: %s\n", filetype, newPath)

		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			fmt.Println("CANT_WRITE_FILE")
			fmt.Println(err)
			return
		}
		defer newFile.Close()

		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			fmt.Println("CANT_WRITE_FILE")
			fmt.Println(err)
			return
		}
		w.Write([]byte("SUCCESS"))
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func createUUIDFileName(fileName string) string {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return fmt.Sprintf("img-%s", uuid)
}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile())
	http.ListenAndServe(":8080", nil)
}

func main() {
	fmt.Println("Hello World!!!")
	setupRoutes()
}
