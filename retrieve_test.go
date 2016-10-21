package superfeedr

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"title":"Example Title","items":[{"id":"tag:1234"},{"id":"tag:5678"}]}`)
	})

	resp, _, err := client.Retrieve.Get("")
	if err != nil {
		t.Errorf("Retrieve.Get returned error: %v", err)
	}

	want := &Feed{Title: "Example Title", Items: []*Item{{ID: "tag:1234"}, {ID: "tag:5678"}}}
	if !reflect.DeepEqual(resp, want) {
		t.Errorf("Retrieve.Get returned %+v, want %+v", resp, want)
	}
}
