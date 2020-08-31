package streamer

import (
	"fmt"
	"net"
	"os"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/rd-eye/config"
)

const (
	// LC - logging category
	LC = "STREAMER"
	// MAXSIZE - max size of bytes message
	MAXSIZE = 2500
	// CONNTYPE - type of server's protocol
	CONNTYPE = "udp"
)

/*
UDPClient - contains data of the connected client via UDP.
	- UDPClient - contains:

		- UDPConn <*net.UDPConn> - data of UDP connection.

		- UDPAddr <*net.UDPAddr> - data of UDP address connection client.
*/
type UDPClient struct {
	UDPConn *net.UDPConn
	UDPAddr *net.UDPAddr
}

// UDPCLients - connected clients via UDP
var UDPCLients []*UDPClient

// Run - starts a streamer to transfer image data.
func Run() {
	logger.Info(LC, "Starting module Streamer")
	go start()
}

// start - starts the UDP server, if something went wrong, os.Exit is called with the code -1.
func start() {
	udpAddr, err := net.ResolveUDPAddr(CONNTYPE, config.Cfg.UDPPath)
	if err != nil {
		logger.Critical(LC, fmt.Sprint("Resolve address error:", err.Error()))
		os.Exit(-1)
	}

	udpConn, err := net.ListenUDP(CONNTYPE, udpAddr)
	if err != nil {
		logger.Critical(LC, fmt.Sprint("UDP connection error:", err.Error()))
		os.Exit(-1)
	}
	defer udpConn.Close()

	logger.Info(LC, fmt.Sprint("UDP server has been started at path ", config.Cfg.UDPPath))
	listenHandler(udpConn)
}

/*
listenHandler - handler connections for UDP server

	- udpConn <*net.UDPConn> - UDPConn's struct contains connected client data.
*/
func listenHandler(udpConn *net.UDPConn) {
	for {
		buf := make([]byte, MAXSIZE)

		_, udpAddr, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			logger.Warning(LC, fmt.Sprint("UDP reader error:", err.Error()))
			continue
		}

		if !isExistClient(udpAddr) {
			logger.Info(LC, fmt.Sprint("[", udpAddr.String(), "] Connected new UDP client"))

			udpClient := &UDPClient{
				UDPConn: udpConn,
				UDPAddr: udpAddr,
			}

			UDPCLients = append(UDPCLients, udpClient)
		}
	}
}

/*
isExistClient - checks is exist connected UDP client.
Returns boolean result. True - connected.

	- udpConn <*net.UDPConn> - UDPConn's struct contains connected client data.
*/
func isExistClient(udpConn *net.UDPAddr) bool {
	for _, udpClient := range UDPCLients {
		if udpClient.UDPAddr.String() == udpConn.String() {
			return true
		}
	}

	return false
}

/*
Send - sends a message over UDP

	- data <[]byte> - data for sending.
*/
func Send(data []byte) {
	for _, udpClient := range UDPCLients {
		// logger.Info(LC, fmt.Sprint("Sending to ", udpClient.UDPAddr.String()))
		if _, err := udpClient.UDPConn.WriteToUDP(data, udpClient.UDPAddr); err != nil {
			logger.Error(LC, fmt.Sprint("Error sending:", err.Error()))
		}
	}
}
