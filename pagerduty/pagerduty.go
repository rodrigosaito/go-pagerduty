package pagerduty

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const BaseEndpoint = "https://events.pagerduty.com/generic/2010-04-15"

type TriggerJSON struct {
	ServiceKey  string                 `json:"service_key"`
	EventType   string                 `json:"event_type"`
	Description string                 `json:"description"`
	IncidentKey string                 `json:"incident_key,omitempty"`
	Client      string                 `json:"client,omitempty"`
	ClientURL   string                 `json:"client_url,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Contexts    []ContextJSON          `json:"contexts,omitempty"`
}

type ContextJSON struct {
	Type string `json:"type"`
	HREF string `json:"href,omitempty"`
	SRC  string `json:"src,omitempty"`
	ALT  string `json:"alt,omitempty"`
	Text string `json:"text,omitempty"`
}

type Trigger struct {
	Description string
	IncidentKey string
	Client      string
	ClientURL   string
	Details     map[string]interface{}
	Contexts    []Context
}

type Context struct {
	Type string
	HREF string
	SRC  string
	ALT  string
	Text string
}

type TriggerResponse struct {
	IncidentKey string `json:"incident_key"`
}

type PagerDuty struct {
	BaseURL    *url.URL
	ServiceKey string
	client     *http.Client
}

func NewPagerDuty(serviceKey string) *PagerDuty {
	baseURL, _ := url.Parse(BaseEndpoint)
	return &PagerDuty{baseURL, serviceKey, &http.Client{}}
}

func (pd *PagerDuty) Trigger(trigger Trigger) (*TriggerResponse, error) {
	contexts := []ContextJSON{}
	for _, c := range trigger.Contexts {
		contexts = append(contexts, ContextJSON{
			Type: c.Type,
			HREF: c.HREF,
			SRC:  c.SRC,
			ALT:  c.ALT,
			Text: c.Text,
		})
	}

	triggerJSON := TriggerJSON{
		ServiceKey:  pd.ServiceKey,
		EventType:   "trigger",
		Description: trigger.Description,
		IncidentKey: trigger.IncidentKey,
		Client:      trigger.Client,
		ClientURL:   trigger.ClientURL,
		Details:     trigger.Details,
		Contexts:    contexts,
	}

	j, err := json.Marshal(triggerJSON)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", pd.BaseURL.String()+"/create_event.json", bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := pd.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	tResp := TriggerResponse{}
	if err := json.Unmarshal(body, &tResp); err != nil {
		return nil, err
	}

	return &tResp, nil
}
