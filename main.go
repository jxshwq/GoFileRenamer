package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Leggo il percorso della directory dalla riga di comando
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <directory_path>")
		return
	}
	directoryPath := os.Args[1]

	// Leggo il contenuto della directory
	files, err := ioutil.ReadDir(directoryPath)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	for _, file := range files {
		// Costruisco il percorso completo del file
		filePath := filepath.Join(directoryPath, file.Name())

		// Eseguo il comando FFmpeg
		cmd := exec.Command("ffmpeg", "-i", filePath)
		output, err := cmd.CombinedOutput()

		// Controllo gli errori
		if err != nil {
			// Stampo l'output di FFmpeg per diagnosticare l'errore
			fmt.Println("Errore durante l'esecuzione di FFmpeg:")
			fmt.Println(string(output))
			fmt.Println("Errore:", err)
			continue // Passo al prossimo file
		}

		// Conversione dei byte di output in stringa
		outputString := string(output)

		// Ricerca della riga contenente la data di acquisizione
		lines := strings.Split(outputString, "\n")
		for _, line := range lines {
			if strings.Contains(line, "creation_time") {
				// Estrazione della data di acquisizione
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					acquisitionDate := strings.TrimSpace(parts[1])
					fmt.Println("Data di acquisizione per", file.Name(), ":", acquisitionDate)
					break
				}
			}
		}
	}
}
