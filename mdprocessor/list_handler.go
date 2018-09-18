package mdprocessor

import (
    "bufio"
    "io"
)

/**
type UnorderedList struct {
    items []ListItem
}

type ListItem struct {
    children []interface
}
*/

func compileList(w *bufio.Writer, line string, r *bufio.Reader) {
    write(w, "<ul><li>")
    compileDecoration(w, getListItemName(line), false)
    sw := make([]int, 16)
    si := 0
    sw[si] = 0
    for {
        line, _, err := r.ReadLine()
        if err == io.EOF {
            break
        } else if err != nil {
            errorExit("Error reading input")
        }
        if len(line) == 0 {
            break
        }
        shiftWidth := getShiftWidth(string(line))
        if shiftWidth == sw[si] {
            write(w, "</li><li>")
            compileDecoration(w, getListItemName(string(line)), false)
        } else if shiftWidth < sw[si] {
            solved := false
            for {
                si--
                write(w, "</li></ul>")
                if shiftWidth < sw[si] {
                    continue
                } else if shiftWidth == sw[si] {
                    solved = true
                    break
                } else {
                    break
                }
            }
            if !solved {
                errorExit("Invalid list indent")
            }
            write(w, "<li>")
            compileDecoration(w, getListItemName(string(line)), false)
        } else {
            si++
            write(w, "<ul><li>")
            compileDecoration(w, getListItemName(string(line)), false)
            sw[si] = shiftWidth
        }
    }
    for i := 0; i <= si ; i++ {
        write(w, "</li></ul>")
    }
    write(w, "</ul>")
}

func getListItemName(line string) string {
    var startIndex int
    hyphenAppeared := false
    for n, c := range line {
        if !hyphenAppeared && c == ' ' {
            continue
        } else if !hyphenAppeared && c == '-' {
            hyphenAppeared = true
        } else if hyphenAppeared && c != ' ' {
            startIndex = n
            break
        }
    }
    return string([]byte(line)[startIndex:])
}

func getShiftWidth(line string) int {
    for n, c := range line {
        if c != ' ' {
            return n
        }
    }
    return len(line)
}
