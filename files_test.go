package mailosaur

import (
	"testing"
    "log"
    "fmt"
    "os"
    "strings"
    assert "github.com/stretchr/testify/require"
)

func init() {
    apiKey := os.Getenv("MAILOSAUR_API_KEY")
    baseUrl := os.Getenv("MAILOSAUR_BASE_URL")
    server = os.Getenv("MAILOSAUR_SERVER")

    if (len(apiKey) == 0 || len(server) == 0) {
        log.Fatal("Missing necessary environment variables - refer to README.md")
    }

    if (len(baseUrl) == 0) {
        baseUrl = "https://next.mailosaur.com"
    }

    client = New(apiKey)
	client.baseUrl = baseUrl

    client.Messages.DeleteAll(server)

    host := os.Getenv("MAILOSAUR_SMTP_HOST")
    if (len(host) == 0) {
        host = "mailosaur.net"
    }

    testEmailAddress := fmt.Sprintf("wait_for_test@%s.%s", server, host)

    sendEmail(client, server, testEmailAddress)

    result, _ := client.Messages.Get(&MessageSearchParams {
        Server: server,
    }, &SearchCriteria {
        SentTo: testEmailAddress,
    })

    email = result
}

func TestFilesGetEmail(t *testing.T) {
    bytes, _ := client.Files.GetEmail(email.Id)

    assert.True(t, len(bytes) > 1)
    assert.True(t, strings.Contains(string(bytes), email.Subject))
}

func TestFilesGetAttachment(t *testing.T) {
    attachment := email.Attachments[0]
    bytes, _ := client.Files.GetAttachment(attachment.Id)

    assert.True(t, len(bytes) > 1)
    assert.Equal(t, attachment.Length, len(bytes))
}
