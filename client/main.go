package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"path"
)

func main() {
	var headerPath, footerPath, pdfPath string
	var bodyPaths []string
	mwkhtmltopdfURL := getEnvOrDefault("MWKHTMLTOPDF_URL", "http://localhost:2777")
	wkArgs := ""
	debug := getEnvOrDefault("MWKHTMLTOPDF_DEBUG", "false")

	if os.Args[1] == "--version" {
		fmt.Println("wkhtmltopdf 0.12.6.1 (with patched qt)")
		os.Exit(0)
	}

	pdfPath = os.Args[len(os.Args) - 1]
	i := len(os.Args) - 2

	for path.Base(os.Args[i])[:11] == "report.body" {
		bodyPaths = append([]string{os.Args[i]}, bodyPaths...)
		i--
	}

	for i >= 1 {
		if os.Args[i] == "--header-html" {
			headerPath = os.Args[i+1]
			i--
			continue
		}
		
		if os.Args[i] == "--footer-html" {
			footerPath = os.Args[i+1]
			i--
			continue
		}

		if len(os.Args[i]) >= 5 && os.Args[i][len(os.Args[i])-5:] == ".html" {
			i--
			continue
		}

		wkArgs = os.Args[i] + " " + wkArgs
		i--
	}

	curlCmd := fmt.Sprintf("curl -X POST %s/generate -H 'Content-Type: multipart/form-data' -F 'args=%s'", mwkhtmltopdfURL, wkArgs)
	if headerPath != "" {
		curlCmd += fmt.Sprintf(" -F 'header_html=@%s'", headerPath)
	}
	if footerPath != "" {
		curlCmd += fmt.Sprintf(" -F 'footer_html=@%s'", footerPath)
	}

	for i, bodyPath := range bodyPaths {
		curlCmd += fmt.Sprintf(" -F 'body_html_%s=@%s'", strconv.Itoa(i), bodyPath)
	}

	curlCmd += fmt.Sprintf(" -o %s", pdfPath)
	curlExec := exec.Command("sh", "-c", curlCmd)

	if debug == "true" {
		fmt.Fprintln(os.Stderr, curlCmd)
	}
	
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
