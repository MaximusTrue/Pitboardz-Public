//go:build windows
// +build windows

package input

import (
	"syscall"
	"unsafe"
)

// Windows API constants
const (
	VK_LBUTTON = 0x01 // Left mouse button
	VK_RBUTTON = 0x02 // Right mouse button
	VK_TAB     = 0x09

	// MX Bikes UI coordinate system assumes 16:9 aspect ratio.
	uiAspectRatio = float32(16.0 / 9.0)

	// Frames to wait before disabling cursor on focus loss (~83ms at 60fps).
	focusDebounceFrames = 5
)

// Windows API DLLs — matches mxbmrp3 InputManager API surface exactly.
var (
	user32                       = syscall.NewLazyDLL("user32.dll")
	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procGetAsyncKeyState         = user32.NewProc("GetAsyncKeyState")
	procGetCursorPos             = user32.NewProc("GetCursorPos")
	procGetForegroundWindow      = user32.NewProc("GetForegroundWindow")
	procGetClientRect            = user32.NewProc("GetClientRect")
	procScreenToClient           = user32.NewProc("ScreenToClient")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procIsWindow                 = user32.NewProc("IsWindow")
	procGetCurrentProcessId      = kernel32.NewProc("GetCurrentProcessId")

	cachedGameWindow uintptr

	// Focus detection state
	cursorEnabled        bool
	myProcessId          uint32
	framesSinceFocusLost int

	// Window bounds in UI-normalized coordinates.
	// On 16:9: [0, 0, 1, 1]. On 21:9 ultrawide: [-0.17, 0, 1.17, 1].
	wbLeft   float32
	wbTop    float32
	wbRight  float32 = 1
	wbBottom float32 = 1
)

// POINT structure for Windows API
type POINT struct {
	X, Y int32
}

// RECT structure for Windows API
type RECT struct {
	Left, Top, Right, Bottom int32
}

// isLeftMouseDown returns true if left mouse button is currently pressed
func IsLeftMouseDown() bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(VK_LBUTTON))
	return (ret & 0x8000) != 0
}

// isRightMouseDown returns true if right mouse button is currently pressed
func IsRightMouseDown() bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(VK_RBUTTON))
	return (ret & 0x8000) != 0
}

// getCursorPos returns the current cursor position in screen coordinates
func getCursorPos() (x, y int32, ok bool) {
	var pt POINT
	ret, _, _ := procGetCursorPos.Call(uintptr(unsafe.Pointer(&pt)))
	if ret == 0 {
		return 0, 0, false
	}
	return pt.X, pt.Y, true
}

func getForegroundWindow() uintptr {
	ret, _, _ := procGetForegroundWindow.Call()
	return ret
}

func getClientRect(hwnd uintptr) (RECT, bool) {
	var rect RECT
	ret, _, _ := procGetClientRect.Call(hwnd, uintptr(unsafe.Pointer(&rect)))
	return rect, ret != 0
}

func screenToClient(hwnd uintptr, pt *POINT) bool {
	ret, _, _ := procScreenToClient.Call(hwnd, uintptr(unsafe.Pointer(pt)))
	return ret != 0
}

func getWindowThreadProcessId(hwnd uintptr) uint32 {
	var pid uint32
	procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))
	return pid
}

func getCurrentProcessId() uint32 {
	ret, _, _ := procGetCurrentProcessId.Call()
	return uint32(ret)
}

func isWindowValid(hwnd uintptr) bool {
	ret, _, _ := procIsWindow.Call(hwnd)
	return ret != 0
}

