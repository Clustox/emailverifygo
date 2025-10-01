package emailverifygo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// EmailAddress represents one email address unit to be validated in a batch
type EmailAddress struct {
	Address string `json:"address"`
}

// EmailBatchResult represents a single result in the batch validation
type EmailBatchResult struct {
	Address   string `json:"address"`
	Status    string `json:"status"`
	SubStatus string `json:"sub_status"`
}

// EmailBatchError an error unit received in the response, that can be associated
// with an email sent to the batch validate endpoint
type EmailBatchError struct {
	Error   string `json:"error"`
	Address string `json:"address"`
}

// BatchValidateRequest represents the structure of the batch validation request
type BatchValidateRequest struct {
	Title      string        `json:"title"`
	Key        string        `json:"key"`
	EmailBatch []EmailAddress `json:"email_batch"`
}

// BatchValidateResultsWrapper wraps the email_batch results from the API
type BatchValidateResultsWrapper struct {
	EmailBatch []EmailBatchResult `json:"email_batch"`
}

// BatchValidateResponse represents the structure of a batch validate response
type BatchValidateResponse struct {
	Status              string            `json:"status,omitempty"`     // For initial submit response
	TaskID              int               `json:"task_id"`              // Task ID for batch validation
	CountSubmitted      int               `json:"count_submitted,omitempty"` // For initial submit response
	CountDuplicatesRemoved     int               `json:"count_duplicates_removed,omitempty"` // For initial submit response
	CountRejected       int               `json:"count_rejected_emails,omitempty"` // For initial submit response
	CountProcessing     int               `json:"count_processing,omitempty"` // For initial submit response
}

type BatchResultResponse struct {
	CountChecked        int                     `json:"count_checked,omitempty"`
	CountTotal          int                     `json:"count_total,omitempty"`
	TaskID              int                     `json:"task_id,omitempty"`
	Name                string                  `json:"name,omitempty"`
	Status              string                  `json:"status,omitempty"`
	ProgressPercentage  float64                 `json:"progress_percentage,omitempty"`
	Results             BatchValidateResultsWrapper `json:"results,omitempty"`
}

// ValidateBatch submits a batch of emails for Verification
//
// Parameters:
//   - title: The name for this batch validation task
//   - emails: A slice of email addresses to validate
//
// Returns:
//   - *BatchValidateResponse: The initial batch submission response
//   - error: Any error that occurred during submission
//
// API Reference: POST /api/v1/validate-batch
func ValidateBatch(title string, emails []string) (*BatchValidateResponse, error) {
	response := &BatchValidateResponse{}
	
	if title == "" {
        return response, fmt.Errorf("Title is required")
    }

	// Check for empty email list
	if len(emails) == 0 {
		return response, fmt.Errorf("Email list cannot be empty")
	}
	
	// Prepare the email batch
	emailBatch := make([]EmailAddress, len(emails))
	for i, email := range emails {
		emailBatch[i] = EmailAddress{Address: email}
	}
	
	// Create the request payload
	requestData := BatchValidateRequest{
		Title:      title,
		Key:        API_KEY,
		EmailBatch: emailBatch,
	}
	
	// Convert the request data to JSON
	requestBody := &strings.Builder{}
	if err := json.NewEncoder(requestBody).Encode(requestData); err != nil {
		return response, fmt.Errorf("failed to encode request data: %w", err)
	}
	
	// Create request URL
	urlToRequest, err := url.JoinPath(URI, ENDPOINT_VALIDATE_BATCH)
	if err != nil {
		return response, fmt.Errorf("invalid URL (%s) or endpoint (%s): %w", URI, ENDPOINT_VALIDATE_BATCH, err)
	}
	
	// Make the POST request
	err = DoPostRequest(urlToRequest, strings.NewReader(requestBody.String()), response)
	return response, err
}

// GetBatchResults retrieves the results of a previously submitted batch validation task
//
// Parameters:
//   - taskID: The ID of the batch validation task to retrieve results for
//
// Returns:
//   - *BatchValidateResponse: The batch validation results
//   - error: Any error that occurred during retrieval
//
// API Reference: GET /api/v1/get-result-bulk-verification-task
func GetBatchResults(taskID int) (*BatchResultResponse, error) {
	response := &BatchResultResponse{}

	if taskID <= 0 {
		return response, fmt.Errorf("TaskID must be greater than 0")
	}
	
	// Prepare the parameters
	params := url.Values{}
	params.Set("task_id", fmt.Sprintf("%d", taskID))
	
	// Prepare URL with API key
	url_to_request, err := PrepareURL(ENDPOINT_BATCH_RESULT, params)
	if err != nil {
		return response, fmt.Errorf("failed to prepare URL: %w", err)
	}
	
	// Make the request
	err = DoGetRequest(url_to_request, response)
	return response, err
}
