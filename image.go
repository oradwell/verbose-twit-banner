package main

import (
	"bufio"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

func GetDrawableFromImagePath(imagePath string) *image.RGBA {
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

func LoadFontFromPath(fontPath string) *truetype.Font {
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

func AddOverlayOnDrawable(drawable *image.RGBA, rectangle image.Rectangle, colour *color.RGBA, opacity *color.Alpha) {
	draw.DrawMask(drawable, rectangle, &image.Uniform{colour}, image.ZP, &image.Uniform{opacity}, image.ZP, draw.Over)
}

func GetFreetypeContext(font *truetype.Font, dpi float64, fontSize float64, drawable *image.RGBA) *freetype.Context {
	context := freetype.NewContext()
	context.SetDPI(dpi)
	context.SetFont(font)
	context.SetFontSize(fontSize)
	context.SetClip(drawable.Bounds())
	context.SetDst(drawable)
	context.SetSrc(image.Black)

	return context
}

func WriteLinesOnRectangle(rectangle image.Rectangle, context *freetype.Context, lines []string, fontSize int, padding int) {
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

func WriteToPngFile(filename string, drawable *image.RGBA) {
	outFile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	png.Encode(outFile, drawable)
}
