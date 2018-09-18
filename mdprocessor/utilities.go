package mdprocessor

import (
    "bufio"
    "fmt"
    "strings"
)

func errorExit(message string) {
    fmt.Println(message)
}

func detectTitle(line string) string {
    if strings.HasPrefix(line, "# ") {
        return getHeaderName(line)
    }
    return ""
}

func write(w *bufio.Writer, s string) {
    _, err := w.WriteString(s)
    if err != nil {
        println(err)
    }
}
