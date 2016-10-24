package superfeedr

import "net/http"

// RetrieveService handles communication for retrieving all recent
// notifications of a topic.
type RetrieveService service

// Get notifications items for the given topic.
func (s *RetrieveService) Get(topic string) (*Feed, *http.Response, error) {
	req, err := s.client.NewRequest("GET", "/", nil)
	if err != nil {
		return nil, nil, err
	}

	req.AddOptions(map[string]string{
		"hub.mode":  "retrieve",
		"hub.topic": topic,
		"format":    "json",
	})

	feed := new(Feed)
	resp, err := s.client.Do(req, feed)
	if err != nil {
		return nil, resp, err
	}

	return feed, resp, err
}
