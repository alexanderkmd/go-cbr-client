package main

import (
	"fmt"
	"time"

	cbr "github.com/ivanglie/go-cbr-client"
)

func main() {
	client := cbr.NewClient()

	// For float64 value:
	rateFloat64, err := client.GetRate("USD", time.Now())

	if err != nil {
		panic(err)
	}
	fmt.Println(rateFloat64)

	// For Decimal value:
	rateDecimal, err := client.GetRateDecimal("USD", time.Now())

	if err != nil {
		panic(err)
	}
	fmt.Println(rateDecimal)

	// For String value with dot as decimal separator:
	rateString, err := client.GetRateString("USD", time.Now())

	if err != nil {
		panic(err)
	}
	fmt.Println(rateString)
}
