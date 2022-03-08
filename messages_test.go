package mailosaur

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	apiKey := os.Getenv("MAILOSAUR_API_KEY")
	baseUrl := os.Getenv("MAILOSAUR_BASE_URL")
	server = os.Getenv("MAILOSAUR_SERVER")
	verifiedDomain = os.Getenv("MAILOSAUR_VERIFIED_DOMAIN")

	if len(verifiedDomain) == 0 {
		verifiedDomain = "mailosaur.net"
	}

	if len(apiKey) == 0 || len(server) == 0 {
		log.Fatal("Missing necessary environment variables - refer to README.md")
	}

	if len(baseUrl) == 0 {
		baseUrl = "https://next.mailosaur.com/"
	}

	client = New(apiKey)
	client.baseUrl = baseUrl

	client.Messages.DeleteAll(server)

	sendEmails(client, server, 5)
}

func TestMessageList(t *testing.T) {
	result, err := client.Messages.List(&MessageListParams{Server: server})
	assert.NoError(t, err)
	assert.Equal(t, 5, len(result.Items))

	for _, email := range result.Items {
		validateEmailSummary(t, email)
	}

	emails = result.Items
}

func TestMessageListReceivedAfter(t *testing.T) {
	pastDate := time.Now().Add(time.Duration(-10) * time.Minute)

	pastEmails, _ := client.Messages.List(&MessageListParams{
		Server:        server,
		ReceivedAfter: pastDate,
	})
	assert.True(t, len(pastEmails.Items) > 0)

	futureEmails, _ := client.Messages.List(&MessageListParams{
		Server:        server,
		ReceivedAfter: time.Now(),
	})
	assert.Equal(t, 0, len(futureEmails.Items))
}

func TestMessageGet(t *testing.T) {
	host := os.Getenv("MAILOSAUR_SMTP_HOST")
	if len(host) == 0 {
		host = "mailosaur.net"
	}

	testEmailAddress := fmt.Sprintf("wait_for_test@%s.%s", server, host)

	sendEmail(client, server, testEmailAddress)

	email, _ := client.Messages.Get(&MessageSearchParams{
		Server: server,
	}, &SearchCriteria{
		SentTo: testEmailAddress,
	})

	validateEmail(t, email)
}

func TestMessageGetById(t *testing.T) {
	emailToRetrieve := emails[0]
	email, _ := client.Messages.GetById(emailToRetrieve.Id)
	validateEmail(t, email)
	validateHeaders(t, email)
}

func TestMessageGetByIdNotFound(t *testing.T) {
	_, err := client.Messages.GetById("efe907e9-74ed-4113-a3e0-a3d41d914765")

	assert.Error(t, err)
	assert.IsType(t, &mailosaurError{}, err)
}

func TestSearchNoCriteriaError(t *testing.T) {
	_, err := client.Messages.Search(&MessageSearchParams{Server: server}, &SearchCriteria{})

	assert.Error(t, err)
	assert.IsType(t, &mailosaurError{}, err)
}

func TestSearchTimeoutErrorSuppressed(t *testing.T) {
	f := false
	result, _ := client.Messages.Search(&MessageSearchParams{
		Server:         server,
		Timeout:        1,
		ErrorOnTimeout: &f,
	}, &SearchCriteria{
		SentFrom: "neverfound@example.com",
	})

	assert.Equal(t, 0, len(result.Items))
}

func TestSearchBySentFrom(t *testing.T) {
	targetEmail := emails[1]

	result, _ := client.Messages.Search(&MessageSearchParams{
		Server: server,
	}, &SearchCriteria{
		SentFrom: targetEmail.From[0].Email,
	})

	assert.Equal(t, 1, len(result.Items))
	assert.Equal(t, targetEmail.From[0].Email, result.Items[0].From[0].Email)
	assert.Equal(t, targetEmail.Subject, result.Items[0].Subject)
}

func TestSearchBySentFromInvalidEmail(t *testing.T) {
	_, err := client.Messages.Search(&MessageSearchParams{
		Server: server,
	}, &SearchCriteria{
		SentFrom: ".not_an_email_address",
	})

	assert.Error(t, err)
	assert.IsType(t, &mailosaurError{}, err)
}

