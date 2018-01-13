package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/feeds"
	"github.com/joho/godotenv"
	"github.com/lestrrat/go-server-starter/listener"
	"github.com/mono0x/my-scraper/lib"
	"github.com/mono0x/my-scraper/lib/facebook"
	"github.com/mono0x/my-scraper/lib/twitter"
	"github.com/pkg/errors"
)

func renderFeed(w http.ResponseWriter, feed *feeds.Feed) {
	w.Header().Set("Content-Type", "application/atom+xml")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	if err := feed.WriteAtom(w); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func sourceRenderer(source scraper.Source) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		feed, err := source.Scrape()
		if err != nil {
			log.Printf("%+v\n", err)
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		renderFeed(w, feed)
	}
}

func run() error {
	listeners, err := listener.ListenAll()
	if err != nil {
		return errors.WithStack(err)
	}

	var l net.Listener
	if len(listeners) > 0 {
		l = listeners[0]
	} else {
		l, err = net.Listen("tcp", ":13000")
		if err != nil {
			return errors.WithStack(err)
		}
	}

	mux := http.NewServeMux()

	entries := []struct {
		Path   string
		Source scraper.Source
	}{
		{"/character-show", scraper.NewCharacterShowSource()},
		{"/fukoku-life", scraper.NewFukokuLifeEventSource()},
		{"/harmonyland-info", scraper.NewHarmonylandInfoSource()},
		{"/kittychan-info", scraper.NewKittychanInfoSource()},
		{"/memoirs-of-shibasaki-saki", scraper.NewMemoirsOfShibasakiSakiSource()},
		{"/prtimes-sanrio", scraper.NewPRTimesSource()},
		{"/puroland-info", scraper.NewPurolandInfoSource()},
		{"/sanrio-news-release", scraper.NewSanrioNewsReleaseSource()},
		{"/seibuen-event", scraper.NewSeibuenEventSource()},
		{"/value-press-sanrio", scraper.NewValuePressSource()},
		{"/yuyakekoyake-news", scraper.NewYuyakekoyakeNewsSource()},
	}
	for _, entry := range entries {
		mux.HandleFunc(entry.Path, sourceRenderer(entry.Source))
	}

	mux.HandleFunc("/facebook", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		source := facebook.NewSource(id)
		sourceRenderer(source)(w, r)
	})

	mux.HandleFunc("/google-calendar", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		source := scraper.NewGoogleCalendarSource(id)
		sourceRenderer(source)(w, r)
	})

	mux.HandleFunc("/instagram", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		id := query.Get("id")
		if id == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		source := scraper.NewInstagramSource(id)
		sourceRenderer(source)(w, r)
	})

	mux.HandleFunc("/twitter", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		idStr := query.Get("id")
		if idStr == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		source := twitter.NewSource(id)
		sourceRenderer(source)(w, r)
	})

	server := http.Server{Handler: mux}

	go func() {
		if err := server.Serve(l); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGTERM)

	for {
		s := <-signalChan
		if s == syscall.SIGTERM {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			return server.Shutdown(ctx)
		}
	}
}

func main() {
	log.SetFlags(log.Lshortfile)

	_ = godotenv.Load()

	if err := run(); err != nil {
		log.Fatalf("%+v\n", err)
	}
}
