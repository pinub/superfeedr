package superfeedr

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testQuery(t, r, "hub.mode=retrieve")
		testQuery(t, r, "hub.topic="+url.QueryEscape("http://www.example.com/test.atom"))
		testQuery(t, r, "format=json")
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"title":"Example Title","items":[{"id":"tag:1234"},{"id":"tag:5678"}]}`)
	})

	resp, _, err := client.Retrieve.Get("http://www.example.com/test.atom")
	if err != nil {
		t.Errorf("Retrieve.Get returned error: %v", err)
	}

	want := &Feed{Title: "Example Title", Items: []*Item{{ID: "tag:1234"}, {ID: "tag:5678"}}}
	if !reflect.DeepEqual(resp, want) {
		t.Errorf("Retrieve.Get returned %+v, want %+v", resp, want)
	}
}
