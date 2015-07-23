package main

import (
	"bufio"
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetPurolandNewsFromDocument(t *testing.T) {
	f, err := os.Open("data/puroland.html")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	feed, err := GetPurolandNewsFromDocument(doc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(feed.Items), 10)
}

func TestGetPurolandInfoFromDocument(t *testing.T) {
	f, err := os.Open("data/puroland.html")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	feed, err := GetPurolandInfoFromDocument(doc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(feed.Items), 49)
}