package screencapturer

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"sync"
	"time"

	"github.com/delaemon/go-uuidv4"

	logger "github.com/dm1trypon/easy-logger"
	"github.com/dm1trypon/rd-eye/config"
	"github.com/dm1trypon/rd-eye/models"
	"github.com/dm1trypon/rd-eye/tools"
	"github.com/vova616/screenshot"
)

// LC - logging category
const LC = "SCREENCAPTURER"

// MAXPCKGSIZE - maximum package size
const MAXPCKGSIZE = 2000

// MAXINSTANCE - maximum instance of screenshoter
const MAXINSTANCE = 2

// ReadySend - notifies when buffering ended
var ReadySend = make(chan bool)

var mutex sync.Mutex

// PackagesBuf - screenshot's data buffer
var PackagesBuf []models.Package

// Run - launches the screen capturing tool
func Run() {
	logger.Info(LC, "Starting module ScreenCapturer")

	start()
}

func screenEncoder() {
	img, err := screenshot.CaptureScreen()
	if err != nil {
		logger.Error(LC, err.Error())
		return
	}

	buf := new(bytes.Buffer)

	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 20}); err != nil {
		logger.Error(LC, err.Error())
		return
	}

	pckg, err := splittingPackage(buf.Bytes())
	if err != nil {
		logger.Error(LC, fmt.Sprint("Packet splitting error"))
		return
	}

	mutex.Lock()
	PackagesBuf = append(PackagesBuf, pckg)
	mutex.Unlock()
}

func splittingPackage(buf []byte) (models.Package, error) {
	pckg := models.Package{}
	part := 0
	idPackage, err := uuidv4.Generate()
	if err != nil {
		logger.Error(LC, fmt.Sprint("Generating package ID failed: ", err.Error()))
		return nil, err
	}

	for {
		if len(buf) < MAXPCKGSIZE {
			partPackage := models.PartPackage{
				ID:   idPackage,
				Part: part,
				Data: buf,
				End:  true,
			}

			pckg = append(pckg, partPackage)
			break
		}

		partBuf := buf[:MAXPCKGSIZE]
		buf = buf[MAXPCKGSIZE:]

		partPackage := models.PartPackage{
			ID:   idPackage,
			Part: part,
			Data: partBuf,
			End:  false,
		}

		pckg = append(pckg, partPackage)
		part++
	}

	return pckg, nil
}

func buffering() {
	for {
		if len(PackagesBuf) < config.Cfg.MaxBufSize {
			continue
		}

		buf, ok := tools.RemoveFromBuffer(PackagesBuf, 0)
		if !ok {
			logger.Error(LC, "Buffering error")
			continue
		}

		mutex.Lock()
		PackagesBuf = buf
		mutex.Unlock()
	}
}

func start() {
	instance := 1

	for {
		if instance > MAXINSTANCE {
			break
		}

		go func() {
			for {
				screenEncoder()
				time.Sleep(time.Duration(config.Cfg.FPS) * time.Millisecond)
			}

		}()

		time.Sleep(time.Duration(config.Cfg.FPS) / MAXINSTANCE * time.Millisecond)
		instance++
	}

	buffering()
}
