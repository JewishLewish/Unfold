package main

import (
	"fmt"
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
	var locX, locY, R, G, B int
	draw.Draw(cimg, img.Bounds(), img, image.Point{}, draw.Over)
	for true {
		fmt.Print("Type the following: locX, locY, R, G, B:")
		fmt.Scan(&locX, &locY, &R, &G, &B)
		if R > 255 || G > 255 || B > 255 {
			fmt.Print("ERROR! RGB max int goes up to 255.")
			continue
		}
		cimg.Set(locX, locY, color.RGBA{uint8(R), uint8(G), uint8(B), 255})
		update(cimg)
	}
	update(cimg)
	//close
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
	png.Encode(outFile, cimg)
	outFile.Close()
}
