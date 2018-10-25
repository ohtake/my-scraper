package instagram

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSource(t *testing.T) {
	source := NewSource(http.DefaultClient, "fukkachan628")
	assert.Equal(t, http.DefaultClient, source.httpClient)
	assert.Equal(t, "fukkachan628", source.userID)
	assert.Equal(t, baseURL, source.baseURL)
}

func TestScrape(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/fukkachan628/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "testdata/www.instagram.com/fukkachan628/index.html")
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	source := NewSource(server.Client(), "fukkachan628")
	source.baseURL = server.URL

	feed, err := source.Scrape()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "ふっかちゃん【公式】", feed.Title)
	assert.Equal(t, "https://www.instagram.com/fukkachan628/", feed.Link.Href)
	assert.Equal(t, 12, len(feed.Items))
	assert.Equal(t, "メル助（@menicon_melsuke）とサニーちゃんとティックトックY(o≧ω≦o)Yたのしす〜♪また遊ぼうねぇY(o0ω★o)Y", feed.Items[0].Title)
	assert.Equal(t, "https://www.instagram.com/p/BpWQzt0FYM1/", feed.Items[0].Link.Href)
	assert.WithinDuration(t, time.Date(2018, 10, 25, 7, 34, 19, 0, time.UTC), feed.Items[0].Created, 0)
}
