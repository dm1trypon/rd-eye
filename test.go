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
const MAXSIZE = 3000

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:44444")
	if err != nil {
		logger.Critical(LC, err.Error())
		return
	}

	for {
		buf := make([]byte, MAXSIZE)

		_, err := conn.Read(buf)
		if err != nil {
			logger.Warning(LC, fmt.Sprint("TCP reader error:", err.Error()))
			continue
		}

		buf = bytes.Trim(buf, "\x00")

		buffer := new(bytes.Buffer)
		if err := json.Compact(buffer, buf); err != nil {
			logger.Error(LC, fmt.Sprint("Compact json error: "))
			continue
		}

		logger.Info(LC, fmt.Sprint("RECV: ", string(buffer.Bytes())))
	}
}
