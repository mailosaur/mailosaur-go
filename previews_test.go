package mailosaur

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	apiKey := os.Getenv("MAILOSAUR_API_KEY")
	baseUrl := os.Getenv("MAILOSAUR_BASE_URL")
	server = os.Getenv("MAILOSAUR_SERVER")

	if len(apiKey) == 0 {
		log.Fatal("Missing necessary environment variables - refer to README.md")
	}

	if len(baseUrl) == 0 {
		baseUrl = "https://next.mailosaur.com/"
	}

	client = New(apiKey)
	client.baseUrl = baseUrl
}

func TestListEmailClients(t *testing.T) {
	result, err := client.Previews.ListEmailClients()
	assert.NoError(t, err)

	assert.True(t, len(result.Items) > 1)
}

func TestGenerateEmailPreviews(t *testing.T) {
	if len(server) == 0 {
		t.Skip()
	}

	randomString := getRandomString()
	host := os.Getenv("MAILOSAUR_SMTP_HOST")
	if len(host) == 0 {
		host = "mailosaur.net"
	}

	testEmailAddress := fmt.Sprintf("%s@%s.%s", randomString, server, host)

	sendEmail(client, server, testEmailAddress)

	email, _ := client.Messages.Get(&MessageSearchParams{
		Server: server,
	}, &SearchCriteria{
		SentTo: testEmailAddress,
	})

	request := &PreviewRequest{EmailClient: "OL2021"}
	options := &PreviewRequestOptions{Previews: []*PreviewRequest{request}}

	result, _ := client.Messages.GeneratePreviews(email.Id, options)
	assert.True(t, len(result.Items) > 0)

	bytes, _ := client.Files.GetPreview(result.Items[0].Id)

	assert.True(t, len(bytes) > 1)
}
