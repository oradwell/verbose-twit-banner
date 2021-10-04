package main

import "bufio"
import "github.com/golang/freetype"
import "image"
import "image/color"
import "image/draw"
import "image/png"
import "io/ioutil"
import "os"
import _ "image/jpeg"

func main() {
	const fontFileName = "fonts/OpenSans-VariableFont_wdth,wght.ttf"

	file, err := os.Open("images/skate-1500x500.jpg")
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	decodedImage, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	decodedImageBounds := decodedImage.Bounds()
	decodedRectangle := image.NewRGBA(image.Rect(0, 0, decodedImageBounds.Dx(), decodedImageBounds.Dy()))
	draw.Draw(decodedRectangle, decodedRectangle.Bounds(), decodedImage, decodedImageBounds.Min, draw.Src)

	greyRectangle := image.Rect(1082, 22, 1482, 272)
	greyBg := image.NewRGBA(greyRectangle)
	grey := color.RGBA{150, 150, 150, 255}
	opacity := color.Alpha{150}
	draw.DrawMask(decodedRectangle, greyBg.Bounds(), &image.Uniform{grey}, image.ZP, &image.Uniform{opacity}, image.ZP, draw.Over)

	ftContext := freetype.NewContext()

	fontBytes, err := ioutil.ReadFile(fontFileName)
	if err != nil {
		panic(err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		panic(err)
	}

	fontSize := 32
	ftContext.SetDPI(72)
	ftContext.SetFont(font)
	ftContext.SetFontSize(float64(fontSize))
	ftContext.SetClip(decodedRectangle.Bounds())
	ftContext.SetDst(decodedRectangle)
	ftContext.SetSrc(image.Black)

	labelPoint := freetype.Pt(greyRectangle.Min.X+10, greyRectangle.Min.Y+fontSize)

	_, err = ftContext.DrawString("204 Following", labelPoint)
	if err != nil {
		panic(err)
	}

	labelPoint = freetype.Pt(greyRectangle.Min.X+10, greyRectangle.Min.Y+fontSize*2+10)

	_, err = ftContext.DrawString("68 Followers", labelPoint)
	if err != nil {
		panic(err)
	}

	outFile, err := os.Create("images/out.png")
	if err != nil {
		panic(err)
	}

	png.Encode(outFile, decodedRectangle)
}
