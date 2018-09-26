package mdprocessor

import (
    "bufio"
)

func preExecute(w *bufio.Writer, title string) {
    write(w, "<!doctype html><html><head><meta charset=\"utf-8\">" +
        "<meta http-equiv=\"X-UA-Compatible\" content=\"IE=edge\">" +
        "<meta name=\"viewport\" content=\"width=device-width\">")
    if len(title) != 0 {
        write(w, "<title>" + title + "</title>")
    }
    write(w, "<style>")
    write(w, defStyle)
    write(w, "</style></head><body><div class=\"markdown-body\">")
}

const MATHJAX = "https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.5/MathJax.js?config=TeX-MML-AM_CHTML"
const MATHJAX_CONFIG = `<script type="text/x-mathjax-config">
MathJax.Hub.Config({tex2jax:{inlineMath:[['$','$'],['\\(','\\(']],
processEscapes:true},CommonHTML:{matchFontHeight:false}});</script>`

func postExecute(w *bufio.Writer) {
    write(w, "</div>")
    if configurationRegistry.hasMath {
        write(w, MATHJAX_CONFIG + "<script src=\"" + MATHJAX +
            "\" async></script>")
    }
    write(w, "</body></html>")
}
