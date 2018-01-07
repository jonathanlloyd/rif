package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const port = 8080

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		_, err := fmt.Fprintf(w, formatRequest(req))
		if err != nil {
			panic("Error writing HTTP response: " + err.Error())
		}
	})
	fmt.Printf("Listening on port %d\n", port)
	listenAddr := fmt.Sprintf(":%d", port)
	err := http.ListenAndServe(listenAddr, mux)
	if err != nil {
		panic("Error starting HTTP server: " + err.Error())
	}
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string

	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)

	// Add the host
	request = append(request, fmt.Sprintf("host: %v", r.Host))

	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// Add body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic("Error reading body")
	}
	if len(body) > 0 {
		request = append(request, fmt.Sprintf("\n%s", string(body)))
	}

	// Return the request as a string
	return strings.Join(request, "\n") + "\r\n"
}