func TestSearchBySentTo(t *testing.T) {
	targetEmail := emails[1]

	result, _ := client.Messages.Search(&MessageSearchParams{
		Server: server,
	}, &SearchCriteria{
		SentTo: targetEmail.To[0].Email,
	})

	assert.Equal(t, 1, len(result.Items))
	assert.Equal(t, targetEmail.To[0].Email, result.Items[0].To[0].Email)
	assert.Equal(t, targetEmail.Subject, result.Items[0].Subject)
}

func TestSearchBySentToInvalidEmail(t *testing.T) {
	_, err := client.Messages.Search(&MessageSearchParams{
		Server: server,
	}, &SearchCriteria{
		SentTo: ".not_an_email_address",
	})

	assert.Error(t, err)
	assert.IsType(t, &mailosaurError{}, err)
}

func TestSearchByBody(t *testing.T) {
	targetEmail := emails[1]
	uniqueString := targetEmail.Subject[:8]

	result, _ := client.Messages.Search(&MessageSearchParams{
		Server: server,
	}, &SearchCriteria{
		Body: uniqueString + " html",
	})

	assert.Equal(t, 1, len(result.Items))
	assert.Equal(t, targetEmail.To[0].Email, result.Items[0].To[0].Email)
	assert.Equal(t, targetEmail.Subject, result.Items[0].Subject)
}

func TestSearchBySubject(t *testing.T) {
	targetEmail := emails[1]
	uniqueString := targetEmail.Subject[:8]

	result, _ := client.Messages.Search(&MessageSearchParams{
		Server: server,
	}, &SearchCriteria{
		Subject: uniqueString,
	})

	assert.Equal(t, 1, len(result.Items))
	assert.Equal(t, targetEmail.To[0].Email, result.Items[0].To[0].Email)
	assert.Equal(t, targetEmail.Subject, result.Items[0].Subject)
}

func TestSearchWithMatchAll(t *testing.T) {
	targetEmail := emails[1]
	uniqueString := targetEmail.Subject[:8]

	result, _ := client.Messages.Search(&MessageSearchParams{
		Server: server,
	}, &SearchCriteria{
		Subject: uniqueString,
		Body:    "this is a link",
		Match:   "ALL",
	})

	assert.Equal(t, 1, len(result.Items))
}

func TestSearchWithMatchAny(t *testing.T) {
	targetEmail := emails[1]
	uniqueString := targetEmail.Subject[:8]

	result, _ := client.Messages.Search(&MessageSearchParams{
		Server: server,
	}, &SearchCriteria{
		Subject: uniqueString,
		Body:    "this is a link",
		Match:   "ANY",
	})

	assert.Equal(t, 6, len(result.Items))
}

func TestSearchWithSpecialCharacters(t *testing.T) {
	result, _ := client.Messages.Search(&MessageSearchParams{
		Server: server,
	}, &SearchCriteria{
		Subject: "Search with ellipsis … and emoji 👨🏿‍🚒",
	})

	assert.Equal(t, 0, len(result.Items))
}

func TestSpamAnalysis(t *testing.T) {
	targetId := emails[0].Id
	result, _ := client.Analysis.Spam(targetId)

	for _, rule := range result.SpamFilterResults.SpamAssassin {
		assert.True(t, len(rule.Rule) != 0)
		assert.True(t, len(rule.Description) != 0)
	}
}

func TestDeleteMessage(t *testing.T) {
	targetEmailId := emails[4].Id

	err := client.Messages.Delete(targetEmailId)
	assert.NoError(t, err)

	// Attempting to delete again should fail
	err = client.Messages.Delete(targetEmailId)

	assert.Error(t, err)
	assert.IsType(t, &mailosaurError{}, err)
}

func TestCreateSendText(t *testing.T) {
	if verifiedDomain == "mailosaur.net" {
		t.Skip()
	}

	subject := "New message"

	var message, err = client.Messages.Create(server, &MessageCreateOptions{
		To:      fmt.Sprintf("anything@%s", verifiedDomain),
		Send:    true,
		Subject: subject,
		Text:    "This is a new email",
	})

	assert.NoError(t, err)
	assert.NotNil(t, message.Id)
	assert.Equal(t, subject, message.Subject)
}

