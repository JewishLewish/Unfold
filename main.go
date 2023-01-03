package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func main() {
	img := canvas() //Fetchs canvas

	cimg := image.NewRGBA(img.Bounds())
	draw.Draw(cimg, img.Bounds(), img, image.Point{}, draw.Over)
	cimg.Set(3, 3, color.RGBA{0, 255, 34, 255})

	update(cimg) //Updates canvas
}

func canvas() image.Image {
	canvas, _ := os.Open("canvas.png") //canvas = Main folder.
	img, _ := png.Decode(canvas)
	canvas.Close()
	return img
}

func update(cimg *image.RGBA) {
	e := os.Remove("canvas.png")
	if e != nil {
		log.Fatal(e)
	}
	outFile, _ := os.Create("canvas.png")
	defer outFile.Close()

	png.Encode(outFile, cimg)
}
