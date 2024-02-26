package main

import (
    "fmt"
    "io/ioutil"
    "os"
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
        fmt.Println(file.Name())
    }
}
