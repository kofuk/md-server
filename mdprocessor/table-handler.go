package mdprocessor

import (
    "bufio"
    "errors"
)

type tableColumn struct {
    align int
}

const (
    ALIGN_NONE = iota
    ALIGN_LEFT
    ALIGN_CENTER
    ALIGN_RIGHT
)

func getTableConfig(line string) ([]tableColumn, error) {
    columns := make([]tableColumn, 0)
    expectPipe := false
    expectHyphenOrColon := true
    isFirstInColumn := true
    var column tableColumn
    for _, c := range line {
        if expectHyphenOrColon {
            if !expectPipe && c != '-' && c != ':' {
                return columns,
                    errors.New("Hyphen(-) or colon(:) expected, but got" +  string(c))
            }
            if c == '-' {
                expectPipe = true
                if isFirstInColumn {
                    column = tableColumn{ALIGN_NONE}
                    isFirstInColumn = false
                }
                continue
            } else if c == ':' {
                if isFirstInColumn {
                    column = tableColumn{ALIGN_LEFT}
                    isFirstInColumn = false
                } else {
                    if column.align == ALIGN_LEFT {
                        column.align = ALIGN_CENTER
                    } else {
                        column.align = ALIGN_RIGHT
                    }
                    expectHyphenOrColon = false
                    expectPipe = true
                }
                continue
            }
        }
        if expectPipe {
            if c != '|' {
                return columns, errors.New("Pipe(|) expected, but got" + string(c))
            }
            expectHyphenOrColon = true
            isFirstInColumn = true
            expectPipe = false
            columns = append(columns, column)
            continue
        }
        return columns, errors.New("Unexpected character")
    }
    columns = append(columns, column)
    return columns, nil
}

func compileTable(headerLine string, config []tableColumn,
        r *bufio.Reader, w *bufio.Writer) {
    write(w, "<table><thead><tr>")
    runes := []rune(headerLine)
    cursor := 0
    for i := 0; i < len(config); i++ {
        write(w, "<th")
        align := config[i].align
        if align == ALIGN_LEFT {
            write(w, ` align="left"`)
        } else if align == ALIGN_CENTER {
            write(w, ` align="center"`)
        } else if align == ALIGN_RIGHT {
            write(w, ` align="right"`)
        }
        write(w, ">")
        data, newCursor := getColumn(runes, cursor)
        cursor = newCursor
        compileDecoration(w, data, true)
        write(w, "</th>")
    }
    write(w, "</tr></thead><tbody>")
    for {
        b, _, err := r.ReadLine()
        if err != nil {
            break
        }
        if len(b) == 0 {
            break
        }
        line := []rune(string(b))
        write(w, "<tr>")
        for i:= 0; i < len(config); i++ {
            cursor = 0
            write(w, "<td")
            align := config[i].align
            if align == ALIGN_LEFT {
                write(w, ` align="left"`)
            } else if align == ALIGN_CENTER {
                write(w, ` align=center`)
            } else if align == ALIGN_RIGHT {
                write(w, ` align="right"`)
            }
            write(w, ">")
            data, newCursor := getColumn(line, cursor)
            cursor = newCursor
            compileDecoration(w, data, true)
        }
        write(w, "</tr>")
    }
    write(w, "</tbody></table>")
}

func getColumn(line []rune, cursor int) (string, int) {
    if cursor == -1 {
        return "", -1
    }
    result := make([]rune, 0)
    for ; cursor < len(line); cursor++ {
        if line[cursor] == '|' {
            return string(result), cursor + 1
        } else {
            result = append(result, line[cursor])
        }
    }
    return string(result), -1
}
