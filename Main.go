package main

import (
	"log"
	"os"
	"runtime"
)

func main() {
	log.Printf("API Host, %s\n", os.Getenv("APISERVERVER"))
	readFlags()
	log.Printf("Configuration: webServer='%s', dataDir='%s'\n", webBinding, dataDirectory)

	// Go should use all CPUs:
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Start the web server:
	startWebServerAndBlockForever()
}
