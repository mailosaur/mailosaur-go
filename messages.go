package mailosaur

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// MessagesService provides operations for finding, retrieving, creating, forwarding,
// replying to, and deleting the email and SMS messages received by your Mailosaur servers.
// Accessed via the Messages field of MailosaurClient.
type MessagesService struct {
	client *MailosaurClient
}

type Attachment struct {
	Id          string `json:"id,omitempty"`
	ContentType string `json:"contentType"`
	FileName    string `json:"fileName"`
	Content     string `json:"content"`
	ContentId   string `json:"contentId,omitempty"`
	Length      int    `json:"length,omitempty"`
	Url         string `json:"url,omitempty"`
}

type Image struct {
	Src string `json:"src"`
	Alt string `json:"alt"`
}

type Code struct {
	Value string `json:"value"`
}

type Link struct {
	Href string `json:"href"`
	Text string `json:"text"`
}

type Message struct {
	Id          string            `json:"id"`
	Type        string            `json:"type"`
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
	Codes  []*Code  `json:"codes"`
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
	Dir           string
}

type MessageListResult struct {
	Items []*MessageSummary `json:"items"`
}

type MessageSummary struct {
	Id          string            `json:"id"`
	Type        string            `json:"type"`
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
	Cc          string       `json:"cc"`
	From        string       `json:"from"`
	Send        bool         `json:"send"`
	Subject     string       `json:"subject"`
	Text        string       `json:"text"`
	Html        string       `json:"html"`
	Attachments []Attachment `json:"attachments"`
}

type MessageForwardOptions struct {
	To   string `json:"to"`
	Cc   string `json:"cc"`
	Text string `json:"text"`
	Html string `json:"html"`
}

type MessageReplyOptions struct {
	Cc          string       `json:"cc"`
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
	Dir            string
}

type Preview struct {
	Id            string `json:"id"`
	EmailClient   string `json:"emailClient"`
	DisableImages bool   `json:"disableImages"`
}

type PreviewListResult struct {
	Items []*Preview `json:"items"`
}

type PreviewRequestOptions struct {
	EmailClients []string `json:"emailClients"`
}

// List returns a list of your messages in summary form, sorted by received date with the
// most recently-received messages appearing first. The params argument specifies the server
// to list messages from along with paging and filtering options. It returns a
// MessageListResult containing the message summaries.
func (s *MessagesService) List(params *MessageListParams) (*MessageListResult, error) {
	u := buildPagePath(
		"api/messages?server="+params.Server,
		params.Page,
		params.ItemsPerPage,
		params.ReceivedAfter,
		params.Dir,
	)

	result, err := s.client.HttpGet(&MessageListResult{}, u)
	return result.(*MessageListResult), err
}

// Get waits for a message to be found, returning as soon as a message matching the given
// search criteria is found. This is the most efficient way to look up a message and is
// recommended wherever possible. The params argument specifies the server to search and
// related options, and criteria specifies what to match. It returns the first matching
// Message. It returns a mailosaurError with error type no_messages_found if no matching
// message exists, or search_timeout if no matching message arrives before the timeout elapses.
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
		return nil, err
	}

	return s.GetById(result.Items[0].Id)
}

// Search returns a list of messages matching the given search criteria, in summary form,
// sorted by received date with the most recently-received messages appearing first. The
// params argument specifies the server to search along with paging and timeout options, and
// criteria specifies what to match. It returns a MessageListResult containing the matching
// message summaries. It returns a mailosaurError with error type search_timeout if no
// matching message is found before the timeout elapses, unless ErrorOnTimeout is set to false.
func (s *MessagesService) Search(params *MessageSearchParams, criteria *SearchCriteria) (*MessageListResult, error) {
	pollCount := 0
	startTime := time.Now()

	u := buildPagePath(
		"api/messages/search?server="+params.Server,
		params.Page,
		params.ItemsPerPage,
		params.ReceivedAfter,
		params.Dir,
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
		}

		var delay int
		if pollCount >= len(delayPatternValues) {
			delay = delayPatternValues[len(delayPatternValues)-1] / 1000
		} else {
			delay = delayPatternValues[pollCount] / 1000
		}

		pollCount++

		// Stop if timeout will be exceeded
		if time.Since(startTime).Seconds()+float64(delay) > float64(params.Timeout) {
			if *params.ErrorOnTimeout == false {
				return result.(*MessageListResult), nil
			}

			criteriaJson, _ := json.Marshal(criteria)
			err := &mailosaurError{
				Message:   "No matching messages found in time. By default, only messages received in the last hour are checked (use receivedAfter to override this). The search criteria used for this query was [" + string(criteriaJson) + "] which timed out after " + fmt.Sprint(params.Timeout) + "s",
				ErrorType: "search_timeout",
			}
			return nil, err
		}

		time.Sleep(time.Duration(delay) * time.Second)
	}
}

// GetById retrieves the detail for a single message. It must be used in conjunction with
// either List or Search in order to obtain the unique identifier for the required message.
// The id argument is the unique identifier of the message to retrieve, and it returns the
// full Message.
func (s *MessagesService) GetById(id string) (*Message, error) {
	result, err := s.client.HttpGet(&Message{}, "api/messages/"+id)
	return result.(*Message), err
}

// Delete permanently deletes a message, along with any attachments related to it. This
// operation cannot be undone. The id argument is the identifier of the message to delete.
func (s *MessagesService) Delete(id string) error {
	return s.client.HttpDelete("api/messages/" + id)
}

// DeleteAll permanently deletes all messages within a server. This operation cannot be
// undone. The server argument is the unique identifier of the server to clear.
func (s *MessagesService) DeleteAll(server string) error {
	return s.client.HttpDelete("api/messages?server=" + server)
}

// Create creates a new message that can be sent to a verified email address. This is useful
// when you want an email to trigger a workflow in your product. The server argument is the
// unique identifier of the server, and messageCreateOptions specifies the message to create.
// It returns the newly-created Message.
func (s *MessagesService) Create(server string, messageCreateOptions *MessageCreateOptions) (*Message, error) {
	result, err := s.client.HttpPost(&Message{}, "api/messages?server="+server, messageCreateOptions)
	return result.(*Message), err
}

// Forward forwards the specified message to a verified email address. This is useful for
// simulating a user forwarding one of your email messages. The id argument is the unique
// identifier of the message to forward, and messageForwardOptions specifies the forwarding
// options. It returns the forwarded Message.
func (s *MessagesService) Forward(id string, messageForwardOptions *MessageForwardOptions) (*Message, error) {
	result, err := s.client.HttpPost(&Message{}, "api/messages/"+id+"/forward", messageForwardOptions)
	return result.(*Message), err
}

// Reply sends a reply to the specified message. This is useful for simulating a user
// replying to one of your email or SMS messages. The id argument is the unique identifier of
// the message to reply to, and messageReplyOptions specifies the reply options. It returns
// the reply Message.
func (s *MessagesService) Reply(id string, messageReplyOptions *MessageReplyOptions) (*Message, error) {
	result, err := s.client.HttpPost(&Message{}, "api/messages/"+id+"/reply", messageReplyOptions)
	return result.(*Message), err
}

// GeneratePreviews generates screenshots of an email rendered in the specified email
// clients. The id argument is the identifier of the email to preview, and options specifies
// which email clients to use. It returns a PreviewListResult containing the generated previews.
func (s *MessagesService) GeneratePreviews(id string, options *PreviewRequestOptions) (*PreviewListResult, error) {
	result, err := s.client.HttpPost(&PreviewListResult{}, "api/messages/"+id+"/screenshots", options)
	return result.(*PreviewListResult), err
}
