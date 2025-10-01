package emailverifygo

import (
	"fmt"
	"net/url"
)

// ValidateResponse contains the fields that are returned by the API
// for single email validation requests.
type ValidateResponse struct {
	Email     string `json:"email"`     // The email address being validated
	Status    string `json:"status"`    // Status of the email (valid, invalid, etc.)
	SubStatus string `json:"sub_status"` // Detailed status information
}

// IsValid returns true if the email status is "valid"
func (v *ValidateResponse) IsValid() bool {
	return v.Status == STATUS_VALID
}

// Validate performs validation on a single email address
// 
// Parameters:
//   - email: The email address to validate
//
// Returns:
//   - *ValidateResponse: The validation result
//   - error: Any error that occurred during validation
//
// API Reference: GET /api/v1/validate
func Validate(email string) (*ValidateResponse, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}
	
	// Prepare the parameters
	params := url.Values{}
	params.Set("email", email)

	response := &ValidateResponse{}

	// Do the request
	url_to_request, err := PrepareURL(ENDPOINT_VALIDATE, params)
	if err != nil {
		return response, fmt.Errorf("failed to prepare URL: %w", err)
	}
	
	err = DoGetRequest(url_to_request, response)
	return response, err
}
