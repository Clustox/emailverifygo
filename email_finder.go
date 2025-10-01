package emailverifygo

import (
	"fmt"
	"net/url"
)

// FindEmailResponse response structure for the Email Finder API
type FindEmailResponse struct {
	Email  string `json:"email"`  // The email address found, or "null" if not found
	Status string `json:"status"` // "found" or "not_found"
}

// IsFound checks if an email was found
func (f *FindEmailResponse) IsFound() bool {
	return f.Status == "found"
}

// FindEmail uses a combined name parameter and domain to find a valid business email
//
// Parameters:
//   - name: The full name of the person to search for (e.g., "John Smith")
//   - domain: The domain to search for the email on
//
// Returns:
//   - *FindEmailResponse: The email finder result
//   - error: Any error that occurred during the request
//
// API Reference: GET /api/v1/finder
func FindEmail(name, domain string) (*FindEmailResponse, error) {
	if name == "" || domain == "" {
		return nil, fmt.Errorf("Both name and domain are required")
	}
	
	response := &FindEmailResponse{}

	request_parameters := url.Values{}
	request_parameters.Set("name", name)
	request_parameters.Set("domain", domain)

	url_to_request, err := PrepareURL(ENDPOINT_EMAIL_FINDER, request_parameters)
	if err != nil {
		return response, fmt.Errorf("failed to prepare URL: %w", err)
	}

	err = DoGetRequest(url_to_request, response)
	return response, err
}
