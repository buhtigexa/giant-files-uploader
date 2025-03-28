package test

import (
	"bugtigexa.giantfilesuploader.com/cmd"
	"testing"
)

// Este test muestra como almacenar un archivo usando el server
func TestServerPing(t *testing.T) {
	server := cmd.NewStreamClient("127.0.0.1:8080")
	server.Stream("file_manager_test.go")

}
