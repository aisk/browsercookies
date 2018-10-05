package browsercookies

import (
	"database/sql"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/go-ini/ini"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mitchellh/go-homedir"
	"golang.org/x/net/publicsuffix"
)

// FireFox is the cookie jar which loads all your cookies in FireFox.
type FireFox struct {
	http.CookieJar
}

func (ff *FireFox) parseProfile(profile string) (string, error) {
	cfg, err := ini.Load(profile)
	if err != nil {
		return "", err
	}
	// TODO: return first profile when no default default profile set.
	for _, section := range cfg.Sections() {
		if section.Key("Default").String() != "1" {
			continue
		}
		path := section.Key("Path").String()
		if section.Key("IsRelative").String() == "1" {
			path = filepath.Join(filepath.Dir(profile), path)
		}
		return path, nil
	}

	return "", errors.New("No default Firefox profile found")
}

func (ff *FireFox) findDefaultProfile() (string, error) {
	homedir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homedir, filepath.FromSlash("Library/Application Support/Firefox/profiles.ini")), nil
	case "linux":
		return filepath.Join(homedir, filepath.FromSlash(".mozilla/firefox/profiles.ini")), nil
	case "win32":
		return filepath.Join(os.Getenv("APPDATA"), filepath.FromSlash("Mozilla/Firefox/profiles.ini")), nil
	}
	return "", errors.New("Unsupported operating system: " + runtime.GOOS)
}

func (ff *FireFox) findCookieFiles() (string, error) {
	profile, err := ff.findDefaultProfile()
	if err != nil {
		return "", nil
	}
	profile, err = ff.parseProfile(profile)
	if err != nil {
		return "", nil
	}
	return filepath.Join(profile, "cookies.sqlite"), nil
}

func (ff *FireFox) getCookiesFromFile(file string) ([]*http.Cookie, error) {
	tempfile, err := ioutil.TempFile("", "go-webcookies")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempfile.Name())
	origin, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(tempfile, origin); err != nil {
		return nil, err
	}
	tempfile.Close()

	db, err := sql.Open("sqlite3", tempfile.Name())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT host, path, isSecure, expiry, name, value FROM moz_cookies")
	if err != nil {
		return nil, err
	}
	cookies := []*http.Cookie{}
	for rows.Next() {
		var host, path, name, value string
		var secure, expiry int
		if err := rows.Scan(&host, &path, &secure, &expiry, &name, &value); err != nil {
			return cookies, err
		}
		cookie := &http.Cookie{
			Name:    name,
			Value:   value,
			Path:    path,
			Domain:  host,
			Expires: time.Unix(int64(expiry), 0),
			Secure:  secure == 1,
		}
		cookies = append(cookies, cookie)
	}
	return cookies, nil
}

func (ff *FireFox) getCookies() ([]*http.Cookie, error) {
	file, err := ff.findCookieFiles()
	if err != nil {
		return nil, err
	}
	return ff.getCookiesFromFile(file)
}

// LoadFireFox will load all your cookies from FireFox.
func LoadFireFox() (*FireFox, error) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	ff := &FireFox{jar}
	cookies, err := ff.getCookies()
	if err != nil {
		return nil, err
	}
	for _, cookie := range cookies {
		scheme := "http://"
		if cookie.Secure {
			scheme = "https://"
		}
		u, err := url.Parse(scheme + cookie.Domain + cookie.Path)
		if err != nil {
			return nil, err
		}
		ff.SetCookies(u, []*http.Cookie{cookie})
	}
	return ff, nil
}
