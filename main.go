package main

import (
	"os"
	"fmt"
	"bufio"
	"regexp"
	"strings"
)

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

	// Regexp
	ipReg := regexp.MustCompile(`([0-9]{1,3}\.){3}[0-9]{1,3}`)
	datetimeReg := regexp.MustCompile(`\d{1,2}\/\w{3}\/\d{1,4}(:[0-9]{1,2}){3} \+([0-9]){1,4}`)
	httpVersionReg := regexp.MustCompile(`HTTP\/\d.\d`)
	accessPathReg := regexp.MustCompile(`[A-Z]{3,6} \/[\s\S]* HTTP`)
	userAgentReg := regexp.MustCompile(`\"Mozilla[\s\S]*\" `)
	statusCodeAndBodySizeReg := regexp.MustCompile(` [0-9]{3} [0-9]{1,10} `)
	renderTimeReg := regexp.MustCompile(`\"[0-9]{1,3}.[0-9]{1,3}\"`)

	if sc.Scan() {
		for sc.Scan() {
			line := sc.Text()

			reqIp := ""
			reqDate := ""
			reqHttpType := ""
			reqPath := ""
			reqUserAgent := ""
			reqStatusCode := ""
			reqBodySize := ""
			reqTime := ""

			// Find values what we need
			accessIp := ipReg.FindAllString(line, -1)
			accessDate := datetimeReg.FindAllString(line, -1)
			httpVer := httpVersionReg.FindAllString(line, -1)
			accessPath := accessPathReg.FindAllString(line, -1)
			userAgent := userAgentReg.FindAllString(line, -1)
			statusCode := statusCodeAndBodySizeReg.FindAllString(line, -1)
			requestTime := renderTimeReg.FindAllString(line, -1)

			// Format them
			if len(accessIp) > 0 {
				reqIp = accessIp[0]
			}

			if len(accessDate) > 0 {
				reqDate = accessDate[0]
			}

			if len(httpVer) > 0 {
				reqHttpType = httpVer[0]
			}

			if len(accessPath) > 0 {
				replacer := strings.NewReplacer("HTTP", "", "GET", "", "POST", "", " ", "")
				reqPath = replacer.Replace(accessPath[0])
			}

			if len(userAgent) > 0 {
				replacer := strings.NewReplacer("\"", "", "-", "")
				reqUserAgent = replacer.Replace(userAgent[0])
			}

			if len(statusCode) > 0 {
				stringArr := strings.Fields(statusCode[0])
				reqStatusCode = stringArr[0]
				reqBodySize = stringArr[1]
			}

			if len(requestTime) > 0 {
				reqTime = strings.Replace(requestTime[0], "\"", "", -1)
			}

			// Separated by a vertical line
			reqLines = append(reqLines, reqIp, "|", reqDate, "|", reqHttpType, "|", reqPath, "|", reqUserAgent, "|",
				reqStatusCode, "|", reqBodySize, "|", reqTime, "\n")


		}

		fmt.Println(reqLines)
		writeLine(reqLines, "test.log")
	}

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