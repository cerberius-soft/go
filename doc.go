// Package cerberusapiclient provides a Go client for interacting with the Cerberius API.
//
// The client is generated from a Swagger/OpenAPI specification and includes
// models for all API requests and responses, as well as operations for each
// API endpoint.
//
// It also includes a custom HTTP transport (`cerberusclient.HMACAuthTransport`)
// to handle the HMAC-SHA256 authentication required by the Cerberius API.
//
// # Installation
//
// To use this client in your Go project, you can fetch it using `go get`:
//
//	go get github.com/your-repo-path/cerberusapiclient
//
// (Replace `github.com/your-repo-path/cerberusapiclient` with the actual import path
// once the repository is available.)
//
// # Usage
//
// Here's a basic example of how to initialize and use the client:
//
//	package main
//
//	import (
//		"fmt"
//		"log"
//		"net/http"
//		"os"
//
//		"cerberusapiclient/cerberusclient" // Authentication transport
//		"cerberusapiclient/gen/client"    // Generated API client
//		"cerberusapiclient/gen/client/operations"
//		"cerberusapiclient/gen/models"
//
//		httptransport "github.com/go-openapi/runtime/client"
//		"github.com/go-openapi/strfmt"
//	)
//
//	func main() {
//		apiKey := os.Getenv("CERBERUS_API_KEY")
//		apiSecret := os.Getenv("CERBERUS_API_SECRET")
//
//		if apiKey == "" || apiSecret == "" {
//			log.Fatal("CERBERUS_API_KEY and CERBERUS_API_SECRET environment variables must be set.")
//		}
//
//		// Create the HMAC authentication transport, wrapping the default HTTP transport
//		authTransport := cerberusclient.NewHMACAuthTransport(apiKey, apiSecret, http.DefaultTransport)
//
//		// Create an HTTP client that uses our custom authentication transport
//		httpClient := &http.Client{
//			Transport: authTransport,
//		}
//
//		// Create the go-openapi runtime transport using our custom HTTP client
//		// client.DefaultHost and client.DefaultBasePath are from the generated client package
//		transport := httptransport.NewWithClient(client.DefaultHost, client.DefaultBasePath, client.DefaultSchemes, httpClient)
//
//		// Create the Cerberius API client
//		apiClient := client.New(transport, strfmt.Default)
//
//		// Example: Call Email Validation
//		emailParams := operations.NewEmailValidationRequestDataParams().
//			WithBody(&models.EmailLookupRequest{Data: []string{"test@example.com"}})
//
//		emailResp, err := apiClient.Operations.EmailValidationRequestData(emailParams)
//		if err != nil {
//			// Handle error (see examples/main.go for detailed error handling)
//			log.Fatalf("Error calling EmailValidationRequestData: %v", err)
//		}
//
//		fmt.Printf("Email Validation Response: %+v\n", emailResp.Payload)
//	}
//
// For more detailed examples, including comprehensive error handling, see the
// `examples/main.go` file in this repository.
//
// # Authentication
//
// Authentication is handled by `cerberusclient.HMACAuthTransport`. This transport
// automatically adds the following headers to each request:
//   - X-API-Key: Your API key.
//   - X-Timestamp: The current UNIX timestamp.
//   - X-Signature: The HMAC-SHA256 signature of (timestamp + api_key) using your API secret.
//
// You must provide your API key and secret when creating the `HMACAuthTransport`.
package cerberusapiclient
