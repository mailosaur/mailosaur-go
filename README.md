# Mailosaur Go Client Library

[Mailosaur](https://mailosaur.com) lets you automate email and SMS tests, like account verification and password resets, and integrate these into your CI/CD pipeline.

[![](https://github.com/mailosaur/mailosaur-go/workflows/CI/badge.svg)](https://github.com/mailosaur/mailosaur-go/actions)

## Installation

### Install Mailosaur

If you're using Go Modules (you have a `go.mod` file in your project root):

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

## Documentation

Please see the [Go client reference](https://mailosaur.com/docs/email-testing/go/) for the most up-to-date documentation.

## Usage

example.go

```golang
mailosaur := mailosaur.New("YOUR_API_KEY")

result, _ := mailosaur.Servers.List()

fmt.Println("Your have a server called: " + result.Items[0].Name)
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
