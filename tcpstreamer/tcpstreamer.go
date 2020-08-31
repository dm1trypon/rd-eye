package tcpstreamer

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/dm1trypon/rd-eye/config"

	logger "github.com/dm1trypon/easy-logger"
)

const (
	// LC - logging category
	LC = "TCPSTREAMER"
	// MAXSIZE - max size of bytes message
	MAXSIZE = 3000
	// CONNTYPE - type of server's protocol
	CONNTYPE = "udp"
)

var clients []net.Conn

// Run - starts a streamer to transfer image data.
func Run() {
	logger.Info(LC, "Starting module TCPStreamer")
	go start()
}

func start() {
	ln, err := net.Listen("tcp", config.Cfg.UDPPath)
	if err != nil {
		logger.Critical(LC, fmt.Sprint("TCP listening error: ", err.Error()))
		os.Exit(-1)
	}

	conn, err := ln.Accept()
	if err != nil {
		logger.Critical(LC, fmt.Sprint("TCP accepting error: ", err.Error()))
		os.Exit(-1)
	}

	clients = append(clients, conn)

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			logger.Error(LC, fmt.Sprint("Reader error: ", err.Error()))
			continue
		}

		logger.Info(LC, fmt.Sprint("RECD: ", string(message)))
	}
}

/*
Send - sends a message over UDP

	- data <[]byte> - data for sending.
*/
func Send(data []byte) {
	data = append(data, []byte("\n")...)

	for {
		// log.Println(len(data))
		if len(data) > MAXSIZE-1 {
			break
		}

		data = append(data, []byte("\x00")...)
	}

	log.Println(len(data))
	for _, conn := range clients {
		logger.Info(LC, fmt.Sprint("Sending to ", conn.RemoteAddr().String()))

		if _, err := conn.Write([]byte(string(data))); err != nil {
			logger.Error(LC, fmt.Sprint("Error sending:", err.Error()))
		}
	}
}
