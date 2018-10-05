package browsercookies

import (
	"errors"
	"net/http"
)

// Chrome is the cookie jar which loads all your cookies in FireFox.
type Chrome struct {
	http.CookieJar
}

// LoadChrome will load all your cookies from FireFox.
func LoadChrome() (*FireFox, error) {
	return nil, errors.New("Load cookies from Chrome is not implemented")
}
