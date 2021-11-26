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
const overlayColourRed = 150
const overlayColourGreen = 150
const overlayColourBlue = 150
const overlayOpacity = 150
const overlayX0 = 1082
const overlayY0 = 22
const overlayX1 = 1482
const overlayY1 = 272
const textPadding = 10

func getBannerDrawable(consumerKey string, consumerSecret string, username string, debug bool) (*image.RGBA, error) {
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

	overlayRectangle := image.Rect(overlayX0, overlayY0, overlayX1, overlayY1)

	AddOverlayOnDrawable(drawable, overlayRectangle, &overlayColour, &color.Alpha{overlayOpacity})

	font := LoadFontFromPath(fontPath)

	ftContext := GetFreetypeContext(font, fontDpi, fontSize, drawable)
	userData, err := GetTwitterUserData(consumerKey, consumerSecret, username)
	if err != nil {
		return image.NewRGBA(image.Rectangle{}), err
	}

	lines := GetTextLines(userData)

	if debug {
		fmt.Printf("Lines to print: %v\n", lines)
	}

	WriteLinesOnRectangle(overlayRectangle, ftContext, lines, int(fontSize), textPadding)

	return drawable, nil
}

func main() {
	const outPath = "out.png"

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
	loopInterval := flag.Int("interval", 0, "Banner update interval in minutes. With value 0, program exits after updating once")
	debug := flag.Bool("debug", os.Getenv("DEBUG") != "", "Print more output")
	dryRun := flag.Bool("dry-run", false, "Write image to out.png instead of updating Twitter banner")

	flag.Parse()

	if *consumerKey == "" || *consumerSecret == "" {
		panic("Twitter consumer key and consumer secret are required. Use '-h' for details")
	}

	if *loopInterval > 0 && *dryRun {
		panic("Dry run doesn't make sense when running unattended")
	}

	for ok := true; ok; ok = (*loopInterval > 0) {
		drawable, err := getBannerDrawable(*consumerKey, *consumerSecret, *username, *debug)
		if err != nil {
			fmt.Printf("Error: %v\n", err)

			continue
		}

		if *dryRun {
			WriteToPngFile(outPath, drawable)

			fmt.Printf("Please see file: %s\n", outPath)

			continue
		} else {
			err = UpdateTwitterBanner(*consumerKey, *consumerSecret, *accessToken, *accessSecret, drawable)
			if err != nil {
				fmt.Printf("Error: %v\n", err)

				continue
			}
			fmt.Println("Updated Twitter banner")
		}

		if *loopInterval > 0 {
			if *debug {
				fmt.Printf("Waiting for %d minute(s)\n", *loopInterval)
			}

			time.Sleep(time.Duration(*loopInterval) * time.Minute)
		}
	}
}
