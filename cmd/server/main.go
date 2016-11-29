package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const (
	// NearestThings is the number of things closest to the tweet to include in the results.
	NearestThings = 3
)

func main() {

	// get the secrets
	// twitter first
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "", "Twitter Access Secret")

	flags.Parse(os.Args[1:])

	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required for twitter")
	}

	// thingful
	thingfulAPIKey := MustFindInEnvironment("THINGFUL_API_KEY")
	thingfulClient := NewClient(thingfulAPIKey)

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)

	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	out := make(chan []item, 5)
	api := &api{
		results: out,
	}

	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {

		if tweet.Coordinates != nil {

			fmt.Println(tweet.Text, " @ ", tweet.Coordinates.Coordinates)

			items := []item{
				item{
					Type: "tweet",
					Location: location{
						Latitude:  tweet.Coordinates.Coordinates[1],
						Longitude: tweet.Coordinates.Coordinates[0],
					},
					Data: tweet,
				},
			}

			searchResults, err := thingfulClient.SearchByLocation(tweet.Coordinates.Coordinates[1], tweet.Coordinates.Coordinates[0], 5000)

			if err != nil {
				log.Fatal(err)
			}

			// access the nearest things
			for i := 0; i < NearestThings; i++ {

				r, err := thingfulClient.Access(searchResults.Data[i].ID)

				if err != nil {
					continue
				}

				items = append(items, item{
					Type: "thing",
					Location: location{
						Latitude:  r.Data[0].Attributes.Location.Latitude,
						Longitude: r.Data[0].Attributes.Location.Longitude,
					},
					Data: r,
				})
			}
			out <- items
		}
	}

	fmt.Println("starting stream...")

	filterParams := &twitter.StreamFilterParams{
		Locations: []string{"-0.489,51.28,0.236,51.686"}, //london
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	go demux.HandleChan(stream.Messages)

	http.Handle("/", http.FileServer(http.Dir("./web")))
	http.HandleFunc("/api/stream/", api.streamer)
	http.ListenAndServe(":3000", nil)
}

type api struct {
	results chan []item
}

func (a *api) streamer(w http.ResponseWriter, r *http.Request) {

	items := <-a.results

	j, err := json.Marshal(items)

	if err != nil {
		log.Fatal(err)

	}
	w.Write(j)
}

// MustFindInEnvironment looks for a value, logging with Panic if not found
func MustFindInEnvironment(envVar string) string {

	v := os.Getenv(envVar)
	if v == "" {
		log.Panicf("$%s environmental variable must be set", envVar)
	}
	return v

}

type item struct {
	Type     string      `json:"type"`
	Location location    `json:"location"`
	Data     interface{} `json:"data"`
}

type location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
