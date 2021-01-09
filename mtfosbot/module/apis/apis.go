package apis

import (
	"io"
	"net/http"
)

// RequestObj -
type RequestObj struct {
	Method  string
	URL     string
	Body    io.Reader
	Headers map[string]string
}

// GetRequest -
func GetRequest(r RequestObj) (req *http.Request, err error) {

	req, err = http.NewRequest(r.Method, r.URL, r.Body)
	if err != nil {
		return
	}

	if len(r.Headers) > 0 {
		for k, v := range r.Headers {
			req.Header.Set(k, v)
		}
	}

	return
}
