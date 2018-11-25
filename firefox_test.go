package browsercookies

import (
	"net/http"
	"testing"
)

func TestFireFoxFindDefaultProfile(t *testing.T) {
	ff := &FireFox{}
	_, err := ff.findDefaultProfile()
	if err != nil {
		t.Fatal(err)
	}
}

func TestFireFoxParseProfile(t *testing.T) {
	ff := &FireFox{}
	profile, err := ff.findDefaultProfile()
	if err != nil {
		t.Fatal(err)
	}
	_, err = ff.parseProfile(profile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFireFoxfindCookieFiles(t *testing.T) {
	ff := &FireFox{}
	file, err := ff.findCookieFiles()
	if err != nil {
		t.Fatal(err)
	}
	if file == "" {
		t.Fatal("Empty cookie file")
	}
}

func TestFireFoxgetCookies(t *testing.T) {
	ff := &FireFox{}
	_, err := ff.getCookies()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadFireFox(t *testing.T) {
	jar, err := LoadFireFox()
	if err != nil {
		t.Fatal(err)
	}
	_ = jar
}

func TestFetchFireFox(t *testing.T) {
	jar, err := LoadFireFox()
	if err != nil {
		t.Fatal(err)
	}
	client := http.Client{
		Jar: jar,
		// disable follow redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
	resp, err := client.Get("https://github.com/settings/profile")
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatal("Load firefox cookies failed, or github.com is not logined")
	}
}
