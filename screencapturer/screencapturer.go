package screencapturer

import (
	"bytes"
	"image/jpeg"
	"time"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/rd-eye/config"
	"github.com/vova616/screenshot"
)

// LC - logging category
const LC = "SCREENCAPTURER"

// StreamBuf - screenshot's data buffer
var StreamBuf = make(chan []byte)

// Run - launches the screen capturing tool
func Run() {
	logger.Info(LC, "Starting module ScreenCapturer")

	go start()
}

func screenEncoder() {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		logger.Error(LC, err.Error())
		return
	}

	buf := new(bytes.Buffer)

	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 10}); err != nil {
		logger.Error(LC, err.Error())
		return
	}

	StreamBuf <- buf.Bytes()
}

func start() {
	for {
		go screenEncoder()
		time.Sleep(time.Duration(config.Cfg.FPS) * time.Millisecond)
	}
}
