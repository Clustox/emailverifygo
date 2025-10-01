package emailverifygo

import (
	"net/http"
)

// HTTPClient interface for making HTTP requests
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Global HTTP client
var (
	// Default HTTP client
	defaultHTTPClient HTTPClient = &http.Client{}
	
	// Current HTTP client (used for actual requests)
	httpClient HTTPClient = defaultHTTPClient
)

// Mock responses for testing (these are just examples and will be replaced by httpmock)
const (
	MOCK_VALID_RESPONSE = `{
		"email": "valid@example.com",
		"status": "valid",
		"sub_status": "permitted"
	}`
	
	MOCK_INVALID_RESPONSE = `{
		"email": "invalid@example.com",
		"status": "invalid",
		"sub_status": "mailbox_not_found"
	}`
	
	MOCK_BATCH_RESPONSE = `{
		"status": "success",
		"task_id": 12345,
		"count_submitted": 3,
		"count_duplicates_removed": 0,
		"count_rejected_emails": 0,
		"count_processing": 3
	}`
	
	MOCK_BATCH_RESULTS_RESPONSE = `{
		"count_checked": 3,
		"count_total": 3,
		"name": "Test Batch",
		"progress_percentage": 100,
		"task_id": 12345,
		"status": "verified",
		"results": {
			"email_batch": [
				{
					"address": "valid@example.com",
					"status": "valid",
					"sub_status": "permitted"
				},
				{
					"address": "invalid@example.com",
					"status": "invalid",
					"sub_status": "mailbox_not_found"
				},
				{
					"address": "test@example.com",
					"status": "invalid", 
					"sub_status": "no_dns_entries"
				}
			]
		}
	}`
	
	MOCK_FINDER_RESPONSE = `{
		"email": "john.doe@example.com",
		"status": "found"
	}`
	
	MOCK_FINDER_NOT_FOUND_RESPONSE = `{
		"email": "null",
		"status": "not_found"
	}`
	
	MOCK_ACCOUNT_BALANCE_RESPONSE = `{
		"api_status": "enabled",
		"daily_credits_limit": 150,
		"remaining_credits": 15000
	}`
	
	MOCK_ERROR_RESPONSE = `{
		"error": "Invalid API key"
	}`
)