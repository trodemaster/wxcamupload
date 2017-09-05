package main

import (
	"os"

	"github.com/dutchcoders/goftp"
)

func main() {
	var err error
	var ftp *goftp.FTP

	// For debug messages: goftp.ConnectDbg("ftp.server.com:21")
	if ftp, err = goftp.ConnectDbg("webcam.wunderground.com:21"); err != nil {
		panic(err)
	}

	defer ftp.Close()
	//fmt.Println("Successfully connected to", server)

	// Username / password authentication
	if err = ftp.Login("username", "pass"); err != nil {
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
