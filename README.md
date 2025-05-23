[![Go Test](https://github.com/cerberius-soft/go/actions/workflows/go-test.yml/badge.svg?branch=main)](https://github.com/cerberius-soft/go/actions/workflows/go-test.yml)

# Cerberius API Go Client

This Go client library allows you to easily interact with the [Cerberius API documentation](https://cerberius.com/docs/). It handles the required HMAC-SHA256 authentication and provides Go-native methods and models for all API endpoints.

## Features

*   **Complete API Coverage:** Provides methods for all Cerberius API operations:
    *   Email Validation
    *   IP Lookup
    *   Prompt Check
*   **Pre-configured Authentication:** Includes built-in HMAC-SHA256 authentication compliant with Cerberius API requirements.
*   **Go-native Experience:** Offers Go models for all API requests and responses, making it easy to integrate into your Go applications.
*   **Clear Usage Patterns:** Comes with straightforward examples and clear error handling patterns (see `examples/main.go`).

## Installation

To use this client in your Go project, you can fetch it using `go get`:

```sh
go get cerberius.com/go-client
```

## Authentication

The Cerberius API requires HMAC-SHA256 authentication. This client simplifies this by providing an `auth.HMACAuthTransport`, which is a standard Go `http.RoundTripper`. You configure it once with your API credentials, and it automatically adds the necessary authentication headers to all outgoing requests.

You need to provide your API Key and API Secret when creating the transport:

```go
import (
    "net/http"
    "cerberius.com/go-client/auth" // For HMACAuthTransport
    "cerberius.com/go-client/generated/client" // For the main API client and default values
    httptransport "github.com/go-openapi/runtime/client"
    "github.com/go-openapi/strfmt"
)

// Best practice: Load APIKey and APISecret from environment variables or a secure config manager.
apiKey := "YOUR_API_KEY"
apiSecret := "YOUR_API_SECRET"

// 1. Create the custom HMAC authentication transport
// This transport wraps a standard HTTP transport (like http.DefaultTransport)
// and adds the HMAC authentication headers.
hmacAuthTransport := auth.NewHMACAuthTransport(apiKey, apiSecret, http.DefaultTransport)

// 2. Create an HTTP client that uses our custom authentication transport
httpClient := &http.Client{
    Transport: hmacAuthTransport,
}

// 3. Create the go-openapi runtime transport using our custom HTTP client
// client.DefaultHost, client.DefaultBasePath, and client.DefaultSchemes are constants
// imported from the cerberius.com/go-client/generated/client package.
runtimeTransport := httptransport.NewWithClient(
    client.DefaultHost,     // Default: "service.cerberius.com"
    client.DefaultBasePath, // Default: "/api"
    client.DefaultSchemes,  // Default: []string{"https"}
    httpClient,
)

// 4. Create the Cerberius API client
// This client provides methods for all API operations.
apiClient := client.New(runtimeTransport, strfmt.Default)
```

With `apiClient` initialized, you can now make calls to the Cerberius API services.

## Usage Example

Here's a basic example of how to call the Email Validation endpoint using the `apiClient` configured above:

```go
import (
    "fmt"
    "log"
    "cerberius.com/go-client/generated/client/operations" // For operation-specific parameters
    "cerberius.com/go-client/generated/models"          // For request and response models
    // ... other imports from authentication setup, and your apiClient instance
)

func main() {
    // ... (apiClient initialization as shown in the Authentication section)

    // Prepare parameters for the email validation request
    emailParams := operations.NewEmailValidationRequestDataParams().
        WithBody(&models.EmailLookupRequest{Data: []string{"test@example.com", "another@example.org"}})

    // Call the API
    emailResp, err := apiClient.Operations.EmailValidationRequestData(emailParams)
    if err != nil {
        // For detailed error handling, including parsing structured error responses
        // from the API, please refer to the handleAPIError function in examples/main.go.
        log.Fatalf("Error calling EmailValidationRequestData: %v", err)
    }

    // Process the successful response
    fmt.Println("Email Validation Response:")
    for _, ed := range emailResp.Payload.Data {
        fmt.Printf("  Email: %s, Domain: %s, Valid Score: %d, Comment: %s\n",
            ed.EmailAddress, ed.Domain, ed.ValidityScore, ed.Comment)
    }
    if emailResp.Payload.ExcessChargesApply {
        fmt.Println("  Note: Excess charges may apply for this request.")
    }
}
```

For more comprehensive examples, including how to call all API operations and detailed error handling, please see the `examples/main.go` file in this repository.

## API Operations

The client provides convenient methods for all Cerberius API operations:

*   **Email Validation**: `apiClient.Operations.EmailValidationRequestData(...)`
    *   Validates a list of email addresses.
*   **IP Lookup**: `apiClient.Operations.IPLookupRequestData(...)`
    *   Looks up information on IP addresses.
*   **Prompt Check**: `apiClient.Operations.PromptCheckRequestData(...)`
    *   Checks if a given prompt is malicious.

For detailed information on request and response structures, refer to the model structs in the `cerberius.com/go-client/generated/models` package and the operation parameter types in `cerberius.com/go-client/generated/client/operations`.

## Error Handling

When an API call fails, the methods return an `error`. If the error originates from the API server (e.g., validation error, unauthorized), it can be type-asserted to access structured details.
The `examples/main.go` file contains a `handleAPIError` function demonstrating best practices for this:

1.  Type-assert the `error` to `*runtime.APIError` (from `github.com/go-openapi/runtime`).
2.  The `APIError.Response` field contains the specific error response object. Type-assert this to the operation's default error type (e.g., `*operations.EmailValidationRequestDataDefault`).
3.  The `Payload` of this specific error type is an `interface{}`. Type-assert this `Payload` to `*models.Response`.
4.  The `*models.Response` struct then contains the `Error` field (of type `*models.Data`), which holds the API's `Code` and `Message`.

This structured error information is crucial for robust error handling in your application.
```
