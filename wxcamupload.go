package main

import (
	"flag"
	"fmt"
	"github.com/dutchcoders/goftp"
	"io"
	"log"
	"net/http"
	"os"
	//"strings"
	//"reflect"
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

	// create a tempfile
	tempfile, err := os.Create("image.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer tempfile.Close()

	// copy contents to file
	filebuffer, err := io.Copy(tempfile, snap.Body)
	if err != nil {
		log.Fatal(err)
	}

	// write debug output
	if wxcamuploadDebug {

		fmt.Println("File size: ", filebuffer)
		fmt.Println("Source Url: ", camUrl)
		fmt.Println("wunderground Url: ", wundergroundUrl)
		fmt.Println("wunderground user: ", wundergroundUser)
		fmt.Println("wunderground pass: ", wundergroundPass)

		// ftp upload
		var ftp *goftp.FTP

		// For debug messages: goftp.ConnectDbg("ftp.server.com:21")
		if ftp, err = goftp.ConnectDbg(wundergroundUrl); err != nil {
			panic(err)
		}

		defer ftp.Close()
		//fmt.Println("Successfully connected to", server)

		// Username / password authentication
		if err = ftp.Login(wundergroundUser, wundergroundPass); err != nil {
			panic(err)
		}

		// Upload a file
		var file *os.File
		if file, err = os.Open("/files/go/src/github.com/trodemaster/wxcamupload/image.jpg"); err != nil {
			panic(err)
		}

		if err := ftp.Stor("/image.jpg", file); err != nil {
			panic(err)
		}

	}

}
