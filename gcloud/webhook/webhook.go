package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Prints the POST body
func callback(resp http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println("Failed reading body")
		resp.WriteHeader(401)
		resp.Write([]byte(fmt.Sprintf(`{"success": false}`)))
		return
	}

	log.Println(string(body))

	resp.WriteHeader(200)
	resp.Write([]byte(fmt.Sprintf(`{"success": true}`)))
	return
}

// Starts a webserver
func webhook() {
	// Use PORT, as defined by Cloud Run specs
	basePort := os.Getenv("PORT")

	if len(basePort) == 0 {
		basePort = "8080"
	}

	ip := "0.0.0.0"

	port := fmt.Sprintf(":%s", basePort)
	log.Printf("Starting webhook on %s%s", ip, port)

	// Routing
	mux := mux.NewRouter()
	mux.SkipClean(true)

	mux.HandleFunc("/", callback).Methods("POST")

	handlers.LoggingHandler(os.Stdout, mux)
	loggedRouter := handlers.LoggingHandler(os.Stdout, mux)

	err := http.ListenAndServe(
		port,
		loggedRouter,
	)

	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

func main() {
	webhook()
}
