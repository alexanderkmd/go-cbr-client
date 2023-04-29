package cbr

import (
	"net/http"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

/*
func Test_Debug(t *testing.T) {
	Debug = true
	getRate("CNY", time.Now(), nil, false)
	assert.True(t, Debug)

	Debug = false
	getRate("CNY", time.Now(), nil, false)
	assert.False(t, Debug)
}
*/

func Test_getRate_Error(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	rate, err := getRate("CNY", time.Now(), nil, false, logger)
	assert.NotNil(t, err)
	assert.Equal(t, float64(0), rate)
}

// Check for the cache functionality
func Test_getRate_Cache(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	loc, _ := time.LoadLocation("Europe/Moscow")
	testDate := time.Date(2022, 12, 16, 1, 1, 1, 0, time.UTC).In(loc)
	tmpHits := cacheHits

	timingDate := time.Now()

	// Get first request uncached
	rate, err := getRate("USD", testDate, http.Get, true, logger)
	assert.Nil(t, err)
	assert.NotEqual(t, float64(0), rate)
	prevTime := logElapsedTime(timingDate, timingDate, logger)

	// Get second request from cache
	rate, err = getRate("EUR", testDate, http.Get, true, logger)
	assert.Nil(t, err)
	assert.NotEqual(t, float64(0), rate)
	assert.Greater(t, cacheHits, tmpHits)
	tmpHits = cacheHits
	prevTime = logElapsedTime(timingDate, prevTime, logger)

	// Get third request from cache
	rate, err = getRate("TMT", testDate, http.Get, true, logger)
	assert.Nil(t, err)
	assert.NotEqual(t, float64(0), rate)
	assert.Greater(t, cacheHits, tmpHits)
	tmpHits = cacheHits
	cacheLen := len(cache)
	prevTime = logElapsedTime(timingDate, prevTime, logger)

	// Get NOT from cache
	testDate = time.Date(2022, 12, 24, 1, 1, 1, 0, time.UTC).In(loc)
	rate, err = getRate("TMT", testDate, http.Get, true, logger)
	assert.Nil(t, err)
	assert.NotEqual(t, float64(0), rate)
	assert.Equal(t, cacheHits, tmpHits)     // no hit to cache was made
	assert.Greater(t, len(cache), cacheLen) // new item appeared in cache
	tmpHits = cacheHits
	cacheLen = len(cache)
	prevTime = logElapsedTime(timingDate, prevTime, logger)

	// Get 4th FROM cache
	rate, err = getRate("USD", testDate, http.Get, true, logger)
	assert.Nil(t, err)
	assert.NotEqual(t, float64(0), rate)
	assert.Greater(t, cacheHits, tmpHits) // hit to cache was made
	assert.Equal(t, len(cache), cacheLen) // new item does not appear in cache
	_ = logElapsedTime(timingDate, prevTime, logger)
}

// Check for the cache functionality
func Test_getRate_CacheDisabled(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	cache = make(map[string]Result)
	loc, _ := time.LoadLocation("Europe/Moscow")
	testDate := time.Date(2022, 12, 10, 1, 1, 1, 0, time.UTC).In(loc)
	tmpHits := cacheHits

	timingDate := time.Now()

	// Get first request uncached
	rate, err := getRate("USD", testDate, http.Get, false, logger)
	assert.Nil(t, err)
	assert.NotEqual(t, float64(0), rate)
	prevTime := logElapsedTime(timingDate, timingDate, logger)

	// Get second request from cache (disabled)
	rate, err = getRate("EUR", testDate, http.Get, false, logger)
	assert.Nil(t, err)
	assert.NotEqual(t, float64(0), rate)
	assert.Equal(t, cacheHits, tmpHits)
	tmpHits = cacheHits
	prevTime = logElapsedTime(timingDate, prevTime, logger)

	// Get third request from cache (disabled)
	rate, err = getRate("TMT", testDate, http.Get, false, logger)
	assert.Nil(t, err)
	assert.NotEqual(t, float64(0), rate)
	assert.Equal(t, cacheHits, tmpHits)
	tmpHits = cacheHits
	cacheLen := len(cache)
	prevTime = logElapsedTime(timingDate, prevTime, logger)

	// Get NOT from cache
	testDate = time.Date(2022, 12, 20, 1, 1, 1, 0, time.UTC).In(loc)
	rate, err = getRate("TMT", testDate, http.Get, false, logger)
	assert.Nil(t, err)
	assert.NotEqual(t, float64(0), rate)
	assert.Equal(t, cacheHits, tmpHits)   // no hit to cache was made
	assert.Equal(t, len(cache), cacheLen) // new item appeared in cache
	tmpHits = cacheHits
	cacheLen = len(cache)
	prevTime = logElapsedTime(timingDate, prevTime, logger)

	// Get 4th FROM cache (disabled)
	rate, err = getRate("USD", testDate, http.Get, false, logger)
	assert.Nil(t, err)
	assert.NotEqual(t, float64(0), rate)
	assert.Equal(t, cacheHits, tmpHits)   // hit to cache was made
	assert.Equal(t, len(cache), cacheLen) // new item does not appear in cache
	_ = logElapsedTime(timingDate, prevTime, logger)
}

func logElapsedTime(start time.Time, prev time.Time, logger *logrus.Logger) time.Time {
	logger.Infof("Elapsed: %v µs", time.Since(start).Microseconds())
	logger.Infof("Elapsed since previous %vµs, total %vµs", time.Since(prev).Microseconds(), time.Since(start).Microseconds())
	return time.Now()
}

func Test_getRate_Decimal(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	loc, _ := time.LoadLocation("Europe/Moscow")
	testDate := time.Date(2022, 12, 16, 1, 1, 1, 0, time.UTC).In(loc)

	// Get first request uncached
	rate, err := getRateDecimal("USD", testDate, http.Get, true, logger)
	assert.Nil(t, err)
	assert.True(t, rate.Equal(decimal.NewFromFloat(64.3015)))
}

func Test_getRate_String(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	loc, _ := time.LoadLocation("Europe/Moscow")
	testDate := time.Date(2022, 12, 16, 1, 1, 1, 0, time.UTC).In(loc)

	// Get first request uncached
	rate, err := getRateString("USD", testDate, http.Get, true, logger)
	assert.Nil(t, err)
	assert.Equal(t, "64.3015", rate)
}
