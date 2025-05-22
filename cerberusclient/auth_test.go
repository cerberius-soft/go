package cerberusclient

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

// mockRoundTripper is a helper/mock for http.RoundTripper.
type mockRoundTripper struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
	called        bool          // To track if RoundTrip was called
	request       *http.Request // To store the request for inspection
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.called = true
	m.request = req
	if m.RoundTripFunc != nil {
		return m.RoundTripFunc(req)
	}
	// Default response if not overridden
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("")), // http.NoBody is deprecated
		Header:     make(http.Header),
	}, nil
}

func TestNewHMACAuthTransportDefaults(t *testing.T) {
	transport := NewHMACAuthTransport("testKey", "testSecret", nil)
	if transport.Transport != http.DefaultTransport {
		t.Errorf("Expected Transport to be http.DefaultTransport, got %T", transport.Transport)
	}
}

func TestNextTransportCalled(t *testing.T) {
	mockNext := &mockRoundTripper{}
	authTransport := NewHMACAuthTransport("testKey", "testSecret", mockNext)

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	_, err := authTransport.RoundTrip(req)
	if err != nil {
		t.Fatalf("RoundTrip failed: %v", err)
	}

	if !mockNext.called {
		t.Error("Expected next transport's RoundTrip method to be called, but it wasn't")
	}
}

func TestSignatureAndHeaders(t *testing.T) {
	apiKey := "testAPIKey"
	apiSecret := "testAPISecret"
	mockNext := &mockRoundTripper{}

	authTransport := NewHMACAuthTransport(apiKey, apiSecret, mockNext)

	// Test case 1: Request without body
	t.Run("RequestWithoutBody", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "http://example.com/no-body", nil)
		_, err := authTransport.RoundTrip(req)
		if err != nil {
			t.Fatalf("RoundTrip failed: %v", err)
		}

		if mockNext.request == nil {
			t.Fatal("Next transport did not receive a request")
		}
		checkHeadersAndSignature(t, mockNext.request, apiKey, apiSecret, false)
	})

	// Test case 2: Request with body and no Content-Type
	t.Run("RequestWithBodyNoContentType", func(t *testing.T) {
		mockNext.request = nil // Reset for this sub-test
		reqBody := []byte(`{"key":"value"}`)
		req, _ := http.NewRequest("POST", "http://example.com/with-body", bytes.NewBuffer(reqBody))
		_, err := authTransport.RoundTrip(req)
		if err != nil {
			t.Fatalf("RoundTrip failed: %v", err)
		}

		if mockNext.request == nil {
			t.Fatal("Next transport did not receive a request")
		}
		checkHeadersAndSignature(t, mockNext.request, apiKey, apiSecret, true)
		if mockNext.request.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", mockNext.request.Header.Get("Content-Type"))
		}
	})

	// Test case 3: Request with body and existing Content-Type
	t.Run("RequestWithBodyExistingContentType", func(t *testing.T) {
		mockNext.request = nil // Reset for this sub-test
		reqBody := []byte(`{"key":"value"}`)
		req, _ := http.NewRequest("POST", "http://example.com/with-body-ct", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/xml") // Different content type
		_, err := authTransport.RoundTrip(req)
		if err != nil {
			t.Fatalf("RoundTrip failed: %v", err)
		}

		if mockNext.request == nil {
			t.Fatal("Next transport did not receive a request")
		}
		checkHeadersAndSignature(t, mockNext.request, apiKey, apiSecret, true) // body exists
		if mockNext.request.Header.Get("Content-Type") != "application/xml" {
			t.Errorf("Expected Content-Type 'application/xml' (to be preserved), got '%s'", mockNext.request.Header.Get("Content-Type"))
		}
	})
}

func checkHeadersAndSignature(t *testing.T, req *http.Request, apiKey, apiSecret string, hasBody bool) {
	// Check X-API-Key
	if req.Header.Get("X-API-Key") != apiKey {
		t.Errorf("Expected X-API-Key '%s', got '%s'", apiKey, req.Header.Get("X-API-Key"))
	}

	// Check X-Timestamp
	timestampStr := req.Header.Get("X-Timestamp")
	if timestampStr == "" {
		t.Fatal("X-Timestamp header is missing")
	}
	timestampVal, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		t.Fatalf("X-Timestamp '%s' is not a valid integer: %v", timestampStr, err)
	}
	// Check if timestamp is recent (e.g., within the last 5 minutes)
	if time.Now().Unix()-timestampVal > 300 {
		t.Errorf("X-Timestamp '%d' is too old", timestampVal)
	}

	// Check X-Signature
	signatureHeader := req.Header.Get("X-Signature")
	if signatureHeader == "" {
		t.Fatal("X-Signature header is missing")
	}

	// Calculate expected signature
	expectedMessage := timestampStr + apiKey
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(expectedMessage))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	if signatureHeader != expectedSignature {
		t.Errorf("Expected X-Signature '%s', got '%s'", expectedSignature, signatureHeader)
	}
}
