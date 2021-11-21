package main

import (
	"errors"
	"fmt"
	"image"
	"os"
	"testing"
)

func TestWritesOntoDrawable(t *testing.T) {
	const filePath = "images/PXL_20210705_154615696.jpg"
	const fontPath = "fonts/OpenSans-VariableFont_wdth,wght.ttf"

	lines := []string{
		"my line",
	}

	font := LoadFontFromPath(fontPath)

	if "*truetype.Font" != fmt.Sprintf("%T", font) {
		t.Errorf("LoadFontFromPath(%s) returned %T want *truetype.Font", fontPath, font)
	}

	drawable := GetDrawableFromImagePath(filePath)

	if "*image.RGBA" != fmt.Sprintf("%T", drawable) {
		t.Errorf("GetDrawableFromImagePath(%s) returned %T want *image.RGBA", filePath, drawable)
	}

	ftContext := GetFreetypeContext(font, 72.0, 16, drawable)

	if "*freetype.Context" != fmt.Sprintf("%T", ftContext) {
		t.Errorf("GetFreetypeContext() returned %T want *freetype.Context", ftContext)
	}

	rectangle := image.Rect(0, 0, 10, 10)

	WriteLinesOnRectangle(rectangle, ftContext, lines, 72, 0)
}

func TestWritesToPngFile(t *testing.T) {
	const outPath = "/var/tmp/test.png"
	const inPath = "images/PXL_20210705_154615696.jpg"

	drawable := GetDrawableFromImagePath(inPath)

	WriteToPngFile(outPath, drawable)

	if _, err := os.Stat(outPath); errors.Is(err, os.ErrNotExist) {
		t.Errorf("WriteToPngFile(%s, drawable) did not write to %s", outPath, outPath)
	}

	os.Remove(outPath)
}
