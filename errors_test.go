package mailosaur

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	baseUrl = os.Getenv("MAILOSAUR_BASE_URL")

	if len(baseUrl) == 0 {
		baseUrl = "https://mailosaur.com/"
	}
}

func TestUnauthorized(t *testing.T) {
	client := New("invalid_key")
	client.baseUrl = baseUrl
	_, err := client.Servers.List()

	assert.Error(t, err)
	assert.Equal(t, "Authentication failed, check your API key.", err.Error())
}

func TestNotFound(t *testing.T) {
	client := New()
	client.baseUrl = baseUrl
	_, err := client.Servers.Get("not_found")

	assert.Error(t, err)
	assert.Equal(t, "Not found, check input parameters.", err.Error())
}

func TestBadRequest(t *testing.T) {
	client := New()
	client.baseUrl = baseUrl
	serverCreateOptions := ServerCreateOptions{}

	_, err := client.Servers.Create(serverCreateOptions)

	assert.Error(t, err)
	assert.Equal(t, "(name) Servers need a name\r\n", err.Error())
}
