package mdprocessor

import (
	"bufio"
)

func compileBlockQuote(w *bufio.Writer, line string, r *bufio.Reader) {
	prevDepth := getQuoteDepth(line)
	for i := 0; i < prevDepth; i++ {
		write(w, "<blockquote><p>")
	}
	compileDecoration(w, getQuoteContent(line), true)
	for {
		line, _, err := r.ReadLine()
		if err != nil || len(line) == 0 || line[0] != '>' {
			break
		}
		depth := getQuoteDepth(string(line))
		if depth < prevDepth {
			for prevDepth != depth {
				write(w, "</p></blockquote>")
				prevDepth--
			}
			compileDecoration(w, getQuoteContent(string(line)), true)
		} else if depth == prevDepth {
			write(w, " ")
			compileDecoration(w, getQuoteContent(string(line)), true)
		} else {
			for depth != prevDepth {
				write(w, "<blockquote><p>")
				prevDepth++
			}
			compileDecoration(w, getQuoteContent(string(line)), true)
		}
	}
	for prevDepth > 0 {
		write(w, "</p></blockquote>")
		prevDepth--
	}
}

func getQuoteDepth(line string) int {
	depth := 0
	for _, c := range []byte(line) {
		if c == '>' {
			depth++
		} else if c == ' ' {
			continue
		} else {
			break
		}
	}
	return depth
}

func getQuoteContent(line string) string {
	contentStart := 0
	for n, c := range []byte(line) {
		if c != '>' && c != ' ' {
			contentStart = n
			break
		}
	}
	return string([]byte(line)[contentStart:])
}
