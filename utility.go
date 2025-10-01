package emailverifygo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// API endpoints
const (
	ENDPOINT_VALIDATE              = "/api/v1/validate"
	ENDPOINT_VALIDATE_BATCH        = "/api/v1/validate-batch"
	ENDPOINT_EMAIL_FINDER          = "/api/v1/finder"
	ENDPOINT_ACCOUNT_BALANCE       = "/api/v1/check-account-balance"
	ENDPOINT_BATCH_RESULT          = "/api/v1/get-result-bulk-verification-task"
)

// Email validation status constants
const (
	STATUS_VALID       = "valid"       // The email is valid and deliverable
	STATUS_INVALID     = "invalid"     // The email is invalid or undeliverable
	STATUS_CATCH_ALL   = "catch_all"   // The domain has a catch-all policy
	STATUS_DO_NOT_MAIL = "do_not_mail" // The email should not be mailed to
	STATUS_UNKNOWN     = "unknown"     // The status could not be determined
	STATUS_ROLE_BASED  = "role_based"  // The email is a role-based address (e.g., info@, support@)
	STATUS_SKIPPED     = "skipped"     // The validation was skipped for this email
)

// Email validation sub-status constants
const (
	SUBSTATUS_PERMITTED            = "permitted"             // Email is permitted for sending
	SUBSTATUS_FAILED_SYNTAX_CHECK  = "failed_syntax_check"   // Email failed syntax validation
	SUBSTATUS_MAILBOX_QUOTA_EXCEEDED = "mailbox_quota_exceeded" // Mailbox is full
	SUBSTATUS_MAILBOX_NOT_FOUND    = "mailbox_not_found"     // Mailbox does not exist
	SUBSTATUS_NO_DNS_ENTRIES       = "no_dns_entries"        // Domain has no DNS entries
	SUBSTATUS_DISPOSABLE           = "disposable"            // Email is from a disposable domain
	SUBSTATUS_NONE                 = "none"                  // No specific sub-status
	SUBSTATUS_OPT_OUT              = "opt_out"               // User has opted out
	SUBSTATUS_BLOCKED_DOMAIN       = "blocked_domain"        // Domain is blocked
)

// APIResponse is an interface for all API response types
type APIResponse interface{}

// ErrMissingAPIKey is returned when the API key is not set
var ErrMissingAPIKey = errors.New("API key not set. Use SetApiKey() or LoadEnvFromFile() to set it")

// Getenv gets an environment variable or returns a default value if it's not set
func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return fallback
}

// Global variables
var (
	// URI used to make requests to the EmailVerify API
	URI = Getenv(`EMAIL_VERIFY_URI`, `https://app.emailverify.io`)

	// API_KEY the API key used in order to make the requests
	API_KEY string = os.Getenv("EMAIL_VERIFY_API_KEY")
)

// GetBaseURI returns the current base URI for the API
func GetBaseURI() string {
	return URI
}

// SetApiKey sets the API key for all future requests
func SetApiKey(newApiKey string) {
	API_KEY = newApiKey
}

// SetURI updates the base URI for the API
// Useful for testing against staging environments
func SetURI(newURI string) {
	if newURI != "" {
		URI = newURI
	}
}

// LoadEnvFromFile loads environment variables from a .env file and sets the API key and URI.
// It returns true if the file was loaded successfully, false otherwise.
// It won't print any messages to stdout, so it's suitable for library use.
func LoadEnvFromFile() bool {
	err := godotenv.Load(".env")
	if err != nil {
		// Just return false without printing anything
		return false
	}

	SetApiKey(os.Getenv("EMAIL_VERIFY_API_KEY"))
	SetURI(os.Getenv("EMAIL_VERIFY_URI"))
	return true
}

// PrepareURL prepares the URL for a request by attaching the API key and params
func PrepareURL(endpoint string, params url.Values) (string, error) {
	if API_KEY == "" {
		return "", ErrMissingAPIKey
	}

	// Set API KEY
	params.Set("key", API_KEY)

	// Create and return the final URL
	finalURL, err := url.JoinPath(URI, endpoint)
	if err != nil {
		return "", fmt.Errorf("failed to join URL paths: %w", err)
	}
	return fmt.Sprintf("%s?%s", finalURL, params.Encode()), nil
}

// ErrorFromResponse parses an error response from the API
// This handles the inconsistent error message formats returned by the API
func ErrorFromResponse(response *http.Response) error {
	// ERROR handling: expect a json payload containing details about the error
	var errorResponse map[string]string

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to read error response: %w", err)
	}
	
	err = json.NewDecoder(strings.NewReader(string(responseBody))).Decode(&errorResponse)

	if err != nil {
		// unexpected non-json payload
		return fmt.Errorf("API error (status %d): %s", response.StatusCode, string(responseBody))
	}

	// return all possible details about the error
	var errorStrings []string
	for _, value := range errorResponse {
		errorStrings = append(errorStrings, value)
	}
	return fmt.Errorf("API error %d: %s", response.StatusCode, strings.Join(errorStrings, ", "))
}

// DoGetRequest performs a GET request to the API
func DoGetRequest(url string, object APIResponse) error {
	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Do the request using the current HTTP client
	response, err := httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	// Close the request
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return ErrorFromResponse(response)
	}

	// Decode JSON Request
	err = json.NewDecoder(response.Body).Decode(&object)
	if err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}
	return nil
}

// DoPostRequest performs a POST request to the API
func DoPostRequest(url string, payload io.Reader, object APIResponse) error {
	// Create a new request
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set content type
	req.Header.Set("Content-Type", "application/json")
	
	// Do the request using the current HTTP client
	response, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}

	// Close the request
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return ErrorFromResponse(response)
	}

	// Decode JSON Request
	err = json.NewDecoder(response.Body).Decode(&object)
	if err != nil {
		return fmt.Errorf("failed to decode JSON response: %w", err)
	}
	return nil
}
