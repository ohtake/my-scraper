package seibuenevent

import (
	"bufio"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/stretchr/testify/assert"
)

var _ scraper.Source = (*SeibuenEventSource)(nil)

func TestSource(t *testing.T) {
	f, err := os.Open("testdata/www.seibu-leisure.co.jp/event/index.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	source := NewSource()
	feed, err := source.ScrapeFromDocument(doc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(feed.Items))
}