func TestCreateSendHtml(t *testing.T) {
	if verifiedDomain == "mailosaur.net" {
		t.Skip()
	}

	subject := "New HTML message"

	var message, err = client.Messages.Create(server, &MessageCreateOptions{
		To:      fmt.Sprintf("anything@%s", verifiedDomain),
		Send:    true,
		Subject: subject,
		Html:    "This is a new email",
	})

	assert.NoError(t, err)
	assert.NotNil(t, message.Id)
	assert.Equal(t, subject, message.Subject)
}

func TestCreateSendWithAttachment(t *testing.T) {

	result, _ := client.Messages.List(&MessageListParams{Server: server})
	emails = result.Items

	if verifiedDomain == "mailosaur.net" {
		t.Skip()
	}

	subject := "New message with attachment"

	catImage, _ := ioutil.ReadFile("testing/cat.png")

	var attachment = Attachment{
		FileName:    "cat.png",
		Content:     base64.StdEncoding.EncodeToString([]byte(catImage)),
		ContentType: "image/png",
	}

	var message, _ = client.Messages.Create(server, &MessageCreateOptions{
		To:          fmt.Sprintf("anything@%s", verifiedDomain),
		Send:        true,
		Subject:     subject,
		Html:        "This is a new email",
		Attachments: []Attachment{attachment},
	})

	assert.Equal(t, 1, len(message.Attachments))
	file1 := message.Attachments[0]
	assert.True(t, len(file1.Id) != 0)
	assert.Equal(t, 82138, file1.Length)
	assert.NotNil(t, file1.Url)
	assert.Equal(t, "cat.png", file1.FileName)
	assert.Equal(t, "image/png", file1.ContentType)
}

func TestForwardText(t *testing.T) {
	if verifiedDomain == "mailosaur.net" {
		t.Skip()
	}

	body := "Forwarded message"
	targetEmailId := emails[0].Id

	var message, err = client.Messages.Forward(targetEmailId, &MessageForwardOptions{
		To:   fmt.Sprintf("anything@%s", verifiedDomain),
		Text: body,
	})

	assert.NoError(t, err)
	assert.NotNil(t, message.Id)
	assert.Contains(t, message.Text.Body, body)
}

func TestForwardHtml(t *testing.T) {
	if verifiedDomain == "mailosaur.net" {
		t.Skip()
	}

	body := "<p>Forwarded <strong>HTML</strong> message.</p>"
	targetEmailId := emails[0].Id

	var message, err = client.Messages.Forward(targetEmailId, &MessageForwardOptions{
		To:   fmt.Sprintf("anything@%s", verifiedDomain),
		Html: body,
	})

	assert.NoError(t, err)
	assert.NotNil(t, message.Id)
	assert.Contains(t, message.Html.Body, body)
}

func TestReplyText(t *testing.T) {
	if verifiedDomain == "mailosaur.net" {
		t.Skip()
	}

	body := "Reply message"
	targetEmailId := emails[0].Id

	var message, err = client.Messages.Reply(targetEmailId, &MessageReplyOptions{
		Text: body,
	})

	assert.NoError(t, err)
	assert.NotNil(t, message.Id)
	assert.Contains(t, message.Text.Body, body)
}

func TestReplyHtml(t *testing.T) {
	if verifiedDomain == "mailosaur.net" {
		t.Skip()
	}

	body := "<p>Reply <strong>HTML</strong> message.</p>"
	targetEmailId := emails[0].Id

	var message, err = client.Messages.Reply(targetEmailId, &MessageReplyOptions{
		Html: body,
	})

	assert.NoError(t, err)
	assert.NotNil(t, message.Id)
	assert.Contains(t, message.Html.Body, body)
}

func TestReplyWithAttachment(t *testing.T) {
	if verifiedDomain == "mailosaur.net" {
		t.Skip()
	}

	body := "<p>Reply with attachment.</p>"
	targetEmailId := emails[0].Id

	catImage, _ := ioutil.ReadFile("testing/cat.png")

	var attachment = Attachment{
		FileName:    "cat.png",
		Content:     base64.StdEncoding.EncodeToString([]byte(catImage)),
		ContentType: "image/png",
	}

	var message, _ = client.Messages.Reply(targetEmailId, &MessageReplyOptions{
		Html:        body,
		Attachments: []Attachment{attachment},
	})

	assert.Equal(t, 1, len(message.Attachments))
	file1 := message.Attachments[0]
	assert.True(t, len(file1.Id) != 0)
	assert.Equal(t, 82138, file1.Length)
	assert.NotNil(t, file1.Url)
	assert.Equal(t, "cat.png", file1.FileName)
	assert.Equal(t, "image/png", file1.ContentType)
}

