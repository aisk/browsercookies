# browsercookies

[![GoDoc](https://godoc.org/github.com/aisk/browsercookies?status.svg)](https://godoc.org/github.com/aisk/browsercookies)
[![Go Report Card](https://goreportcard.com/badge/github.com/aisk/browsercookies)](https://goreportcard.com/report/github.com/aisk/browsercookies)

Make HTTP requests with cookies from your browsers!

![cookie jar](http://howthoughtful.co.za/wp-content/uploads/2017/07/Cookie-Jar-4-600x600.png)

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
