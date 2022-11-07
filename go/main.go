package main

import (
	"log"
	"net/http"
	"os"
	"fmt"
	"strconv"

	"github.com/gorilla/mux"
)

const STATIC_DIR string = "/docker/cache"

var baseUrl string
var heatMins int

func main() {
	log.Println("Application started")
	baseUrl = os.Getenv("PROXY_HOST")
	if baseUrl == "" {
		log.Fatalln("environment PROXY_HOST missing")
	}
	heatMins, _ = strconv.Atoi(os.Getenv("HEAT_TIME_MINUTES"))
	if heatMins < 0 {
		heatMins = 15
	}
	log.Printf("Heat Time is set to %d mins", heatMins)

	runApi(8080)
}

func runApi(port int) {
	log.Println("Initializing API")
	router := mux.NewRouter()

	log.Println("Adding Logging")
	router.Use(loggingMiddleware)
	router.Use(cachingMiddleware)

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(STATIC_DIR))))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

