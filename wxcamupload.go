package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/dutchcoders/goftp"
	"github.com/oliamb/cutter"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"reflect"
)

func main() {

	// Setup defaults
	camUrl := "http://192.150.23.219/snap.jpeg"
	wundergroundUrl := "webcam.wunderground.com:21"
	wundergroundUser := "username"
	wundergroundPass := "password"
	var wxcamuploadDebug bool = false

	// override default from environment variable if they exist
	if os.Getenv("CAM_URL") != "" {
		camUrl = os.Getenv("CAM_URL")
	}

	if os.Getenv("WUNDERGROUND_URL") != "" {
		wundergroundUrl = os.Getenv("WUNDERGROUND_URL")
	}

	if os.Getenv("WUNDERGROUND_USER") != "" {
		wundergroundUser = os.Getenv("WUNDERGROUND_USER")
	}

	if os.Getenv("WUNDERGROUND_PASS") != "" {
		wundergroundPass = os.Getenv("WUNDERGROUND_PASS")
	}

	if os.Getenv("WXCAMUPLOAD_DEBUG") == "true" {
		wxcamuploadDebug = true
	}

	// override default and env variable with option flag from command line
	// flag.StringVar( pointer for the variable, flag name, default value, description)
	flag.StringVar(&camUrl, "cam", camUrl, "Url for the source camera")
	flag.StringVar(&wundergroundUrl, "wuurl", wundergroundUrl, "Url for weatherunderground ftp")
	flag.StringVar(&wundergroundUser, "user", wundergroundUser, "username for weatherunderground")
	flag.StringVar(&wundergroundPass, "pass", wundergroundPass, "password for weatherunderground")

	// boolean flags
	// name of flag, default, description
	// cant set wxcamuploadDebug directly for some reason? seems like Parse needs to happen first.
	wxcamuploadDebugFlag := flag.Bool("debug", wxcamuploadDebug, "show debug output")

	flag.Parse()

	// read wxcamuploadDebugFlag and then set wxcamuploadDebug
	if *wxcamuploadDebugFlag == true {
		wxcamuploadDebug = true
	}

	// download the image from camera
	snap, err := http.Get(camUrl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("snap type:", reflect.TypeOf(snap))
	fmt.Println("snap.body type:", reflect.TypeOf(snap.Body))

	// attempt conversion
	img, _, err := image.Decode(snap.Body)
	if err != nil {
		log.Fatal("Cannot decode image:", err)
	}
	fmt.Println("img type:", reflect.TypeOf(img))

	// crop to just the upper left corner
	cImg, err := cutter.Crop(img, cutter.Config{
		Height:  800,               // height in pixel or Y ratio(see Ratio Option below)
		Width:   1920,              // width in pixel or X ratio
		Mode:    cutter.TopLeft,    // Accepted Mode: TopLeft, Centered
		Anchor:  image.Point{0, 0}, // Position of the top left point
		Options: 0,                 // Accepted Option: Ratio
	})

	if err != nil {
		log.Fatal("Cannot crop image:", err)
	}
	fmt.Println("cImg dimension:", cImg.Bounds())

	// convert from image.Image to []byte
	buf := &bytes.Buffer{}
	if err := jpeg.Encode(buf, cImg, nil); err != nil {
		log.Fatalf("Error converting: %s\n", err)
	}

	// write debug output
	if wxcamuploadDebug {
		fmt.Println("Source Url: ", camUrl)
		fmt.Println("wunderground Url: ", wundergroundUrl)
		fmt.Println("wunderground user: ", wundergroundUser)
		fmt.Println("wunderground pass: ", wundergroundPass)

	}

	// ftp upload
	var ftp *goftp.FTP

	// For debug messages: goftp.ConnectDbg("ftp.server.com:21")
	if ftp, err = goftp.ConnectDbg(wundergroundUrl); err != nil {
		panic(err)
	}

	defer ftp.Close()

	// Username / password authentication
	if err = ftp.Login(wundergroundUser, wundergroundPass); err != nil {
		panic(err)
	}

	// Upload a file
	if err := ftp.Stor("/image.jpg", buf); err != nil {
		panic(err)
	}
}
