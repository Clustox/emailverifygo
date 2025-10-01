# EmailVerify.io Go Package

A Go package for interacting with EmailVerify.io services. This package provides a simple and efficient way to validate email addresses, perform batch validations, find email addresses by name and domain, and check your account balance.

## Installation

```bash
go get https://github.com/Clustox/emailverifygo
```

## Features

1. **Email Validation**: Verify the validity of individual email addresses
2. **Batch Email Validation**: Validate multiple email addresses at once
3. **Email Finder**: Find email addresses by name and domain
4. **Account Balance**: Check your account balance and credits

## API Key Configuration

This package uses the EmailVerify API which requires an API key. This key can be provide in three ways:
1. Through an environment variable EMAIL_VERIFY_API_KEY (loaded automatically in code)
2. Through an .env file that contains EMAIL_VERIFY_API_KEY and then calling following method before usage:

```go
emailverifygo.LoadEnvFromFile()
```

3. By settings explicitly in code, using the following method:

```go
emailverifygo.SetApiKey("<YOUR_API_KEY>")
```

## Usage Examples


### Check Account Balance

Check your account balance and credit information.

```go
package main

import (
	"fmt"

	"github.com/Clustox/emailverifygo"
)

func main() {
    emailverifygo.SetApiKey("<YOUR_API_KEY>")

	response, error_ := emailverifygo.GetAccountBalance()

	if error_ != nil {
		fmt.Println("error occurred: ", error_.Error())
	} else {
		fmt.Printf("Full response: %+v\n", response)
		fmt.Println("API Status:", response.APIStatus)
		fmt.Println("Daily Credits Limit:", response.DailyCreditsLimit)
		fmt.Println("Referral Credits:", response.ReferralCredits) // For appsumo accounts, ReferralCredits will be 0
		fmt.Println("Remaining Credits:", response.RemainingCredits)
		fmt.Println("Remaining Daily Credits:", response.RemainingDailyCredits) // For non-appsumo accounts, RemainingDailyCredits will be 0
		fmt.Println("Bonus Credits:", response.BonusCredits) // For non-appsumo accounts, BonusCredits will be 0
	}
}
```

### Email Validation

Validate a single email address to check if it's valid, invalid, or has other status flags.

```go
package main

import (
	"fmt"
	"github.com/Clustox/emailverifygo"
)

func main() {
	emailverifygo.SetApiKey("<YOUR_API_KEY>")
	
	response, error_ := emailverifygo.Validate("possible_typo@example.com")
	
	if error_ != nil {
		fmt.Println("error occurred: ", error_.Error())
	} else {
		// Now you can check status
		fmt.Println("Response", response)
		fmt.Println("Email", response.Email)
		fmt.Println("Status", response.Status)
		fmt.Println("Substatus", response.SubStatus)

		if response.Status == emailverifygo.STATUS_DO_NOT_MAIL {
			fmt.Println("This email's status is do not email")
		}

		if response.SubStatus == emailverifygo.SUBSTATUS_MAILBOX_QUOTA_EXCEEDED {
			fmt.Println("This email's sub status is mailbox quota exceeded")
		}
	}
}
```

We have below mentioned Constants you can use to check variable status and substatus

```go
// Status constants
const (
	STATUS_VALID       = "valid"       // The email is valid and deliverable
	STATUS_INVALID     = "invalid"     // The email is invalid or undeliverable
	STATUS_CATCH_ALL   = "catch_all"   // The domain has a catch-all policy
	STATUS_DO_NOT_MAIL = "do_not_mail" // The email should not be mailed to
	STATUS_UNKNOWN     = "unknown"     // The status could not be determined
	STATUS_ROLE_BASED  = "role_based"  // The email is a role-based address (e.g., info@, support@)
	STATUS_SKIPPED     = "skipped"     // The validation was skipped for this email
)

// Sub-status constants
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
```

### Batch Email Validation

Submit multiple emails for validation in a single batch operation.

```go
package main

import (
	"fmt"

	"github.com/Clustox/emailverifygo"
)

func main() {
	emailverifygo.SetApiKey("<YOUR_API_KEY>")

	emails := []string{
		"user1@example.com",
		"user2@example.com",
		"user3@example.com",
		"user4@example.com",
	}
	
	response, error_ := emailverifygo.ValidateBatch("<Title>", emails) //Tittle and emails are required fields

	if error_ != nil {
		fmt.Println("error occurred: ", error_.Error())
	} else {
		fmt.Println("Response", response)
		fmt.Println("Status", response.Status)
		fmt.Println("TaskID", response.TaskID) // IMPORTANT! SAVE IT to later fetch the results
		fmt.Println("CountSubmitted", response.CountSubmitted)
		fmt.Println("CountRejected", response.CountRejected)
		fmt.Println("CountProcessing", response.CountProcessing)
		fmt.Println("Count Duplicate Removed", response.CountDuplicatesRemoved)
	}
}
```

### Retrieve Batch Validation Results

Retrieve the results of a previously submitted batch validation using the Task ID.

```go
package main
import (
	"fmt"
	"github.com/Clustox/emailverifygo"
)
	
func main() {
	emailverifygo.SetApiKey("<YOUR_API_KEY>")

	response, error_ := emailverifygo.GetBatchResults(<TASK_ID>) // TASK ID received in bulk validate, must be integer

	if error_ != nil {
		fmt.Println("error occurred: ", error_.Error())
	} else {
		fmt.Println("Response", response)
		fmt.Println("CountChecked", response.CountChecked)
		fmt.Println("Count Total", response.CountTotal)
		fmt.Println("Title", response.Name)
		fmt.Println("Status", response.Status)
		fmt.Println("Task Id", response.TaskID)
		fmt.Println("Progress Percentage", response.ProgressPercentage)
		fmt.Println("Results", response.Results) // results can be {[]} when verification under process

		// Way to fetch emails and their result
		for _, res := range response.Results.EmailBatch {
     	   fmt.Printf("Email: %s, Status: %s, SubStatus: %s\n", res.Address, res.Status, res.SubStatus)
        }
	}
}
```

### Find Email by Name and Domain

Find email addresses associated with a person at a specific domain.

```go
package main

import (
	"fmt"

	"github.com/Clustox/emailverifygo"
)

func main() {
	emailverifygo.SetApiKey("<YOUR_API_KEY>")

	response, error_ := emailverifygo.FindEmail("<NAME>", "DOMAIN.COM")

	if error_ != nil {
		fmt.Println("error occurred: ", error_.Error())
	} else {
		fmt.Println("Response", response)
		fmt.Println("Email", response.Email) // WILL BE 'null' WHEN EMAIL NOT FOUND
		fmt.Println("Status", response.Status) // STATUS CAN BE found or not_found
	}
}

```

## Testing

The package includes several levels of tests to ensure everything works correctly:

### Unit Tests

Run the basic unit tests (uses mocked responses, no API key required):

```bash
go test -v
```