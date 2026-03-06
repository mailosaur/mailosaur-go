package mailosaur

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	baseUrl := os.Getenv("MAILOSAUR_BASE_URL")

	if len(baseUrl) == 0 {
		baseUrl = "https://mailosaur.com/"
	}

	client = New()
	client.baseUrl = baseUrl
}

func TestLimits(t *testing.T) {
	result, err := client.Usage.Limits()
	assert.NoError(t, err)

	assert.True(t, result.Servers.Limit > 0)
	assert.True(t, result.Users.Limit > 0)
	assert.True(t, result.Email.Limit > 0)
	assert.True(t, result.Sms.Limit > 0)
}

func TestTransactions(t *testing.T) {
	result, err := client.Usage.Transactions()
	assert.NoError(t, err)

	assert.True(t, len(result.Items) > 1)
}
