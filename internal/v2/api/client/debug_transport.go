package client

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// DebugTransport wraps an http.RoundTripper and logs response bodies on error.
type DebugTransport struct {
	Transport http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface.
func (t *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Check if we're running acceptance tests
	isAccTest := os.Getenv("TF_ACC") != ""

	// Perform the request
	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		// Log transport-level errors
		logMessage := fmt.Sprintf("HTTP Transport Error: %v", err)
		tflog.Error(req.Context(), logMessage)

		// When running acceptance tests, output to stderr for test visibility
		if isAccTest {
			fmt.Fprintf(os.Stderr, "\n[DEBUG] %s\n", logMessage)
		}

		return resp, err
	}

	// If the status code indicates an error, log the response body
	if resp.StatusCode >= 400 {
		// Read the response body
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			// If we can't read the body, log that issue and return the original response
			logMessage := fmt.Sprintf("Failed to read error response body: %v", readErr)
			tflog.Error(req.Context(), logMessage)

			if isAccTest {
				fmt.Fprintf(os.Stderr, "\n[DEBUG] %s\n", logMessage)
			}

			return resp, readErr
		}

		// Close the original body
		_ = resp.Body.Close()

		// Create a new body with the same content for the response
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Trim any extremely long bodies to avoid flooding logs
		bodyStr := string(bodyBytes)
		if len(bodyStr) > 1000 {
			bodyStr = bodyStr[:1000] + "... [truncated]"
		}

		// Clean up the body string (remove excessive whitespace)
		bodyStr = strings.TrimSpace(bodyStr)

		// Create a detailed error message
		errorDetail := fmt.Sprintf(
			"API Error Response:\nURL: %s\nMethod: %s\nStatus: %s\nBody: %s",
			req.URL.String(), req.Method, resp.Status, bodyStr,
		)

		// Log with tflog
		tflog.Error(req.Context(), errorDetail)

		// When running acceptance tests, ensure errors are visible in test output
		if isAccTest {
			fmt.Fprintf(os.Stderr, "\n[DEBUG] %s\n", errorDetail)
		}
	}

	return resp, err
}
