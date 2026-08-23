// Package ui renders the PitBoardzz HUD and owns its editable layout.
package ui

import (
	"github.com/MaximusTrue/PitBoardzz/src/config"
	"github.com/MaximusTrue/PitBoardzz/src/plugin"
)

const (
	SpriteIndexPointer       = 1
	SpriteIndexPointerShadow = 2
)

// Renderer is implemented by the CGO bridge so UI code never depends on C types.
type Renderer interface {
	Quad(x0, y0, x1, y1 float32, color Color)
	Sprite(x0, y0, x1, y1 float32, spriteIndex int, color Color)
	Text(text string, x, y, size float32, justify int, color Color)
	Triangle(x, y, height float32, color Color)
	PyramidTriangle(x, y, height float32, color Color)
	Exclamation(x, y, height, thickness float32, color Color)
}

type LogWriter func(format string, args ...any)
type DebugMessages func(maxMessages int) []string

var (
	activeState               *plugin.State
	activeRenderer            Renderer
	configuration             *config.Store
	logWriter                 LogWriter
	debugMessages             DebugMessages
	showStopwatch             = true
	showLeaderboard           = true
	showLeaderboardBackground = true
	showSpeedometer           = true
	showFuelInfo              = true
	showSuspension            = true
	showDebugPanel            = false
	showDebugConsole          = false
)

// Initialize connects the HUD to plugin state and its runtime services.
func Initialize(state *plugin.State, store *config.Store, logger LogWriter, messages DebugMessages) {
	activeState = state
	configuration = store
	logWriter = logger
	debugMessages = messages
	initLayoutConfig()
}

// Periodic refreshes input and configuration state once per draw frame.
func Periodic() {
	periodic()
}

// Draw renders all enabled HUD panels through the SDK-backed renderer.
func Draw(renderer Renderer) {
	activeRenderer = renderer
	defer func() { activeRenderer = nil }()

	if showSpeedometer {
		drawSpeedometer()
	}
	if showStopwatch {
		drawTimingPanel()
	}
	if showFuelInfo {
		drawFuelPanel()
	}
	if showSuspension {
		drawSuspensionPanel()
	}
	if showDebugPanel {
		drawDebugPanel()
	}
	if showDebugConsole {
		drawDebugConsole()
	}
	if showLeaderboard {
		drawLeaderboardPanel()
	}
	if editModeActive {
		drawEditModeOverlay()
	}
}

func writeLog(format string, args ...any) {
	if logWriter != nil {
		logWriter(format, args...)
	}
}

func getDebugConsoleMessages(maxMessages int) []string {
	if debugMessages == nil {
		return nil
	}
	return debugMessages(maxMessages)
}

func addQuadCCW(x0, y0, x1, y1 float32, color Color) {
	if activeRenderer != nil {
		activeRenderer.Quad(x0, y0, x1, y1, color)
	}
}

func addSpriteQuadCCW(x0, y0, x1, y1 float32, spriteIndex int, color Color) {
	if activeRenderer != nil {
		activeRenderer.Sprite(x0, y0, x1, y1, spriteIndex, color)
	}
}

func addText(text string, x, y, size float32, justify int, color Color) {
	if activeRenderer != nil {
		activeRenderer.Text(text, x, y, size, justify, color)
	}
}

func addTriangle(x, y, height float32, color Color) {
	if activeRenderer != nil {
		activeRenderer.Triangle(x, y, height, color)
	}
}

func addPyramidTriangle(x, y, height float32, color Color) {
	if activeRenderer != nil {
		activeRenderer.PyramidTriangle(x, y, height, color)
	}
}

func addExclamation(x, y, height, thickness float32, color Color) {
	if activeRenderer != nil {
		activeRenderer.Exclamation(x, y, height, thickness, color)
	}
}

func min(first, second int) int {
	if first < second {
		return first
	}
	return second
}

func max(first, second int) int {
	if first > second {
		return first
	}
	return second
}
