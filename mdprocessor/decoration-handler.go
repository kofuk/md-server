package mdprocessor

import (
    "bufio"
    "errors"
    "fmt"
)

const (
    D_BOLD = iota
    D_ITALIC
    D_STRIKE
    D_CODE
    D_MATH
    D_LINK
    D_IMAGE
    D_TOOLTIP
)

func compileDecoration(w *bufio.Writer, expr string, returnAllowed bool) {
    line, hasReturn := getLineContent(expr)
    decorations := make([]int, 4)
    di := 0
    cs := []rune(line)
    cursor := 0
    length := len(cs)
    for cursor < length {
        c := cs[cursor]
        if c == '*' {
            cursor++
            if cursor < length && cs[cursor] == '*' {
                if di != 0 && decorations[di - 1] == D_BOLD {
                    write(w, "</strong>")
                    di--
                } else {
                    write(w, "<strong>")
                    decorations[di] = D_BOLD
                    di++
                }
            } else {
                cursor--
                if di != 0 && decorations[di - 1] == D_ITALIC {
                    write(w, "</em>")
                    di--
                } else {
                    write(w, "<em>")
                    decorations[di] = D_ITALIC
                    di++
                }
            }
        } else if c == '~' {
            cursor++
            if cs[cursor] == '~' {
                if di != 0 && decorations[di - 1] == D_STRIKE {
                    write(w, "</del>")
                    di--
                } else {
                    write(w, "<del>")
                    decorations[di] = D_STRIKE
                    di++
                }
            } else {
                cursor--
                write(w, "~")
            }
        } else if c == '`' {
            if di != 0 && decorations[di - 1] == D_CODE {
                write(w, "</code>")
                di--
            } else {
                write(w, "<code>")
                decorations[di] = D_CODE
                di++
            }
        } else if c == '$' {
            configurationRegistry.hasMath = true
            decorations[di] = D_MATH
            di++
            write(w, string(c))
            cursor++
            for {
                if cursor >= len(cs) {
                    break
                } else {
                    if cs[cursor] == '$' {
                        di--
                        write(w, string(cs[cursor]))
                        break
                    } else {
                        if cs[cursor] == '\\' {
                            cursor++
                            if cursor >= len(cs) {
                                break
                            }
                            write(w, "\\")
                        }
                        write(w, string(cs[cursor]))
                    }
                }
                cursor++
            }
        } else if c == '[' {
            title, content, newCursor, err := getEmbedded(cs, cursor)
            if err == nil {
                write(w, fmt.Sprintf(`<a href="%s">%s</a>`, content, title))
                cursor = newCursor
            } else {
                write(w, "[")
            }
        } else if c == '!' {
            cursor++
            if cursor >= len(cs) {
                write(w, "!")
                break
            }
            if cs[cursor] != '[' {
                continue
            }
            title, content, newCursor, err := getEmbedded(cs, cursor)
            if err == nil {
                write(w, fmt.Sprintf(`<img src="%s", alt="%s">`, content, title))
                cursor = newCursor
            }
        } else if c == '?' {
            cursor++
            if cursor >= len(cs) {
                write(w, "?")
                break
            }
            if cs[cursor] != '[' {
                continue
            }
            title, content, newCursor, err := getEmbedded(cs, cursor)
            if err == nil {
                if title == "!" || title == "?" {
                    var displayContent string
                    if title == "!" {
                        displayContent = ICON_INFO
                    } else if title == "?" {
                        displayContent = ICON_HELP
                    }
                    write(w, fmt.Sprintf(`<span class="tooltip" data-tooltip="%s">` +
                            `<img src="%s"></span>`,
                            content, displayContent))
                } else {
                    write(w, fmt.Sprintf(`<span data-tooltip="%s"` +
                            ` class="tooltip text-tooltip">%s</span>`,
                            content, title))
                }
                cursor = newCursor
            }
        } else if c == '\\' {
            cursor++
            if cursor < len(cs) {
                write(w, string(cs[cursor]))
            }
        } else {
            write(w, string(c))
        }
        cursor++
    }
    di--
    for di >= 0 {
        dec := decorations[di]
        if dec == D_BOLD {
            write(w, "</strong>")
        } else if dec == D_ITALIC {
            write(w, "</em>")
        } else if dec == D_STRIKE {
            write(w, "</del>")
        } else if dec == D_CODE {
            write(w, "</code>")
        } else if dec == D_MATH {
            write(w, "$")
        }
        di--
    }
    if returnAllowed && hasReturn {
        write(w, "<br>")
    }
}

func getLineContent(line string) (string, bool) {
    spaceSeqLen := 0
    expr := []rune(line)
    result := ""
    for _, c := range expr {
        if c == rune(' ') {
            spaceSeqLen++
        } else {
            if spaceSeqLen > 0 {
                result += " "
                spaceSeqLen = 0
            }
            result += string(c)
        }
    }
    return result, spaceSeqLen >= 2
}

func getEmbedded(line []rune, cursor int) (string, string, int, error) {
    initialCursor := cursor
    cursor++
    title := make([]rune, 0)
    content := make([]rune, 0)
    isInTitle := true
    terminated := false
    for ; cursor < len(line); cursor++ {
        if isInTitle {
            if line[cursor] == ']' {
                isInTitle = false
                continue
            }
            title = append(title, line[cursor])
        } else {
            if line[cursor] == '(' {
                continue
            } else if line[cursor] == ')' {
                terminated = true
                break
            }
            content = append(content, line[cursor])
        }
    }
    if !terminated {
        return "", "", initialCursor, errors.New("Unterminated embedded content")
    }
    return string(title), string(content), cursor, nil
}

func renderIfHr(w *bufio.Writer, line string) bool {
    if len(line) < 3 {
        return false
    }
    for _, c := range line {
        if c != '=' {
            return false
        }
    }
    write(w, "<hr>")
    return true
}
