// Package superfeedr provides basic methods for the superfeedr api.
package superfeedr

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

// Client manages cummunication with the superfeedr API.
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

// NewClient returns a superfeedr API client. If no httpClient is provided,
// the default http.DefaultClient will be used. Authentication requires a
// non default client.
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

// NewRequest created a new API request. A relative url can be provides in the
// urlString to be resolved relative to the baseURL of the Client.
func (c *Client) NewRequest(method, urlString string, body interface{}) (*Request, error) {
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

	return newRequest(req), nil
}

// Do sends the API request and returns the API response. The response is
// JSON decoded into the given v interface. If v implements the io.Writer
// interface, the raw response will be written to v, without decoding it.
func (c *Client) Do(req *Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req.Request)
	if err != nil {
		return nil, err
	}

	defer func() {
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	err = checkResponse(resp)
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

// Request is used for abstracting the http.Request to provide a method to add
// query parameters to a request (AddOptions).
type Request struct {
	*http.Request
}

func newRequest(r *http.Request) *Request {
	return &Request{Request: r}
}

// AddOptions adds the given associative array to the request as a query
// parameter.
func (r *Request) AddOptions(options map[string]string) {
	q := r.URL.Query()

	for k, v := range options {
		q.Add(k, v)
	}

	r.URL.RawQuery = q.Encode()
}

// ErrorResponse is used for abstracting the http.Response to provide a more
// easier way to access errors caused by the API request.
type ErrorResponse struct {
	Response *http.Response
}

// Error returns the error caused by the API.
func (r *ErrorResponse) Error() string {
	return fmt.Sprintf(
		"%v %v: %d",
		r.Response.Request.Method,
		r.Response.Request.URL,
		r.Response.StatusCode,
	)
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}

	return errorResponse
}
