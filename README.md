# [Mailosaur - Go library](https://mailosaur.com/) &middot; [![](https://github.com/mailosaur/mailosaur-go/workflows/CI/badge.svg)](https://github.com/mailosaur/mailosaur-go/actions)

Mailosaur lets you automate email and SMS tests as part of software development and QA.

- **Unlimited test email addresses for all**  - every account gives users an unlimited number of test email addresses to test with.
- **End-to-end (e2e) email and SMS testing** Allowing you to set up end-to-end tests for password reset emails, account verification processes and MFA/one-time passcodes sent via text message.
- **Fake SMTP servers** Mailosaur also provides dummy SMTP servers to test with; allowing you to catch email in staging environments - preventing email being sent to customers by mistake.

## Get Started

You can find the full [Mailosaur documentation](https://mailosaur.com/docs/) on the website.

If you get stuck, just contact us at support@mailosaur.com.

### Installation

```sh
go mod init
```

Then, reference mailosaur-go via `import`:

```golang
import (
    "github.com/mailosaur/mailosaur-go"
)
```

Alternatively, you can also `go get` the package into your project:

```sh
go get -u github.com/mailosaur/mailosaur-go
```

### API Reference

This library is powered by the Mailosaur [email & SMS testing API](https://mailosaur.com/docs/api/). You can easily check out the API itself by looking at our [API reference documentation](https://mailosaur.com/docs/api/) or via our Postman or Insomnia collections:

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/6961255-6cc72dff-f576-451a-9023-b82dec84f95d?action=collection%2Ffork&collection-url=entityId%3D6961255-6cc72dff-f576-451a-9023-b82dec84f95d%26entityType%3Dcollection%26workspaceId%3D386a4af1-4293-4197-8f40-0eb49f831325)
 [![Run in Insomnia}](https://insomnia.rest/images/run.svg)](https://insomnia.rest/run/?label=Mailosaur&uri=https%3A%2F%2Fmailosaur.com%2Finsomnia.json)

## Creating an account

Create a [free trial account](https://mailosaur.com/app/signup) for Mailosaur via the website.

Once you have this, navigate to the [API tab](https://mailosaur.com/app/project/api) to find the following values:

- **Server ID** - Servers act like projects, which group your tests together. You need this ID whenever you interact with a server via the API.
- **Server Domain** - Every server has its own domain name. You'll need this to send email to your server.
- **API Key** - You can create an API key per server (recommended), or an account-level API key to use across your whole account. [Learn more about API keys](https://mailosaur.com/docs/managing-your-account/api-keys/).

## Test email addresses with Mailosaur

Mailosaur gives you an **unlimited number of test email addresses** - with no setup or coding required!

Here's how it works:

* When you create an account, you are given a server.
* Every server has its own **Server Domain** name (e.g. `abc123.mailosaur.net`)
* Any email address that ends with `@{YOUR_SERVER_DOMAIN}` will work with Mailosaur without any special setup. For example:
  * `build-423@abc123.mailosaur.net`
  * `john.smith@abc123.mailosaur.net`
  * `rAnDoM63423@abc123.mailosaur.net`
* You can create more servers when you need them. Each one will have its own domain name.

***Can't use test email addresses?** You can also [use SMTP to test email](https://mailosaur.com/docs/email-testing/sending-to-mailosaur/#sending-via-smtp). By connecting your product or website to Mailosaur via SMTP, Mailosaur will catch all email your application sends, regardless of the email address.*

## Usage

example_test.go

```golang
package emailtests

import (
  "fmt"
  "testing"
  "github.com/mailosaur/mailosaur-go"
)

func TestExample(t *testing.T) {
  // Available in the API tab of a server
  apiKey := "YOUR_API_KEY";
  serverId := "SERVER_ID";
  serverDomain := "SERVER_DOMAIN";

  m := mailosaur.New(apiKey);

  params := &mailosaur.MessageSearchParams {
    Server: serverId,
  }

  criteria := &mailosaur.SearchCriteria {
    SentTo: "anything@" + serverDomain,
  }

  email, err := m.Messages.Get(params, criteria)

  if (err != nil) {
    t.Error(err)
  }

  // If we have an email, print the subject
  fmt.Println("Subject: " + email.Subject)
}
```

## Development

The test suite requires the following environment variables to be set:

```sh
export MAILOSAUR_API_KEY=your_api_key
export MAILOSAUR_SERVER=server_id
export MAILOSAUR_VERIFIED_DOMAIN=server_domain
```

Run all tests:

```sh
go test -v
```

## Contacting us

You can get us at [support@mailosaur.com](mailto:support@mailosaur.com)
