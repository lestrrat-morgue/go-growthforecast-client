go-growthforecast-client
========================

[![Build Status](https://travis-ci.org/lestrrat/go-growthforecast-client.png)](https://travis-ci.org/lestrrat/go-growthforecast-client)

[![Coverage Status](https://coveralls.io/repos/lestrrat/go-growthforecast-client/badge.png?branch=HEAD)](https://coveralls.io/r/lestrrat/go-growthforecast-client?branch=HEAD)

GrowthForecast Client In Golang

```go
import (
    "log"
    gf "github.com/lestrrat/go-growthforecast-client"
)

func main() {
    client := gf.NewClient("http://gf.mycompany.com")
    err := client.Post("service/section/graph", &Graph{ Number: 1 })
    if err != nil {
        log.Printf("Failed to post to GrowthForecast: %s", err)
    }
}
```

## Auto-Generated API Docs

[http://godoc.org/github.com/lestrrat/go-growthforecast-client](http://godoc.org/github.com/lestrrat/go-growthforecast-client)
