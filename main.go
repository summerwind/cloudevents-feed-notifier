package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var (
	VERSION = "0.0.1"
	COMMIT  = "HEAD"
)

const (
	UserAgent = "cloudevents-feed-notifier"
)

var lastFetchedTime map[string]time.Time

func init() {
	lastFetchedTime = map[string]time.Time{}
}

func loadConfig(configPath string) (*Config, error) {
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	c := NewConfig()
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func fetch(feedURL string) ([]byte, error) {
	res, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return ioutil.ReadAll(res.Body)
}

func parse(data []byte) (*Feed, error) {
	var feed Feed

	err := xml.Unmarshal(data, &feed)
	if err != nil {
		return nil, err
	}

	return &feed, nil
}

func publish(url string, event *Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	t := time.Now()
	id := uuid.NewV4()

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Add("ce-specversion", "1.0")
	req.Header.Add("ce-type", "feed.publish")
	req.Header.Add("ce-source", event.Entry.URL)
	req.Header.Add("ce-id", id.String())
	req.Header.Add("ce-time", t.Format(time.RFC3339))
	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("User-Agent", fmt.Sprintf("%s/%s", UserAgent, VERSION))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func notify(c *Config) error {
	events := []Event{}

	for _, feedURL := range c.Feeds {
		now := time.Now()
		since, ok := lastFetchedTime[feedURL]
		if !ok {
			since = now
		}

		body, err := fetch(feedURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to fetch feed: %v\n", err)
			continue
		}

		lastFetchedTime[feedURL] = now
		log.Printf("fetched: %s", feedURL)

		f, err := parse(body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to parse feed: %v\n", err)
			continue
		}

		// Workaround for behavior of XML decoder
		// See: https://github.com/golang/go/issues/8535
		if f.Channel != nil && f.Channel.Link == "" {
			f.Channel.Link = feedURL
		}

		ev := f.NewEvents(since)
		events = append(events, ev...)
	}

	for _, event := range events {
		err := publish(c.Endpoint, &event)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to publish event: %v\n", err)
			continue
		}

		log.Printf("notified: %s", event.Entry.URL)
	}

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	v, err := cmd.Flags().GetBool("version")
	if err != nil {
		return err
	}

	if v {
		fmt.Printf("%s (%s)\n", VERSION, COMMIT)
		return nil
	}

	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	c, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	if c.Endpoint == "" {
		return fmt.Errorf("invalid endpoint: %s", c.Endpoint)
	}

	interval, err := time.ParseDuration(c.Interval)
	if err != nil {
		return fmt.Errorf("invalid interval: %v", err)
	}

	go func() {
		err := notify(c)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Notify error: %v\n", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-sigCh:
			return nil
		case <-ticker.C:
			go func() {
				err := notify(c)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Notify error: %v\n", err)
				}
			}()
		}
	}
}

func main() {
	var cmd = &cobra.Command{
		Use:   "cloudevents-feed-notifier",
		Short: "Fetch feeds and notify a entry as an event",
		RunE:  run,

		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.Flags().StringP("config", "c", "config.yml", "Path to the configuration file")
	cmd.Flags().BoolP("version", "v", false, "Display version information and exit")

	err := cmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}
