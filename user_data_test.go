package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
