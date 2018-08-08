package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// Filelist : the filename identified
type Filelist struct {
	Filename string `json:"filename"`
}

// FullFileList : The complete list of filenames identified
type FullFileList struct {
	FullResult []Filelist `json:"reports"`
}

func handlerList(response http.ResponseWriter, request *http.Request) {

	var myfilenames FullFileList
	var FileMatched int

	root := request.FormValue(`env`)
	recorddate := request.FormValue(`mydate`)

	d, err := time.Parse("20060102", recorddate)
	if err != nil {
		log.Printf("ERROR: invalid date specified! %s\n %s\n", d, err)
		return
	}

	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("INFO: date selected is '%+v", d)

	for i, file := range files {
		if strings.Contains(file.Name(), recorddate) {
			FileMatched++
			fmt.Println(file.Name())
			var result Filelist
			result.Filename = file.Name()
			myfilenames.FullResult = append(myfilenames.FullResult, result)
		} else if len(files) == i+1 && FileMatched < 1 {
			log.Println("ERROR: Unable to match date - no files found")
			return
		}
	}

	b, err := json.MarshalIndent(myfilenames, "", "\t")
	if err != nil {
		log.Printf("ERROR: Couldn't marshal JSON '%s'", myfilenames)
	}

	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response.Write(b)
}
