package cmd

import (
	"context"
	"database/sql"
	"encoding/binary"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"io"
	"os"
	"strconv"

	"log"
	"net"
	"time"
)

var ErrDbConn error = fmt.Errorf("Unable to connect to database\n")

type StreamServer struct {
	readTimeOut time.Duration
	sm          *FileManager
	host        string
	port        string
	db          *DbConfig
}

func NewStreamServer(host, port string) *StreamServer {
	b, err := strconv.ParseInt(os.Getenv("READ_TIME_OUT_CONNECTION"), 10, 64)
	if err != nil {
		log.Printf(err.Error())
		b = 0
	}

	db, err := connectDb()
	if err != nil {
		log.Printf(err.Error())
		return nil
	}
	return &StreamServer{
		readTimeOut: time.Duration(b),
		sm:          NewFileManager(),
		host:        host,
		port:        port,
		db:          db,
	}
}

type DbConfig struct {
	user   string
	pass   string
	host   string
	port   string
	dbName string
	sqlDB  *sql.DB
}

func connectDb() (*DbConfig, error) {
	dbConn := &DbConfig{}
	if dbConn.user == "" || dbConn.pass == "" || dbConn.port == "" || dbConn.dbName == "" {
		return dbConn, ErrDbConn
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbConn.user, dbConn.pass, dbConn.host, dbConn.port, dbConn.dbName)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	dbConn.sqlDB = db
	return dbConn, nil
}

func (s *StreamServer) Start(ctx context.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf(" Error while executing server: %v", err)
			return
		}

	}()

	listener, err := net.Listen("tcp", s.host+":"+s.port)
	if err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}

	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:

				if err != nil {
					log.Fatalf("Error while accepting client: %v", err)
				}
			}
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
		conn.SetReadDeadline(time.Now().Add(s.readTimeOut * time.Second))
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
