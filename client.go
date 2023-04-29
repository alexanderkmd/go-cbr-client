package cbr

import (
	"errors"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

// fetchFunction is a function that mimics http.Get() method
type fetchFunction func(url string) (resp *http.Response, err error)

// Client is a currency rates service client... what else?
type Client interface {
	// Returns currency rate in float64
	GetRate(string, time.Time) (float64, error)

	// Returns currency rate in Decimal
	//
	// Rationale: https://pkg.go.dev/github.com/shopspring/decimal - FAQ section
	GetRateDecimal(string, time.Time) (decimal.Decimal, error)

	// Returns currency rate string with dot as decimal separator
	GetRateString(string, time.Time) (string, error)

	// Returns currency struct
	GetCurrencyInfo(string, time.Time) (Currency, error)

	SetFetchFunction(fetchFunction)
	SetBaseUrl(string) error
	SetLogLevel(logrus.Level)
}

type client struct {
	fetch    fetchFunction
	UseCache bool
	logLevel logrus.Level
	logger   logrus.Logger
}

// Returns currency rate in float64
func (s *client) GetRate(currency string, t time.Time) (float64, error) {
	return getRate(currency, t, s.fetch, s.UseCache, &s.logger)
}

// Returns currency rate in Decimal
//
// Rationale: https://pkg.go.dev/github.com/shopspring/decimal - FAQ section
func (s *client) GetRateDecimal(currency string, t time.Time) (decimal.Decimal, error) {
	return getRateDecimal(currency, t, s.fetch, s.UseCache, &s.logger)
}

// Returns currency rate string with dot as decimal separator
func (s *client) GetRateString(currency string, t time.Time) (string, error) {
	return getRateString(currency, t, s.fetch, s.UseCache, &s.logger)
}

// Returns currency struct
func (s *client) GetCurrencyInfo(currency string, t time.Time) (Currency, error) {
	return getCurrency(currency, t, s.fetch, s.UseCache, &s.logger)
}

func (s *client) SetFetchFunction(f fetchFunction) {
	s.fetch = f
}

func (s *client) SetLogLevel(logLevel logrus.Level) {
	s.logger.SetLevel(logLevel)
	s.logLevel = logLevel
}

// Sets alternative baseUrl for compatible API
func (s *client) SetBaseUrl(url string) error {
	if url == "" {
		return errors.New("empty base URL was provided")
	}
	baseURL = url
	return nil
}

// NewClient creates a new rates service instance
func NewClient() Client {
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	return &client{
		fetch:    http.Get,
		UseCache: true,
		logger:   *logger,
		logLevel: logrus.WarnLevel,
	}
}
