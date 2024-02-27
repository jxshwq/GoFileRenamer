package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory_path>")
		return
	}
	directoryPath := os.Args[1]

	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, file := range files {
		filePath := filepath.Join(directoryPath, file.Name())

		cmd := exec.Command("ffmpeg", "-i", filePath)
		stderr, _ := cmd.StderrPipe()
		cmd.Start()

		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "creation_time") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					acquisitionDate := strings.TrimSpace(strings.Join(parts[1:], ":"))
					date, err := time.Parse("2006-01-02T15:04:05.000000Z", acquisitionDate)
					if err != nil {
						fmt.Println("Error parsing date:", err)
						continue
					}
					fmt.Println("Acquisition date for", file.Name(), ":", date.Year(), date.Month(), date.Day())

					// Rename the file
					newFileName := fmt.Sprintf("%d-%02d-%02d-%02d-%02d%s", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), filepath.Ext(file.Name()))
					newFilePath := filepath.Join(directoryPath, newFileName)
					err = os.Rename(filePath, newFilePath)
					if err != nil {
						fmt.Println("Error renaming file:", err)
					} else {
						fmt.Println("Renamed", file.Name(), "to", newFileName)
					}
					break
				}
			}
		}

		cmd.Wait()
	}
}