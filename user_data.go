package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type userLookup struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type publicMetrics struct {
	FollowerCount  int `json:"followers_count"`
	FollowingCount int `json:"following_count"`
	TweetCount     int `json:"tweet_count"`
	ListedCount    int `json:"listed_count"`
}

type user struct {
	PublicMetrics publicMetrics `json:"public_metrics"`
	Location      string
}

func getTwitterOauth2Client(consumerKey string, consumerSecret string) *http.Client {
	config := &clientcredentials.Config{
		ClientID:     consumerKey,
		ClientSecret: consumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	return config.Client(oauth2.NoContext)
}

func fetchIdForUsername(client *http.Client, username string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/by/username/"+username, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	dec := json.NewDecoder(resp.Body)
	var data = make(map[string]userLookup)
	err = dec.Decode(&data)
	if err != nil {
		return "", err
	}

	return data["data"].Id, nil
}

func fetchUserData(client *http.Client, userId string) (map[string]user, error) {
	req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/"+userId+"?user.fields=public_metrics,location", nil)
	if err != nil {
		return make(map[string]user), err
	}

	resp, err := client.Do(req)
	if err != nil {
		return make(map[string]user), err
	}

	dec := json.NewDecoder(resp.Body)
	userData := make(map[string]user)
	err = dec.Decode(&userData)
	if err != nil {
		return make(map[string]user), err
	}

	return userData, nil
}

func GetTextLines(metrics map[string]string, promotionalLine string) []string {
	lines := []string{
		fmt.Sprintf("%s Following", metrics["following_count"]),
		fmt.Sprintf("%s Followers", metrics["followers_count"]),
		fmt.Sprintf("%s Tweets", metrics["tweet_count"]),
		fmt.Sprintf("In %s lists", metrics["listed_count"]),
	}

	if metrics["location"] != "" {
		lines = append(lines, fmt.Sprintf("Located in %s", metrics["location"]))
	}

	lines = append(lines, fmt.Sprintf("%s", time.Now().UTC().Format("2006-01-02 15:04 MST")))
	if promotionalLine != "" {
		lines = append(lines, promotionalLine)
	}

	return lines
}

func GetTwitterUserData(consumerKey string, consumerSecret string, username string) (map[string]string, error) {
	client := getTwitterOauth2Client(consumerKey, consumerSecret)

	userId, err := fetchIdForUsername(client, username)
	if err != nil {
		return make(map[string]string), err
	}

	userData, err := fetchUserData(client, userId)
	if err != nil {
		return make(map[string]string), err
	}

	return map[string]string{
		"followers_count": fmt.Sprintf("%d", userData["data"].PublicMetrics.FollowerCount),
		"following_count": fmt.Sprintf("%d", userData["data"].PublicMetrics.FollowingCount),
		"tweet_count":     fmt.Sprintf("%d", userData["data"].PublicMetrics.TweetCount),
		"listed_count":    fmt.Sprintf("%d", userData["data"].PublicMetrics.ListedCount),
		"location":        userData["data"].Location,
	}, nil
}
