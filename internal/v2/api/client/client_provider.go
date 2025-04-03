package client

// ClientProvider is an interface for accessing the API client.
// This avoids cyclic dependencies between packages.
type ClientProvider interface {
	// V2Client returns the V2 API client
	V2Client() *Client
}
