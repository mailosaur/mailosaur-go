package mailosaur

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

func (s *PreviewsService) ListEmailClients() (*EmailClientListResult, error) {
	result, err := s.client.HttpGet(&EmailClientListResult{}, "api/screenshots/clients")
	return result.(*EmailClientListResult), err
}
