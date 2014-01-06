package main

import (
	"fmt"
	"github.com/alfonsokim/go-tweet-mining/tweet"
	"github.com/kurrik/twittergo"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func testSearch(client *twittergo.Client) {
	query := url.Values{}
	query.Set("q", "iphone")
	url := fmt.Sprintf("/1.1/search/tweets.json?%v", query.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}
	resp, err := client.SendRequest(req)
	if err != nil {
		fmt.Printf("Could not send request: %v\n", err)
		os.Exit(1)
	}
	results := &twittergo.SearchResults{}
	err = resp.Parse(results)
	if err != nil {
		fmt.Printf("Problem parsing response: %v\n", err)
		os.Exit(1)
	}
	for i, tweet := range results.Statuses() {
		user := tweet.User()
		fmt.Printf("%v.) %v\n", i+1, tweet.Text())
		fmt.Printf("From %v (@%v) ", user.Name(), user.ScreenName())
		fmt.Printf("at %v\n\n", tweet.CreatedAt().Format(time.RFC1123))
	}
	if resp.HasRateLimit() {
		fmt.Printf("Rate limit:           %v\n", resp.RateLimit())
		fmt.Printf("Rate limit remaining: %v\n", resp.RateLimitRemaining())
		fmt.Printf("Rate limit reset:     %v\n", resp.RateLimitReset())
	} else {
		fmt.Printf("Could not parse rate limit from response.\n")
	}
}

func testStreaming(client *twittergo.Client) {
	query := url.Values{}
	query.Set("track", "iphone")
	client.Host = "stream.twitter.com"
	url := fmt.Sprintf("/1.1/statuses/filter.json?%v", query.Encode())
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Printf("Could not parse request: %v\n", err)
		os.Exit(1)
	}
	urlTemp := fmt.Sprintf("https://stream.twitter.com%v", url)
	req, err = http.NewRequest("POST", urlTemp, nil)
	//client.Sign(req)
	client.FetchAppToken()
	h := fmt.Sprintf("Bearer %v", client.AppToken.AccessToken)
	req.Header.Set("Authorization", h)
	fmt.Println(req.Header)
	r, err := client.HttpClient.Do(req)
	fmt.Println(r)
	fmt.Println(err)
	/*
		resp, err := client.SendRequest(req)
		if err != nil {
			fmt.Printf("Could not send request: %v\n", err)
			os.Exit(1)
		}
		results := &twittergo.SearchResults{}
		err = resp.Parse(results)
		if err != nil {
			fmt.Printf("Problem parsing response: %v\n", err)
			os.Exit(1)
		}
		for i, tweet := range results.Statuses() {
			user := tweet.User()
			fmt.Printf("%v.) %v\n", i+1, tweet.Text())
			fmt.Printf("From %v (@%v) ", user.Name(), user.ScreenName())
			fmt.Printf("at %v\n\n", tweet.CreatedAt().Format(time.RFC1123))
		}
	*/
}

func main() {
	client, err := tweet.GetClient()
	if err != nil {
		log.Fatal(err)
	}
	//testSearch(client)
	testStreaming(client)
}
