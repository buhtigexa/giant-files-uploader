package cmd

import (
	"bugtigexa.giantfilesuploader.com/model"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const TO = 30 * time.Second

type StreamServer struct {
	sm      *model.FileManager
	address string
}

func NewStreamServer(addr string) *StreamServer {
	return &StreamServer{
		sm:      model.NewFileManager(),
		address: addr,
	}
}

func (s *StreamServer) Start() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf(" Error while executing server: %v", err)
			return
		}

	}()

	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Error while accepting client: %v", err)
		}
		go s.processStream(conn)
	}

}

func (s *StreamServer) processStream(conn net.Conn) error {
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
