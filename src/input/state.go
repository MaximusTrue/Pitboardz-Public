package input

// CursorEnabled reports whether MX Bikes currently owns the focused window.
func CursorEnabled() bool {
	return cursorEnabled
}

// WindowBounds returns the usable UI bounds in normalized coordinates.
func WindowBounds() (left, top, right, bottom float32) {
	return wbLeft, wbTop, wbRight, wbBottom
}

func clampFloat32(value, minimum, maximum float32) float32 {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}
