package mailosaur

import (
	"time"
)

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

func (s *UsageService) Limits() (*UsageAccountLimits, error) {
	result, err := s.client.HttpGet(&UsageAccountLimits{}, "api/usage/limits")
	return result.(*UsageAccountLimits), err
}

func (s *UsageService) Transactions() (*UsageTransactionListResult, error) {
	result, err := s.client.HttpGet(&UsageTransactionListResult{}, "api/usage/transactions")
	return result.(*UsageTransactionListResult), err
}
