// Package superfeedr provides basic methods for the superfeedr api.
package superfeedr

import (
	"net/http"
	"net/url"
	"time"
)

// Feed contains the title and multiple Items
type Feed struct {
	Title string  `json:"title"`
	Items []*Item `json:"items"`
}

// Item contains all infos of an item.
type Item struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Content   string    `json:"content"`
	Language  string    `json:"language"`
	Published time.Time `json:"published"`
	Updated   time.Time `json:"updated"`
	Links     []*Link
	// standardLinks interface{} `json:"standardLinks"`
}

// Link represents a structured link.
type Link struct {
	Title string `json:"title"`
	Href  string `json:"href"`
	Rel   string `json:"rel"`
	Type  string `json:"alternate"`
}

// Config struct for all possible configuration parameters.
type Config struct {
	Username string
	Password string
	URL      string
}

const (
	defaultBaseURL = "https://push.superfeedr.com"
)

type (
	Client struct {
		client *http.Client

		BaseURL *url.URL

		// Reuse a single struct.
		common service

		// Services used for talking to different parts of the API.
		Retrieve    *RetrieveService
		Subscribe   *SubscribeService
		Unsubscribe *UnsubscribeService
		List        *ListService
	}

	service struct {
		client *Client
	}
)

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL}
	c.common.client = c
	c.Retrieve = (*RetrieveService)(&c.common)
	c.Subscribe = (*SubscribeService)(&c.common)
	c.Unsubscribe = (*UnsubscribeService)(&c.common)
	c.List = (*ListService)(&c.common)

	return c
}

// Superfeedr represents the object used to work with.
type Superfeedr struct {
	config Config
}

// NewSuperfeedr creates and sets the basic attributes.
func NewSuperfeedr(config Config) *Superfeedr {
	if config.URL == "" {
		config.URL = "https://push.superfeedr.com"
	}

	return &Superfeedr{config: config}
}

func (f *Superfeedr) client(method string) (*http.Request, error) {
	if method == "" {
		method = "GET"
	}

	req, err := http.NewRequest(method, f.config.URL, nil)
	if err != nil {
		return nil, err
	}

	if f.config.Username != "" && f.config.Password != "" {
		req.SetBasicAuth(f.config.Username, f.config.Password)
	}

	return req, nil
}

// Retrieve entries for the given topic. You must be a subscriber of the
// given topic.
func (f *Superfeedr) Retrieve(topic string) (*Feed, error) {
	req, err := f.client("GET")
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("hub.mode", "retrieve")
	q.Add("hub.topic", topic)
	q.Add("format", "json")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// return nil, err
	// }

	return &Feed{}, nil
}
