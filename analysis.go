package mailosaur

type AnalysisService struct {
	client *MailosaurClient
}

type SpamAssassinRule struct {
	Score       float64 `json:"score"`
	Rule        string  `json:"rule"`
	Description string  `json:"description"`
}

type SpamFilterResults struct {
	SpamAssassin []*SpamAssassinRule `json:"spamAssassin"`
}

type SpamAnalysisResult struct {
	SpamFilterResults *SpamFilterResults `json:"spamFilterResults"`
	Score             float64            `json:"score"`
}

type EmailAuthenticationResult struct {
	Result      string            `json:"result"`
	Description string            `json:"description"`
	RawValue    string            `json:"rawValue"`
	Tags        map[string]string `json:"tags"`
}

type BlockListResult struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Result string `json:"result"`
}

type Content struct {
	Embed                  bool `json:"embed"`
	Iframe                 bool `json:"iframe"`
	Object                 bool `json:"object"`
	Script                 bool `json:"script"`
	ShortUrls              bool `json:"shortUrls"`
	TextSize               int  `json:"textSize"`
	TotalSize              int  `json:"totalSize"`
	MissingAlt             bool `json:"missingAlt"`
	MissingListUnsubscribe bool `json:"missingListUnsubscribe"`
}

type DnsRecords struct {
	A   []string `json:"a"`
	Mx  []string `json:"mx"`
	Ptr []string `json:"ptr"`
}

type SpamAssassinResult struct {
	Score  int                 `json:"score"`
	Result string              `json:"result"`
	Rules  []*SpamAssassinRule `json:"rules"`
}

type DeliverabilityReport struct {
	Spf          *EmailAuthenticationResult   `json:"spf"`
	Dkim         []*EmailAuthenticationResult `json:"dkim"`
	Dmarc        *EmailAuthenticationResult   `json:"dmarc"`
	BlockLists   []*BlockListResult           `json:"blockLists"`
	Content      *Content                     `json:"content"`
	DnsRecords   *DnsRecords                  `json:"dnsRecords"`
	SpamAssassin *SpamAssassinResult          `json:"spamAssassin"`
}

func (s *AnalysisService) Spam(id string) (*SpamAnalysisResult, error) {
	result, err := s.client.HttpGet(&SpamAnalysisResult{}, "api/analysis/spam/"+id)
	return result.(*SpamAnalysisResult), err
}

func (s *AnalysisService) Deliverability(id string) (*DeliverabilityReport, error) {
	result, err := s.client.HttpGet(&DeliverabilityReport{}, "api/analysis/deliverability/"+id)
	return result.(*DeliverabilityReport), err
}
