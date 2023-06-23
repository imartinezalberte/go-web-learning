package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
)

// Custom counter, awesome :)
type counter int64

func (h *counter)ServeHTTP(w http.ResponseWriter, r *http.Request) {	
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintf(w, "Hello from counter number %d", atomic.AddInt64((*int64)(h), 1))
}

func ViewSnippet(w http.ResponseWriter, r *http.Request) {
	idQueryParam := "id"

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		// Exactly the same as w.WriteHeader(http.StatusMethodNotAllowed) and w.Write([]byte("Method Not Allowed"))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idQueryParamV, err := strconv.Atoi(r.URL.Query().Get(idQueryParam))
	if err != nil || idQueryParamV <= 0 {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello world from view snippet number " + strconv.Itoa(idQueryParamV)))
}

func CreateSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Allow header must be sent back to the client with http.StatusMethodNotAllowed(405).
		// More info: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Allow
		w.Header().Set("Allow", http.MethodPost)
		w.WriteHeader(http.StatusMethodNotAllowed)
		// Headers that are set after a WriteHeader or Write method call, does not take any effect on the response
		// Code: w.Header().Set("Content-Type", "application/json") // Does not take any effect
		// It does not make sense to call WriteHeader twice in a request, because after the first one, the status is not modified
		// Code: w.WriteHeader(http.StatusOK) // 2023/06/21 23:39:22 http: superfluous response.WriteHeader call from main.book.func2 (main.go:49)
		w.Write([]byte("method not allowed"))
		// Code: w.Header().Set("Content-Type", "application/json") // Does not take any effect
		return
	}

	w.Write([]byte("Hello world from create snippet"))
}

func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts, err := template.ParseFiles("./ui/html/pages/home.gohtml", "./ui/html/partials/nav.gohtml", "./ui/html/base.gohtml")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Optional(w http.ResponseWriter, r *http.Request) {
	// Order matters, we have to set up as a first arguments our definitions
	ts, err := template.ParseFiles("./ui/html/pages/optional.gohtml", "./ui/html/partials/nav.gohtml", "./ui/html/base.gohtml")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func book(mux *http.ServeMux) {
	mux.HandleFunc("/", Home)
	// mux.Handle("/static/", middlewareDoNotAllowFetchingFolders(http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static")))))
	mux.Handle("/static", http.NotFoundHandler())
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(CustomSystem{ fs: http.Dir(cfg.staticDir) })))
	mux.Handle("/counter", new(counter))
	mux.HandleFunc("/optional", Optional)
	mux.HandleFunc("/snippet/view", ViewSnippet)
	mux.HandleFunc("/snippet/create", CreateSnippet)
}

// This is not the best solution, because if we have an index.html file inside one of these folders, attacking {endpoint_folder}/ is not going
// to redirect us to {endpoint_folder}/index.html
func middlewareDoNotAllowFetchingFolders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

type CustomSystem struct {
	fs http.FileSystem
}

func (c CustomSystem) Open(path string) (http.File, error) {
	f, err := c.fs.Open(path)
	if err != nil {
		return nil, err
	}

	os, err := f.Stat()
	if os.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err = c.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, err
			}
			return nil, err
		}
	}

	return f, nil
}
