package mailosaur

type AnalysisService struct {
    client *MailosaurClient
}

type SpamAssassinRule struct {
    Score               float64                 `json:"score"`
    Rule                string                  `json:"rule"`
    Description         string                  `json:"description"`
}

type SpamFilterResults struct {
    SpamAssassin        []*SpamAssassinRule     `json:"spamAssassin"`
}

type SpamAnalysisResult struct {
    SpamFilterResults   *SpamFilterResults      `json:"spamFilterResults"`
    Score               float64                 `json:"score"`
}

func (s *AnalysisService) Spam(id string) (*SpamAnalysisResult, error) { 
    result, err := s.client.HttpGet(&SpamAnalysisResult{}, "/api/analysis/spam/" + id)
    return result.(*SpamAnalysisResult), err
}
