package mailosaur

import (
	"time"
)

// UsageService provides operations for inspecting your account's usage limits and recent
// transactional usage. These endpoints require authentication with an account-level API key.
// Accessed via the Usage field of MailosaurClient.
type UsageService struct {
	client *MailosaurClient
}

type UsageAccountLimit struct {
	Limit   int `json:"limit"`
	Current int `json:"current"`
}

type UsageAccountLimits struct {
	Servers *UsageAccountLimit `json:"servers"`
	Users   *UsageAccountLimit `json:"users"`
	Email   *UsageAccountLimit `json:"email"`
	Sms     *UsageAccountLimit `json:"sms"`
}

type UsageTransaction struct {
	Timestamp time.Time `json:"timestamp"`
	Email     int       `json:"email"`
	Sms       int       `json:"sms"`
}

type UsageTransactionListResult struct {
	Items []*UsageTransaction `json:"items"`
}

// Limits retrieves account usage limits, detailing the current limits and usage for your
// account. This endpoint requires authentication with an account-level API key. It returns
// the UsageAccountLimits for your account.
func (s *UsageService) Limits() (*UsageAccountLimits, error) {
	result, err := s.client.HttpGet(&UsageAccountLimits{}, "api/usage/limits")
	return result.(*UsageAccountLimits), err
}

// Transactions retrieves the last 31 days of transactional usage. This endpoint requires
// authentication with an account-level API key. It returns a UsageTransactionListResult for
// the last 31 days.
func (s *UsageService) Transactions() (*UsageTransactionListResult, error) {
	result, err := s.client.HttpGet(&UsageTransactionListResult{}, "api/usage/transactions")
	return result.(*UsageTransactionListResult), err
}
