package mailosaur

import (
	"log"
	"strconv"
	"strings"
	"time"
)

type MessagesService struct {
	client *MailosaurClient
}

type Attachment struct {
	Id          string `json:"id"`
	ContentType string `json:"contentType"`
	FileName    string `json:"fileName"`
	Content     string `json:"content"`
	ContentId   string `json:"contentId"`
	Length      int    `json:"length"`
	Url         string `json:"url"`
}

type Image struct {
	Src string `json:"src"`
	Alt string `json:"alt"`
}

type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type Message struct {
	Id          string            `json:"id"`
	From        []*MessageAddress `json:"from"`
	To          []*MessageAddress `json:"to"`
	Cc          []*MessageAddress `json:"cc"`
	Bcc         []*MessageAddress `json:"bcc"`
	Received    time.Time         `json:"received"`
	Subject     string            `json:"subject"`
	Html        *MessageContent   `json:"html"`
	Text        *MessageContent   `json:"text"`
	Attachments []*Attachment     `json:"attachments"`
	Metadata    *Metadata         `json:"metadata"`
	Server      string            `json:"server"`
}

type MessageAddress struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type MessageContent struct {
	Links  []*Link  `json:"links"`
	Images []*Image `json:"images"`
	Body   string   `json:"body"`
}

type MessageHeader struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

type MessageListParams struct {
	Server        string
	ReceivedAfter time.Time
	Page          int
	ItemsPerPage  int
}

type MessageListResult struct {
	Items []*MessageSummary `json:"items"`
}

type MessageSummary struct {
	Id          string            `json:"id"`
	Server      string            `json:"server"`
	From        []*MessageAddress `json:"from"`
	To          []*MessageAddress `json:"to"`
	Cc          []*MessageAddress `json:"cc"`
	Bcc         []*MessageAddress `json:"bcc"`
	Received    time.Time         `json:"received"`
	Subject     string            `json:"subject"`
	Summary     string            `json:"summary"`
	Attachments int               `json:"attachments"`
}

type MessageCreateOptions struct {
	To          string       `json:"to"`
	Send        bool         `json:"send"`
	Subject     string       `json:"subject"`
	Text        string       `json:"text"`
	Html        string       `json:"html"`
	Attachments []Attachment `json:"attachments"`
}

type MessageForwardOptions struct {
	To   string `json:"to"`
	Text string `json:"text"`
	Html string `json:"html"`
}

type MessageReplyOptions struct {
	Text        string       `json:"text"`
	Html        string       `json:"html"`
	Attachments []Attachment `json:"attachments"`
}

type Metadata struct {
	Headers  []*MessageHeader  `json:"headers"`
	MailFrom string            `json:"mailFrom"`
	RcptTo   []*MessageAddress `json:"rcptTo"`
	Ehlo     string            `json:"ehlo"`
}

type SearchCriteria struct {
	SentFrom string `json:"sentFrom"`
	SentTo   string `json:"sentTo"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Match    string `json:"match"`
}

type MessageSearchParams struct {
	Server         string
	ReceivedAfter  time.Time
	Page           int
	ItemsPerPage   int
	Timeout        int
	ErrorOnTimeout *bool
}

func (s *MessagesService) List(params *MessageListParams) (*MessageListResult, error) {
	u := buildPagePath(
		"api/messages?server="+params.Server,
		params.Page,
		params.ItemsPerPage,
		params.ReceivedAfter,
	)

	result, err := s.client.HttpGet(&MessageListResult{}, u)
	return result.(*MessageListResult), err
}

func (s *MessagesService) Get(params *MessageSearchParams, criteria *SearchCriteria) (*Message, error) {
	// Timeout defaulted to 10s, receivedAfter to 1h
	if params.ReceivedAfter.IsZero() {
		params.ReceivedAfter = time.Now().Add(-(1 * time.Hour))
	}

	if params.Timeout == 0 {
		params.Timeout = 10
	}

	params.Page = 0
	params.ItemsPerPage = 1

	result, err := s.Search(params, criteria)
	if err != nil {
		log.Fatal(err)
	}

	return s.GetById(result.Items[0].Id)
}

func (s *MessagesService) Search(params *MessageSearchParams, criteria *SearchCriteria) (*MessageListResult, error) {
	pollCount := 0
	startTime := time.Now()

	u := buildPagePath(
		"api/messages/search?server="+params.Server,
		params.Page,
		params.ItemsPerPage,
		params.ReceivedAfter,
	)

	// Default value for Match
	if len(criteria.Match) == 0 {
		criteria.Match = "ALL"
	}

	// Default value for ErrorOnTimeout
	if params.ErrorOnTimeout == nil {
		t := true
		params.ErrorOnTimeout = &t
	}

	for {
		result, delayHeader, err := s.client.executeRequestWithDelayHeader(&MessageListResult{}, "POST", u, criteria, 200)

		if err != nil {
			return nil, err
		}

		if params.Timeout == 0 || len(result.(*MessageListResult).Items) != 0 {
			return result.(*MessageListResult), nil
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
			if err != nil {
				log.Fatal(err)
			}
		}

		var delay int
		if pollCount >= len(delayPatternValues) {
			delay = delayPatternValues[len(delayPatternValues)-1] / 1000
		} else {
			delay = delayPatternValues[pollCount] / 1000
		}

		pollCount++

		// Stop if timeout will be exceeded
		if time.Now().Sub(startTime).Seconds()+float64(delay) > float64(params.Timeout) {
			if *params.ErrorOnTimeout == false {
				return result.(*MessageListResult), nil
			}

			err := &mailosaurError{
				Message:   "No matching messages found in time. By default, only messages received in the last hour are checked (use receivedAfter to override this).",
				ErrorType: "search_timeout",
			}
			return nil, err
		}

		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func (s *MessagesService) GetById(id string) (*Message, error) {
	result, err := s.client.HttpGet(&Message{}, "api/messages/"+id)
	return result.(*Message), err
}

func (s *MessagesService) Delete(id string) error {
	return s.client.HttpDelete("api/messages/" + id)
}

func (s *MessagesService) DeleteAll(server string) error {
	return s.client.HttpDelete("api/messages?server=" + server)
}

func (s *MessagesService) Create(server string, messageCreateOptions *MessageCreateOptions) (*Message, error) {
	result, err := s.client.HttpPost(&Message{}, "api/messages?server="+server, messageCreateOptions)
	return result.(*Message), err
}

func (s *MessagesService) Forward(id string, messageForwardOptions *MessageForwardOptions) (*Message, error) {
	result, err := s.client.HttpPost(&Message{}, "api/messages/"+id+"/forward", messageForwardOptions)
	return result.(*Message), err
}

func (s *MessagesService) Reply(id string, messageReplyOptions *MessageReplyOptions) (*Message, error) {
	result, err := s.client.HttpPost(&Message{}, "api/messages/"+id+"/reply", messageReplyOptions)
	return result.(*Message), err
}
