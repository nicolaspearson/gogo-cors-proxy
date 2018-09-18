package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

var (
	// host:port to proxy requests to
	target = "localhost:8080"
	// host:port to listen on
	listen = "localhost:8181"
	// protocol used by the target
	protocol = "http"
	// host header to be used for the proxy request
	host = "localhost:3000"
	// origin header to be used for the proxy request
	origin = "http://localhost:3000"
	// enable / disable default access control methods
	methods = true
	// enable / disable debug messages
	debug = false
)

// handleReverseRequest writes back the server response to the client.
// If an "OPTIONS" request is called, we only return Access-Control-Allow-*
func handleReverseRequest(w http.ResponseWriter, r *http.Request) {
	// build the url
	toCall := fmt.Sprintf("%s://%s%s", protocol, target, r.URL.String())
	logger("Create request for ", toCall)

	if len(origin) == 0 {
		for n, h := range r.Header {
			// get the origin from the request
			if strings.Contains(n, "Origin") {
				for _, h := range h {
					origin = h
				}
			}
		}
	}

	// always allow access origin
	w.Header().Add("Access-Control-Allow-Origin", origin)
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	if methods {
		w.Header().Add("Access-Control-Allow-Methods", "GET, PUT, POST, HEAD, TRACE, DELETE, PATCH, COPY, HEAD, LINK, OPTIONS")
	}

	if r.Method == "OPTIONS" {
		logger("CORS asked for ", toCall)
		for n, h := range r.Header {
			if strings.Contains(n, "Access-Control-Request") {
				for _, h := range h {
					k := strings.Replace(n, "Request", "Allow", 1)
					w.Header().Add(k, h)
				}
			}
		}
		return
	}

	// create the request to server
	req, err := http.NewRequest(r.Method, toCall, r.Body)

	// add ALL headers to the connection
	for n, h := range r.Header {
		for _, h := range h {
			req.Header.Add(n, h)
		}
	}

	// use the host provided by the flag
	if len(host) > 0 {
		req.Host = host
	}

	// create a basic client to send the request
	client := http.Client{}
	if r.TLS != nil {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	for h, v := range resp.Header {
		for _, v := range v {
			w.Header().Add(h, v)
		}
	}
	// copy the response from the server to the connected client request
	w.WriteHeader(resp.StatusCode)

	wr, err := io.Copy(w, resp.Body)
	if err != nil {
		log.Println(wr, err)
	} else {
		logger("Written", wr, "bytes")
	}

}

// validateFlags checks if host:port format is correct
func validateFlags() {
	for _, f := range []string{target, listen} {
		if !strings.Contains(f, ":") {
			log.Fatalf("%s is not incorrect, you must use a colon (:) to separate the host and port", f)
		}
	}

	parts := strings.Split(target, ":")
	if parts[0] == "" {
		log.Println("You didn't set the target host to connect, using localhost:" + parts[1])
		target = "localhost:" + parts[1]
	}

	if protocol != "http" && protocol != "https" {
		log.Fatalf(`Protocol can only be "http" or "https", not %q`, protocol)
	}
}

// logger prints messages to the console when the debug flag is set to true
func logger(v ...interface{}) {
	if debug {
		log.Println(v...)
	}
}

func main() {
	flag.StringVar(&target, "target", target, "host:port to proxy requests to")
	flag.StringVar(&listen, "listen", listen, "host:port to listen on")
	flag.StringVar(&protocol, "protocol", protocol, "protocol used by the target")
	flag.StringVar(&host, "host", host, "host header to be used for the proxy request")
	flag.StringVar(&origin, "origin", origin, "origin header to be used for the proxy request")
	flag.BoolVar(&methods, "methods", methods, "enable / disable default access control methods")
	flag.BoolVar(&debug, "debug", debug, "enable / disable debug messages")
	flag.Parse()
	validateFlags()

	http.HandleFunc("/", handleReverseRequest)
	log.Println(listen, "-->", target)
	http.ListenAndServe(listen, nil)
}
