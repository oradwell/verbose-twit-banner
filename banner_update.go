package main

import (
	"bytes"
	"github.com/dghubble/oauth1"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
)

func UpdateTwitterBanner(consumerKey string, consumerSecret string, accessToken string, accessSecret string, drawable *image.RGBA) {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	client := config.Client(oauth1.NoContext, token)

	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	fw, err := writer.CreateFormField("banner")
	if err != nil {
		panic(err)
	}

	png.Encode(fw, drawable)

	writer.Close()

	req, err := http.NewRequest("POST", "https://api.twitter.com/1.1/account/update_profile_banner.json", &body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}
