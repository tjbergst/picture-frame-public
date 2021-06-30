package frame

import "testing"

func TestMoveMouse(t *testing.T) {
	moveMouse(MouseCoordinates{x: 400, y: 35}, false)
}

func TestMoveMouseClick(t *testing.T) {
	moveMouse(MouseCoordinates{x: 300, y: 300}, true)
}

func TestCalibrationCheck(t *testing.T) {
	positions := []string{"top-left", "top-right", "bottom-left", "bottom-right"}
	var color string

	for _, pos := range positions {
		moveMouse(navMap[pos], false)
		color = getColor()
		t.Logf("got color at %s: %s", pos, color)
	}
}
