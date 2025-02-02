# Golang client for the Central Bank of the Russian Federation currency rates API

[![Go Reference](https://pkg.go.dev/badge/github.com/alexanderkmd/go-cbr-client.svg)](https://pkg.go.dev/github.com/alexanderkmd/go-cbr-client)
[![Test](https://github.com/alexanderkmd/go-cbr-client/actions/workflows/test.yml/badge.svg)](https://github.com/alexanderkmd/go-cbr-client/actions/workflows/test.yml)
[![Codecov](https://codecov.io/gh/alexanderkmd/go-cbr-client/branch/master/graph/badge.svg?token=46HUJQAM56)](https://codecov.io/gh/alexanderkmd/go-cbr-client)

go-cbr-client is a fork of [matperez's](https://github.com/matperez) [client](https://github.com/matperez/go-cbr-client) and  [ivangile's](https://github.com/ivanglie) [client](https://github.com/ivanglie/go-cbr-client) for [CBRF API](http://www.cbr.ru/development/).

Request cache and different output types added.

## Example

First, ensure the library is installed and up to date by running ```go get -u github.com/alexanderkmd/go-cbr-client```.

This is a very simple app that just displays exchange rate of US dollar.

```golang
package main

import (
    "fmt"
    "time"

    cbr "github.com/alexanderkmd/go-cbr-client"
)

func main() {
    client := cbr.NewClient()
    
    // For float64 value:
    rateFloat64, err := client.GetRate("USD", time.Now())

    // For Decimal value:
    rateDecimal, err := client.GetRateDecimal("USD", time.Now())

    // For String value with dot as decimal separator:
    rateString, err := client.GetRateString("USD", time.Now())

    if err != nil {
        panic(err)
    }
    fmt.Println(rateFloat64)
}
```

See [main.go](./_example/main.go).

## Set Alternative API point

Due to lots of IP bans from CBR (HTTP error 403) - you can change the different/alternative compatible API URL.

For example:: [www.cbr-xml-daily.ru](https://www.cbr-xml-daily.ru/) provide one.

```golang
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
```

## References

For more information check out the following links:

* [CBRF API](http://www.cbr.ru/development/SXML/)
* [CBRF technical resources](http://www.cbr.ru/eng/development/) (EN)
