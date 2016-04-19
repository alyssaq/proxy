package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

var (
	// flagPort is the open port the application listens on
	flagPort = flag.String("port", "9000", "Port to listen on")
	// flagProxyURL is the url to proxy
	flagProxyURL = flag.String("proxyURL", "http://localhost", "Proxy URL")
)

type proxyTransport struct {
	Transport *http.Transport
}

func (t *proxyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	// Add additional headers to the proxied response
	resp.Header.Set("Vary", "Origin, Accept-Encoding")
	resp.Header.Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	resp.Header.Set("Access-Control-Allow-Origin", "*")
	return resp, nil
}

func proxyHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Override host to the host of the target server not the proxy
		r.Host = r.URL.Host
		handler.ServeHTTP(w, r)
	})
}

func init() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	flag.Parse()
}

func main() {
	target, err := url.Parse(*flagProxyURL)
	if err != nil {
		log.Fatal(err)
	}
	reverseProxy := httputil.NewSingleHostReverseProxy(target)
	reverseProxy.Transport = &proxyTransport{&http.Transport{}}

	mux := http.NewServeMux()
	mux.Handle("/", proxyHandler(reverseProxy))
	log.Printf("listening on port %s", *flagPort)
	log.Fatal(http.ListenAndServe(":"+*flagPort, mux))
}
