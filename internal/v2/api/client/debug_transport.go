package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// DebugTransport is an HTTP transport wrapper that logs detailed error information,
// including response bodies for HTTP errors, to aid in debugging and testing.
type DebugTransport struct {
	Transport http.RoundTripper
}

func (t *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	isAcceptanceTest := os.Getenv("TF_ACC") != ""

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		logMessage := fmt.Sprintf("HTTP Transport Error: %v", err)
		tflog.Error(req.Context(), logMessage)

		if isAcceptanceTest {
			fmt.Fprintf(os.Stderr, "\n[DEBUG] %s\n", logMessage)
		}

		return resp, err
	}

	if resp.StatusCode >= 400 {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			logMessage := fmt.Sprintf("Failed to read error response body: %v", readErr)
			tflog.Error(req.Context(), logMessage)

			if isAcceptanceTest {
				fmt.Fprintf(os.Stderr, "\n[DEBUG] %s\n", logMessage)
			}

			return resp, readErr
		}

		_ = resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		bodyStr := string(bodyBytes)
		if len(bodyStr) > 4096 {
			bodyStr = bodyStr[:4096] + "... [truncated]"
		}
		bodyStr = strings.TrimSpace(bodyStr)

		errorDetails := fmt.Sprintf(
			"API Error Response:\nURL: %s\nMethod: %s\nStatus: %s\nBody: %s",
			req.URL.String(), req.Method, resp.Status, bodyStr,
		)

		tflog.Error(req.Context(), errorDetails)

		return resp, errors.New(errorDetails)
	}

	return resp, err
}
