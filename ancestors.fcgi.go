package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/fcgi"
	"net/url"
	"os"
	"runtime"

	"gopkg.in/resty.v0"
)

var app_addr string
var index int

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app_addr = os.Getenv("APP_ADDR") // e.g. "0.0.0.0:8080" or ""
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()

	headers.Add("Content-Type", "application/json; charset=utf-8")
	headers.Add("Access-Control-Allow-Origin", "https://www.pbcote.com")

	resp, _ := resty.R().Get(fmt.Sprintf("https://apps.wikitree.com/api.php?action=getAncestors&format=json&depth=%s&key=%s", r.FormValue("depth"), url.QueryEscape(r.FormValue("wid"))))
	io.WriteString(w, resp.String())
}

func main() {
	http.HandleFunc("/", ServeHTTP)

	var err error
	if app_addr != "" { // Run as a local web server
		err = http.ListenAndServe(app_addr, nil)
	} else { // Run as FCGI via standard I/O
		err = fcgi.Serve(nil, nil)
	}
	if err != nil {
		log.Fatal(err)
	}
}
