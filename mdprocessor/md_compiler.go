package mdprocessor

import (
    "bufio"
    "net/http"
    "io"
    "strings"
    "os"
)

var configurationRegistry struct {
    hasMath bool
}

func Process(input *os.File, output *http.ResponseWriter) {
    r := bufio.NewReader(input)
    w := bufio.NewWriter(*output)
    linebytes, _, err := r.ReadLine()
    if err == io.EOF {
        errorExit("Input is empty")
        return
    } else if err != nil {
        errorExit("Error reading input")
        return
    }
    firstLine := string(linebytes)
    title := detectTitle(firstLine)
    preExecute(w, title)
    w.Flush()
    processLine(firstLine, r, w)
    w.Flush()
    for {
        linebytes, _, err := r.ReadLine()
        if err == io.EOF {
            break
        } else if err != nil {
            errorExit("Error reading file")
            return
        }
        if len(linebytes) == 0 {
            continue
        }
        line := string(linebytes)
        processLine(line, r, w)
    }
    postExecute(w)
    w.Flush()
}

func processLine(line string, r *bufio.Reader, w *bufio.Writer) {
    if renderIfHr(w, line) {
        w.Flush()
        return
    }
    if strings.HasPrefix(line, "#") {
        compileHeader(w, line)
    } else if strings.HasPrefix(line, "```") {
        compileCodeBlock(w, line, r)
    } else if strings.HasPrefix(line, "- ") {
        compileList(w, line, r)
    } else if strings.HasPrefix(line, ">") {
        compileBlockQuote(w, line, r)
    } else {
        w.WriteString(" ")
        compileDecoration(w, line, true)
    }
    w.Flush()
}
