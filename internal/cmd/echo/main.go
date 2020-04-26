package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func echo(w http.ResponseWriter, req *http.Request) {
	log.Printf("[method: %7v][url: %v]", req.Method, req.URL.String())
	io.WriteString(os.Stderr, "\n")
	req.Header.Write(os.Stderr)
	io.WriteString(os.Stderr, "\n")
	delay, err := time.ParseDuration(req.URL.Query().Get("delay"))
	if err != nil {
		delay = 0
	} else if delay > time.Minute {
		delay = time.Minute
	}
	time.Sleep(delay)
	w.WriteHeader(http.StatusOK)
	if req.Method == "GET" {
		io.WriteString(w, req.URL.Query().Encode())
	} else {
		io.Copy(w, req.Body)
	}
	return
}

func main() {
	bind := flag.String("bind", "127.0.0.1:8081", "Address to bind and listen for requests")
	flag.Parse()

	http.HandleFunc("/", echo)

	log.Printf("Server available at %v", *bind)
	err := http.ListenAndServe(*bind, nil)
	log.Print(err)
}
