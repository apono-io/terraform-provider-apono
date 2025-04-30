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

const MAX_DEBUG_BODY_LENGTH = 32000

// DebugTransport is an HTTP transport wrapper that logs detailed error information,
// including request and response bodies for HTTP errors, to aid in debugging and testing.
// During acceptance tests, it also logs the request body to stderr.
type DebugTransport struct {
	Transport http.RoundTripper
}

func (t *DebugTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	isAcceptanceTest := os.Getenv("TF_ACC") != ""

	var requestBodyStr string
	if isAcceptanceTest && req.Body != nil {
		bodyBytes, _ := io.ReadAll(req.Body)
		req.Body.Close()

		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		requestBodyStr = string(bodyBytes)
		requestBodyStr = truncateString(requestBodyStr, MAX_DEBUG_BODY_LENGTH)
	}

	resp, err := t.Transport.RoundTrip(req)
	if err != nil {
		werr := fmt.Errorf("HTTP Transport Error: %w", err)
		tflog.Error(req.Context(), werr.Error())

		return resp, werr
	}

	if resp.StatusCode >= 400 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			werr := fmt.Errorf("Failed to read error response body: %w", err)
			tflog.Error(req.Context(), werr.Error())

			return resp, werr
		}

		_ = resp.Body.Close()
		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		bodyStr := string(bodyBytes)
		bodyStr = truncateString(bodyStr, MAX_DEBUG_BODY_LENGTH)

		errorDetails := fmt.Sprintf(
			"API Error Response:\nURL: %s\nMethod: %s\nStatus: %s\nBody: %s",
			req.URL.String(), req.Method, resp.Status, bodyStr,
		)

		tflog.Error(req.Context(), errorDetails)

		if isAcceptanceTest {
			if requestBodyStr != "" {
				fmt.Fprintf(os.Stderr, "\n[DEBUG] Request Body: %s\n", requestBodyStr)
			}
		}

		return resp, errors.New(errorDetails)
	}

	return resp, err
}

func truncateString(s string, maxLength int) string {
	s = strings.TrimSpace(s)
	if len(s) > maxLength {
		return s[:maxLength] + "... [truncated]"
	}
	return s
}
