package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dghubble/oauth1"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func getJpegPathInDirectory(directory string) string {
	var files []string

	filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}

		fileExt := filepath.Ext(path)

		if !info.IsDir() && (fileExt == ".jpg" || fileExt == ".jpeg") {
			files = append(files, path)
		}

		return nil
	})

	rand.Seed(time.Now().Unix())

	return files[rand.Intn(len(files))]
}

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

func updateBanner(consumerKey string, consumerSecret string, accessToken string, accessSecret string, drawable *image.RGBA) {
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
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}

func getTwitterUserData(consumerKey string, consumerSecret string, username string) map[string]string {
	config := &clientcredentials.Config{
		ClientID:     consumerKey,
		ClientSecret: consumerSecret,
		TokenURL:     "https://api.twitter.com/oauth2/token",
	}

	client := config.Client(oauth2.NoContext)

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

	req, err = http.NewRequest("GET", "https://api.twitter.com/2/users/"+dat["data"].Id+"?user.fields=public_metrics,location", nil)
	if err != nil {
		panic(err)
	}
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

	return lines
}

func main() {
	const fontPath = "fonts/OpenSans-VariableFont_wdth,wght.ttf"
	const outPath = "out.png"
	const imageDir = "images"
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

	godotenv.Load()

	defaultUsername := os.Getenv("TWITTER_USERNAME")
	if defaultUsername == "" {
		defaultUsername = "oliverradwell"
	}

	consumerKey := flag.String("consumer-key", os.Getenv("TWITTER_APP_CONSUMER_KEY"), "Twitter App consumer key")
	consumerSecret := flag.String("consumer-secret", os.Getenv("TWITTER_APP_CONSUMER_SECRET"), "Twitter App consumer secret")
	accessToken := flag.String("access-token", os.Getenv("TWITTER_ACCESS_TOKEN"), "Twitter User access token")
	accessSecret := flag.String("access-secret", os.Getenv("TWITTER_ACCESS_SECRET"), "Twitter User access secret")
	username := flag.String("username", defaultUsername, "Twitter username")
	debug := flag.Bool("debug", os.Getenv("DEBUG") != "", "If enabled, image is written to out.png")

	flag.Parse()

	if *consumerKey == "" || *consumerSecret == "" {
		panic("Twitter consumer key and consumer secret are required to be passed in")
	}

	lines := getLines(getTwitterUserData(*consumerKey, *consumerSecret, *username))

	overlayColour := color.RGBA{overlayColourRed, overlayColourBlue, overlayColourGreen, 255}

	srcPath := getJpegPathInDirectory(imageDir)
	if *debug {
		fmt.Printf("Using image: %s\n", srcPath)
	}

	drawable := getDrawableFromImagePath(srcPath)

	overlayRectangle := image.Rect(overlayX0, overlayY0, overlayX1, overlayY1)

	addOverlay(drawable, overlayRectangle, &overlayColour, &color.Alpha{overlayOpacity})

	font := getFontFromFontPath(fontPath)

	ftContext := getFreetypeContext(font, fontDpi, fontSize, drawable)

	addLines(overlayRectangle, ftContext, lines, int(fontSize), textPadding)

	if *debug {
		writeToPng(outPath, drawable)
	}

	updateBanner(*consumerKey, *consumerSecret, *accessToken, *accessSecret, drawable)
}
