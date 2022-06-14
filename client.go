package mailosaur

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type MailosaurClient struct {
	baseUrl    string
	apiKey     string
	userAgent  string
	httpClient *http.Client

	Servers  *ServersService
	Messages *MessagesService
	Analysis *AnalysisService
	Files    *FilesService
	Usage    *UsageService
	Devices  *DevicesService
}

type mailosaurError struct {
	Message          string
	ErrorType        string
	HttpStatusCode   int
	HttpResponseBody string
}

type ErrorDetail struct {
	Description string `json:"description"`
}

type Error struct {
	Field  string        `json:"field"`
	Detail []ErrorDetail `json:"detail"`
}

type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

func (e *mailosaurError) Error() string {
	return e.Message
}

func New(apiKey string) *MailosaurClient {
	return NewWithClient(apiKey, &http.Client{Timeout: time.Minute})
}

func NewWithClient(apiKey string, httpClient *http.Client) *MailosaurClient {
	c := &MailosaurClient{
		baseUrl:    "https://mailosaur.com/",
		apiKey:     apiKey,
		httpClient: httpClient,
		userAgent:  "mailosaur-go/1.0.0",
	}

	c.Servers = &ServersService{client: c}
	c.Messages = &MessagesService{client: c}
	c.Analysis = &AnalysisService{client: c}
	c.Files = &FilesService{client: c}
	c.Usage = &UsageService{client: c}
	c.Devices = &DevicesService{client: c}

	return c
}

func (c *MailosaurClient) httpRequest(method, path string, body interface{}) (*http.Request, error) {
	u := c.baseUrl + path

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	return req, nil
}

func (c *MailosaurClient) executeRequestWithDelayHeader(result interface{}, method string, path string, body interface{}, expectedStatus int) (interface{}, string, error) {
	req, err := c.httpRequest(method, path, body)

	if err != nil {
		return result, "", err
	}

	req.SetBasicAuth(c.apiKey, "")
	resp, err := c.httpClient.Do(req)

	if err != nil {
		return result, "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		err := &mailosaurError{}
		err.HttpStatusCode = resp.StatusCode

		var bodyBytes []byte
		if err.HttpStatusCode != 204 {
			bodyBytes, _ = ioutil.ReadAll(resp.Body)
			err.HttpResponseBody = string(bodyBytes)
		}

		message := ""
		switch resp.StatusCode {
		case 400:
			var jsonResult ErrorResponse
			json.Unmarshal(bodyBytes, &jsonResult)
			for _, e := range jsonResult.Errors {
				message += fmt.Sprintf("(%s) %s\r\n", e.Field, e.Detail[0].Description)
			}
			err.Message = message
			err.ErrorType = "invalid_request"
		case 401:
			err.Message = "Authentication failed, check your API key."
			err.ErrorType = "authentication_error"
		case 403:
			err.Message = "Insufficient permission to perform that task."
			err.ErrorType = "permission_error"
		case 404:
			err.Message = "Not found, check input parameters."
			err.ErrorType = "invalid_request"
		default:
			err.Message = "An API error occurred, see httpResponse for further information."
			err.ErrorType = "api_error"
			break
		}

		return result, "", err
	}

	// If no result type is being marshalled, just return the bytes
	if result == nil {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		return bodyBytes, resp.Header.Get("x-ms-delay"), err
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, resp.Header.Get("x-ms-delay"), err
}

func (c *MailosaurClient) executeRequest(result interface{}, method string, path string, body interface{}, expectedStatus int) (interface{}, error) {
	result, _, err := c.executeRequestWithDelayHeader(result, method, path, body, expectedStatus)
	return result, err
}

func (c *MailosaurClient) HttpPost(result interface{}, path string, body interface{}) (interface{}, error) {
	return c.executeRequest(result, "POST", path, body, 200)
}

func (c *MailosaurClient) HttpGet(result interface{}, path string) (interface{}, error) {
	return c.executeRequest(result, "GET", path, nil, 200)
}

func (c *MailosaurClient) HttpPut(result interface{}, path string, body interface{}) (interface{}, error) {
	return c.executeRequest(result, "PUT", path, body, 200)
}

func (c *MailosaurClient) HttpDelete(path string) error {
	_, err := c.executeRequest(nil, "DELETE", path, nil, 204)
	return err
}

func buildPagePath(path string, page int, itemsPerPage int, receivedAfter time.Time) string {
	if page > 0 {
		path += "&page=" + fmt.Sprint(page)
	}

	if itemsPerPage > 0 {
		path += "&itemsPerPage=" + fmt.Sprint(itemsPerPage)
	}

	if !receivedAfter.IsZero() {
		path += "&receivedAfter=" + url.QueryEscape(receivedAfter.Format(time.RFC3339))
	}

	return path
}
