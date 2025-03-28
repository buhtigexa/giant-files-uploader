package cmd

import (
	"bugtigexa.giantfilesuploader.com/model"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const TO = 10 * time.Second

type StreamServer struct {
	sm *model.FileManager
}

func NewStreamServer() *StreamServer {
	return &StreamServer{
		sm: model.NewFileManager(),
	}
}

func (s *StreamServer) processStream(conn net.Conn) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	defer func() {
		log.Printf("Closing connection to %s", conn.RemoteAddr())
		conn.Close()
	}()
	for {
		var size int64
		conn.SetReadDeadline(time.Now().Add(TO))
		if err := binary.Read(conn, binary.BigEndian, &size); err != nil {
			if os.IsTimeout(err) {
				return err
			}
			if err == io.EOF {
				return nil
			}
			return err
		}

		if size == 0 {
			return nil
		}

		buff := make([]byte, size)
		n, err := io.ReadFull(conn, buff)
		if err != nil {
			if os.IsTimeout(err) {
				break
			}
			if err == io.EOF {
				return nil
			}
			return err
		}

		// if I'm here is because we could read from conn .
		if n == 0 {
			return nil
		}
		if _, err := s.sm.Store(buff); err != nil {
			return err
		}
	}
	return nil
}

func (s *StreamServer) Start() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf(" Error while executing server: %v", err)
			return
		}

	}()

	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error while accepting client: %v", err)
		}
		s.processStream(conn)
	}

}
