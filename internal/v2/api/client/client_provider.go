package client

// ClientProvider is an interface for accessing the API client.
// This avoids cyclic dependencies between packages.
type ClientProvider interface {
	PublicClient() *Client
}
