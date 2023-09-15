package mailosaur

type FilesService struct {
	client *MailosaurClient
}

func (s *FilesService) GetAttachment(id string) ([]byte, error) {
	result, err := s.client.HttpGet(nil, "api/files/attachments/"+id)
	return result.([]byte), err
}

func (s *FilesService) GetEmail(id string) ([]byte, error) {
	result, err := s.client.HttpGet(nil, "api/files/email/"+id)
	return result.([]byte), err
}

func (s *FilesService) GetPreview(id string) ([]byte, error) {
	result, err := s.client.HttpGet(nil, "api/files/previews/"+id)
	return result.([]byte), err
}
