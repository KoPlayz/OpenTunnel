package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AddResponse struct {
	Status string `json:"status"`
	File   string `json:"file"`
}

const apiKey = "ADadm7QY50k2rySCj0Nang69hhZne8SR" // Replace with the same API key as the server

func fetchFileList(server string) ([]string, error) {
	req, err := http.NewRequest("GET", server+"/list", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", apiKey) // Add the API key to the request header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check for 401 Unauthorized
	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("Error 401: Unauthorized")
	}

	// Check for other non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server error: %s", string(body))
	}

	var files []string
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}
	return files, nil
}

func addDownload(server, url string) error {
	payload, _ := json.Marshal(map[string]string{"url": url})
	req, err := http.NewRequest("POST", server+"/add", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", apiKey) // Add the API key to the request header
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error: %s", string(body))
	}
	return nil
}

func downloadAndDecode(server, filename string) error {
	// Decode the filename to handle URL-encoded characters (e.g., %20 -> space)
	decodedFilename, err := url.QueryUnescape(filename)
	if err != nil {
		return fmt.Errorf("failed to decode filename: %v", err)
	}

	// Create a new GET request
	req, err := http.NewRequest("GET", server+"/"+filename, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-API-Key", apiKey) // Add the API key to the request header

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error: %s", string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Check if the data is valid base64
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return fmt.Errorf("failed to decode base64 data: %v", err)
	}

	newName := strings.TrimSuffix(decodedFilename, ".ote")
	return os.WriteFile(filepath.Base(newName), decoded, 0644)
}

func main() {
	server := "http://web.koplayz.dev:8080"

	a := app.New()
	w := a.NewWindow("OpenTunnel V4 Client")

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter download URL...")

	fileList := container.NewVBox()

	refreshFiles := func() {
		files, err := fetchFileList(server)
		fileList.Objects = nil
		if err != nil {
			fileList.Add(widget.NewLabel(fmt.Sprintf("Error: %v", err)))
		} else {
			for _, f := range files {
				file := strings.ReplaceAll(f, "%20", " ") // Replace %20 with spaces
				btn := widget.NewButton(file, func() {
					if err := downloadAndDecode(server, f); err != nil { // Use original filename for download
						fmt.Println("Error:", err)
					} else {
						fmt.Println("Downloaded & decoded:", file)
					}
				})
				fileList.Add(btn)
			}
		}
		fileList.Refresh()
	}

	addButton := widget.NewButton("Add", func() {
		url := urlEntry.Text
		if url == "" {
			return
		}
		if err := addDownload(server, url); err != nil {
			fmt.Println("Add failed:", err)
		} else {
			fmt.Println("URL added:", url)
			refreshFiles()
		}
	})

	content := container.NewVBox(
		urlEntry,
		addButton,
		widget.NewLabel("Files on server:"),
		fileList,
	)

	w.SetContent(content)
	refreshFiles()
	w.Resize(fyne.NewSize(400, 600))
	w.ShowAndRun()
}
