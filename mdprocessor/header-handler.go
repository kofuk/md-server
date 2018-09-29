package mdprocessor

import (
    "bufio"
    "net/url"
    "strconv"
)

func compileHeader(w *bufio.Writer, line string) {
    nsharp := 0
    for _, c := range line {
        if nsharp > 6 {
            errorExit("Header level can be [1..6]")
        }
        if c == '#' {
            nsharp++
        } else {
            break
        }
    }
    name := getHeaderName(line)
    write(w, "<h" + strconv.Itoa(nsharp) +
        " id=\""+ url.PathEscape(name) + "\">")
    compileDecoration(w, name, false)
    write(w, "</h" + strconv.Itoa(nsharp) + ">")
}

func getHeaderName(line string) string {
    sharpAppeared := false
    var titleIndex int
    for n, c := range []byte(line) {
        if c == '#' {
            sharpAppeared = true
        } else if sharpAppeared && c != '#' && c != ' ' {
            titleIndex = n
            break
        }
    }
    return string([]byte(line)[titleIndex:])
}
