package mailosaur

import (
	"strconv"
	"strings"
	"time"
)

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
	timeout := 120
	pollCount := 0
	startTime := time.Now()

	for {
		result, delayHeader, err := s.client.executeRequestWithDelayHeader(nil, "GET", "api/files/screenshots/"+id, nil, 200)

		if err == nil {
			return result.([]byte), nil
		}

		// Check if it's a mailosaur error and if the status code is 202 (still processing)
		if mailosaurErr, ok := err.(*mailosaurError); ok {
			if mailosaurErr.HttpStatusCode == 202 {
				// Continue polling
			} else if mailosaurErr.HttpStatusCode == 410 {
				return nil, &mailosaurError{
					Message:   "Permanently expired or deleted.",
					ErrorType: "gone",
					HttpStatusCode: 410,
					HttpResponseBody: mailosaurErr.HttpResponseBody,
				}
			} else {
				// Other errors should be returned immediately
				return nil, err
			}
		} else {
			return nil, err
		}

		delayPattern := "1000"
		if len(delayHeader) != 0 {
			delayPattern = delayHeader
		}
		delayPatternSplit := strings.Split(delayPattern, ",")

		var delayPatternValues []int

		for _, v := range delayPatternSplit {
			var n int
			n, _ = strconv.Atoi(strings.TrimSpace(v))
			delayPatternValues = append(delayPatternValues, n)
		}

		var delay int
		if pollCount >= len(delayPatternValues) {
			delay = delayPatternValues[len(delayPatternValues)-1] / 1000
		} else {
			delay = delayPatternValues[pollCount] / 1000
		}

		pollCount++

		// Stop if timeout will be exceeded
		if time.Since(startTime).Seconds()+float64(delay) > float64(timeout) {
			err := &mailosaurError{
				Message:   "An email preview was not generated in time. The email client may not be available, or the preview ID [" + id + "] may be incorrect.",
				ErrorType: "preview_timeout",
			}
			return nil, err
		}

		time.Sleep(time.Duration(delay) * time.Second)
	}
}
