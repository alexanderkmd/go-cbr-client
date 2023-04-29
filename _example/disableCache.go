package main

import (
	"fmt"
	"time"

	cbr "github.com/alexanderkmd/go-cbr-client"
)

func main() {
	client := cbr.NewClient()

	client.UseCache = false // disables cache for all requests

	rate, err := client.GetRate("USD", time.Now())
	if err != nil {
		panic(err)
	}
	fmt.Println(rate)
}
