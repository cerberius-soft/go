[![Go Test](https://github.com/YOUR_OWNER/YOUR_REPO/actions/workflows/go-test.yml/badge.svg)](https://github.com/YOUR_OWNER/YOUR_REPO/actions/workflows/go-test.yml)

# Cerberius API Go Client

This repository contains a Go client library for interacting with the [Cerberius API](https://service.cerberius.com/api/schema).
The client is generated from the official OpenAPI v2 schema and includes support for the required HMAC-SHA256 authentication.

## Features

- Client code generated using `go-swagger`.
- All API operations supported:
    - Email Validation (`/email-lookup`)
    - IP Lookup (`/ip-lookup`)
    - Prompt Check (`/prompt-check`)
- Custom HMAC-SHA256 authentication transport.
- Models for all API requests and responses.
- Example usage and comprehensive error handling demonstration.

## Installation

To use this client in your Go project, you can fetch it using `go get`:

```sh
go get cerberius.com/go-client
```

## Authentication

The Cerberius API requires HMAC-SHA256 authentication. This client provides an `http.RoundTripper` implementation (`auth.HMACAuthTransport`) that handles this automatically.

You need to provide your API Key and API Secret when creating the transport:

```go
import (
    "net/http"
    "cerberius.com/go-client/auth" // For HMACAuthTransport
    "cerberius.com/go-client/generated/client" // For the main API client
    httptransport "github.com/go-openapi/runtime/client"
    "github.com/go-openapi/strfmt"
)

// Best practice: Load from environment variables or a secure config manager.
apiKey := "YOUR_API_KEY"
apiSecret := "YOUR_API_SECRET"

// Create the custom HMAC authentication transport
hmacAuthTransport := auth.NewHMACAuthTransport(apiKey, apiSecret, http.DefaultTransport)

// Create an HTTP client that uses our custom authentication transport
httpClient := &http.Client{
    Transport: hmacAuthTransport,
}

// Create the go-openapi runtime transport using our custom HTTP client
// client.DefaultHost and client.DefaultBasePath are from cerberius.com/go-client/generated/client
runtimeTransport := httptransport.NewWithClient(
    client.DefaultHost,     // e.g., "service.cerberius.com"
    client.DefaultBasePath, // e.g., "/api"
    client.DefaultSchemes,  // e.g., []string{"https"}
    httpClient,
)

// Create the Cerberius API client
apiClient := client.New(runtimeTransport, strfmt.Default)
```

## Usage Example

Here's a basic example of how to call the Email Validation endpoint:

```go
import (
    "fmt"
    "log"
    "cerberius.com/go-client/generated/client/operations"
    "cerberius.com/go-client/generated/models"
    // ... other imports from authentication setup
)

func main() {
    // ... (apiClient initialization as shown in Authentication section)

    emailParams := operations.NewEmailValidationRequestDataParams().
        WithBody(&models.EmailLookupRequest{Data: []string{"test@example.com"}})

    emailResp, err := apiClient.Operations.EmailValidationRequestData(emailParams)
    if err != nil {
        // See examples/main.go for detailed error handling, including how to
        // parse structured error responses from the API.
        log.Fatalf("Error calling EmailValidationRequestData: %v", err)
    }

    fmt.Println("Email Validation Response:")
    for _, ed := range emailResp.Payload.Data {
        fmt.Printf("  Email: %s, Domain: %s, Valid Score: %d, Comment: %s\n",
            ed.EmailAddress, ed.Domain, ed.ValidityScore, ed.Comment)
    }
}
```

For more detailed examples, including how to call all API operations and perform comprehensive error handling, please see the `examples/main.go` file in this repository.

## API Operations

The client provides methods for the following API operations:

*   **Email Validation**: `apiClient.Operations.EmailValidationRequestData(...)`
    *   Validates a list of email addresses.
*   **IP Lookup**: `apiClient.Operations.IPLookupRequestData(...)`
    *   Looks up information on IP addresses.
*   **Prompt Check**: `apiClient.Operations.PromptCheckRequestData(...)`
    *   Checks if a given prompt is malicious.

Refer to the parameters and response model structs in the `cerberius.com/go-client/generated/models` package and the operation parameters in `cerberius.com/go-client/generated/client/operations` for more details on request and response structures.

## Error Handling

API errors are returned as `error` types. To access structured error details from the Cerberius API (like error code and message), you'll need to type-assert the error.
The `examples/main.go` file contains a `handleAPIError` function that demonstrates how to do this for each operation. Typically, this involves:
1. Type-asserting the error to `*runtime.APIError`.
2. Type-asserting `APIError.Response` to the specific operation's default error type (e.g., `*operations.EmailValidationRequestDataDefault`).
3. Accessing the `Payload` of this default error type, which should be `*models.Response`, containing the API's `Code` and `Message`.

```
