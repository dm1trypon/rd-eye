package main

import (
	"fmt"
	"os"
	"time"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/rd-eye/config"
	"github.com/dm1trypon/rd-eye/screencapturer"
	"github.com/dm1trypon/rd-eye/streamer"
)

// LC - logging category
const LC = "MAIN"

// MAXSIZE - max size package
const MAXSIZE = 2048

var bufFrames [][]byte

func setLogger() {
	cfg := logger.Cfg{
		AppName: "RD_EYE",
		LogPath: "",
		Level:   0,
	}

	logger.SetConfig(cfg)
}

func screenCapturerEvents() {
	for {
		bufFrames = append(bufFrames, <-screencapturer.StreamBuf)
	}
}

func looper() {
	for {
		time.Sleep(200 * time.Millisecond)
	}
}

func payLoader(data []byte) {
	logger.Info(LC, fmt.Sprint("Length: ", len(data)))
	for {
		if len(data) < MAXSIZE {
			data = append(data, "\n"...)
			streamer.Send(data)

			break
		}

		part := data[:MAXSIZE]
		data = data[MAXSIZE:]
		streamer.Send(part)
	}

}

func main() {
	setLogger()
	logger.Info(LC, "[STARTING SERVICE]")

	if !config.Run("/home/dmitry/Projects/ProjectsGo/src/github.com/dm1trypon/rd-eye/config.json",
		"/home/dmitry/Projects/ProjectsGo/src/github.com/dm1trypon/rd-eye/config.schema.json") {
		os.Exit(-1)
	}

	streamer.Run()
	screencapturer.Run()
	screenCapturerEvents()
}
