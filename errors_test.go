package mailosaur

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnauthorized(t *testing.T) {
	client := New(os.Getenv("invalid_key"))
	_, err := client.Servers.List()

	assert.Error(t, err)
	assert.Equal(t, "Authentication failed, check your API key.", err.Error())
}

func TestNotFound(t *testing.T) {
	client := New(os.Getenv("MAILOSAUR_API_KEY"))
	_, err := client.Servers.Get("not_found")

	assert.Error(t, err)
	assert.Equal(t, "Not found, check input parameters.", err.Error())
}

func TestBadRequest(t *testing.T) {
	client := New(os.Getenv("MAILOSAUR_API_KEY"))
	serverCreateOptions := ServerCreateOptions{}

	_, err := client.Servers.Create(serverCreateOptions)

	assert.Error(t, err)
	assert.Equal(t, "(name) Servers need a name\r\n", err.Error())
}
