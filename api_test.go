package emailverifygo

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	
	// Set a test API key
	SetApiKey("test_api_key")
	
	// Mock responses for validation endpoint
	httpmock.RegisterResponder("GET", `=~^(.*)`+ENDPOINT_VALIDATE+`(.*)\z`,
		func(r *http.Request) (*http.Response, error) {
			args := r.URL.Query()
			if args.Get("key") == "" {
				return httpmock.NewStringResponse(400, `{"error": "Missing parameter: key."}`), nil
			}
			
			email := args.Get("email")
			if email == "valid@example.com" {
				return httpmock.NewStringResponse(200, MOCK_VALID_RESPONSE), nil
			} else {
				return httpmock.NewStringResponse(200, MOCK_INVALID_RESPONSE), nil
			}
		},
	)
	
	t.Run("TestValidateValidEmail", func(t *testing.T) {
		result, err := Validate("valid@example.com")
		
		assert.Nil(t, err, "Expected no error")
		assert.True(t, result.IsValid(), "Expected email to be valid")
		assert.Equal(t, "valid@example.com", result.Email, "Expected email to match")
		assert.Equal(t, "valid", result.Status, "Expected status to be 'valid'")
		assert.Equal(t, "permitted", result.SubStatus, "Expected sub_status to be 'permitted'")
	})
	
	t.Run("TestValidateInvalidEmail", func(t *testing.T) {
		result, err := Validate("invalid@example.com")
		
		assert.Nil(t, err, "Expected no error")
		assert.False(t, result.IsValid(), "Expected email to be invalid")
		assert.Equal(t, "invalid", result.Status, "Expected status to be 'invalid'")
		assert.Equal(t, "mailbox_not_found", result.SubStatus, "Expected sub_status to be 'mailbox_not_found'")
	})
	
	t.Run("TestValidateWithInvalidAPIKey", func(t *testing.T) {
		// Set an invalid API key
		SetApiKey("")
		
		_, err := Validate("test@example.com")
		assert.NotNil(t, err, "Expected error with missing API key")
		assert.Contains(t, err.Error(), "API key not set", "Expected error about missing API key")
		
		// Reset API key for other tests
		SetApiKey("test_api_key")
	})
}
