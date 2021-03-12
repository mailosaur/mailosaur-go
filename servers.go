package mailosaur

import (
    "os"
    "fmt"
    "math/rand"
)

type ServersService struct {
    client *MailosaurClient
}

type Server struct {
	Id          string      `json:"id"`
	Password    string      `json:"password"`
	Name        string      `json:"name"`
	Messages    int         `json:"messages"`
}

type ServerListResult struct {
	Items       []*Server   `json:"items"`
}

type ServerCreateOptions struct {
    Name        string      `json:"name"`
}

func (s *ServersService) List() (*ServerListResult, error) { 
    result, err := s.client.HttpGet(&ServerListResult{}, "api/servers")
    return result.(*ServerListResult), err
}

func (s *ServersService) Create(serverCreateOptions ServerCreateOptions) (*Server, error) {
    result, err := s.client.HttpPost(&Server{}, "api/servers", serverCreateOptions)
    return result.(*Server), err
}

func (s *ServersService) Get(id string) (*Server, error) {
    result, err := s.client.HttpGet(&Server{}, "api/servers/" + id)
    return result.(*Server), err
}

func (s *ServersService) Update(id string, server *Server) (*Server, error) {
    result, err := s.client.HttpPut(&Server{}, "api/servers/" + id, server)
    return result.(*Server), err
}

func (s *ServersService) Delete(id string) (error) {
    return s.client.HttpDelete("api/servers/" + id)
}

func (s *ServersService) GenerateEmailAddress(id string) (string) {
    host := os.Getenv("MAILOSAUR_SMTP_HOST")
    if (len(host) == 0) {
        host = "mailosaur.net"
    }

    return fmt.Sprintf("%s@%s.%s", getRandomString(), id, host)
}

func getRandomString() (string) {
    var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
 
    s := make([]rune, 8)
    for i := range s {
        s[i] = letters[rand.Intn(len(letters))]
    }
    return string(s)
}