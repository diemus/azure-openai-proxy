package openai

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func NewOpenAIReverseProxy() *httputil.ReverseProxy {
	remote, _ := url.Parse("https://api.openai.com")
	director := func(req *http.Request) {
		// Set the Host, Scheme, Path, and RawPath of the request to the remote host and path
		originURL := req.URL.String()
		req.Host = remote.Host
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host

		log.Printf("proxying request %s -> %s", originURL, req.URL.String())
	}
	return &httputil.ReverseProxy{Director: director}
}
