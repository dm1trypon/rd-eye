package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/rd-eye/config"
	"github.com/dm1trypon/rd-eye/screencapturer"
	"github.com/dm1trypon/rd-eye/tcpstreamer"
)

// LC - logging category
const LC = "MAIN"

// MAXSIZE - max size package
const MAXSIZE = 2500

var bufFrames [][]byte

func setLogger() {
	cfg := logger.Cfg{
		AppName: "RD_EYE",
		LogPath: "",
		Level:   0,
	}

	logger.SetConfig(cfg)
}

func sender() {
	for {
		// logger.Info(LC, fmt.Sprint("Length: ", len(screencapturer.PackagesBuf)))
		if len(screencapturer.PackagesBuf) < 1 {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		pckg := screencapturer.PackagesBuf[0]

		for _, part := range pckg {
			// logger.Info(LC, fmt.Sprint("UUID part of package: ", part.ID, "; Part: ", part.Part, "; END: ", part.End))

			bData, err := json.Marshal(part)
			if err != nil {
				logger.Error(LC, fmt.Sprint("An error occurred while converting the structure to JSON: ", err.Error()))
				continue
			}

			tcpstreamer.Send(bData)
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func main() {
	setLogger()
	logger.Info(LC, "[STARTING SERVICE]")

	if !config.Run("E:/PROJECTS/Go/src/github.com/dm1trypon/rd-eye/config.json",
		"E:/PROJECTS/Go/src/github.com/dm1trypon/rd-eye/config.schema.json") {
		os.Exit(-1)
	}

	tcpstreamer.Run()
	go sender()
	screencapturer.Run()
}
