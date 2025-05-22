// Package auth provides an HTTP transport that handles HMAC authentication
// for the Cerberus API. This is used in conjunction with the generated API client.
package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

// HMACAuthTransport is an http.RoundTripper that injects Cerberus API specific
// HMAC authentication headers into each outgoing request.
// It wraps an existing http.RoundTripper, typically http.DefaultTransport or
// the transport used by the go-openapi runtime client.
type HMACAuthTransport struct {
	APIKey    string             // APIKey is the Cerberus API Key.
	APISecret string             // APISecret is the Cerberus API Secret.
	Transport http.RoundTripper // Transport is the underlying transport to delegate requests to.
}

// NewHMACAuthTransport creates a new HMACAuthTransport.
// It takes the API key, API secret, and an optional next http.RoundTripper.
// If next is nil, http.DefaultTransport will be used as the underlying transport.
// This transport is then used to build the go-openapi runtime client.
func NewHMACAuthTransport(apiKey, apiSecret string, next http.RoundTripper) *HMACAuthTransport {
	if next == nil {
		next = http.DefaultTransport
	}
	return &HMACAuthTransport{
		APIKey:    apiKey,
		APISecret: apiSecret,
		Transport: next,
	}
}

// RoundTrip implements the http.RoundTripper interface.
// It calculates the HMAC signature, adds the required authentication headers
// (X-API-Key, X-Timestamp, X-Signature) to the request, and then delegates
// the request to the underlying Transport.
// It also ensures that the Content-Type header is set to "application/json"
// for requests that have a body and don't have this header set.
func (t *HMACAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone the request to avoid modifying the original request.
	reqClone := req.Clone(req.Context())

	// Get current UNIX timestamp as a string.
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Construct the message for HMAC signature: timestamp + apiKey.
	message := timestamp + t.APIKey

	// Calculate HMAC-SHA256 signature.
	mac := hmac.New(sha256.New, []byte(t.APISecret))
	mac.Write([]byte(message))
	signature := hex.EncodeToString(mac.Sum(nil))

	// Add authentication headers to the cloned request.
	reqClone.Header.Set("X-API-Key", t.APIKey)
	reqClone.Header.Set("X-Timestamp", timestamp)
	reqClone.Header.Set("X-Signature", signature)

	// Set Content-Type to application/json if there's a body and it's not already set.
	// The generated client should typically handle this for POST/PUT requests with a body,
	// but it's a good safeguard.
	if reqClone.Body != nil && reqClone.Header.Get("Content-Type") == "" {
		reqClone.Header.Set("Content-Type", "application/json")
	}

	// Delegate the request to the nested RoundTripper.
	return t.Transport.RoundTrip(reqClone)
}
