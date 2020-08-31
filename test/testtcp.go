package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"

	logger "github.com/dm1trypon/easy-logger"
)

// LC - logging category
const LC = "MAIN"

// MAXSIZE - max size package
const MAXSIZE = 10000

var resPackage []byte

func main() {
	logger.Info(LC, "STARTING SERVICE")
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:44444")
	if err != nil {
		logger.Critical(LC, fmt.Sprint("Error: ", err.Error()))
		return
	}

	localUDPAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		logger.Critical(LC, fmt.Sprint("Error: ", err.Error()))
		return
	}

	udpConn, err := net.DialUDP("udp", localUDPAddr, udpAddr)
	if err != nil {
		logger.Critical(LC, fmt.Sprint("Error: ", err.Error()))
		return
	}

	defer udpConn.Close()

	udpConn.Write([]byte("sa"))

	for {
		buf := make([]byte, MAXSIZE)

		_, _, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			logger.Warning(LC, fmt.Sprint("UDP reader error:", err.Error()))
			continue
		}

		buf = bytes.Trim(buf, "\x00")

		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, buf); err != nil {
			logger.Error(LC, fmt.Sprint("Compact json error: ", err.Error()))
			continue
		}

		logger.Info(LC, fmt.Sprint("RECV: ", string(buffer.Bytes())))
		combainPackage(buf)
	}
}

func combainPackage(buf []byte) {
	resPackage = append(resPackage, string(buf)...)

	// logger.Info(LC, fmt.Sprint("Length of package: ", len(buf)))

	if bytes.HasSuffix(buf, []byte("\n")) {
		logger.Info(LC, fmt.Sprint("Length of package: ", len(resPackage)))
		resPackage = resPackage[:0]
		logger.Info(LC, fmt.Sprint("Length: ", resPackage))
	}
}
