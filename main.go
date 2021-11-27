package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const fontPath = "fonts/OpenSans-VariableFont_wdth,wght.ttf"
const imageDir = "images"
const fontDpi = 72.0
const fontSize = 32.0

// Overlay background RGBA values
const overlayColourRed = 150
const overlayColourGreen = 150
const overlayColourBlue = 150
const overlayOpacity = 150

// Overlay rectangle coordinates
const overlayX0 = 1082
const overlayY0 = 22
const overlayX1 = 1482

// Amount of pixels to allocate
// to each line when calculating Y1
// Should take the font size
// and padding into consideration
const lineHeight = 46

// Gap between lines
const textPadding = 10

func getBannerDrawable(consumerKey string, consumerSecret string, username string, promotionalLine string, debug bool) (*image.RGBA, error) {
	overlayColour := color.RGBA{
		overlayColourRed,
		overlayColourBlue,
		overlayColourGreen,
		uint8(color.Opaque.A),
	}

	srcPath, err := GetJpegPathInDirectory(imageDir)
	if debug {
		fmt.Printf("Using base image: %s\n", srcPath)
	}
	if err != nil {
		return image.NewRGBA(image.Rectangle{}), err
	}

	drawable := GetDrawableFromImagePath(srcPath)

	userData, err := GetTwitterUserData(consumerKey, consumerSecret, username)
	if err != nil {
		return image.NewRGBA(image.Rectangle{}), err
	}

	lines := GetTextLines(userData, promotionalLine)
	numLines := len(lines)
	overlayY1 := numLines * lineHeight

	overlayRectangle := image.Rect(overlayX0, overlayY0, overlayX1, overlayY1)

	AddOverlayOnDrawable(drawable, overlayRectangle, &overlayColour, &color.Alpha{overlayOpacity})

	font := LoadFontFromPath(fontPath)

	ftContext := GetFreetypeContext(font, fontDpi, fontSize, drawable)

	if debug {
		fmt.Printf("Lines to print: %v\n", lines)
	}

	WriteLinesOnRectangle(overlayRectangle, ftContext, lines, int(fontSize), textPadding)

	return drawable, nil
}

func main() {
	const outPath = "out.png"

	var loopDuration time.Duration
	var err error

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
	promotionalLine := flag.String("promotional-line", os.Getenv("PROMOTIONAL_LINE"), "Extra line to add to the bottom of the information overlay")
	loopInterval := flag.String("interval", "", "Banner update interval (e.g. 5m). With empty value, program exits after updating once")
	debug := flag.Bool("debug", os.Getenv("DEBUG") != "", "Print more output")
	dryRun := flag.Bool("dry-run", false, "Write image to out.png instead of updating Twitter banner")

	flag.Parse()

	if *consumerKey == "" || *consumerSecret == "" {
		panic("Twitter consumer key and consumer secret are required. Use '-h' for details")
	}

	if *loopInterval != "" {
		loopDuration, err = time.ParseDuration(*loopInterval)
		if err != nil {
			panic(fmt.Sprintf("Error parsing interval: %v\n", err))
		}

		if *dryRun {
			panic("Dry run doesn't make sense when running unattended")
		}
	}

	for ok := true; ok; ok = (*loopInterval != "") {
		drawable, err := getBannerDrawable(*consumerKey, *consumerSecret, *username, *promotionalLine, *debug)
		if err != nil {
			fmt.Printf("Error: %v\n", err)

			continue
		}

		if *dryRun {
			WriteToPngFile(outPath, drawable)

			fmt.Printf("Please see file: %s\n", outPath)

			continue
		} else {
			err = UpdateTwitterBanner(*consumerKey, *consumerSecret, *accessToken, *accessSecret, drawable, *debug)
			if err != nil {
				fmt.Printf("Error: %v\n", err)

				continue
			}
			fmt.Println("Updated Twitter banner")
		}

		if *loopInterval != "" {
			if *debug {
				fmt.Printf("Waiting for %s\n", *loopInterval)
			}

			time.Sleep(loopDuration)
		}
	}
}
