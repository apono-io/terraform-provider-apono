package client

import "net/http"

// UserAgentTransport adds User-Agent header to requests.
type UserAgentTransport struct {
	UserAgent string
	Transport http.RoundTripper
}

func (t *UserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.UserAgent)
	return t.Transport.RoundTrip(req)
}
