package scraper

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

const (
	SeibuenEventUrl = "http://www.seibuen-yuuenchi.jp/event/index.html?category=e1"
)

type SeibuenEventSource struct {
}

func NewSeibuenEventSource() *SeibuenEventSource {
	return &SeibuenEventSource{}
}

func (s *SeibuenEventSource) Scrape() (*feeds.Feed, error) {
	doc, err := goquery.NewDocument(SeibuenEventUrl)
	if err != nil {
		return nil, err
	}
	return s.ScrapeFromDocument(doc)
}

func (s *SeibuenEventSource) ScrapeFromDocument(doc *goquery.Document) (*feeds.Feed, error) {
	titleReplacer := strings.NewReplacer("『", "", "』", "")
	textReplacer := strings.NewReplacer("\n", "", "\t", "")

	var items []*feeds.Item
	var (
		title string
	)
	doc.Find(".elem-section > div > div > div > div > div").Each(func(_ int, s *goquery.Selection) {
		switch {
		case s.HasClass("elem-heading-lv3"):
			title = titleReplacer.Replace(s.Find("h3").Text())
		case s.HasClass("elem-pic-block"):
			properties := map[string]string{}
			s.Find("table tr").Each(func(_ int, s *goquery.Selection) {
				key := textReplacer.Replace(strings.TrimSpace(s.Find("th").Text()))
				value := textReplacer.Replace((strings.TrimSpace(s.Find("td").Text())))
				properties[key] = value
			})
			if len(properties) == 0 {
				return
			}

			summary := textReplacer.Replace(s.Find(".txt-box > div > .txt-body > div > .elem-paragraph").Text())

			description := fmt.Sprintf("%s<br /><br />日程: %s<br />時間: %s<br />場所: %s<br />その他: %s", summary, properties["日程"], properties["時間"], properties["場所"], properties["その他"])

			hash := fmt.Sprintf("%x", sha256.Sum256([]byte(title+properties["日程"])))

			items = append(items, &feeds.Item{
				Title:       title,
				Description: description,
				Link:        &feeds.Link{Href: SeibuenEventUrl},
				Id:          hash,
			})
		}
	})

	feed := &feeds.Feed{
		Title: "西武園ゆうえんち メルヘンタウン",
		Link:  &feeds.Link{Href: SeibuenEventUrl},
		Items: items,
	}

	return feed, nil
}