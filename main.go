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
	"strconv"
)

var img = canvas()
var cimg = image.NewRGBA(img.Bounds())

func main() {
	go web() //Website operates async.
	fmt.Print("Website is up. \n")

	cimg := image.NewRGBA(img.Bounds())
	fmt.Print("Image has been created! \n")
	draw.Draw(cimg, img.Bounds(), img, image.Point{}, draw.Over)
	var user, action string
	for true {
		fmt.Print("$terminal =>")
		fmt.Scan(&user, &action) //Admin Rectangle x y x2 y2
	}
}

// Website
func web() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form data", http.StatusBadRequest)
				return
			}

			var inputs [5]int
			for i := 0; i < 5; i++ {
				inputStr := r.Form.Get(fmt.Sprintf("input%d", i))
				input, err := strconv.Atoi(inputStr)
				if err != nil {
					http.Error(w, "Error parsing int inputs", http.StatusBadRequest)
					return
				}
				inputs[i] = input
			}

			fmt.Print(inputs)
			pixelplace(inputs[0], inputs[1], uint8(inputs[2]), uint8(inputs[3]), uint8(inputs[4])) //LocX LocY R G B

		} else {
			w.Header().Set("Content-Type", "image/png")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "upgrade")
			w.Header().Set("Upgrade", "websocket")

			//If the user is NOT POST-ing then they will just see the picture of the canvas.
			png.Encode(w, sitecanvas())

			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

		}
	})

	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Pixel Placing Mechanism
func pixelplace(locX int, locY int, R, G, B uint8) {
	draw.Draw(cimg, img.Bounds(), img, image.Point{}, draw.Over)

	if R > 255 || G > 255 || B > 255 {
		fmt.Print("ERROR! RGB max int goes up to 255.")
		return
	}
	cimg.Set(locX, locY, color.RGBA{uint8(R), uint8(G), uint8(B), 255})
	update(cimg)
	return
}

// Canvas Manipulation and Data Fetching
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

func sitecanvas() image.Image {
	canvas, _ := os.Open("main.png") //canvas = Main folder.
	img, _ := png.Decode(canvas)
	canvas.Close()
	return img
}
