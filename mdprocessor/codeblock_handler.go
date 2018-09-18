package mdprocessor

import (
    "bufio"
    "strings"
)

func compileCodeBlock(w *bufio.Writer, line string, r *bufio.Reader) {
    lang := getCodeBlockLang(line)
    isMath := lang == "math"
    if isMath {
        configurationRegistry.hasMath = true
        write(w, "<div style=\"margin:20px\">$")
    } else {
        write(w, "<pre><code>")
    }
    for {
        line, _, err := r.ReadLine()
        if err != nil {
            errorExit("Error reading input file")
        }
        if len(line) == 3 && string(line) == "```" {
            break
        }
        line2 := string(line)
        if isMath {
            line2 = strings.Replace(line2, "$", "\\$", -1)
        } else {
            line2 = strings.Replace(line2, "<", "&lt;", -1)
            line2 = strings.Replace(line2, ">", "&gt;", -1)
        }
        write(w, line2 + "\n")
    }
    if isMath {
        write(w, "$</div>")
    } else {
        write(w, "</code></pre>")
    }
}

func getCodeBlockLang(line string) string {
    return string([]byte(line)[3:])
}
