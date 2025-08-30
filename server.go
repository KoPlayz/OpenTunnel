package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func fetchFileList(server string) ([]string, error) {
	resp, err := http.Get(server + "/list")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var files []string
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}
	return files, nil
}

func downloadAndDecode(server, filename string) error {
	// Download file
	resp, err := http.Get(server + "/" + filename)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Decode base64
	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}

	// Save without ".ote"
	newName := strings.TrimSuffix(filename, ".ote")
	return os.WriteFile(filepath.Base(newName), decoded, 0644)
}

func main() {
	server := "http://127.0.0.1:8080"

	a := app.New()
	w := a.NewWindow("OpenTunnel V3 (Client)")

	files, err := fetchFileList(server)
	if err != nil {
		w.SetContent(widget.NewLabel(fmt.Sprintf("Error: %v", err)))
		w.ShowAndRun()
		return
	}

	// List of buttons for each file
	var items []fyne.CanvasObject
	for _, f := range files {
		file := f
		btn := widget.NewButton(file, func() {
			err := downloadAndDecode(server, file)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Downloaded & decoded:", file)
			}
		})
		items = append(items, btn)
	}

	w.SetContent(container.NewVBox(items...))
	w.ShowAndRun()
}
