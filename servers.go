package mailosaur

import (
	"fmt"
	"math/rand"
	"os"
)

// ServersService provides operations for creating and managing your Mailosaur servers — the
// virtual inboxes that group your tests together, each with its own domain and
// SMTP/POP3/IMAP credentials. Accessed via the Servers field of MailosaurClient.
type ServersService struct {
	client *MailosaurClient
}

type Server struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Messages int    `json:"messages"`
}

type ServerListResult struct {
	Items []*Server `json:"items"`
}

type ServerCreateOptions struct {
	Name string `json:"name"`
}

// List returns a list of your virtual servers, sorted in alphabetical order. It returns a
// ServerListResult containing your servers.
func (s *ServersService) List() (*ServerListResult, error) {
	result, err := s.client.HttpGet(&ServerListResult{}, "api/servers")
	return result.(*ServerListResult), err
}

// Create creates a new virtual server. The serverCreateOptions argument specifies the
// options used to create the server. It returns the newly-created Server.
func (s *ServersService) Create(serverCreateOptions ServerCreateOptions) (*Server, error) {
	result, err := s.client.HttpPost(&Server{}, "api/servers", serverCreateOptions)
	return result.(*Server), err
}

// Get retrieves the detail for a single server. The id argument is the unique identifier of
// the server, and it returns the Server.
func (s *ServersService) Get(id string) (*Server, error) {
	result, err := s.client.HttpGet(&Server{}, "api/servers/"+id)
	return result.(*Server), err
}

// GetPassword retrieves the password for a server, which can be used for SMTP, POP3, and
// IMAP connectivity. The id argument is the unique identifier of the server, and it returns
// the server's password.
func (s *ServersService) GetPassword(id string) (string, error) {
	type Result struct {
		Value string `json:"value"`
	}

	result, err := s.client.HttpGet(&Result{}, "api/servers/"+id+"/password")
	parsed := result.(*Result)

	return parsed.Value, err
}

// Update updates the attributes of a server. The id argument is the unique identifier of the
// server, and server is the updated server. It returns the updated Server.
func (s *ServersService) Update(id string, server *Server) (*Server, error) {
	result, err := s.client.HttpPut(&Server{}, "api/servers/"+id, server)
	return result.(*Server), err
}

// Delete permanently deletes a server, along with all messages and associated attachments
// within it. This operation cannot be undone. The id argument is the unique identifier of
// the server to delete.
func (s *ServersService) Delete(id string) error {
	return s.client.HttpDelete("api/servers/" + id)
}

// GenerateEmailAddress generates a random email address by prefixing a random string to the
// server's domain name. The id argument is the identifier of the server, and it returns a
// random email address ending in the server's domain.
func (s *ServersService) GenerateEmailAddress(id string) string {
	host := os.Getenv("MAILOSAUR_SMTP_HOST")
	if len(host) == 0 {
		host = "mailosaur.net"
	}

	return fmt.Sprintf("%s@%s.%s", getRandomString(), id, host)
}

func getRandomString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := make([]rune, 8)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
