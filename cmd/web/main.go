package main

import (
	"flag"
	"net/http"
)

var cfg Config

func main() {
	// addr := flag.String("addr", ":http-alt", "HTTP network address")
	flag.StringVar(&cfg.addr, "addr", ":http-alt", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "staticdir", "./ui/static", "Static dir to serve")

	flag.Parse()

	mux := http.NewServeMux()

	// Private handlers
	// Fixed path
	mux.HandleFunc("/greeting", Greeting)

	// Subtree path
	mux.HandleFunc("/staticsubtree/", StaticSubTree)

	mux.HandleFunc("/headers", HeadersTask)

	mux.HandleFunc("/serving_just_one_file", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./ui/dummy/index.html")
	})

	// Book handlers
	book(mux)

	svc := http.Server{
		Addr:         cfg.addr,
		Handler:      mux,
		ReadTimeout:  30,
		WriteTimeout: 30,
		ErrorLog:     errorLog,
	}

	InfoF("Starting server on %s", cfg.addr)
	ErrorLn(svc.ListenAndServe())
}

/**
When sending a response in Go, three headers are set by default:
	- Date
	- Content-Length
	- Content-Type (text/plain; chartset=utf-8) Go uses http.DetectContentType() function to detect our content-Type. By default if it is not able to detect the MIME type, then it is set to application/octect-stream

How to disable the folder architecture when fetching a folder resource? (Help: https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings)
- Create an index.html file in each of them: find ./ui/static -type d -exec touch index.html\;
- Creating a middleware and each request that ends with a trailing slash, is returned a 404 not found.
*/
