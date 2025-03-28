package test

import (
	"bugtigexa.giantfilesuploader.com/cmd"
	"testing"
)

func TestServerPing(t *testing.T) {
	server := cmd.NewStreamClient("localhost:8080")
	server.Stream("file_manager_test.go")

}
