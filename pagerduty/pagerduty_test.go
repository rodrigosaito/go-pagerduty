package pagerduty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	server *httptest.Server
	mux    *http.ServeMux
	pd     *PagerDuty
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	baseURL, _ := url.Parse(server.URL + "/generic/2010-04-15")
	pd = NewPagerDuty("e93facc04764012d7bfb002500d5d1a6")
	pd.BaseURL = baseURL
}

func teardown() {
	server.Close()
}

func assertMethod(t *testing.T, req *http.Request, expected string) {
	assert.Equal(t, expected, req.Method)
}

func assertContentType(t *testing.T, req *http.Request) {
	assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
}

func assertBody(t *testing.T, req *http.Request, expected map[string]interface{}) {
	body, _ := ioutil.ReadAll(req.Body)
	var actual map[string]interface{}
	json.Unmarshal(body, &actual)
	for key, value := range expected {
		assert.Equal(t, value, actual[key])
	}
}

func TestTrigger_WithDescription(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/generic/2010-04-15/create_event.json", func(w http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, "POST")
		assertContentType(t, req)
		assertBody(t, req, map[string]interface{}{
			"service_key": "e93facc04764012d7bfb002500d5d1a6",
			"event_type":  "trigger",
			"description": "Something bad has happened",
		})
		fmt.Fprint(w, `
{
	"status": "success",
	"message": "Event processed",
	"incident_key": "srv01/HTTP"
}
		`)
	})

	resp, err := pd.Trigger(Trigger{Description: "Something bad has happened"})
	if assert.Nil(t, err) {
		assert.Equal(t, "srv01/HTTP", resp.IncidentKey)
	}
}

func TestTrigger_Complete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/generic/2010-04-15/create_event.json", func(w http.ResponseWriter, req *http.Request) {
		assertMethod(t, req, "POST")
		assertContentType(t, req)
		assertBody(t, req, map[string]interface{}{
			"service_key":  "e93facc04764012d7bfb002500d5d1a6",
			"event_type":   "trigger",
			"description":  "Something bad has happened",
			"incident_key": "e93facc04764012d7bfb002500d5d1a6",
			"client":       "Sample Monitoring Service",
			"client_url":   "https://monitoring.service.com",
			"details": map[string]interface{}{
				"ping time": "1500ms",
				"load avg":  0.75,
			},
			"contexts": []interface{}{
				map[string]interface{}{
					"type": "link",
					"href": "http://acme.pagerduty.com",
				},
				map[string]interface{}{
					"type": "link",
					"href": "http://acme.pagerduty.com",
					"text": "View the incident on PagerDuty",
				},
				map[string]interface{}{
					"type": "image",
					"src":  "https://chart.googleapis.com/chart?chs=600x400&chd=t:6,2,9,5,2,5,7,4,8,2,1&cht=lc&chds=a&chxt=y&chm=D,0033FF,0,0,5,1",
				},
				map[string]interface{}{
					"type": "image",
					"src":  "https://chart.googleapis.com/chart?chs=600x400&chd=t:6,2,9,5,2,5,7,4,8,2,1&cht=lc&chds=a&chxt=y&chm=D,0033FF,0,0,5,1",
					"href": "https://google.com",
				},
			},
		})
		fmt.Fprint(w, `
{
	"status": "success",
	"message": "Event processed",
	"incident_key": "srv01/HTTP"
}
		`)
	})

	resp, err := pd.Trigger(Trigger{
		Description: "Something bad has happened",
		IncidentKey: "e93facc04764012d7bfb002500d5d1a6",
		Client:      "Sample Monitoring Service",
		ClientURL:   "https://monitoring.service.com",
		Details: map[string]interface{}{
			"ping time": "1500ms",
			"load avg":  0.75,
		},
		Contexts: []Context{
			Context{
				Type: "link",
				HREF: "http://acme.pagerduty.com",
			},
			Context{
				Type: "link",
				HREF: "http://acme.pagerduty.com",
				Text: "View the incident on PagerDuty",
			},
			Context{
				Type: "image",
				SRC:  "https://chart.googleapis.com/chart?chs=600x400&chd=t:6,2,9,5,2,5,7,4,8,2,1&cht=lc&chds=a&chxt=y&chm=D,0033FF,0,0,5,1",
			},
			Context{
				Type: "image",
				SRC:  "https://chart.googleapis.com/chart?chs=600x400&chd=t:6,2,9,5,2,5,7,4,8,2,1&cht=lc&chds=a&chxt=y&chm=D,0033FF,0,0,5,1",
				HREF: "https://google.com",
			},
		},
	})
	if assert.Nil(t, err) {
		assert.Equal(t, "srv01/HTTP", resp.IncidentKey)
	}
}
