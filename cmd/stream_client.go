package cmd

import (
	"bufio"
	"bugtigexa.giantfilesuploader.com/model"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"net"
	"os"
	"time"
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

func (s *StreamClient) Stream(fname string) error {
	conn, err := net.Dial(NETWORK, s.addr)
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(fname, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	defer conn.Close()

	reader := bufio.NewReaderSize(f, 10*Kb)
	buff := make([]byte, 2*Kb)

	i := 0
	for {
		n, err := reader.Read(buff)
		if err != nil {
			break
		}
		if n >= 1 {
			data := model.Data{
				Id:    fname,
				Part:  i,
				Value: buff[:n],
				Time:  time.Now(),
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				return err
			}
			if _, err = s.send(jsonData, conn); err != nil {
				return err
			}
			i++
		}
	}
	return nil
}

func (s *StreamClient) send(data []byte, conn net.Conn) (int, error) {
	var size int64
	size = int64(len(data))
	buffer := new(bytes.Buffer)
	if err := binary.Write(buffer, binary.BigEndian, size); err != nil {
		return 0, err
	}

	buffer.Write(data)
	n, err := conn.Write(buffer.Bytes())
	if err != nil {
		return n, err
	}
	return n, err
}
