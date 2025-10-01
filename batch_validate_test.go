package emailverifygo

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestBatchValidation(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	
	// Set a test API key
	SetApiKey("test_api_key")
	
	// Mock responses for batch validation endpoint
	httpmock.RegisterResponder("POST", `=~^(.*)`+ENDPOINT_VALIDATE_BATCH+`(.*)\z`,
		func(r *http.Request) (*http.Response, error) {
			// Validate the request body
			decoder := json.NewDecoder(r.Body)
			var requestBody struct {
				Key        string        `json:"key"`
				EmailBatch []interface{} `json:"email_batch"`
			}
			if err := decoder.Decode(&requestBody); err != nil {
				return httpmock.NewStringResponse(400, `{"error": "Invalid request body"}`), nil
			}
			
			// Check if API key is missing
			if requestBody.Key == "" {
				return httpmock.NewStringResponse(400, `{"error": "Missing parameter: key."}`), nil
			}
			
			// Check if email batch is empty
			if len(requestBody.EmailBatch) == 0 {
				return httpmock.NewStringResponse(400, `{"error": "No emails provided"}`), nil
			}
			
			// Return successful response
			return httpmock.NewStringResponse(200, MOCK_BATCH_RESPONSE), nil
		},
	)
	
	// Mock responses for batch results endpoint
	httpmock.RegisterResponder("GET", `=~^(.*)`+ENDPOINT_BATCH_RESULT+`(.*)\z`,
		func(r *http.Request) (*http.Response, error) {
			args := r.URL.Query()
			if args.Get("key") == "" {
				return httpmock.NewStringResponse(400, `{"error": "Missing parameter: key."}`), nil
			}
			
			return httpmock.NewStringResponse(200, MOCK_BATCH_RESULTS_RESPONSE), nil
		},
	)

	t.Run("TestValidateBatchSuccess", func(t *testing.T) {
		emails := []string{
			"valid@example.com",
			"invalid@example.com",
			"test@example.com",
		}
		
		result, err := ValidateBatch("Test Batch", emails)
		
		assert.Nil(t, err, "Expected no error")
		assert.Equal(t, "success", result.Status, "Expected status to be 'success'")
		assert.Equal(t, 12345, result.TaskID, "Expected task ID to be 12345")
		assert.Equal(t, 3, result.CountSubmitted, "Expected 3 emails submitted")
	})
	
	t.Run("TestGetBatchResultsSuccess", func(t *testing.T) {
		result, err := GetBatchResults(12345)
		
		assert.Nil(t, err, "Expected no error")
		assert.Equal(t, 12345, result.TaskID, "Expected task ID to be 12345")
		assert.Equal(t, "verified", result.Status, "Expected status to be 'verified'")
		assert.Equal(t, "Test Batch", result.Name, "Expected name to be 'Test Batch'")
		assert.Equal(t, 3, len(result.Results.EmailBatch), "Expected 3 email results")
		
		// Check first email result
		assert.Equal(t, "valid@example.com", result.Results.EmailBatch[0].Address, "Expected first email to be valid@example.com")
		assert.Equal(t, "valid", result.Results.EmailBatch[0].Status, "Expected first email status to be valid")
		assert.Equal(t, "permitted", result.Results.EmailBatch[0].SubStatus, "Expected first email sub_status to be permitted")
	})
	
	t.Run("TestValidateBatchWithInvalidAPIKey", func(t *testing.T) {
		// Store original API key
		originalKey := API_KEY
		
		// Set an invalid API key
		SetApiKey("")
		
		emails := []string{"test@example.com"}
		_, err := ValidateBatch("Test Batch", emails)
		
		assert.NotNil(t, err, "Expected error with missing API key")
		// Just check that we get any error, don't assert on the specific message
		// since it can vary between "API key not set" and "Missing parameter: key"
		
		// Reset API key for other tests
		SetApiKey(originalKey)
	})
	
	t.Run("TestValidateBatchWithEmptyList", func(t *testing.T) {
		emails := []string{}
		_, err := ValidateBatch("Empty Batch", emails)
		
		assert.NotNil(t, err, "Expected error with empty email list")
		assert.Contains(t, err.Error(), "Email", "Expected error about empty email list")
	})
	
	t.Run("TestGetBatchResultsWithInvalidAPIKey", func(t *testing.T) {
		// Store original API key
		originalKey := API_KEY
		
		// Set an invalid API key
		SetApiKey("")
		
		_, err := GetBatchResults(12345)
		
		assert.NotNil(t, err, "Expected error with missing API key")
		// Just check that we get any error, don't assert on the specific message
		
		// Reset API key for other tests
		SetApiKey(originalKey)
	})
}
