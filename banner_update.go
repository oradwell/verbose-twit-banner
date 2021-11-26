package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"

	"github.com/dghubble/oauth1"
)

func getTwitterOauth1Client(consumerKey string, consumerSecret string, accessToken string, accessSecret string) *http.Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	return config.Client(oauth1.NoContext, token)
}

func writeBannerForm(body *bytes.Buffer, drawable *image.RGBA) (*multipart.Writer, error) {
	multipartWriter := multipart.NewWriter(body)

	fwriter, err := multipartWriter.CreateFormField("banner")
	if err != nil {
		return multipartWriter, err
	}

	err = png.Encode(fwriter, drawable)
	if err != nil {
		return multipartWriter, err
	}

	multipartWriter.Close()

	return multipartWriter, nil
}

func doTwitterUploadRequest(client *http.Client, multipartWriter *multipart.Writer, body *bytes.Buffer, debug bool) error {
	const apiUrl = "https://api.twitter.com/1.1/account/update_profile_banner.json"

	req, err := http.NewRequest("POST", apiUrl, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 201 {
		return errors.New(fmt.Sprintf("Failed upload request. Status: %s", res.Status))
	}

	if debug {
		fmt.Printf("Upload request returned: %s\n", res.Status)
	}

	return nil
}

func UpdateTwitterBanner(consumerKey string, consumerSecret string, accessToken string, accessSecret string, drawable *image.RGBA, debug bool) error {
	var body bytes.Buffer

	client := getTwitterOauth1Client(consumerKey, consumerSecret, accessToken, accessSecret)

	multipartWriter, err := writeBannerForm(&body, drawable)
	if err != nil {
		return err
	}

	err = doTwitterUploadRequest(client, multipartWriter, &body, debug)
	if err != nil {
		return err
	}

	return nil
}
