package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
)

func getDrawableFromImagePath(imagePath string) *image.RGBA {
	file, err := os.Open(imagePath)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	decodedImage, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	decodedImageBounds := decodedImage.Bounds()
	drawable := image.NewRGBA(image.Rect(0, 0, decodedImageBounds.Dx(), decodedImageBounds.Dy()))
	draw.Draw(drawable, drawable.Bounds(), decodedImage, decodedImageBounds.Min, draw.Src)

	return drawable
}

func getFontFromFontPath(fontPath string) *truetype.Font {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		panic(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	return font
}

func addOverlay(drawable *image.RGBA, rectangle image.Rectangle, colour *color.RGBA, opacity *color.Alpha) {
	draw.DrawMask(drawable, rectangle, &image.Uniform{colour}, image.ZP, &image.Uniform{opacity}, image.ZP, draw.Over)
}

func getFreetypeContext(font *truetype.Font, dpi float64, fontSize float64, drawable *image.RGBA) *freetype.Context {
	context := freetype.NewContext()
	context.SetDPI(dpi)
	context.SetFont(font)
	context.SetFontSize(fontSize)
	context.SetClip(drawable.Bounds())
	context.SetDst(drawable)
	context.SetSrc(image.Black)

	return context
}

func addLines(rectangle image.Rectangle, context *freetype.Context, lines []string, fontSize int, padding int) {
	pointX := rectangle.Min.X + padding
	pointY := rectangle.Min.Y + fontSize

	for _, text := range lines {
		labelPoint := freetype.Pt(pointX, pointY)
		_, err := context.DrawString(text, labelPoint)
		if err != nil {
			panic(err)
		}

		pointY += fontSize + padding
	}
}

func writeToPng(filename string, drawable *image.RGBA) {
	outFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	png.Encode(outFile, drawable)
}

func getTwitterApiClient(bearerToken string) *http.Client {
	return &http.Client{}
}

func getTwitterUserData(bearerToken string, username string) map[string]string {
	client := getTwitterApiClient(bearerToken)

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

	req, err := http.NewRequest("GET", "https://api.twitter.com/2/users/by/username/"+username, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	dec := json.NewDecoder(resp.Body)
	var dat = make(map[string]userLookup)
	err = dec.Decode(&dat)
	if err != nil {
		panic(err)
	}
	fmt.Println(dat["data"].Id)

	req, err = http.NewRequest("GET", "https://api.twitter.com/2/users/"+dat["data"].Id+"?user.fields=public_metrics,location", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	resp, err = client.Do(req)
	if err != nil {
		panic(err)
	}

	dec = json.NewDecoder(resp.Body)
	userDat := make(map[string]user)
	err = dec.Decode(&userDat)
	if err != nil {
		panic(err)
	}

	return map[string]string{
		"followers_count": fmt.Sprintf("%d", userDat["data"].PublicMetrics.FollowerCount),
		"following_count": fmt.Sprintf("%d", userDat["data"].PublicMetrics.FollowingCount),
		"tweet_count":     fmt.Sprintf("%d", userDat["data"].PublicMetrics.TweetCount),
		"listed_count":    fmt.Sprintf("%d", userDat["data"].PublicMetrics.ListedCount),
		"location":        userDat["data"].Location,
	}
}

func getLines(metrics map[string]string) []string {
	return []string{
		fmt.Sprintf("%s Following", metrics["following_count"]),
		fmt.Sprintf("%s Followers", metrics["followers_count"]),
		fmt.Sprintf("%s Tweets", metrics["tweet_count"]),
		fmt.Sprintf("In %s lists", metrics["listed_count"]),
		fmt.Sprintf("Located in %s", metrics["location"]),
	}
}

func main() {
	const fontPath = "fonts/OpenSans-VariableFont_wdth,wght.ttf"
	const outPath = "images/out.png"
	const srcPath = "images/skate-1500x500.jpg"
	const fontDpi = 72.0
	const fontSize = 32.0
	const overlayColourRed = 150
	const overlayColourGreen = 150
	const overlayColourBlue = 150
	const overlayOpacity = 150
	const overlayX0 = 1082
	const overlayY0 = 22
	const overlayX1 = 1482
	const overlayY1 = 272
	const textPadding = 10

	bearerToken := flag.String("bearer", "", "Twitter Bearer Token")
	username := flag.String("username", "oliverradwell", "Twitter username")

	flag.Parse()

	lines := getLines(getTwitterUserData(*bearerToken, *username))

	overlayColour := color.RGBA{overlayColourRed, overlayColourBlue, overlayColourGreen, 255}

	drawable := getDrawableFromImagePath(srcPath)

	overlayRectangle := image.Rect(overlayX0, overlayY0, overlayX1, overlayY1)

	addOverlay(drawable, overlayRectangle, &overlayColour, &color.Alpha{overlayOpacity})

	font := getFontFromFontPath(fontPath)

	ftContext := getFreetypeContext(font, fontDpi, fontSize, drawable)

	addLines(overlayRectangle, ftContext, lines, int(fontSize), textPadding)

	writeToPng(outPath, drawable)
}
