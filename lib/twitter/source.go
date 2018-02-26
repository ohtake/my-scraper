package twitter

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"

	"github.com/ChimeraCoder/anaconda"
	"github.com/gorilla/feeds"
	"github.com/pkg/errors"
)

type TwitterSource struct {
	userId int64
}

var (
	once       sync.Once
	twitterApi *anaconda.TwitterApi
)

func getTwitterApi() *anaconda.TwitterApi {
	once.Do(func() {
		twitterApi = anaconda.NewTwitterApiWithCredentials(
			os.Getenv("TWITTER_OAUTH_TOKEN"),
			os.Getenv("TWITTER_OAUTH_TOKEN_SECRET"),
			os.Getenv("TWITTER_CONSUMER_KEY"),
			os.Getenv("TWITTER_CONSUMER_SECRET"))
	})
	return twitterApi
}

func NewSource(userId int64) *TwitterSource {
	return &TwitterSource{
		userId: userId,
	}
}

func (s *TwitterSource) Scrape() (*feeds.Feed, error) {
	api := getTwitterApi()

	values := url.Values{}
	values.Set("user_id", strconv.FormatInt(s.userId, 10))
	values.Set("count", "100")
	timeline, err := api.GetUserTimeline(values)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return s.Render(timeline)
}

func (s *TwitterSource) Render(timeline []anaconda.Tweet) (*feeds.Feed, error) {
	if len(timeline) == 0 {
		return nil, errors.New("timeline is empty")
	}
	user := timeline[0].User
	userURL := fmt.Sprintf("https://twitter.com/%s", user.ScreenName)
	items := make([]*feeds.Item, 0, len(timeline))
	for _, tweet := range timeline {
		created, err := tweet.CreatedAtTime()
		if err != nil {
			continue
		}
		items = append(items, &feeds.Item{
			Title:   tweet.Text,
			Created: created,
			Link:    &feeds.Link{Href: fmt.Sprintf("%s/status/%s", userURL, tweet.IdStr)},
		})
	}
	return &feeds.Feed{
		Title: fmt.Sprintf("%s (@%s)", user.Name, user.ScreenName),
		Link:  &feeds.Link{Href: userURL},
		Items: items,
	}, nil
}
