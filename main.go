package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/pborman/getopt"
)

var (
	csvFile     = getopt.StringLong("input", 'i', "", "The input CSV file")
	outFile     = getopt.StringLong("out", 'o', "", "The output CSV file (default: STDOUT)")
	webUrl      = getopt.StringLong("web", 'w', "", "The url to fetch from")
	search      = getopt.StringLong("search", 's', "", "The search type: twitter, email")
	workers     = getopt.IntLong("workers-count", 'c', 50, "The number of wokers")
	printErrors = getopt.BoolLong("errors", 'e', "Print the errors")
	showVersion = getopt.BoolLong("version", 'v', "Prints the version")
	checkUpdate = getopt.BoolLong("update", 'u', "Update the version")
	help        = getopt.BoolLong("help", 'h', "Show this message")

	timeout    = time.Duration(30 * time.Second)
	httpClient = http.Client{
		Timeout: timeout,
	}
	jobsChannel = make(chan string)
	jobsWorking = make(map[int]int)
	mutex       = &sync.Mutex{}
	searchRe    = map[string]*regexp.Regexp{
		"twitter": regexp.MustCompile(`https://twitter.com/(\w+)`),
		//"email":   regexp.MustCompile(`([^@\s][!#$%&'*+-\/=?^{|}~\w]+@[^@\s][!#$%&'*+-\/=?^{|}~\w]+)`),
		//"email": regexp.MustCompile(`([^@\s][!#$%&'*+-\/=?^{|}~\w]+@[^@\s](?:(?:[^.][A-Za-z\d\-]+)|(?:\[(?:\d{1,3}?\.){3}\d{1,3}\])))`),
		//"email": regexp.MustCompile(`([\w+\-]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+)`),
		"email": regexp.MustCompile(`((?:[a-z0-9!#$%&'*+/=?^_{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\]))`),
	}
	validSearch = regexp.MustCompile(`twitter|email`)
)

func checkError(msg string, err error) {
	if err != nil {
		log.Fatal(msg, err)
	}
}

func main() {
	getopt.Parse()

	if *help {
		getopt.PrintUsage(os.Stdout)
		return
	}

	if *showVersion {
		printVersion()
		return
	}

	if *checkUpdate {
		err := doUpdate()
		checkError("Error updating: ", err)
		return
	}

	err := validateRequiredOpts()

	checkError("", err)

	pwd, err := os.Getwd()
	checkError("Error getting the pwd:", err)

	// Open input File
	f, err := os.Open(path.Join(pwd, *csvFile))
	checkError("Error opening input file", err)

	out := getOutFile(pwd)
	defer out.Close()

	r := csv.NewReader(bufio.NewReader(f))
	w := csv.NewWriter(out)

	initializeWorkers(w)

	if len(*webUrl) == 0 {
		startReading(r)
	} else {
		jobsChannel <- *webUrl
	}

	time.Sleep(1 * time.Second)

	for len(jobsWorking) != 0 {
		time.Sleep(5 * time.Second)
	}

	close(jobsChannel)

	w.Flush()
}

func fetchAndWrite(jobs <-chan string, w *csv.Writer, id int) {
	for u := range jobs {

		mutex.Lock()
		jobsWorking[id] = 1
		mutex.Unlock()

		row := []string{u}
		handlers, err := extractInformation(u)
		if err != nil && *printErrors {
			row = append(row, err.Error())
		} else {
			row = append(row, handlers...)
		}
		if err := w.Write(row); err != nil {
			log.Fatalln("Error writing record to csv:", err, row)
		}

		mutex.Lock()
		delete(jobsWorking, id)
		mutex.Unlock()

	}
}

func extractInformation(u string) ([]string, error) {
	urlObject, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	if len(urlObject.Scheme) == 0 {
		urlObject.Scheme = "http"
	}
	resp, err := httpClient.Get(urlObject.String())
	if err != nil {
		urlObject.Path = "www." + urlObject.Path
		resp, err = httpClient.Get(urlObject.String())
		if err != nil {
			return nil, err
		}
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return applyRe(string(body))
}

func applyRe(body string) ([]string, error) {
	matches := make(map[string]int)
	for _, match := range searchRe[*search].FindAllStringSubmatch(body, -1) {
		if len(match) > 1 {
			matches[match[1]] = 0
		}
	}
	result := make([]string, 0)
	for k, _ := range matches {
		result = append(result, k)
	}
	return result, nil
}

func getOutFile(pwd string) *os.File {
	if len(*outFile) == 0 {
		return os.Stdout
	}
	// Create output File
	outFilePath := path.Join(pwd, *outFile)

	// Open output File
	out, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	//fw, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_WRONLY, 0666)
	checkError("Error opening the otput file", err)
	return out
}

func initializeWorkers(w *csv.Writer) {
	for i := 1; i < *workers; i++ {
		go fetchAndWrite(jobsChannel, w, i)
	}
}

func startReading(r *csv.Reader) {
	for {
		record, err := r.Read()

		// Stop at EOF
		if err == io.EOF {
			break
		}

		checkError("Error reading the line", err)

		if u := record[0]; len(u) != 0 {
			jobsChannel <- u
		}
	}
}

func validateRequiredOpts() error {
	if len(*csvFile) == 0 && len(*webUrl) == 0 {
		return errors.New("The --csv flag must hava the path to the file or the --url must have a value")
	}

	if len(*search) == 0 || !validSearch.MatchString(*search) {
		return errors.New("Not a supported --search type")
	}
	return nil
}
