package pancors

import (
	"bufio"
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/AlchemyTelcoSolutions/safehttp"
	"github.com/jellydator/ttlcache/v3"
)

type corsTransport struct {
	referer     string
	origin      string
	credentials string
	cache       *ttlcache.Cache[string, string]
}

func NewCorsTransport(referer string, origin string, credentials string, cache *ttlcache.Cache[string, string]) corsTransport {
	return corsTransport{
		referer,
		origin,
		credentials,
		cache,
	}
}

func (t corsTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	item := t.cache.Get(r.URL.String())
	if item != nil {
		reader := bufio.NewReader(bytes.NewReader([]byte(item.Value())))
		response, err := http.ReadResponse(reader, r)
		if err != nil {
			return nil, err
		}
		return response, nil
	}

	// Put in the Referer if specified
	if t.referer != "" {
		r.Header.Add("Referer", t.referer)
	}

	client := safehttp.NewClient(safehttp.Options{})

	// Do the actual request
	res, err := client.Transport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	res.Header.Set("Access-Control-Allow-Origin", t.origin)
	res.Header.Set("Access-Control-Allow-Credentials", t.credentials)

	buffer := bytes.NewBuffer([]byte{})
	writer := bufio.NewWriter(buffer)
	err = res.Write(writer)
	if err != nil {
		return nil, err
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	t.cache.Set(r.URL.String(), buffer.String(), ttlcache.DefaultTTL)

	reader := bufio.NewReader(bytes.NewReader(buffer.Bytes()))
	response, err := http.ReadResponse(reader, r)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func handleProxy(w http.ResponseWriter, r *http.Request, origin string, credentials string, cache *ttlcache.Cache[string, string]) {
	// Check for the User-Agent header
	if r.Header.Get("User-Agent") == "" {
		http.Error(w, "Missing User-Agent header", http.StatusBadRequest)
		return
	}

	// Get the optional Referer header
	referer := r.URL.Query().Get("referer")
	if referer == "" {
		referer = r.Header.Get("Referer")
	}

	// Get the URL
	urlParam := r.URL.Query().Get("url")
	// Validate the URL
	urlParsed, err := url.Parse(urlParam)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	// Check if HTTP(S)
	if urlParsed.Scheme != "http" && urlParsed.Scheme != "https" {
		http.Error(w, "The URL scheme is neither HTTP nor HTTPS", http.StatusBadRequest)
		return
	}

	// Setup for the proxy
	proxy := httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL = urlParsed
			r.Host = urlParsed.Host
		},
		Transport: NewCorsTransport(referer, origin, credentials, cache),
	}

	// Execute the request
	proxy.ServeHTTP(w, r)
}

// HandleProxy is a handler which passes requests to the host and returns their
// responses with CORS headers
func HandleProxy(w http.ResponseWriter, r *http.Request) {
	cache := ttlcache.New(
		ttlcache.WithTTL[string, string](2 * time.Minute),
	)

	go cache.Start() // starts automatic expired item deletion

	handleProxy(w, r, "*", "true", cache)
}

// HandleProxyFromHosts is a handler which passes requests only from specified to the host
func HandleProxyWith(origin string, credentials string) func(http.ResponseWriter, *http.Request) {
	if !(credentials == "true" || credentials == "false") {
		log.Panicln("Access-Control-Allow-Credentials can only be 'true' or 'false'")
	}

	cache := ttlcache.New(
		ttlcache.WithTTL[string, string](2 * time.Minute),
	)

	go cache.Start() // starts automatic expired item deletion

	return func(w http.ResponseWriter, r *http.Request) {
		handleProxy(w, r, origin, credentials, cache)
	}
}
