package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"

	"github.com/gocontrib/hub"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("home.html"))
var h = hub.NewHub()

func main() {
	flag.Parse()

	go h.Run()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", h.ServeWs)

	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTempl.Execute(w, r.Host)
}
