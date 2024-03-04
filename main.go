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
	// Check if the directory path is provided as a command-line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory_path>")
		return
	}
	directoryPath := os.Args[1]

	// Read the files in the specified directory
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	// Iterate over each file in the directory
	for _, file := range files {
		filePath := filepath.Join(directoryPath, file.Name())

		// Run FFmpeg command to get the creation time of the file
		cmd := exec.Command("ffmpeg", "-i", filePath)
		stderr, _ := cmd.StderrPipe()
		cmd.Start()

		// Read the output of the FFmpeg command line by line
		scanner := bufio.NewScanner(stderr)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "creation_time") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					acquisitionDate := strings.TrimSpace(strings.Join(parts[1:], ":"))

					// Parse the acquisition date using the specified format
					date, err := time.Parse("2006-01-02T15:04:05.000000Z", acquisitionDate)
					if err != nil {
						fmt.Println("Error parsing date:", err)
						continue
					}

					// Print the acquisition date
					fmt.Println("Acquisition date for", file.Name(), ":", date.Year(), date.Month(), date.Day())

					// Rename the file with the new format
					newFileName := fmt.Sprintf("%d-%02d-%02d-%02d-%02d%s", date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), filepath.Ext(file.Name()))
					newFilePath := filepath.Join(directoryPath, newFileName)
					err = os.Rename(filePath, newFilePath)
					if err != nil {
						fmt.Println("Error renaming file:", err)
					} else {
						fmt.Println("Renamed", file.Name(), "to", newFileName)
					}
					fmt.Println()
					break
				}
			}
		}

		cmd.Wait()
	}
}