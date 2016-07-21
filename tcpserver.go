package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8888"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	tcpAddr, _ := net.ResolveTCPAddr(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	l, err := net.ListenTCP(CONN_TYPE, tcpAddr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		conn.SetKeepAlive(true)
		conn.SetKeepAlivePeriod(300 * time.Second)
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

const http200PUT = "HTTP/1.0 200 OK\r\n" +
	"Connection: Keep-Alive\r\n" +
	"Keep-Alive: max=5, timeout=300\r\n" +
	"Content-Length: 10\r\n\r\n" +
	`{"Code":0}`
const http200GET = "HTTP/1.0 200 OK\r\n" +
	"Connection: Keep-Alive\r\n" +
	"Keep-Alive: max=5, timeout=300\r\n" +
	"Content-Length: 22\r\n\r\n" +
	`{"Code":0,"Result":[]}`

// Handles incoming requests.
func handleRequest(conn *net.TCPConn) {
	// Read the incoming connection into the buffer.
	buf := make([]byte, 200)
	reader := bufio.NewReader(conn)
	// writer := bufio.NewWriter(conn)
	// writer.Write([]byte(http200))
	// fmt.Println("len:", len(http200))
	// writer.Flush()
	for {
		// will listen for message to process ending in newline (\n)
		_, err := reader.Read(buf)
		if err != nil {
			// fmt.Println("conn close")
			conn.Close()
			break
		}
		http := string(buf)
		// fmt.Println("message:", http)
		if http[0:3] == "GET" {
			conn.Write([]byte(http200GET))
		} else if http[0:3] == "PUT" {
			conn.Write([]byte(http200PUT))
		}
		// send new string back to client
		// writer.Write([]byte(http200))
		// fmt.Println("conn:", conn)

		// conn.CloseWrite()
	}
}
