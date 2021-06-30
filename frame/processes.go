package frame

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
)

func killChromium() {
	log.Info("killing chromium process")

	cmd := exec.Command("pkill", "chromium")
	if err := cmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if code := exitError.ExitCode(); code != 0 {
				log.Warnf("chromium not running [%d]\n", code)
			} else {
				log.Errorf("error killing chromium\n%s [%d]\n", err, code)
			}
		}
	}
}
