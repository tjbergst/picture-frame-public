package frame

import (
	"time"

	"github.com/go-vgo/robotgo"
	log "github.com/sirupsen/logrus"
)

// MouseCoordinates represent x, y positions on the screen
type MouseCoordinates struct {
	x int
	y int
}

// NavigationMap represents a map of name: coordinates to click
type NavigationMap map[string]MouseCoordinates

// CalibrationColors represents a map of name: expected color for
// calibration image
type CalibrationColors map[string]string

var navMap NavigationMap
var calColors CalibrationColors

func init() {
	// load nav map
	navMap = NavigationMap{
		"center":       MouseCoordinates{x: 400, y: 300},
		"menu":         MouseCoordinates{x: 595, y: 32},
		"slideshow":    MouseCoordinates{x: 560, y: 30},
		"corner":       MouseCoordinates{x: 0, y: 0},
		"top-left":     MouseCoordinates{x: 300, y: 100},
		"top-right":    MouseCoordinates{x: 600, y: 100},
		"bottom-left":  MouseCoordinates{x: 300, y: 300},
		"bottom-right": MouseCoordinates{x: 600, y: 300},
	}

	// load calibration colors
	calColors = CalibrationColors{
		"top-left":     "fffd54",
		"top-right":    "eb3223",
		"bottom-left":  "0020f5",
		"bottom-right": "75fa4c",
	}
}

func moveMouse(coords MouseCoordinates, click bool) {
	log.Debugf("moving mouse to %d, %d\n", coords.x, coords.y)
	robotgo.MoveMouse(coords.x, coords.y)

	if click {
		time.Sleep(500 * time.Millisecond)
		clickMouse()
	}
}

func clickMouse() {
	log.Debug("clicking mouse")
	robotgo.MouseClick("left", false)
}

func getColor() string {
	color := robotgo.GetMouseColor()
	log.Debugf("got mouse color: %s", color)
	return color
}

func checkColors(threshold int) bool {
	log.Info("checking calibration colors")
	var wrongColors int

	for pos, expectedColor := range calColors {
		moveMouse(navMap[pos], false)
		if mouseColor := getColor(); mouseColor != expectedColor {
			log.Warnf("found incorrect color at %s: %s [%s]", pos, mouseColor, expectedColor)
			wrongColors++
		}
	}

	if wrongColors >= threshold {
		log.Errorf("found %d wrong colors, fail", wrongColors)
		return false
	}
	return true
}
