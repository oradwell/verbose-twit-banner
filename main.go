package main

import jpeg "image/jpeg"
import "image"

import "os"
import "bufio"

func main() {
	file, err := os.Open("images/skate-1500x500.jpg")

	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(file)

	decodedImage, _, err := image.Decode(reader)

	if err != nil {
		panic(err)
	}

	outFile, err := os.Create("images/out.jpg")

	if err != nil {
		panic(err)
	}

	jpeg.Encode(outFile, decodedImage, &jpeg.Options{Quality: 90})
}
