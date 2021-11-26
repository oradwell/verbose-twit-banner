package main

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"testing"
)

func TestGetTwitterOauth1Client(t *testing.T) {
	const consumerKey = "v"
	const consumerSecret = "x"
	const accessToken = "y"
	const accessSecret = "z"

	client := getTwitterOauth1Client(consumerKey, consumerSecret, accessToken, accessSecret)

	if "*http.Client" != fmt.Sprintf("%T", client) {
		t.Errorf("getTwitterOauth1Client(%s, %s, %s, %s) returned type %T want *http.Client", consumerKey, consumerSecret, accessToken, accessSecret, client)
	}
}

func TestWritesBannerForm(t *testing.T) {
	const filePath = "images/PXL_20210705_154615696.jpg"
	var body bytes.Buffer

	drawable := GetDrawableFromImagePath(filePath)

	multipartWriter, _ := writeBannerForm(&body, drawable)

	if "*multipart.Writer" != fmt.Sprintf("%T", multipartWriter) {
		t.Errorf("writeBannerForm(body, drawable) returned %T want *multipart.Writer", multipartWriter)
	}
}

func TestTwitterUploadRequestReturnsRequestError(t *testing.T) {
	const errorMessage = "Failed upload request. Status: 400 Bad Request"

	var body bytes.Buffer

	client := &http.Client{}
	multipartWriter := multipart.NewWriter(&body)

	err := doTwitterUploadRequest(client, multipartWriter, &body)

	if fmt.Sprintf("%v", err) != errorMessage {
		t.Errorf("doTwitterUploadRequest(client, multipartWriter, body) returned %#v want %q", fmt.Sprintf("%v", err), errorMessage)
	}
}
