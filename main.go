package main

import _ "image/jpeg"
import "image/png"
import "image/draw"
import "image/color"
import "image"

import "os"

// import "bufio"

func main() {
	// file, err := os.Open("images/skate-1500x500.jpg")

	// if err != nil {
	// 	panic(err)
	// }

	// reader := bufio.NewReader(file)

	// decodedImage, _, err := image.Decode(reader)

	// myRectangle := image.Rect(10, 20, 30, 40)

	m := image.NewRGBA(image.Rect(0, 0, 640, 480))
	grey := color.RGBA{150, 150, 150, 255}
	draw.Draw(m, m.Bounds(), &image.Uniform{grey}, image.ZP, draw.Src)

	// if err != nil {
	// 	panic(err)
	// }

	outFile, err := os.Create("images/out.png")

	if err != nil {
		panic(err)
	}

	png.Encode(outFile, m)
}
