package cbr

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestGetRate(t *testing.T) {
	client := NewClient()
	client.SetLogLevel(logrus.DebugLevel)
	rate, err := client.GetRate("USD", time.Now())
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, rate, float64(1))
}

func TestGetRate_Error(t *testing.T) {
	client := NewClient()
	rate, err := client.GetRate(" ", time.Now())
	assert.Error(t, err)
	assert.Equal(t, float64(0), rate)
}

func TestSetFetchFunction(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	c := client{nil, false, logger.GetLevel(), *logger}
	c.SetFetchFunction(func(url string) (resp *http.Response, err error) { return http.Get(url) })
	assert.Equal(t, reflect.Func, reflect.TypeOf(c.fetch).Kind())
}

func TestSetBaseUrl(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	c := client{nil, false, logger.GetLevel(), *logger}

	// test for error with empty URL
	err := c.SetBaseUrl("")
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("empty base URL was provided"), err)
	}

	// test no error for URL
	err = c.SetBaseUrl("http://example.com")
	assert.Equal(t, nil, err)
}
