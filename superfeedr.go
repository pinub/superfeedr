// Package superfeedr provides basic methods for the superfeedr api.
package superfeedr

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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

type Client struct {
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

type service struct {
	client *Client
}

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

func (c *Client) NewRequest(method, urlString string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				err = nil
			}
		}
	}

	return resp, err
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	return nil
}
