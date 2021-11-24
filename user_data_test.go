package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

func TestGetsTwitterClient(t *testing.T) {
	const consumerKey = "myKey"
	const consumerSecret = "mySecret"

	client := getTwitterOauth2Client(consumerKey, consumerSecret)

	if "*http.Client" != fmt.Sprintf("%T", client) {
		t.Errorf("getTwitterOauth2Client(%s, %s) returned %T want *http.Client", consumerKey, consumerSecret, client)
	}

	if "*oauth2.Transport" != fmt.Sprintf("%T", client.Transport) {
		t.Errorf("getTwitterOauth2Client(%s, %s) returned Transport %T want *oauth2.Transport", consumerKey, consumerSecret, client.Transport)
	}
}

func TestFetchIdForUsername(t *testing.T) {
	_, err := fetchIdForUsername(&http.Client{}, "")

	if err == nil {
		t.Errorf("fetchIdForUsername(client, username) succeess want error")
	}

	if _, ok := err.(*json.SyntaxError); !ok {
		t.Errorf("fetchIdForUsername(client, username) returned error %T want *json.SyntaxError", err)
	}
}

func TestFetchesUserData(t *testing.T) {
	const userId = "123"

	client := &http.Client{}

	_, err := fetchUserData(client, userId)

	if err == nil {
		t.Errorf("fetchUserData(client, userId) succeess want error")
	}

	if _, ok := err.(*json.UnmarshalTypeError); !ok {
		t.Errorf("fetchUserData(client, userId) returned error %T want *json.UnmarshalTypeError", err)
	}
}

func TestGetsTextLines(t *testing.T) {
	var metrics map[string]string
	var lines []string
	var expected []string

	metrics = map[string]string{
		"following_count": "5",
		"followers_count": "10",
		"tweet_count":     "101",
		"listed_count":    "2",
		"location":        "London, UK",
	}

	lines = GetTextLines(metrics)
	lines = lines[:len(lines)-1]

	expected = []string{
		"5 Following",
		"10 Followers",
		"101 Tweets",
		"In 2 lists",
		"Located in London, UK",
	}

	if !reflect.DeepEqual(lines, expected) {
		t.Errorf("GetTextLines(%#v) returned %#v want %#v", metrics, lines, expected)
	}

	delete(metrics, "location")

	lines = GetTextLines(metrics)
	lines = lines[:len(lines)-1]

	expected = []string{
		"5 Following",
		"10 Followers",
		"101 Tweets",
		"In 2 lists",
	}

	if !reflect.DeepEqual(lines, expected) {
		t.Errorf("GetTextLines(%#v) returned %#v want %#v", metrics, lines, expected)
	}
}

func TestGetsTwitterUserData(t *testing.T) {
	const consumerKey = "asff"
	const consumerSecret = "afas"
	const username = "bob"

	_, err := GetTwitterUserData(consumerKey, consumerSecret, username)

	if err == nil {
		t.Errorf("GetTwitterUserData(%s, %s, %s) returned success want error", consumerKey, consumerSecret, username)
	}

	if _, ok := err.(*url.Error); !ok {
		t.Errorf("GetTwitterUserData(%s, %s, %s) returned %T as error want *url.Error", consumerKey, consumerSecret, username, err)
	}
}
