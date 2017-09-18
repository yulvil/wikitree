package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/fcgi"
	"net/url"
	"os"
	"runtime"
	"time"

	"gopkg.in/resty.v0"
)

var app_addr string
var index int

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	app_addr = os.Getenv("APP_ADDR") // e.g. "0.0.0.0:8080" or ""
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func mkdir(path string) {
	os.MkdirAll(path, 0700)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	headers := w.Header()

	headers.Add("Content-Type", "application/json; charset=utf-8")
	headers.Add("Access-Control-Allow-Origin", "https://www.pbcote.com")

	var wid = r.FormValue("wid")
	var depth = r.FormValue("depth")
	var cache = r.FormValue("cache")

	if wid == "" {
		io.WriteString(w, `{"ancestors":{}}`)
		return
	}

	var t = time.Now()
	var dirname = t.Format("20060102")
	var filename = fmt.Sprintf("cache/%s/%s-%02s.json", dirname, wid, depth)

	if cache != "bust" {

		cached, err := os.Open(filename)
		if err == nil {
			defer cached.Close()
			io.Copy(w, cached)
			//fmt.Println("Using cache")
			return
		}
	}

	//fmt.Println("Using service")
	resp, _ := resty.R().Get(fmt.Sprintf("https://apps.wikitree.com/api.php?action=getAncestors&format=json&depth=%s&key=%s", depth, url.QueryEscape(wid)))
	io.WriteString(w, resp.String())

	mkdir("cache/" + dirname)

	ioutil.WriteFile(filename, []byte(resp.String()), 0644)
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
