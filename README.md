# [Mailosaur - Go library](https://mailosaur.com/) &middot; [![](https://github.com/mailosaur/mailosaur-go/workflows/CI/badge.svg)](https://github.com/mailosaur/mailosaur-go/actions)

Mailosaur lets you automate email and SMS tests as part of software development and QA.

- **Unlimited test email addresses for all**  - every account gives users an unlimited number of test email addresses to test with.
- **End-to-end (e2e) email and SMS testing** Allowing you to set up end-to-end tests for password reset emails, account verification processes and MFA/one-time passcodes sent via text message.
- **Fake SMTP servers** Mailosaur also provides dummy SMTP servers to test with; allowing you to catch email in staging environments - preventing email being sent to customers by mistake.

## Get Started

This guide provides several key sections:

- [Mailosaur - Go library · ](#mailosaur---go-library--)
  - [Get Started](#get-started)
    - [Installation](#installation)
    - [Set your API key](#set-your-api-key)
    - [Create your code](#create-your-code)
    - [API Reference](#api-reference)
  - [Creating an account](#creating-an-account)
  - [Test email addresses with Mailosaur](#test-email-addresses-with-mailosaur)
  - [Find an email](#find-an-email)
    - [What is this code doing?](#what-is-this-code-doing)
    - [My email wasn't found](#my-email-wasnt-found)
  - [Find an SMS message](#find-an-sms-message)
  - [Testing plain text content](#testing-plain-text-content)
    - [Extracting verification codes from plain text](#extracting-verification-codes-from-plain-text)
  - [Testing HTML content](#testing-html-content)
    - [Working with HTML using goquery](#working-with-html-using-goquery)
  - [Working with hyperlinks](#working-with-hyperlinks)
    - [Links in plain text (including SMS messages)](#links-in-plain-text-including-sms-messages)
  - [Working with attachments](#working-with-attachments)
    - [Writing an attachment to disk](#writing-an-attachment-to-disk)
  - [Working with images and web beacons](#working-with-images-and-web-beacons)
    - [Remotely-hosted images](#remotely-hosted-images)
    - [Triggering web beacons](#triggering-web-beacons)
  - [Spam checking](#spam-checking)
  - [Development](#development)
  - [Contacting us](#contacting-us)

You can find the full [Mailosaur documentation](https://mailosaur.com/docs/) on the website.

If you get stuck, just contact us at support@mailosaur.com.

### Installation

Install the Mailosaur Go library via `go get`:

```sh
go get -u github.com/mailosaur/mailosaur-go
```

### Set your API key

Get your API key from the Mailosaur Dashboard and set it as an environment variable:

```sh
export MAILOSAUR_API_KEY='your-api-key-here'
```

### Create your code

Now import the library and create a client:

```go
import (
    "github.com/mailosaur/mailosaur-go"
)

m := mailosaur.New()
```

### API Reference

This library is powered by the Mailosaur [email & SMS testing API](https://mailosaur.com/docs/api/). You can easily check out the API itself by looking at our [API reference documentation](https://mailosaur.com/docs/api/) or via our Postman or Insomnia collections:

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/6961255-6cc72dff-f576-451a-9023-b82dec84f95d?action=collection%2Ffork&collection-url=entityId%3D6961255-6cc72dff-f576-451a-9023-b82dec84f95d%26entityType%3Dcollection%26workspaceId%3D386a4af1-4293-4197-8f40-0eb49f831325)
 [![Run in Insomnia](https://insomnia.rest/images/run.svg)](https://insomnia.rest/run/?label=Mailosaur&uri=https%3A%2F%2Fmailosaur.com%2Finsomnia.json)

## Creating an account

Create a [free trial account](https://mailosaur.com/app/signup) for Mailosaur via the website.

To start testing you'll need three things:

- **API key** - [manage API keys within the Mailosaur Dashboard](https://mailosaur.com/app/keys). You can scope a key to a single inbox (server) — recommended — or create an account-level key. [Learn more about API keys](https://mailosaur.com/docs/managing-your-account/api-keys/).
- **Inbox (server) domain** - open [your inbox](https://mailosaur.com/app/servers/default) (server) within the Mailosaur Dashboard to see its domain name (e.g. `abc123.mailosaur.net`). You'll need this to send email to the inbox (server).
- **Inbox (server) ID** - the first part of the inbox (server) domain. For `abc123.mailosaur.net` the ID is `abc123`. You need this whenever you interact with the inbox (server) via the API.

## Test email addresses with Mailosaur

Mailosaur gives you an **unlimited number of test email addresses** - with no setup or coding required!

Here's how it works:

* When you create an account, you are given an inbox (server).
* Every inbox (server) has its own domain name (e.g. `abc123.mailosaur.net`)
* Any email address that ends with `@{YOUR_SERVER_DOMAIN}` will work with Mailosaur without any special setup. For example:
  * `build-423@abc123.mailosaur.net`
  * `john.smith@abc123.mailosaur.net`
  * `rAnDoM63423@abc123.mailosaur.net`
* You can create more inboxes (servers) when you need them. Each one will have its own domain name.

***Can't use test email addresses?** You can also [use SMTP to test email](https://mailosaur.com/docs/email-testing/sending-to-mailosaur/#sending-via-smtp). By connecting your product or website to Mailosaur via SMTP, Mailosaur will catch all email your application sends, regardless of the email address.*

## Find an email

In automated tests you will want to wait for a new email to arrive. This library makes that easy with the `Messages.Get` method. Here's how you use it:

```go
package emailtests

import (
    "fmt"
    "testing"

    "github.com/mailosaur/mailosaur-go"
)

func TestExample(t *testing.T) {
    m := mailosaur.New()

    serverId := "abc123"
    serverDomain := "abc123.mailosaur.net"

    params := &mailosaur.MessageSearchParams{
        Server: serverId,
    }

    criteria := &mailosaur.SearchCriteria{
        SentTo: "anything@" + serverDomain,
    }

    email, err := m.Messages.Get(params, criteria)
    if err != nil {
        t.Error(err)
    }

    fmt.Println(email.Subject) // "Hello world!"
}
```

### What is this code doing?

1. Sets up an instance of `MailosaurClient`, reading the API key from the `MAILOSAUR_API_KEY` environment variable.
2. Waits for an email to arrive at the inbox (server) with ID `abc123`.
3. Outputs the subject line of the email.

### My email wasn't found

First, check that the email you sent is visible in the [Mailosaur Dashboard](https://mailosaur.com/app/project/messages).

If it is, the likely reason is that by default, `Messages.Get` only searches emails received by Mailosaur in the last 1 hour. You can override this behavior (see the `ReceivedAfter` field below), however we only recommend doing this during setup, as your tests will generally run faster with the default settings:

```go
params := &mailosaur.MessageSearchParams{
    Server: serverId,
    // Override ReceivedAfter to search all messages since Jan 1st
    ReceivedAfter: time.Date(2021, time.January, 1, 0, 0, 0, 0, time.UTC),
}

email, err := m.Messages.Get(params, criteria)
```

## Find an SMS message

**Important:** Trial accounts do not automatically have SMS access. Please contact our support team to enable a trial of SMS functionality.

If your account has [SMS testing](https://mailosaur.com/sms-testing/) enabled, you can reserve phone numbers to test with, then use the Mailosaur API in a very similar way to when testing email:

```go
m := mailosaur.New()

serverId := "abc123"

params := &mailosaur.MessageSearchParams{
    Server: serverId,
}

criteria := &mailosaur.SearchCriteria{
    SentTo: "4471235554444",
}

sms, err := m.Messages.Get(params, criteria)
if err != nil {
    t.Error(err)
}

fmt.Println(sms.Text.Body)
```

## Testing plain text content

Most emails, and all SMS messages, should have a plain text body. Mailosaur exposes this content via the `Text.Body` property on an email or SMS message:

```go
fmt.Println(message.Text.Body) // "Hi Jason, ..."

if strings.Contains(message.Text.Body, "Jason") {
    fmt.Println(`Email contains "Jason"`)
}
```

### Extracting verification codes from plain text

You may have an email or SMS message that contains an account verification code, or some other one-time passcode. Mailosaur automatically extracts these for you and makes them available via the `Codes` array on the message content:

```go
fmt.Println(message.Text.Body) // "Your access code is 243546."

fmt.Println(message.Text.Codes[0].Value) // "243546"
```

[Read more](https://mailosaur.com/docs/automation/codes)

## Testing HTML content

Most emails also have an HTML body, as well as the plain text content. You can access HTML content in a very similar way to plain text:

```go
fmt.Println(message.Html.Body) // "<html><head ..."
```

### Working with HTML using goquery

If you need to traverse the HTML content of an email — for example, finding an element via a CSS selector — you can use the [goquery](https://github.com/PuerkitoBio/goquery) library.

```sh
go get github.com/PuerkitoBio/goquery
```

```go
import (
    "strings"

    "github.com/PuerkitoBio/goquery"
)

// ...

doc, err := goquery.NewDocumentFromReader(strings.NewReader(message.Html.Body))
if err != nil {
    t.Error(err)
}

verificationCode := doc.Find(".verification-code").Text() // "542163"
```

[Read more](https://mailosaur.com/docs/test-cases/html-content/)

## Working with hyperlinks

When an email is sent with an HTML body, Mailosaur automatically extracts any hyperlinks found within anchor (`<a>`) and area (`<area>`) elements and makes these viable via the `Html.Links` array.

Each link has a `Text` property, representing the display text of the hyperlink within the body, and an `Href` property containing the target URL:

```go
// How many links?
fmt.Println(len(message.Html.Links)) // 2

firstLink := message.Html.Links[0]
fmt.Println(firstLink.Text) // "Google Search"
fmt.Println(firstLink.Href) // "https://www.google.com/"
```

**Important:** To ensure you always have valid emails, Mailosaur only extracts links that have been correctly marked up with `<a>` or `<area>` tags.

### Links in plain text (including SMS messages)

Mailosaur auto-detects links in plain text content too, which is especially useful for SMS testing:

```go
// How many links?
fmt.Println(len(message.Text.Links)) // 2

firstLink := message.Text.Links[0]
fmt.Println(firstLink.Href) // "https://www.google.com/"
```

## Working with attachments

If your email includes attachments, you can access these via the `Attachments` property:

```go
// How many attachments?
fmt.Println(len(message.Attachments)) // 2
```

Each attachment contains metadata on the file name and content type:

```go
firstAttachment := message.Attachments[0]
fmt.Println(firstAttachment.FileName)    // "contract.pdf"
fmt.Println(firstAttachment.ContentType) // "application/pdf"
```

The `Length` property returns the size of the attached file (in bytes):

```go
firstAttachment := message.Attachments[0]
fmt.Println(firstAttachment.Length) // 4028
```

### Writing an attachment to disk

```go
import "os"

// ...

firstAttachment := message.Attachments[1]

fileBytes, err := m.Files.GetAttachment(firstAttachment.Id)
if err != nil {
    t.Error(err)
}

err = os.WriteFile(firstAttachment.FileName, fileBytes, 0644)
if err != nil {
    t.Error(err)
}
```

## Working with images and web beacons

The `Html.Images` property of a message contains an array of images found within the HTML content of an email. The length of this array corresponds to the number of images found within an email:

```go
// How many images in the email?
fmt.Println(len(message.Html.Images)) // 1
```

### Remotely-hosted images

Emails will often contain many images that are hosted elsewhere, such as on your website or product. It is recommended to check that these images are accessible by your recipients.

All images should have an alternative text description, which can be checked using the `Alt` attribute.

```go
image := message.Html.Images[0]
fmt.Println(image.Alt) // "Hot air balloon"
```

### Triggering web beacons

A web beacon is a small image that can be used to track whether an email has been opened by a recipient.

Because a web beacon is simply another form of remotely-hosted image, you can use the `Src` attribute to perform an HTTP request to that address:

```go
import "net/http"

// ...

image := message.Html.Images[0]
fmt.Println(image.Src) // "https://example.com/s.png?abc123"

// Make an HTTP call to trigger the web beacon
resp, err := http.Get(image.Src)
if err != nil {
    t.Error(err)
}
defer resp.Body.Close()

fmt.Println(resp.StatusCode) // 200
```

## Spam checking

You can perform a [SpamAssassin](https://spamassassin.apache.org/) check against an email. The structure returned matches the [spam test object](https://mailosaur.com/docs/api/#spam):

```go
result, err := m.Analysis.Spam(message.Id)
if err != nil {
    t.Error(err)
}

fmt.Println(result.Score) // 0.5

for _, r := range result.SpamFilterResults.SpamAssassin {
    fmt.Println(r.Rule)
    fmt.Println(r.Description)
    fmt.Println(r.Score)
}
```

## Development

The test suite requires the following environment variables to be set:

```sh
export MAILOSAUR_API_KEY=your_api_key
export MAILOSAUR_SERVER=server_id
```

Run all tests:

```sh
go test -v
```

## Contacting us

You can get us at [support@mailosaur.com](mailto:support@mailosaur.com)
