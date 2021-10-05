package main

import "bufio"
import "github.com/golang/freetype"
import "github.com/golang/freetype/truetype"
import "image"
import "image/color"
import "image/draw"
import "image/png"
import "io/ioutil"
import "os"
import _ "image/jpeg"

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

	lines := []string{"204 Following", "68 Followers", "405 Tweets"}

	overlayColour := color.RGBA{overlayColourRed, overlayColourBlue, overlayColourGreen, 255}

	drawable := getDrawableFromImagePath(srcPath)

	overlayRectangle := image.Rect(overlayX0, overlayY0, overlayX1, overlayY1)

	addOverlay(drawable, overlayRectangle, &overlayColour, &color.Alpha{overlayOpacity})

	font := getFontFromFontPath(fontPath)

	ftContext := getFreetypeContext(font, fontDpi, fontSize, drawable)

	addLines(overlayRectangle, ftContext, lines, int(fontSize), textPadding)

	writeToPng(outPath, drawable)
}
