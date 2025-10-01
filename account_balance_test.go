package emailverifygo

import (
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAccountBalance(t *testing.T) {
	// Activate httpmock
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	
	// Set a test API key
	SetApiKey("test_api_key")
	
	// Mock responses for account balance endpoint
	httpmock.RegisterResponder("GET", `=~^(.*)`+ENDPOINT_ACCOUNT_BALANCE+`(.*)\z`,
		func(r *http.Request) (*http.Response, error) {
			args := r.URL.Query()
			if args.Get("key") == "" {
				return httpmock.NewStringResponse(400, `{"error": "Missing parameter: key."}`), nil
			}
			
			return httpmock.NewStringResponse(200, MOCK_ACCOUNT_BALANCE_RESPONSE), nil
		},
	)

	t.Run("TestGetAccountBalanceSuccess", func(t *testing.T) {
		result, err := GetAccountBalance()
		
		assert.Nil(t, err, "Expected no error")
		assert.Equal(t, "enabled", result.APIStatus, "Expected API status to be 'enabled'")
		assert.Equal(t, 150, result.DailyCreditsLimit, "Expected daily credits limit to be 150")
		assert.Equal(t, 15000, result.RemainingCredits, "Expected remaining credits to be 15000")
	})
	
	t.Run("TestGetAccountBalanceWithInvalidAPIKey", func(t *testing.T) {
		// Set an invalid API key
		SetApiKey("")
		
		_, err := GetAccountBalance()
		assert.NotNil(t, err, "Expected error with missing API key")
		assert.Contains(t, err.Error(), "API key not set", "Expected error about missing API key")
		
		// Reset API key for other tests
		SetApiKey("test_api_key")
	})
}
