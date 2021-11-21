package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"os"

	"github.com/joho/godotenv"
)

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
		panic("Twitter consumer key and consumer secret are required. Use '-h' for details")
	}

	userData, err := GetTwitterUserData(*consumerKey, *consumerSecret, *username)
	if err != nil {
		errors.Unwrap(fmt.Errorf("%w", err))
	}

	lines := GetTextLines(userData)

	overlayColour := color.RGBA{
		overlayColourRed,
		overlayColourBlue,
		overlayColourGreen,
		uint8(color.Opaque.A),
	}

	srcPath, err := GetJpegPathInDirectory(imageDir)
	if *debug {
		fmt.Printf("Using image: %s\n", srcPath)
	}
	if err != nil {
		errors.Unwrap(fmt.Errorf("%w", err))
	}

	drawable := GetDrawableFromImagePath(srcPath)

	overlayRectangle := image.Rect(overlayX0, overlayY0, overlayX1, overlayY1)

	AddOverlayOnDrawable(drawable, overlayRectangle, &overlayColour, &color.Alpha{overlayOpacity})

	font := LoadFontFromPath(fontPath)

	ftContext := GetFreetypeContext(font, fontDpi, fontSize, drawable)

	WriteLinesOnRectangle(overlayRectangle, ftContext, lines, int(fontSize), textPadding)

	if *debug {
		WriteToPngFile(outPath, drawable)
		fmt.Printf("Please see file: %s\n", outPath)
	} else {
		UpdateTwitterBanner(*consumerKey, *consumerSecret, *accessToken, *accessSecret, drawable)
	}
}
