package superfeedr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrieve(t *testing.T) {
	topic := "https://pinub.github.io/superfeedr/github.com-blog.atom"
	s := NewSuperfeedr(Config{
		Username: "username",
		Password: "password",
		URL:      "https://pinub.github.io/superfeedr/github.com-blog.atom",
	})

	feed, err := s.Retrieve(topic)
	assert.Nil(t, err)
	assert.NotNil(t, feed)
}

func Test(t *testing.T) {
	s := NewSuperfeedr(Config{})

	assert.NotNil(t, s)
	assert.Equal(t, "*superfeedr.Superfeedr", fmt.Sprintf("%T", s))
}
