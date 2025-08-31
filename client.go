package main

import (
	"bufio"
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

var apiKey string = "" // Global API key variable
var configFilePath string = "otv4.conf"

// Reads server, port, and API key configuration from the config file
func readConfig() (string, string, string, error) {
	file, err := os.Open(configFilePath)
	if err != nil {
		return "", "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var server, port, apiKey string
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
		switch key {
		case "server":
			server = value
		case "port":
			port = value
		case "apiKey":
			apiKey = value
		}
	}
	if err := scanner.Err(); err != nil {
		return "", "", "", err
	}
	return server, port, apiKey, nil
}

// Writes server, port, and API key configuration to the config file
func writeConfig(server, port, apiKey string) error {
	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("server=%s\nport=%s\napiKey=%s\n", server, port, apiKey))
	return err
}

func fetchFileList(server string) ([]string, error) {
	// Update and trim the API key dynamically
	apiKey = strings.TrimSpace(apiKey)

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
	// Update and trim the API key dynamically
	apiKey = strings.TrimSpace(apiKey)

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

func refreshFiles(serverEntry, portEntry, apiKeyEntry *widget.Entry, fileList *fyne.Container, defaultServer, defaultPort string) {
	server := serverEntry.Text
	port := portEntry.Text
	apiKey = strings.TrimSpace(apiKeyEntry.Text) // Dynamically update the API key

	if server == "" {
		server = defaultServer
	}
	if port == "" {
		port = defaultPort
	}
	server = fmt.Sprintf("%s:%s", server, port)

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

// GUI
func main() {
	a := app.New()
	w := a.NewWindow("OpenTunnel V4 Client")

	// Default values for server, port, and API key
	defaultServer := ""
	defaultPort := "8080"
	defaultAPIKey := ""

	// Try to read configuration from the file
	server, port, apiKey, err := readConfig()
	if err != nil {
		server = defaultServer
		port = defaultPort
		apiKey = defaultAPIKey
	}

	// Entry for Server
	serverEntry := widget.NewEntry()
	serverEntry.SetPlaceHolder("Enter server address")
	serverEntry.SetText(server)

	// Entry for Port
	portEntry := widget.NewEntry()
	portEntry.SetPlaceHolder("Port")
	portEntry.SetText(port)

	// Entry for API Key
	apiKeyEntry := widget.NewPasswordEntry()
	apiKeyEntry.SetPlaceHolder("Enter API key")
	apiKeyEntry.SetText(apiKey)

	// Group Server and Port entries horizontally with proportional widths
	serverPortContainer := container.NewGridWithColumns(2,
		container.NewVBox(widget.NewLabel("Server:"), serverEntry), // Server takes 2/3 of the width
		container.NewVBox(widget.NewLabel("Port:"), portEntry),     // Port takes 1/3 of the width
	)

	// Entry for URL
	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter download URL...")

	fileList := container.NewVBox()

	addButton := widget.NewButton("Add", func() {
		url := urlEntry.Text
		server := serverEntry.Text
		port := portEntry.Text
		apiKey = strings.TrimSpace(apiKeyEntry.Text) // Dynamically update the API key

		if server == "" {
			server = defaultServer
		}
		if port == "" {
			port = defaultPort
		}
		server = fmt.Sprintf("%s:%s", server, port)

		if url == "" {
			fmt.Println("Error: URL must be specified.")
			return
		}
		if err := addDownload(server, url); err != nil {
			fmt.Println("Add failed:", err)
		} else {
			fmt.Println("URL added:", url)
			refreshFiles(serverEntry, portEntry, apiKeyEntry, fileList, defaultServer, defaultPort)
		}

		// Save the configuration
		if err := writeConfig(serverEntry.Text, portEntry.Text, apiKeyEntry.Text); err != nil {
			fmt.Println("Failed to save configuration:", err)
		}
	})

	content := container.NewVBox(
		serverPortContainer, // Use the grid container for server and port
		apiKeyEntry,         // Add the API key entry box
		urlEntry,
		addButton,
		widget.NewLabel("Files on server:"),
		fileList,
	)

	w.SetContent(content)
	refreshFiles(serverEntry, portEntry, apiKeyEntry, fileList, defaultServer, defaultPort)
	w.Resize(fyne.NewSize(400, 600))
	w.ShowAndRun()
}
