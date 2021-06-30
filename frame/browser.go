package frame

import (
	"os/exec"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

// Browser represents a web browser process to be started and kept alive
type Browser struct {
	Executable string
	Args       []string
	URL        string
}

// SetDefaults makes the nil value useful by loading defaults
func (b *Browser) SetDefaults() {
	if b.Executable == "" {
		b.Executable = "chromium-browser"
	}
	if len(b.Args) == 0 {
		b.Args = []string{"--kiosk", "--disable-restore-session-state"}
	}
	if b.URL == "" {
		b.URL = "https://photos.app.goo.gl/4kFDGY4tzzjTYr1k7"
	}
}

func (b *Browser) startBrowser(restart chan bool) {
	// kill existing browsers
	killChromium()

	// start the browser
	log.Infof("starting browser with url:\n\t%s\n", b.URL)
	process := exec.Command(b.Executable, append(b.Args, b.URL)...)
	if err := process.Start(); err != nil {
		log.Fatalf("error starting process\n%s", err)
	}

	process.Wait()
	log.Warn("browser process ended")
	restart <- true
}

func (b *Browser) startSlideshow(initDelay int) {
	time.Sleep(time.Duration(initDelay) * time.Second)

	log.Info("starting slideshow")
	moveMouse(navMap["menu"], true)

	time.Sleep(5 * time.Second)
	moveMouse(navMap["slideshow"], true)
}

func (b *Browser) run(isRestart bool) {
	moveMouse(navMap["center"], false)
	restart := make(chan bool)

	// start browser
	wg.Add(1)
	go b.startBrowser(restart)

	// start slideshow
	b.startSlideshow(10)
	wg.Add(1)
	go b.checkRunning(4, restart)

	// attempt to hide mouse
	time.Sleep(2 * time.Second)
	moveMouse(navMap["corner"], false)

	// start refresher if not restarting
	if !isRestart {
		wg.Add(1)
		go refresh(180, restart)
	}

	// watch restart channel
	select {
	case <-restart:
		log.Warn("restart signal received")
		time.Sleep(10 * time.Second)
		b.run(true)
	}
	wg.Done()
}

// ensures slideshow is running by checking colors at 4 quadrants using
// special calibration image (data/calibration-squares.png)
func (b *Browser) checkRunning(initDelay int, restart chan bool) {
	time.Sleep(time.Duration(initDelay) * time.Second)
	log.Info("ensuring slideshow is running")

	var threshold int = 3
	if ok := checkColors(threshold); !ok {
		// sleep for a bit to let run() reach select statement
		time.Sleep(5 * time.Second)
		log.Debug("sending signal to restart channel")
		restart <- true
		wg.Done()
		return
	}
	log.Debug("successfully validated slideshow is running")
	wg.Done()
}

// refresh takes advantage of the run restarter by killing chromium
// periodically to force a restart, loading fresh images
func refresh(minutes int, restart chan bool) {
	log.Infof("starting refresh every %d minutes\n", minutes)
	ticker := time.NewTicker(time.Duration(minutes) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Warn("time to refresh")
			restart <- true
		}
	}
	// wg.Done()
}

// RunForever starts the browser, slideshow and monitors chromium process
// restarting if necessary
func (b *Browser) RunForever() {
	b.run(false)
	wg.Wait()
}
