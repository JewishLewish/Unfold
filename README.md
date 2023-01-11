# Unfold
![Unfold](https://user-images.githubusercontent.com/65754609/211249822-1a11e1cc-da8a-4566-9220-d299a98578c3.png)

## Websocket Subdirectories
![Frame 2](https://user-images.githubusercontent.com/65754609/211909555-5fa4b496-7ce1-4297-851b-9deb02216e4f.png)


## HTTP Methods
![Frame 3 (1)](https://user-images.githubusercontent.com/65754609/211909724-7252ad03-347e-429e-9a5b-972da513b19f.png)

## Settings.json
```json
{
  "frames": false,
  "update_duration_seconds":2,
  "framerate":60,
  "port":25567,
  "address":"0.0.0.0",
  "ratelimit":180,
  "sitefiles": "static/index.html",
  "useplacemap": false
}
```

## Setup

### Requirements

* Download Go Lang (1.19 preferred)

* Git Clone ffmpeg (optional for timelapse) 

### Instructions

1) Download Go Lang (1.19 preferred) 

2) Git clone this github repo 

3) From there use Git Bash and direct to your directory.
  
    -> use `` go build `` to build the GO code to your device's OS and CPU Archectecture 

    -> for different operating systems / CPU archiecture then look at Go's Cross Compiling System
  
4) Place the code into your server. Run it once and the settings.json file will be created.

5) You just need one png image (default). Put down a `` canvas.png `` and the code will use it as the main canvas.

  **SIDENOTE**, if you are using a "placing map" then make a Black-and-White Pixel Only Canvas call it "placeable.png." 
  
      White Pixel = User cannot place. Black Pixel = User Can Place

ex. (left is "canvas.png" and right is "placeable.png")

![canvas](https://user-images.githubusercontent.com/65754609/211696350-cb089955-7aeb-4db8-b2b0-09992349309d.png)
![placeable](https://user-images.githubusercontent.com/65754609/211696355-de09a2c9-9918-48a0-89c1-acb663f90180.png)

 6. Have fun.
 
 ### Timelapse System
 

https://user-images.githubusercontent.com/65754609/211737795-b4606719-2aa5-4cce-be7c-b75eb70625c9.mp4

Still in development, use ffmpeg. Command use:

```ffmpeg -framerate 30 -pattern_type glob -i 'timelapse/frame%d.png' -r 120 -vcodec libx264 timelapse.mp4```
