package main

import (
	"os"
	"fmt"
	"bufio"
	"regexp"
	"strings"
)

type ReqData struct {
	Ip, Date, HttpMethod, Path, UserAgent, StatusCode, RenderTime []string
}

func main() {

	var reqLines []string

	// Open testing log
	logFile := "../log_nginx/access.log"
	fl, err := os.Open(logFile)
	if err != nil {
		fmt.Println(logFile, err)
	}
	defer fl.Close()

	// Initialize the scanner
	sc := bufio.NewScanner(fl)

	if sc.Scan() {
		for sc.Scan() {
			line := sc.Text()

			var reqIp, reqDate, reqHttpType, reqPath, reqUserAgent, reqStatusCode, reqBodySize, reqTime string
			req := readLine(line)

			// Format them
			if len(req.Ip) > 0 {
				reqIp = req.Ip[0]
			}

			if len(req.Date) > 0 {
				reqDate = req.Date[0]
			}

			if len(req.HttpMethod) > 0 {
				reqHttpType = req.HttpMethod[0]
			}

			if len(req.Path) > 0 {
				replacer := strings.NewReplacer("HTTP", "", "GET", "", "POST", "", " ", "")
				reqPath = replacer.Replace(req.Path[0])
			}

			if len(req.UserAgent) > 0 {
				replacer := strings.NewReplacer("\"", "", "-", "")
				reqUserAgent = replacer.Replace(req.UserAgent[0])
			}

			if len(req.StatusCode) > 0 {
				stringArr := strings.Fields(req.StatusCode[0])
				reqStatusCode = stringArr[0]
				reqBodySize = stringArr[1]
			}

			if len(req.RenderTime) > 0 {
				reqTime = strings.Replace(req.RenderTime[0], "\"", "", -1)
			}

			// Separated by a vertical line
			reqLines = append(reqLines, reqIp, "|", reqDate, "|", reqHttpType, "|", reqPath, "|", reqUserAgent, "|",
				reqStatusCode, "|", reqBodySize, "|", reqTime, "\n")


		}

		fmt.Println(reqLines)
		writeLine(reqLines, "test.log")
	}

}

// Read line and find the values I need
func readLine(line string) ReqData {

	req := ReqData{
		Ip: regexp.MustCompile(`([0-9]{1,3}\.){3}[0-9]{1,3}`).FindAllString(line, -1),
		Date: regexp.MustCompile(`\d{1,2}\/\w{3}\/\d{1,4}(:[0-9]{1,2}){3} \+([0-9]){1,4}`).FindAllString(line, -1),
		HttpMethod: regexp.MustCompile(`HTTP\/\d.\d`).FindAllString(line, -1),
		Path: regexp.MustCompile(`[A-Z]{3,6} \/[\s\S]* HTTP`).FindAllString(line, -1),
		UserAgent: regexp.MustCompile(`\"Mozilla[\s\S]*\" `).FindAllString(line, -1),
		StatusCode: regexp.MustCompile(` [0-9]{3} [0-9]{1,10} `).FindAllString(line, -1),
		RenderTime: regexp.MustCompile(`\"[0-9]{1,3}.[0-9]{1,3}\"`).FindAllString(line, -1),
	}

	return req
}

// Flush write into the file
func writeLine(lines []string, path string) error {
	file, err := os.OpenFile(path, os.O_APPEND | os.O_WRONLY | os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprint(w, line)
	}
	return w.Flush()
}