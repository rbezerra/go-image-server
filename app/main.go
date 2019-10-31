package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"

	"./handlers"

	c "./config"
	"./db"
)

func setupRoutes() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/upload", handlers.UploadFile).Methods("POST")
	router.HandleFunc("/imagens", handlers.ListImages).Methods("GET")
	router.HandleFunc("/imagem-info/{uuid}/{tamanho}", handlers.GetImageInfo).Methods("GET")
	router.HandleFunc("/imagem-info/{uuid}/", handlers.GetImageInfo).Methods("GET")
	router.HandleFunc("/imagem/{uuid}/{tamanho}", handlers.GetImage).Methods("GET")
	router.HandleFunc("/imagem/{uuid}/", handlers.GetImage).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	var configuration c.Configurations

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err := viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	datasource := "postgres://" +
		viper.GetString("database.dbuser") + ":" +
		viper.GetString("database.password") + "@" +
		viper.GetString("database.hostname") + ":" +
		viper.GetString("database.port") + "/" +
		viper.GetString("database.dbname") + "?sslmode=" +
		viper.GetString("database.sslmode")

	db.InitDB(datasource)
	fmt.Println("Ready to work")
	setupRoutes()
}
