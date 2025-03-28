package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"net"
	"os"
)

const (
	_  = iota
	Kb = 1 << (10 * iota)
	Mb
	Gb
	NETWORK = "tcp"
)

type StreamClient struct {
	addr string
}

func NewStreamClient(address string) *StreamClient {
	return &StreamClient{address}
}

func (s *StreamClient) Stream(fname string) {
	conn, err := net.Dial(NETWORK, s.addr)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(fname, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	reader := bufio.NewReaderSize(f, 10*Kb)
	buff := make([]byte, 2*Kb)
	defer conn.Close()
	for {
		n, err := reader.Read(buff)
		if err != nil {
			break
		}
		if n >= 1 {
			s.send(buff[:n], conn)

		}
	}
}

func (s *StreamClient) send(data []byte, conn net.Conn) (int, error) {
	var size int64
	size = int64(len(data))
	buffer := new(bytes.Buffer)
	if err := binary.Write(buffer, binary.BigEndian, size); err != nil {
		panic(err)
	}

	buffer.Write(data)
	n, err := conn.Write(buffer.Bytes())
	if err != nil {
		panic(err)
	}
	return n, err
}

func (s *StreamClient) Receive(conn net.Conn) ([]byte, error) {
	var size int64
	if err := binary.Read(conn, binary.BigEndian, &size); err != nil {
		panic(err)
	}

	buf := make([]byte, size)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
