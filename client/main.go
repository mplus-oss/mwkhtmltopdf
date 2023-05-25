package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	var headerPath, footerPath, pdfPath string
	var bodyPaths []string
	mwkhtmltopdfURL := getEnvOrDefault("MWKHTMLTOPDF_URL", "http://localhost:2777")
	fileArgReached := false
	wkArgs := ""

	if os.Args[1] == "--version" {
		curlCmd := fmt.Sprintf("curl %s", mwkhtmltopdfURL)
		curlExec, err := exec.Command("sh", "-c", curlCmd).Output()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(curlExec))
		os.Exit(0)
	}

	for i := 1; i < len(os.Args); i++ {
		if fileArgReached {
			bodyPaths = append(bodyPaths, os.Args[i])
			continue
		}

		if os.Args[i] == "--header-html" {
			headerPath = os.Args[i+1]
			i++
			continue
		}
		
		if os.Args[i] == "--footer-html" {
			fileArgReached = true
			footerPath = os.Args[i+1]
			i++
			continue
		}

		wkArgs += " " + os.Args[i]
	}

	pdfPath = bodyPaths[len(bodyPaths)-1]
	bodyPaths = bodyPaths[:len(bodyPaths)-1]
	curlCmd := fmt.Sprintf("curl -X POST %s/generate -H 'Content-Type: multipart/form-data'", mwkhtmltopdfURL)

	curlCmd += fmt.Sprintf(" -F 'args=%s'", wkArgs)
	curlCmd += fmt.Sprintf(" -F 'header_html=@%s'", headerPath)
	curlCmd += fmt.Sprintf(" -F 'footer_html=@%s'", footerPath)

	for i, bodyPath := range bodyPaths {
		curlCmd += fmt.Sprintf(" -F 'body_html_%s=@%s'", strconv.Itoa(i), bodyPath)
	}

	curlCmd += fmt.Sprintf(" -o %s", pdfPath)
	curlExec := exec.Command("sh", "-c", curlCmd)
	err := curlExec.Run()
	if err != nil {
		panic(err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
