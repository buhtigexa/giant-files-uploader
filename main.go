package main

import "bugtigexa.giantfilesuploader.com/cmd"

func main() {
	server := cmd.NewStreamServer("127.0.0.1:8080")
	server.Start()
}
