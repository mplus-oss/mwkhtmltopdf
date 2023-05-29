package main

import (
	"log"
	"fmt"
	"io"
	"os"
	"os/exec"
	"net/http"
	"github.com/labstack/echo/v4"
	"strings"
)

func main() {
	e := echo.New()
	e.HideBanner = true
	e.GET("/", func(c echo.Context) error {
		cmd := exec.Command("sh", "-c", "wkhtmltopdf --version")
		out, err := cmd.Output()
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, string(out))
	})
	e.POST("/generate", func(c echo.Context) error {
		headerEnabled := true
		headerArgs := ""
		footerEnabled := true
		footerArgs := ""

		headerFileMultiPart, err := c.FormFile("header_html")
		if err != nil {
			headerEnabled = false
		}

		footerFileMultiPart, err := c.FormFile("footer_html")
		if err != nil {
			footerEnabled = false
		}

		wkArgs := c.FormValue("args")

		dir, err := os.MkdirTemp("", "pdfgen")
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to create temp dir")
		}
		
		defer os.RemoveAll(dir)

		if headerEnabled {
			headerFile, err := headerFileMultiPart.Open()
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to open header file")
			}
			defer headerFile.Close()
			headerFilePath := fmt.Sprintf("%s/%s", dir, headerFileMultiPart.Filename)
			headerFileOutput, err := os.Create(headerFilePath)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to create header file")
			}
			defer headerFileOutput.Close()
			if _, err := io.Copy(headerFileOutput, headerFile); err != nil {
				return c.String(http.StatusInternalServerError, "Failed to copy header file")
			}
			headerArgs = fmt.Sprintf("--header-html %s", headerFilePath)
		}

		if footerEnabled {
			footerFile, err := footerFileMultiPart.Open()
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to open footer file")
			}
			defer footerFile.Close()
			footerFilePath := fmt.Sprintf("%s/%s", dir, footerFileMultiPart.Filename)
			footerFileOutput, err := os.Create(footerFilePath)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Failed to create footer file")
			}
			defer footerFileOutput.Close()
			if _, err := io.Copy(footerFileOutput, footerFile); err != nil {
				return c.String(http.StatusInternalServerError, "Failed to copy footer file")
			}
			footerArgs = fmt.Sprintf("--footer-html %s", footerFilePath)
		}

		bodyFilesList, err := c.MultipartForm()
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to get body files list")
		}
		
		bodyFileArgs := []string{}

		for key, _ := range bodyFilesList.File {
			if key[:9] == "body_html" {
				bodyFileMultiPart, err := c.FormFile(key)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Failed to get body file")
				}
				bodyFile, err := bodyFileMultiPart.Open()
				if err != nil {
					return c.String(http.StatusInternalServerError, "Failed to open body file")
				}
				defer bodyFile.Close()
				bodyFilePath := fmt.Sprintf("%s/%s", dir, bodyFileMultiPart.Filename)
				bodyFileOutput, err := os.Create(bodyFilePath)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Failed to create body file")
				}
				defer bodyFileOutput.Close()
				if _, err := io.Copy(bodyFileOutput, bodyFile); err != nil {
					return c.String(http.StatusInternalServerError, "Failed to copy body file")
				}
				bodyFileArgs = append(bodyFileArgs, bodyFilePath)
			}
		}

		bodyFileArgsStr := strings.Join(bodyFileArgs, " ")

		pdfPath := fmt.Sprintf("%s/output.pdf", dir)
		pdfCmd := fmt.Sprintf("wkhtmltopdf %s %s %s %s %s", wkArgs, headerArgs, footerArgs, bodyFileArgsStr, pdfPath)
		log.Println(fmt.Sprintf("Running command: %s", pdfCmd))
		out, err := exec.Command("sh", "-c", pdfCmd).CombinedOutput()
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to generate PDF: %s", string(out)))
		}

		pdfFile, err := os.Open(pdfPath)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to open PDF file")
		}
		defer pdfFile.Close()

		return c.Stream(http.StatusOK, "application/pdf", pdfFile)
	})
	e.Logger.Fatal(e.Start(":2777"))
}