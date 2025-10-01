package emailverifygo

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestEmailFinder(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	
	// Set a test API key
	SetApiKey("test_api_key")
	
	// Mock responses for the email finder endpoint
	httpmock.RegisterResponder("GET", `=~^(.*)`+ENDPOINT_EMAIL_FINDER+`(.*)\z`,
		func(r *http.Request) (*http.Response, error) {
			args := r.URL.Query()
			if args.Get("key") == "" {
				return httpmock.NewStringResponse(400, `{"error": "Missing parameter: key."}`), nil
			}
			
			// Check if domain and name parameters are present
			domain := args.Get("domain")
			name := args.Get("name")
			
			if domain == "" {
				return httpmock.NewStringResponse(400, `{"error": "Missing parameter: domain."}`), nil
			}
			
			if name == "" {
				return httpmock.NewStringResponse(400, `{"error": "Missing parameter: name."}`), nil
			}
			
			// If all parameters are valid, return a successful response
			if domain == "example.com" && name == "John Doe" {
				return httpmock.NewStringResponse(200, MOCK_FINDER_RESPONSE), nil
			}
			
			// Default to not found response
			return httpmock.NewStringResponse(200, MOCK_FINDER_NOT_FOUND_RESPONSE), nil
		},
	)

	t.Run("TestFindEmailSuccess", func(t *testing.T) {
		result, err := FindEmail("John Doe", "example.com")
		
		assert.Nil(t, err, "Expected no error")
		assert.Equal(t, "john.doe@example.com", result.Email, "Expected email to be john.doe@example.com")
		assert.Equal(t, "found", result.Status, "Expected status to be 'found'")
		assert.True(t, result.IsFound(), "Expected IsFound() to return true")
	})
	
	t.Run("TestFindEmailNotFound", func(t *testing.T) {
		result, err := FindEmail("Jane Smith", "unknown.com")
		
		assert.Nil(t, err, "Expected no error")
		assert.Equal(t, "null", result.Email, "Expected email to be null")
		assert.Equal(t, "not_found", result.Status, "Expected status to be 'not_found'")
		assert.False(t, result.IsFound(), "Expected IsFound() to return false")
	})
	
	t.Run("TestFindEmailWithInvalidAPIKey", func(t *testing.T) {
		// Set an invalid API key
		SetApiKey("")
		
		_, err := FindEmail("John Doe", "example.com")
		
		assert.NotNil(t, err, "Expected error with missing API key")
		assert.Contains(t, err.Error(), "API key not set", "Expected error about missing API key")
		
		// Reset API key for other tests
		SetApiKey("test_api_key")
	})
	
	t.Run("TestFindEmailWithMissingParameters", func(t *testing.T) {
		// Test with missing domain
		_, err := FindEmail("John Doe", "")
		
		assert.NotNil(t, err, "Expected error with missing domain")
		assert.Contains(t, err.Error(), "domain", "Expected error about missing domain")
		
		// Test with missing name
		_, err = FindEmail("", "example.com")
		
		assert.NotNil(t, err, "Expected error with missing name")
		assert.Contains(t, err.Error(), "name", "Expected error about missing name")
	})
}