// updateFocusState checks window focus with debounce and updates cachedGameWindow.
// Sets cursorEnabled to true when a window belonging to this process is focused.
// Matches mxbmrp3 InputManager::updateCursorEnabled() logic exactly.
func UpdateFocusState() {
	if myProcessId == 0 {
		myProcessId = getCurrentProcessId()
	}

	fg := getForegroundWindow()
	if fg == 0 {
		framesSinceFocusLost++
		if framesSinceFocusLost >= focusDebounceFrames {
			cursorEnabled = false
		}
		return
	}

	fgPid := getWindowThreadProcessId(fg)
	if fgPid == myProcessId {
		framesSinceFocusLost = 0
		cursorEnabled = true

		// Update cached game window if this is a reasonably sized window.
		if fg != cachedGameWindow {
			rect, ok := getClientRect(fg)
			if ok {
				w := rect.Right - rect.Left
				h := rect.Bottom - rect.Top
				if (w >= 640 && h >= 480) || cachedGameWindow == 0 || !isWindowValid(cachedGameWindow) {
					cachedGameWindow = fg
				}
			}
		}
	} else {
		framesSinceFocusLost++
		if framesSinceFocusLost >= focusDebounceFrames {
			cursorEnabled = false
		}
	}
}

// getNormalizedMousePos converts screen mouse position to normalized UI coordinates (0.0-1.0).
// Uses pillarbox/letterbox math to correctly handle any aspect ratio including ultrawide.
// Matches mxbmrp3 InputManager::updateCursorPosition() logic exactly.
func NormalizedMousePosition() (x, y float32, ok bool) {
	hwnd := cachedGameWindow
	if hwnd == 0 {
		return 0, 0, false
	}

	screenX, screenY, cursorOk := getCursorPos()
	if !cursorOk {
		return 0, 0, false
	}

	// Convert screen coordinates to client-area coordinates.
	pt := POINT{X: screenX, Y: screenY}
	if !screenToClient(hwnd, &pt) {
		return 0, 0, false
	}

	// Get client area dimensions (excludes title bar and borders).
	clientRect, rectOk := getClientRect(hwnd)
	if !rectOk {
		return 0, 0, false
	}
	clientWidth := float32(clientRect.Right - clientRect.Left)
	clientHeight := float32(clientRect.Bottom - clientRect.Top)
	if clientWidth <= 0 || clientHeight <= 0 {
		return 0, 0, false
	}

	// Calculate the 16:9 UI area within the client area.
	windowAspect := clientWidth / clientHeight
	var uiWidth, uiHeight, uiOffsetX, uiOffsetY float32

	if windowAspect > uiAspectRatio {
		// Pillarboxed (ultrawide) — black bars on sides.
		uiHeight = clientHeight
		uiWidth = clientHeight * uiAspectRatio
		uiOffsetX = (clientWidth - uiWidth) / 2
		uiOffsetY = 0
	} else {
		// Letterboxed (narrow) — black bars on top/bottom.
		uiWidth = clientWidth
		uiHeight = clientWidth / uiAspectRatio
		uiOffsetX = 0
		uiOffsetY = (clientHeight - uiHeight) / 2
	}

	if uiWidth <= 0 || uiHeight <= 0 {
		return 0, 0, false
	}

	// Update window bounds in UI-normalized coordinates.
	wbLeft = -uiOffsetX / uiWidth
	wbTop = -uiOffsetY / uiHeight
	wbRight = (clientWidth - uiOffsetX) / uiWidth
	wbBottom = (clientHeight - uiOffsetY) / uiHeight

	normalizedX := (float32(pt.X) - uiOffsetX) / uiWidth
	normalizedY := (float32(pt.Y) - uiOffsetY) / uiHeight

	// Clamp to window bounds, not [0,1] — allows ultrawide placement.
	normalizedX = clampFloat32(normalizedX, wbLeft, wbRight)
	normalizedY = clampFloat32(normalizedY, wbTop, wbBottom)
	return normalizedX, normalizedY, true
}

// getCursorXScale returns the width scale needed for square sprites
// in normalized UI coordinates. Since the UI is always 16:9, this is constant.
func CursorXScale() float32 {
	return 1.0 / uiAspectRatio // 9/16 ≈ 0.5625
}

// isKeyDown returns true if the given virtual key is currently pressed
func isKeyDown(vk uintptr) bool {
	ret, _, _ := procGetAsyncKeyState.Call(vk)
	return (ret & 0x8000) != 0
}

func IsTabDown() bool { return isKeyDown(VK_TAB) }

// inputAvailable returns true if Windows input APIs are available
func Available() bool {
	return true
}
