package mailosaur

// PreviewsService provides operations for discovering the email clients available for
// generating email previews (screenshots of an email rendered in real clients). Accessed via
// the Previews field of MailosaurClient.
type PreviewsService struct {
	client *MailosaurClient
}

type EmailClient struct {
	Label string `json:"label"`
	Name  string `json:"name"`
}

type EmailClientListResult struct {
	Items []*EmailClient `json:"items"`
}

// ListEmailClients lists all email clients that can be used to generate email previews. It
// returns an EmailClientListResult of available email clients.
func (s *PreviewsService) ListEmailClients() (*EmailClientListResult, error) {
	result, err := s.client.HttpGet(&EmailClientListResult{}, "api/screenshots/clients")
	return result.(*EmailClientListResult), err
}