func validateEmail(t *testing.T, email *Message) {
	validateMetadata(t, &MessageSummary{
		From:     email.From,
		To:       email.To,
		Cc:       email.Cc,
		Bcc:      email.Bcc,
		Subject:  email.Subject,
		Received: email.Received,
	})
	validateAttachmentMetadata(t, email)
	validateHtml(t, email)
	validateText(t, email)
	assert.NotNil(t, email.Metadata.Ehlo)
	assert.NotNil(t, email.Metadata.MailFrom)
	assert.Equal(t, 1, len(email.Metadata.RcptTo))
}

func validateEmailSummary(t *testing.T, email *MessageSummary) {
	validateMetadata(t, email)
	assert.True(t, len(email.Summary) != 0)
	assert.Equal(t, 2, email.Attachments)
}

func validateHtml(t *testing.T, email *Message) {
	// Html.Body
	assert.True(t, strings.HasPrefix(email.Html.Body, "<div dir=\"ltr\">"))

	// Html.Links
	assert.Equal(t, 3, len(email.Html.Links))
	assert.Equal(t, "https://mailosaur.com/", email.Html.Links[0].Href)
	assert.Equal(t, "mailosaur", email.Html.Links[0].Text)
	assert.Equal(t, "https://mailosaur.com/", email.Html.Links[1].Href)
	assert.Equal(t, "", email.Html.Links[1].Text)
	assert.Equal(t, "http://invalid/", email.Html.Links[2].Href)
	assert.Equal(t, "invalid", email.Html.Links[2].Text)

	// Html.Codes
	assert.Equal(t, 2, len(email.Html.Codes))
	assert.Equal(t, "123456", email.Html.Codes[0].Value)
	assert.Equal(t, "G3H1Y2", email.Html.Codes[1].Value)

	// Html.Images
	assert.True(t, strings.HasPrefix(email.Html.Images[1].Src, "cid:"))
	assert.Equal(t, "Inline image 1", email.Html.Images[1].Alt)
}

func validateText(t *testing.T, email *Message) {
	// Text.Body
	assert.True(t, strings.HasPrefix(email.Text.Body, "this is a test"))

	// Text.Links
	assert.Equal(t, 2, len(email.Text.Links))
	assert.Equal(t, "https://mailosaur.com/", email.Text.Links[0].Href)
	assert.Equal(t, email.Text.Links[0].Href, email.Text.Links[0].Text)
	assert.Equal(t, "https://mailosaur.com/", email.Text.Links[1].Href)
	assert.Equal(t, email.Text.Links[1].Href, email.Text.Links[1].Text)

	// Text.Codes
	assert.Equal(t, 2, len(email.Text.Codes))
	assert.Equal(t, "654321", email.Text.Codes[0].Value)
	assert.Equal(t, "5H0Y2", email.Text.Codes[1].Value)
}

func validateHeaders(_ *testing.T, _ *Message) {
	// Not implemented
}

func validateMetadata(t *testing.T, email *MessageSummary) {
	assert.Equal(t, 1, len(email.From))
	assert.Equal(t, 1, len(email.To))

	assert.True(t, len(email.From[0].Email) != 0)
	assert.True(t, len(email.From[0].Name) != 0)
	assert.True(t, len(email.To[0].Email) != 0)
	assert.True(t, len(email.To[0].Name) != 0)
	assert.True(t, len(email.Subject) != 0)

	assert.Equal(t, time.Now().Format("2006-01-02"), email.Received.Format("2006-01-02"))
}

func validateAttachmentMetadata(t *testing.T, email *Message) {
	assert.Equal(t, 2, len(email.Attachments))

	file1 := email.Attachments[0]
	assert.True(t, len(file1.Id) != 0)
	assert.Equal(t, 82138, file1.Length)
	assert.NotNil(t, file1.Url)
	assert.Equal(t, "cat.png", file1.FileName)
	assert.Equal(t, "image/png", file1.ContentType)

	file2 := email.Attachments[1]
	assert.True(t, len(file2.Id) != 0)
	assert.Equal(t, 212080, file2.Length)
	assert.NotNil(t, file2.Url)
	assert.Equal(t, "dog.png", file2.FileName)
	assert.Equal(t, "image/png", file2.ContentType)
}
