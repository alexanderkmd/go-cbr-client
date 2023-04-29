package cbr

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/encoding/charmap"
)

var baseURL = "http://www.cbr.ru/scripts/XML_daily_eng.asp"

const (
	dateFormat = "02/01/2006"
)

// Cache for requests
var cache map[string]Result
var cacheHits int

// Currency is a currency item
type Currency struct {
	ID       string `xml:"ID,attr"`
	NumCode  uint   `xml:"NumCode"`
	CharCode string `xml:"CharCode"`
	Nom      uint   `xml:"Nominal"`
	Name     string `xml:"Name"`
	Value    string `xml:"Value"`
}

// Returns properly formatted currency Value string
func (cur Currency) ValueString() string {
	return strings.Replace(cur.Value, ",", ".", -1)
}

// Returns currency Value in float64, without nominal correction
func (cur Currency) ValueFloatRaw() (float64, error) {
	return strconv.ParseFloat(cur.ValueString(), 64)
}

// Returns currency Value in float64, corrected by nominal
func (cur Currency) ValueFloat() (float64, error) {
	res, err := cur.ValueFloatRaw()
	if err != nil {
		return res, err
	}
	return res / float64(cur.Nom), nil
}

// Returns currency Value in Decimal, without nominal correction
//
// Rationale: https://pkg.go.dev/github.com/shopspring/decimal - FAQ section
func (cur Currency) ValueDecimalRaw() (decimal.Decimal, error) {
	return decimal.NewFromString(cur.ValueString())
}

// Returns currency Value in Decimal, corrected by nominal
//
// Rationale: https://pkg.go.dev/github.com/shopspring/decimal - FAQ section
func (cur Currency) ValueDecimal() (decimal.Decimal, error) {
	res, err := cur.ValueDecimalRaw()
	if err != nil {
		return res, err
	}
	nominal := decimal.NewFromInt(int64(cur.Nom))
	return res.Div(nominal), nil
}

// Result is a result representation
type Result struct {
	XMLName    xml.Name   `xml:"ValCurs"`
	Date       string     `xml:"Date,attr"`
	Currencies []Currency `xml:"Valute"`
}

func getRate(currency string, t time.Time, fetch fetchFunction, useCache bool, logger *logrus.Logger) (float64, error) {
	logger.Debugf("Fetching the currency rate for %s at %v", currency, t.Format("02.01.2006"))

	curr, err := getCurrency(currency, t, fetch, useCache, logger)
	if err != nil {
		return 0, err
	}
	return curr.ValueFloat()
}

func getRateDecimal(currency string, t time.Time, fetch fetchFunction, useCache bool, logger *logrus.Logger) (decimal.Decimal, error) {
	logger.Debugf("Fetching the currency rate for %s at %v\n  in Decimal", currency, t.Format("02.01.2006"))

	curr, err := getCurrency(currency, t, fetch, useCache, logger)
	if err != nil {
		return decimal.Zero, err
	}
	return curr.ValueDecimal()
}

func getRateString(currency string, t time.Time, fetch fetchFunction, useCache bool, logger *logrus.Logger) (string, error) {
	logger.Debugf("Fetching the currency rate string for %s at %v", currency, t.Format("02.01.2006"))

	curr, err := getCurrency(currency, t, fetch, useCache, logger)
	if err != nil {
		return "", err
	}
	return curr.ValueString(), nil
}

func getCurrency(currency string, t time.Time, fetch fetchFunction, useCache bool, logger *logrus.Logger) (Currency, error) {
	result, err := getCurrenciesCacheOrRequest(t, fetch, useCache, logger)
	if err != nil {
		return Currency{}, err
	}

	for _, v := range result.Currencies {
		if v.CharCode == currency {
			return v, nil
		}
	}
	return Currency{}, fmt.Errorf("unknown currency: %s", currency)
}

func getCurrenciesCacheOrRequest(t time.Time, fetch fetchFunction, useCache bool, logger *logrus.Logger) (Result, error) {
	formatedDate := t.Format(dateFormat)

	result := Result{}

	// if currencies were already requested for this date - return from cache, if it is used
	if cachedResult, exist := cache[formatedDate]; exist && useCache {
		logger.Info("Got currency data from cache!")
		result = cachedResult
		cacheHits += 1
	} else {
		err := getCurrencies(&result, t, fetch, logger)
		if err != nil {
			logger.Errorf("Error getiing Currencies: %s", err)
			logger.Debug(result)
			return result, err
		}
		if useCache {
			// if cache is used - put result to cache
			if len(cache) == 0 {
				// if cache is empty - initialize
				cache = make(map[string]Result)
			}
			cache[formatedDate] = result
		}
	}
	return result, nil
}

func getCurrencies(v *Result, t time.Time, fetch fetchFunction, logger *logrus.Logger) error {
	url := baseURL + "?date_req=" + t.Format(dateFormat)
	logger.Debug(url)
	if fetch == nil {
		logger.Error("Empty fetch function provided")
		return errors.New("fetch is empty")
	}
	resp, err := fetch(url)
	if err != nil {
		logger.Errorf("Error fetching URL: %s", err)
		return err
	}
	status := resp.StatusCode
	logger.Debugf("Response status code: %v:", status)
	if status != 200 {
		logger.Errorf("Request returned abnormal status code: %v", status)
		logger.Debug(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Error reading response Body: %s", err)
		return err
	}
	defer resp.Body.Close()

	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}
	err = decoder.Decode(&v)
	if err != nil {
		logger.Errorf("Error decoding XML: %s", err)
		logger.Debugf("%s", body)
		return err
	}

	return nil
}
