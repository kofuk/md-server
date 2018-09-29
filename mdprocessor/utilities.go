package mdprocessor

import (
    "bufio"
    "fmt"
    "strings"
)

func errorExit(message string) {
    fmt.Println(message)
}

func detectTitle(line string) (string, bool) {
    if strings.HasPrefix(line, "# ") {
        return getHeaderName(line), true
    } else if strings.HasPrefix(line, "Title:") {
        titleOffset := 6
        for line[titleOffset] == ' ' {
            titleOffset++
        }
        return line[titleOffset:], false
    }
    return "", false
}

func write(w *bufio.Writer, s string) {
    _, err := w.WriteString(s)
    if err != nil {
        println(err)
    }
}
