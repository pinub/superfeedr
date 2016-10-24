package superfeedr

import "net/http"

type RetrieveService service

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
