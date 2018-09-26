// +build cli

package main

import (
    "bufio"
    "log"
    "os"
    "strings"

    "github.com/KoFuk/md-server/mdprocessor"
)

func main() {
    if len(os.Args) != 2 {
        log.Fatal("Just 1 argument required.")
    }
    file, err := os.Open(os.Args[1])
    if err != nil {
        log.Fatal("Cannot open file for read")
    }
    defer file.Close()
    outFile, err := os.Create(outFileName(os.Args[1]))
    if err != nil {
        log.Fatal("Cannot create output file", outFileName(os.Args[1]))
    }
    defer outFile.Close()
    mdprocessor.Process(file, bufio.NewWriter(outFile))
}

func outFileName(inName string) string {
    return strings.TrimSuffix(inName, ".md") + ".html"
}
