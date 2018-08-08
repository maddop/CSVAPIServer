package main

import (
	"flag"
)

func readFlags() {
	flag.StringVar(&webBinding, `webBinding`, `:8080`, `How to bind the web server, i.e. interface and port.`)
	flag.StringVar(&dataDirectory, `dataDirectory`, `data`, `Where to store the data?`)
	flag.Parse()
}
