package main

import (
	"fmt"
	"time"

	cbr "github.com/alexanderkmd/go-cbr-client"
)

func main() {
	client := cbr.NewClient()
	client.SetBaseUrl("http://new-base-url.com")

	// For float64 value:
	rateFloat64, err := client.GetRate("USD", time.Now())

	if err != nil {
		panic(err)
	}
	fmt.Println(rateFloat64)
}
