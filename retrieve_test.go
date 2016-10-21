package superfeedr

import (
	"net/http"
	"testing"
)

func TestGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
	})

}
