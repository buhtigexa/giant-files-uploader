package main

import (
	"testing"
)

func TestServerPing(t *testing.T) {
	server := NewStreamClient()
	server.Stream("../random.bin")

}
