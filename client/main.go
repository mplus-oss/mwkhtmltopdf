package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
)

func main() {
	var headerPath, footerPath, pdfPath string
	var bodyPaths []string
	mwkhtmltopdfURL := getEnvOrDefault("MWKHTMLTOPDF_URL", "http://localhost:2777")
	wkArgs := ""

	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Println("wkhtmltopdf 0.12.6.1 (with patched qt)")
		return
	}

	if len(os.Args) < 3 {
		fmt.Println("Not enough arguments")
		return
	}

	pdfPath = os.Args[len(os.Args)-1]
	i := len(os.Args) - 2

	for path.Base(os.Args[i])[:11] == "report.body" {
		bodyPaths = append(bodyPaths, os.Args[i])
		i--
	}

	for i >= 1 {
		arg := os.Args[i]
		if arg == "--header-html" {
			headerPath = os.Args[i+1]
			i--
			continue
		}

		if arg == "--footer-html" {
			footerPath = os.Args[i+1]
			i--
			continue
		}

		if len(arg) >= 5 && arg[len(arg)-5:] == ".html" {
			i--
			continue
		}

		wkArgs = arg + " " + wkArgs
		i--
	}

	hc := &http.Client{}
	form := url.Values{}
	form.Add("args", wkArgs)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	if headerPath != "" {
		headerFile, err := os.Open(headerPath)
		if err != nil {
			fmt.Println("Failed to open header file:", err)
			return
		}
		defer headerFile.Close()
		part, err := writer.CreateFormFile("header_html", filepath.Base(headerFile.Name()))
		if err != nil {
			fmt.Println("Failed to create header part:", err)
			return
		}
		if _, err := io.Copy(part, headerFile); err != nil {
			fmt.Println("Failed to copy header file:", err)
			return
		}
	}

	if footerPath != "" {
		footerFile, err := os.Open(footerPath)
		if err != nil {
			fmt.Println("Failed to open footer file:", err)
			return
		}
		defer footerFile.Close()
		part, err := writer.CreateFormFile("footer_html", filepath.Base(footerFile.Name()))
		if err != nil {
			fmt.Println("Failed to create footer part:", err)
			return
		}
		if _, err := io.Copy(part, footerFile); err != nil {
			fmt.Println("Failed to copy footer file:", err)
			return
		}
	}

	for i, bodyPath := range bodyPaths {
		bodyFile, err := os.Open(bodyPath)
		if err != nil {
			fmt.Println("Failed to open body file:", err)
			return
		}
		defer bodyFile.Close()
		part, err := writer.CreateFormFile(fmt.Sprintf("body_html_%d", i), filepath.Base(bodyFile.Name()))
		if err != nil {
			fmt.Println("Failed to create body part:", err)
			return
		}
		if _, err := io.Copy(part, bodyFile); err != nil {
			fmt.Println("Failed to copy body file:", err)
			return
		}
	}

	req, err := http.NewRequest("POST", mwkhtmltopdfURL+"/generate", bytes.NewReader(body.Bytes()))
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := hc.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to generate PDF:", resp.Status)
		return
	}

	out, err := os.Create(pdfPath)
	if err != nil {
		fmt.Println("Failed to create output file:", err)
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, resp.Body); err != nil {
		fmt.Println("Failed to write output file:", err)
		return
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
