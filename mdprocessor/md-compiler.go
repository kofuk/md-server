package mdprocessor

import (
    "bufio"
    "net/url"
    "io"
    "strings"
    "os"
)

var configurationRegistry struct {
    hasMath bool
}

func Process(input *os.File, w *bufio.Writer) {
    r := bufio.NewReader(input)
    indexMd, topic := createIndexMd(r)
    input, err := os.Open(input.Name())
    if err != nil {
        errorExit("Cannot reopen file")
    }
    r = bufio.NewReader(input)
    linebytes, _, err := r.ReadLine()
    if err == io.EOF {
        errorExit("Input is empty")
        return
    } else if err != nil {
        errorExit("Error reading input")
        return
    }
    firstLine := string(linebytes)
    title, isH1 := detectTitle(firstLine)
    preExecute(w, title, topic, !isH1)
    w.Flush()
    if !strings.HasPrefix(firstLine, "Title:") {
        processLine(firstLine, r, w)
    }
    w.Flush()
    indexReader := bufio.NewReader(strings.NewReader(indexMd))
    indexLine, _, err := indexReader.ReadLine()
    if err == nil {
        write(w, "<details><summary>目次</summary><div>")
        compileList(w, string(indexLine), indexReader)
        write(w, "</div></details>")
    }
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

func createIndexMd(r *bufio.Reader) (string, string) {
    result := ""
    topic := ""
    for {
        bytes, _, err := r.ReadLine()
        if err != nil {
            return result, topic
        }
        line := string(bytes)
        if strings.HasPrefix(line, "#") {
            level := 0
            for _, c := range line {
                if level > 6 {
                    errorExit("Header level can be [1..6]")
                }
                if c == '#' {
                    level++
                } else {
                    break
                }
            }
            space := make([]byte, (level - 1) * 2)
            for i := 0; i < (level - 1) * 2; i++ {
                space[i] = ' '
            }
            name := getHeaderName(line)
            if level == 1 {
                if topic != "" {
                    topic += "・"
                }
                topic += name
            }
            result = strings.Join([]string{
                result,
                string(space),
                `- <a href="#`,
                url.PathEscape(name),
                `">`,
                name,
                "</a>\n",
            }, "")
        }
    }
    return result, topic
}
