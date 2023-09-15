package mailosaur

type PreviewsService struct {
	client *MailosaurClient
}

type PreviewPreviewEmailClient struct {
	Id               string `json:"id"`
	Name             string `json:"name"`
	PlatformGroup    string `json:"platformGroup"`
	PlatformType     string `json:"platformType"`
	PlatformVersion  string `json:"platformVersion"`
	CanDisableImages bool   `json:"canDisableImages"`
	Status           string `json:"status"`
}

type PreviewEmailClientListResult struct {
	Items []*PreviewPreviewEmailClient `json:"items"`
}

func (s *PreviewsService) ListEmailClients() (*PreviewEmailClientListResult, error) {
	result, err := s.client.HttpGet(&PreviewEmailClientListResult{}, "api/previews/clients")
	return result.(*PreviewEmailClientListResult), err
}
