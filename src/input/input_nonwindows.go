//go:build !windows

package input

const uiAspectRatio = float32(16.0 / 9.0)

var (
	cursorEnabled bool
	wbLeft        float32
	wbTop         float32
	wbRight       float32 = 1
	wbBottom      float32 = 1
)

func Available() bool {
	return false
}

func UpdateFocusState() {}

func IsTabDown() bool {
	return false
}

func IsLeftMouseDown() bool {
	return false
}

func IsRightMouseDown() bool {
	return false
}

func NormalizedMousePosition() (float32, float32, bool) {
	return 0, 0, false
}

func CursorXScale() float32 {
	return 1.0 / uiAspectRatio
}
