package cmd

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

const (
	_  = iota
	Kb = 1 << (10 * iota)
	Mb
	Gb
	NETWORK    = "tcp"
	bufferSize = 500
)

type StreamClient struct {
	addr string
	conn net.Conn
}

func NewStreamClient(address string) *StreamClient {
	conn, err := net.Dial(NETWORK, address)
	if err != nil {
		fmt.Printf("stream client err:%v\n", err)
		panic(err)
	}

	return &StreamClient{address, conn}

}

func (s *StreamClient) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *StreamClient) Stream(fname string) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("stream client err:%v\n", err)
		}
	}()

	f, err := os.OpenFile(fname, os.O_RDONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}

	if stat.Size() == 0 {
		return fmt.Errorf("file is empty")
	}
	reader := bufio.NewReaderSize(f, bufferSize)
	buff := make([]byte, 200)

	i := 1
	for {
		n, err := reader.Read(buff)
		if n >= 1 {
			data := &Data{
				FileName: fname,
				Part:     i,
				Value:    buff[:n],
				Time:     time.Now(),
				Total:    float32(int64(n*i) / stat.Size()),
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				return err
			}
			if _, err = s.send(jsonData, s.conn); err != nil {
				return err
			}
			i++
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
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

	n, err := buffer.Write(data)
	if err != nil {
		return 0, err
	}
	n, err = conn.Write(buffer.Bytes())
	if err != nil {
		return n, err
	}
	return n, err
}
