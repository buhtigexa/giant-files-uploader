package main

import "bugtigexa.giantfilesuploader.com/cmd"

func main() {
	server := cmd.NewStreamServer()
	server.Start()

}
