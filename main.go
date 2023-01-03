package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
)

func unfold(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func web() {
	http.HandleFunc("/", unfold)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	go web()
	fmt.Print("Website is up.")

	img := canvas() //Fetchs canvas
	cimg := image.NewRGBA(img.Bounds())
	fmt.Print("Image has been founded! \n")

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

	outFile, _ := os.Create("main.png")
	png.Encode(outFile, img)
	outFile.Close()

	return img
}

func update(cimg *image.RGBA) {
	e := os.Remove("main.png")
	if e != nil {
		log.Fatal(e)
	}

	outFile, _ := os.Create("main.png")
	png.Encode(outFile, cimg)
	outFile.Close()
}
