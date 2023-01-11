# Unfold
![Unfold](https://user-images.githubusercontent.com/65754609/211249822-1a11e1cc-da8a-4566-9220-d299a98578c3.png)

## Sites
``/canvas`` -> Directs you to the canvas of the project

``/pixel?x=0&y=0`` -> Gives JSON information about certain canvas locations

``/`` -> Main Page. You can setup your own custom page via settings.json. 

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

``frames (false/true)`` -> When true, it will start recording "frames" of the canvas.

``update_duration_seconds (0-...)`` -> Frames must be true. Duration between each frame saves.

``framerate (1-...)`` -> Framerate for Timelapse 

``port (int)`` -> Port for your server

``address (string)`` -> Address for your server (0.0.0.0 for default)

``ratelimit (int)`` -> Ratelimits pixels per second. It's recommended below 20.

``sitefiles (string, relative path)`` -> This will open the index.html file and use it as it's main site for ``/``

``sitefiles (bool)`` -> 2nd canvas that determines which areas can a user place. 

## Plans
The following features are still in development:

Setup Tools -> Allows users not to manually have to install every file needed

Built-In Terminal -> Allows users to type in certain commands

Admin Panel(?) -> Allows admins / managers to see audit logs of servers

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
