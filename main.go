package main

import (
    "flag"
    "fmt"
    "html/template"
    "net/http"
    "io"
    "io/ioutil"
    "log"
    "os"

    "github.com/KoFuk/md-server/mdprocessor"
)

var addr, fileDir, username, passwd *string
var hasAuth bool

//TODO: Implement basic auth.

func main() {

    addr = flag.String("addr", ":80", "Address and port number to bind.")
    fileDir = flag.String("datadir", "pages/",
        "Location to save and load data.")
    username = flag.String("username", "", "If not empty, enforce basic auth" +
        " with specified password.")
    passwd = flag.String("passwd", "", "If username is not empty, use this" +
        " value as password.")
    flag.Parse()
    hasAuth = *username != ""
    if hasAuth {
        fmt.Println("NOTE: auth functionally is not implemented yet.")
    }
    http.HandleFunc("/", indexHandler)
    http.HandleFunc("/pages/", mdHandler)
    http.Handle("/raw/", http.StripPrefix(
        "/raw/", http.FileServer(http.Dir(*fileDir))))
    http.Handle("/static/", http.FileServer(http.Dir(".")))
    http.HandleFunc("/edit", editHandler)
    http.HandleFunc("/save", saveHandler)
    if http.ListenAndServe(*addr, nil) != nil {
        log.Fatal("Unable to listen '", *addr,
            "'; The port may already be used or have to be run as root.")
    }
}

var indexTemplate = template.Must(
    template.ParseFiles("template/index.html", "template/edit.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        query := r.URL.Query()["delete"]
        if len(query) != 0 && query[0] != "" {
            os.Remove(*fileDir + query[0])
        }
        ls := ls()
        indexTemplate.ExecuteTemplate(w, "index.html", ls)
    } else if r.Method == "POST" {
        if !auth(w, r) {
            return
        }
        uploadHandler(w, r)
    }
}

func mdHandler(w http.ResponseWriter, r *http.Request) {
    file, err := os.Open(string([]rune(r.URL.Path)[1:]))
    if err != nil {
        http.NotFound(w, r)
    }
    defer file.Close()
    mdprocessor.Process(file, &w)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", 405)
        return
    }
    if !auth(w, r) {
        return
    }
    file, header, err := r.FormFile("markdown")
    if err != nil {
        fmt.Println(err)
        http.Error(w, "Bad request", 400)
        return
    }
    if header.Filename == "" {
        fmt.Println("Filename mustn't be empty")
        http.Error(w, "Filename mustn't be empty", 400)
        return
    }
    newFile, err := os.OpenFile(*fileDir + header.Filename,
        os.O_WRONLY | os.O_CREATE, 0666)
    if err != nil {
        http.Error(w, "Server error", 500)
        return
    }
    defer newFile.Close()
    newFile.Truncate(0)
    _, err = io.Copy(newFile, file)
    if err != nil {
        http.Error(w, "Internal server error", 500)
        return
    }
    r.Method = "GET"
    indexHandler(w, r)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
    if !auth(w, r) {
        return
    }
    name := r.URL.Query()["name"]
    if len(name) == 0 || name[0] == "" {
        http.Error(w, "Name is not specified", 400)
        return
    }
    indexTemplate.ExecuteTemplate(w, "edit.html", name[0])
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Method not allowed", 405)
        return
    }
    if !auth(w, r) {
        return
    }
    name := r.URL.Query()["name"]
    if len(name) == 0 || name[0] == "" {
        http.Error(w, "Name is not specified", 400)
        return
    }
    newFile, err := os.OpenFile(*fileDir + name[0],
        os.O_WRONLY | os.O_CREATE, 0666)
    if err != nil {
        http.Error(w, "Server error", 500)
        return
    }
    defer newFile.Close()
    _, err = io.Copy(newFile, r.Body)
    if err != nil {
        http.Error(w, "Internal server error", 500)
        return
    }
    fmt.Sprintln(w, "Success")
}

func ls() []string {
    files, err := ioutil.ReadDir(*fileDir)
    if err != nil {
        panic(err)
    }
    result := make([]string, len(files))
    for i, file := range files {
        result[i] = file.Name()
    }
    return result
}

func auth(w http.ResponseWriter, r *http.Request) bool {
    return true
}
