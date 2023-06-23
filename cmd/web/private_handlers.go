package main

import (
	"net/http"
	"strings"
)

func General(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from subtree path!"))
}

func Greeting(w http.ResponseWriter, r *http.Request) {
	name := "World"

	if possibleName, ok := r.URL.Query()["name"]; ok {
		name = possibleName[0]
	}

	w.Write([]byte("Hello " + name))
}

func StaticSubTree(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/static/notfound" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte(strings.TrimPrefix(r.URL.Path, "/static")))
}

func HeadersTask(w http.ResponseWriter, r *http.Request) {
	var response string
	header := "Cache-Control"

	// When calling: Set, Add, Del and Get methods we are parsing the header string using textproto.CanonicalMIMEHeaderKey()
	// This is useful, because we can be case-insensitive.
	w.Header().Set(header, "public, max-age=31536000")
	response += "Get(Cache-Control) = " + w.Header().Get(header) + "\n\r"
	response += "Values(Cache-Control) = " + strings.Join(w.Header().Values(header), ", ") + "\n\r"

	w.Header().Del(header)
	// Del method doesn't remove system-generated headers
	w.Header().Del("Date")
	// If you want to suppres some system-generated headers, then you must set its value to nil directly
	w.Header()["Date"] = nil

	w.Header().Add(header, "public")
	w.Header().Add(header, "max-age=31536000")
	response += "Get(Cache-Control) = " + w.Header().Get(header) + "\n\r"
	response += "Values(Cache-Control) = " + strings.Join(w.Header().Values(header), ", ") + "\n\r"

	w.Header().Del(header)

	w.Write([]byte(response))
}
