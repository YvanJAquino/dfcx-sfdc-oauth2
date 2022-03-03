package dfcx

import (
	"encoding/json"
	"net/http"
)

// Structs
type WebhookRequest struct {
	DetectIntentResponseID string          `json:"detectIntentResponseId,omitempty"`
	IntentInfo             IntentInfo      `json:"intentInfo,omitempty"`
	PageInfo               PageInfo        `json:"pageInfo,omitempty"`
	SessionInfo            SessionInfo     `json:"sessionInfo,omitempty"`
	FulfillmentInfo        FulfillmentInfo `json:"fulfillmentInfo,omitempty"`
	Messages               []Messages      `json:"messages,omitempty"`
	Text                   string          `json:"text,omitempty"`
	LanguageCode           string          `json:"languageCode,omitempty"`
}

func (wr *WebhookRequest) FromRequest(r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(wr)
	if err != nil {
		return err
	}
	return nil
}

func FromRequest(r *http.Request) (*WebhookRequest, error) {
	var wr WebhookRequest
	err := json.NewDecoder(r.Body).Decode(&wr)
	if err != nil {
		return nil, err
	}
	return &wr, nil
}

type IntentInfo struct {
	LastMatchedIntent string  `json:"lastMatchedIntent,omitempty"`
	DisplayName       string  `json:"displayName,omitempty"`
	Confidence        float64 `json:"confidence,omitempty"`
}

type PageInfo struct {
	CurrentPage string `json:"currentPage,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

type SessionInfo struct {
	Session    string            `json:"session,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

type FulfillmentInfo struct {
	Tag string `json:"tag,omitempty"`
}

type Messages struct {
	Text         Text   `json:"text,omitempty"`
	ResponseType string `json:"responseType,omitempty"`
	Source       string `json:"source,omitempty"`
}

type Text struct {
	Text                      []string `json:"text,omitempty"`
	RedactedText              []string `json:"redactedText,omitempty"`
	AllowPlaybackInterruption bool     `json:"allowPlaybackInterruption,omitempty"`
}
