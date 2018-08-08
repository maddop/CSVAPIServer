package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

func startWebServerAndBlockForever() {

	// Set up the HTTP server:
	serverMUX := http.NewServeMux()
	serverMUX.HandleFunc("/list", handlerList)
	serverMUX.HandleFunc("/jsonreports", handlerJSON)

	server := &http.Server{}
	server.Addr = webBinding
	server.Handler = serverMUX
	server.SetKeepAlivesEnabled(true)
	server.ReadTimeout = 60 * 10 * time.Second  // 10 minutes
	server.WriteTimeout = 60 * 10 * time.Second // 10 minutes

	// Start the server:
	log.Printf("Starting HTTP server on '%s'.\n", webBinding)
	if errHTTP := server.ListenAndServe(); errHTTP != nil {
		log.Printf("ERROR: Failed to start HTTP server: %s\n", errHTTP.Error())
		os.Exit(2)
	}
}
