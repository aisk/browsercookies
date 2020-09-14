# browsercookies

[![GoDoc](https://godoc.org/github.com/aisk/browsercookies?status.svg)](https://godoc.org/github.com/aisk/browsercookies)
[![Go Report Card](https://goreportcard.com/badge/github.com/aisk/browsercookies)](https://goreportcard.com/report/github.com/aisk/browsercookies)
[![Maintainability](https://api.codeclimate.com/v1/badges/ad3073a115dfe893f2b8/maintainability)](https://codeclimate.com/github/aisk/browsercookies/maintainability)

Make HTTP requests with cookies from your browsers!

![cookie jar](https://www.kitchenistic.com/media/2019/08/best-cookie-jars.jpg)

*This is a Go port of richardpenman's [browsercookie](https://bitbucket.org/richardpenman/browsercookie).*

---

## Supported browsers:

- [x] FireFox
- [ ] Chrome

## Example:

```go
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/aisk/browsercookies"
)

func main() {
	jar, err := browsercookies.LoadFireFox()
	if err != nil {
		panic(err)
	}

	httpclient := http.Client{Jar: jar}
	resp, err := httpclient.Get("https://github.com/settings/profile")
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(`<meta name="user-login" content="(\w*?)">`)
	fmt.Println(re.FindAllStringSubmatch(string(body), -1)[0][1]) // => aisk
}
```

## License:

LGPL
