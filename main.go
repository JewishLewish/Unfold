package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	rate "golang.org/x/time/rate"
)

// Global Variables since it's going to be used a lot.
var origimg = canvas()
var cimg = image.NewRGBA(canvas().Bounds())

// This is used to ratelimit each user individually.
var rateLimits = make(map[string]*rate.Limiter)

type settings struct {
	Frbool bool   `json:"frames"`
	Fps    int    `json:"framerate"`
	Update int    `json:"update_duration_seconds"`
	Port   int    `json:"port"`
	Addr   string `json:"address"`
	Rlim   int    `json:"ratelimit"`
	SFil   string `json:"sitefiles"`
	Pmap   bool   `json:"useplacemap"`
}

func main() {
	_, err := os.Stat("settings.json")
	if os.IsNotExist(err) {
		fmt.Print("settings.json was not founded. Creating file...")
		setup()
	}

	file, _ := os.Open("settings.json")
	defer file.Close()

	var set settings
	if json.NewDecoder(file).Decode(&set) != nil {
		fmt.Println(err)
		return
	}

	img := canvas()
	go web(set.Port, set.Addr, set.Rlim) //Website operates async.
	fmt.Print("Website is being operated!\n")

	draw.Draw(cimg, img.Bounds(), img, image.Point{}, draw.Over)

	if set.Pmap {
		mask()
	}
	fmt.Print("Image has been created! \n")

	if set.Frbool {
		fmt.Print("Frames system is up! \n")
		go frames(set.Update)
	}

	defer turnoff() //When the code ends, we can shut it off.

	var act1 string

	for {
		fmt.Print("$terminal =>")
		fmt.Scan(&act1) //Admin or User

		if act1 == "ban" {
			var ip string
			fmt.Print("IP to ban: ")
			fmt.Scan(&ip)
			rateLimits[ip] = rate.NewLimiter(rate.Limit(0), 0)
			fmt.Print("Banned " + ip)
		} else if act1 == "backup" {
			backup()
		} else if act1 == "rectangle" {
			var x1, y1, x2, y2 int
			fmt.Print("$Declare Location (x, y, x2, y2) => ")
			fmt.Scan(&x1, &y1, &x2, &y2)
			rectangle(x1, y2, x2, y2)
		} else if act1 == "setup_frames" {
			dlffpg()
		} else if act1 == "timelapse" {
			timelapse(set.Fps)
		} else if act1 == "place" {
			var x1, y1 int
			var r, g, b uint8
			fmt.Print("$Declare Location (x, y, r, g, b) => ")
			fmt.Scan(&x1, &y1, &r, &g, &b)
			pixelplace(x1, y1, uint8(r), uint8(g), uint8(b))
		} else {
			continue
		}
	}
}

type info struct {
	R uint8 `json:"R"`
	G uint8 `json:"G"`
	B uint8 `json:"B"`
	T bool  `json:"placeable"`
}

func getpixel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	queryParams := r.URL.Query()
	locX, err := strconv.Atoi(queryParams.Get("x"))
	if err != nil {
		http.Error(w, "Location X cannot be properly parsed.", http.StatusBadRequest)
		return
	}

	locY, err := strconv.Atoi(queryParams.Get("y"))
	if err != nil {
		http.Error(w, "Location Y cannot be properly parsed.", http.StatusBadRequest)
		return
	}

	if erstring := boundcheck(locX, locY); erstring != "n.a" {
		http.Error(w, erstring, http.StatusForbidden)
	}

	re, g, b, a := cimg.At(locX, locY).RGBA()

	if a == 0 {
		json.NewEncoder(w).Encode(info{T: false})
	} else {
		json.NewEncoder(w).Encode(info{R: uint8(re), G: uint8(g), B: uint8(b), T: true})
	}

}

func homepage(w http.ResponseWriter, r *http.Request) {
	file, _ := os.Open("settings.json")
	defer file.Close()

	var set settings
	json.NewDecoder(file).Decode(&set)
	if set.SFil != "none" {
		http.ServeFile(w, r, set.SFil)
	}
}

// Website
type Payload struct {
	UInput []int `json:"data"`
}

