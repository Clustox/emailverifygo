package emailverifygo

import (
	"net/url"
)

// AccountBalanceResponse represents the response from the account balance API
// Different fields may be present depending on the account type.
type AccountBalanceResponse struct {
	APIStatus            string `json:"api_status"`
	DailyCreditsLimit    int    `json:"daily_credits_limit"`
	ReferralCredits      int    `json:"referral_credits,omitempty"` // absent for appsumo users
	RemainingCredits     int    `json:"remaining_credits,omitempty"` // absent for appsumo users
	RemainingDailyCredits int   `json:"remaining_daily_credits,omitempty"`
	BonusCredits         int    `json:"bonus_credits,omitempty"`
}

// GetAccountBalance gets the current account balance and credits information
//
// Returns:
//   - *AccountBalanceResponse: The account balance information
//   - error: Any error that occurred during the request
//
// API Reference: GET /api/v1/check-account-balance


func GetAccountBalance() (*AccountBalanceResponse, error) {
	var error_ error
	response := &AccountBalanceResponse{}

	// Prepare URL with API key
	url_to_request, error_ := PrepareURL(ENDPOINT_ACCOUNT_BALANCE, url.Values{})
	if error_ != nil {
		return response, error_
	}
	
	// Make the request
	error_ = DoGetRequest(url_to_request, response)
	return response, error_
}
