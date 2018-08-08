package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sys/windows"
)

// PodStats : details of each
type PodStats struct {
	ServerName string
	DateTime   string
	Project    string
	PodName    string
	CPUReq     string
	CPUPercent string
	CPULim     string
	MemReq     string
	MemPercent string
	MemLimit   string
	CPUUsage   string
}

//var Matched int
var root string
var recorddate string

// Level : used for capturing loglevel
var Level string

func init() {
	//log.SetFormatter(&log.JSONFormatter{})
	Formatter := new(log.TextFormatter)
	Formatter.TimestampFormat = "01-02-2006 15:04:05"
	Formatter.FullTimestamp = true
	log.SetFormatter(Formatter)
}

func loglevel(Level string) {
	if (Level == "Debug") || (Level == "debug") {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func handlerJSON(response http.ResponseWriter, request *http.Request) {

	//detect if Windows, if so then enable colours for logging!
	//https://github.com/sirupsen/logrus/issues/172
	if runtime.GOOS == "windows" {

		var originalMode uint32
		stdout := windows.Handle(os.Stdout.Fd())
		windows.GetConsoleMode(stdout, &originalMode)
		windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
		defer windows.SetConsoleMode(stdout, originalMode)
	}

	//pull in the values provided in the uri
	root := request.FormValue(`root`)
	recorddate := request.FormValue(`recorddate`)
	loglevel(request.FormValue(`log`))

	var pods []PodStats

	log.WithFields(log.Fields{
		"root":       root,
		"recorddate": recorddate,
	}).Debug("Captured user input")

	//if not present, add the required slash for the path and fix for the OS
	if strings.Contains(root, "/") || strings.Contains(root, "\\") {
		root = filepath.FromSlash(root)
	} else {
		root = root + "/"
		root = filepath.FromSlash(root)
		log.WithFields(log.Fields{
			"root": root,
		}).Debug("Changed root to suit OS")
	}

	//Verify the date format is correct
	d, err := time.Parse("20060102", recorddate)
	if err != nil {
		log.Printf("ERROR: invalid date specified! %s\n %s\n", d, err)
		return
	}

	log.WithFields(log.Fields{
		"recorddate": d,
	}).Debug("Parsed time")

	//Pull in the list of files within the directory specified
	files, err := ioutil.ReadDir(root)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"count": len(files),
	}).Info("Total number of files to check")

	Matched := 0

	//parse each of the files
	for i, file := range files {

		//Check if filename contains the record date
		if strings.Contains(file.Name(), recorddate) {
			Matched++

			log.WithFields(log.Fields{
				"filename": file.Name(),
			}).Info("Successful match")

			//Pull back csv data from file
			csvData := reportcontents(root + file.Name())

			//call the Split function below to get the hostname
			PodServerName := strings.FieldsFunc(file.Name(), Split)
			if err != nil {
				fmt.Println(err)
			}

			//populate the fields of the struct pods from csvData
			for _, each := range csvData {
				var pod PodStats
				pod.ServerName = PodServerName[0]
				pod.DateTime = each[0] + " " + each[1]
				pod.Project = each[2]
				pod.PodName = each[3]
				pod.CPUReq = each[4]
				pod.CPUPercent = each[5]
				pod.CPULim = each[6]
				pod.MemReq = each[7]
				pod.MemPercent = each[8]
				pod.MemLimit = each[9]
				pod.CPUUsage = each[10]
				pods = append(pods, pod)
			}

		} else if len(files) == i+1 && Matched < 1 {
			response.Header().Set("Content-Type", "text/html; charset=UTF-8")
			response.Write([]byte("ERROR: Unable to match date - no files found"))

			log.WithFields(log.Fields{
				"date": recorddate,
			}).Error("No files found for requested date")

			return
		}

	}

	log.WithFields(log.Fields{
		"count": Matched,
	}).Debug("Number of files matched")

	//Create JSON output from the pods
	jsonData, err := json.MarshalIndent(pods, "", "\t")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//Output data to http
	response.Header().Set("Content-Type", "application/json; charset=UTF-8")
	response.Write(jsonData)

	log.WithFields(log.Fields{
		"jsonData": string(jsonData),
	}).Debug("JSON output")

}

// Split : a string using any of the delimeters below
func Split(r rune) bool {
	return r == '.' || r == '\\' || r == '/'
}

//parse the file as csv (use reader.Comma to change delimeter from comma)
func reportcontents(source string) [][]string {
	sourcefilename := source
	csvFile, err := os.Open(sourcefilename)
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.Comma = ' '

	csvData, err := reader.ReadAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return csvData

}