func web(port int, addr string, ratelim int) {

	mux := http.NewServeMux()
	mux.HandleFunc("/", homepage)
	mux.HandleFunc("/pixel", getpixel)
	mux.HandleFunc("/canvas", func(w http.ResponseWriter, r *http.Request) {

		clientIP := strings.Split(r.RemoteAddr, ":")[0]
		if rateLimits[clientIP] == nil {
			print(clientIP)
			rateLimits[clientIP] = rate.NewLimiter(rate.Limit(ratelim), ratelim) //Ratelimits ratelim (default: 180) pixels per second per user of request.
		}

		if !rateLimits[clientIP].Allow() {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		if r.Method == http.MethodPost {
			r.Header.Set("Content-Type", "application/json")

			var uin Payload
			if err := json.NewDecoder(r.Body).Decode(&uin); err != nil {
				http.Error(w, "Error decoding JSON payload", http.StatusBadRequest)
				return
			}

			erstring := boundcheck(uin.UInput[0], uin.UInput[1])
			if erstring != "n.a" {
				http.Error(w, erstring, http.StatusForbidden)
			}

			go pixelplace(uin.UInput[0], uin.UInput[1], uint8(uin.UInput[2]), uint8(uin.UInput[3]), uint8(uin.UInput[4])) //LocX LocY R G B
			w.Write([]byte("Pixel successfully placed at: " + fmt.Sprint(uin.UInput[0]) + "," + fmt.Sprint(uin.UInput[1])))

		} else if r.Method == "Jimp" {
			fmt.Print(r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
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

	target := addr + ":" + strconv.Itoa(port)
	fmt.Print(target + "\n")
	log.Fatal(http.ListenAndServe(target, mux))
}

// Pixel Placing Mechanism
func pixelplace(locX int, locY int, R, G, B uint8) {
	cimg.Set(locX, locY, color.RGBA{uint8(R), uint8(G), uint8(B), 255})
}

// Admin Pixel Placing
func rectangle(lX, lY, lX2, lY2 int) {
	fmt.Print("Drawing White Recetangle... \n")
	draw.Draw(cimg, image.Rect(lX, lY, lX2, lY2), &image.Uniform{color.White}, image.Point{lX, lX2}, draw.Over)
	fmt.Print("Rectangle completed! \n")
}

// Canvas Updating - Constantly operating

func frames(delay int) {
	os.Mkdir("timelapse", 0777)
	var i int
	sleep := time.Duration(delay) * time.Second

	for {
		time.Sleep(sleep)
		i++
		file, _ := os.Create(fmt.Sprintf("timelapse/frame%06d.png", i))
		png.Encode(file, cimg)
		file.Close()
	}
}

// Canvas Manipulation and Data Fetching
func canvas() image.Image {
	canvas, _ := os.Open("canvas.png") //canvas = Main folder.
	img, _ := png.Decode(canvas)
	defer canvas.Close()
	return img
}

func sitecanvas() image.Image {
	//Create Image to merge.
	return merge()
}

func backup() {
	fmt.Print("Backing up canvas...\n")

	outFile, _ := os.Create("backup.png")
	png.Encode(outFile, merge())
	outFile.Close()

	fmt.Print("Backup is complete. backup.png is made!\n")
}

func merge() image.Image {
	bounds := origimg.Bounds()
	newImg := image.NewRGBA(bounds)
	draw.Draw(newImg, bounds, origimg, image.Point{0, 0}, draw.Over)
	draw.Draw(newImg, bounds, cimg, image.Point{0, 0}, draw.Over)
	return newImg
}

// setup
func setup() {
	fmt.Print("Creating Default Settings.json...")
	defaultset := settings{
		Frbool: false,
		Fps:    15,
		Update: 60,
		Port:   8080,
		Addr:   "0.0.0.0",
		Rlim:   30,
		SFil:   "none",
		Pmap:   false,
	}
	data, err := json.Marshal(defaultset)
	if err != nil {
		fmt.Println(err)
		return
	}

	outFile, _ := os.Create("settings.json")

	_, err = outFile.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	outFile.Close()
	fmt.Print("Settings.json is made!")
}

//Timelapse

func dlffpg() {
	fmt.Print("Cloning ffmpeg. Give it a moment or two...")
	cmd := exec.Command("git", "clone", "https://git.ffm")
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print("ffmpeg is successfully installed!")
}

func timelapse(frameRate int) {
	fmt.Print("Processing timelapse...")
	cmd := exec.Command("ffmpeg", "-framerate", fmt.Sprintf("%d", frameRate), "-i", "timelapse/frames%06d.png", "-c:v", "libx264", "-r", fmt.Sprintf("%d", frameRate), "-pix_fmt", "yuv420p", "out.mp4")

	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Print("Your video is completed!")
}

// Mask System
func mask() {
	fmt.Print("Applying mask... \n")
	file2, _ := os.Open("placeable.png")

	Pmapc, _, _ := image.Decode(file2)
	file2.Close()

	b := Pmapc.Bounds()
	transparentImg := image.NewRGBA(Pmapc.Bounds())

	//draw.Draw(cimg, img.Bounds(), img, image.Point{}, draw.Over)
	draw.Draw(transparentImg, b, Pmapc, image.Point{}, draw.Over)
	for y := 0; y < b.Max.Y; y++ {
		for x := 0; x < b.Max.X; x++ {
			r, g, b, _ := transparentImg.At(x, y).RGBA()
			if r == 0xffff && g == 0xffff && b == 0xffff {
				cimg.Set(x, y, color.Transparent)
			}
		}
	}
}

// Checking
func boundcheck(lx, ly int) string {
	if lx > cimg.Bounds().Max.X {
		return ("X location is outside of Canvas range.")
	}
	if ly > cimg.Bounds().Max.Y {
		return ("X location is outside of Canvas range.")
	}
	return "n.a"

}

// Server Shutdown -> This would save the file and the last state of the image.
func turnoff() {
	fmt.Print("Shutting down...")
	file, _ := os.Create("canvas.png")
	png.Encode(file, cimg)
	file.Close()
}
