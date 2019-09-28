package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/lestrrat-go/server-starter/listener"
	"github.com/mono0x/my-scraper/lib/server"
	"golang.org/x/sync/errgroup"
)

func run() error {
	listeners, err := listener.ListenAll()
	if err != nil && err != listener.ErrNoListeningTarget {
		return fmt.Errorf("%w", err)
	}

	var l net.Listener
	if len(listeners) > 0 {
		l = listeners[0]
	} else {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		l, err = net.Listen("tcp", ":"+port)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	}

	handler, err := server.NewHandler()
	if err != nil {
		return err
	}
	s := http.Server{Handler: handler}

	eg := errgroup.Group{}
	eg.Go(func() error {
		if err := s.Serve(l); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("%w", err)
		}
		return nil
	})
	eg.Go(func() error {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGTERM, os.Interrupt)
		<-signalChan

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.Shutdown(ctx); err != nil {
			return fmt.Errorf("%w", err)
		}
		return nil
	})
	return eg.Wait()
}

func main() {
	log.SetFlags(log.Lshortfile)

	_ = godotenv.Load()

	if err := run(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
