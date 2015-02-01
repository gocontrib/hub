package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/gocontrib/sock"
)

var addr = flag.String("addr", ":8080", "http service address")
var homeTempl = template.Must(template.ParseFiles("home.html"))
var hub = sock.NewHub()

func main() {
	flag.Parse()

	go hub.Run()

	// bot
	var stop = every(time.Second, func() {
		hub.Send("bot message")
	})
	defer func() {
		stop <- true
		fmt.Println("bot killed")
	}()

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", hub.ServeHTTP)

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

// Returns `chan bool` to stop the execution via: stop <- true
func every(d time.Duration, fn func()) chan bool {
	var stop = make(chan bool, 1)

	go func() {
		var ticker = time.NewTicker(d)

		for {
			select {
			case <-ticker.C:
				fn()

			case <-stop:
				ticker.Stop()
				close(stop)
				return
			}
		}
	}()

	return stop
}
