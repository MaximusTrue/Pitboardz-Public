package ui

import (
	"fmt"

	"github.com/MaximusTrue/PitBoardzz/src/input"
	"gopkg.in/ini.v1"
)

// UIElement represents a UI component that can be toggled
type UIElement struct {
	ID      string // Unique identifier
	Name    string // Display name
	Enabled *bool  // Pointer to enabled flag
}

// GUI Panel state
var (
	editModeActive bool
	panelX         float32 = 0.35
	panelY         float32 = 0.25
	wasTabDown     bool
	wasRightDown   bool
	wasLeftDown    bool
	uiElements     []UIElement

	// Drag state: draggingElement identifies what is being dragged.
	// "" = nothing, "panel" = settings panel, or a UIElement ID.
	draggingElement string
	dragOffsetX     float32
	dragOffsetY     float32
)

// Panel dimensions
const (
	panelWidth        float32 = 0.30
	panelHeaderHeight float32 = 0.04
	panelRowHeight    float32 = 0.035
	panelFooterHeight float32 = 0.03
	panelPadding      float32 = 0.01
	buttonWidth       float32 = 0.06
	buttonHeight      float32 = 0.025
	updateButtonWidth float32 = 0.10
)

// initUIElements sets up the element registry
func initUIElements() {
	uiElements = []UIElement{
		{ID: "leaderboard", Name: "Leaderboard", Enabled: &showLeaderboard},
		{ID: "leaderboard_bg", Name: "Leaderboard Background", Enabled: &showLeaderboardBackground},
		{ID: "stopwatch", Name: "Stopwatch", Enabled: &showStopwatch},
		{ID: "speedometer", Name: "Speedometer", Enabled: &showSpeedometer},
		{ID: "fuel", Name: "Fuel Info", Enabled: &showFuelInfo},
		{ID: "suspension", Name: "Suspension", Enabled: &showSuspension},
		{ID: "debug_panel", Name: "Debug Panel", Enabled: &showDebugPanel},
		{ID: "debug_console", Name: "Debug Console", Enabled: &showDebugConsole},
	}
}

// getPanelBounds returns the full panel bounds
func getPanelBounds() (x0, y0, x1, y1 float32) {
	x0 = panelX
	y0 = panelY
	x1 = panelX + panelWidth
	// Height = header + rows + optional action rows + footer + padding.
	numRows := len(uiElements)
	extraRows := 0
	if activeState.TestBikeEventType == 1 {
		extraRows++
	}
	if activeState.UpdateAvailable() {
		extraRows++
	}
	panelHeight := panelHeaderHeight + float32(numRows+extraRows)*panelRowHeight + panelFooterHeight + panelPadding*2
	y1 = panelY + panelHeight
	return
}

func isPointInRect(px, py, x0, y0, x1, y1 float32) bool {
	return px >= x0 && px <= x1 && py >= y0 && py <= y1
}

func clampFloat32(v, minV, maxV float32) float32 {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

// getButtonBounds returns the bounds for a toggle button at row index
func getButtonBounds(rowIndex int) (x0, y0, x1, y1 float32) {
	// Button is on the right side of each row
	x1 = panelX + panelWidth - panelPadding
	x0 = x1 - buttonWidth
	y0 = panelY + panelHeaderHeight + panelPadding + float32(rowIndex)*panelRowHeight + (panelRowHeight-buttonHeight)/2
	y1 = y0 + buttonHeight
	return
}

// getResetButtonBounds returns the bounds for the "Reset Session Time" button
func getResetButtonBounds() (x0, y0, x1, y1 float32) {
	resetBtnWidth := float32(0.14)
	rowIndex := len(uiElements)
	cx := panelX + panelWidth/2
	x0 = cx - resetBtnWidth/2
	x1 = cx + resetBtnWidth/2
	y0 = panelY + panelHeaderHeight + panelPadding + float32(rowIndex)*panelRowHeight + (panelRowHeight-buttonHeight)/2
	y1 = y0 + buttonHeight
	return
}

// getUpdateButtonBounds returns the bounds for the update action.
func getUpdateButtonBounds() Bounds {
	rowIndex := len(uiElements)
	if activeState.TestBikeEventType == 1 {
		rowIndex++
	}
	centerX := panelX + panelWidth/2
	buttonY := panelY + panelHeaderHeight + panelPadding + float32(rowIndex)*panelRowHeight + (panelRowHeight-buttonHeight)/2
	return Bounds{
		X0: centerX - updateButtonWidth/2,
		Y0: buttonY,
		X1: centerX + updateButtonWidth/2,
		Y1: buttonY + buttonHeight,
	}
}

// getElementBounds returns the screen bounds for a UI element by ID.
func getElementBounds(id string) (x0, y0, x1, y1 float32, ok bool) {
	switch id {
	case "stopwatch":
		if !showStopwatch {
			return 0, 0, 0, 0, false
		}
		// Width covers the widest possible line: "Last: 00:00.000" + gap + delta
		w := textWidthMono("Last: 00:00.000", stopwatchTextSize) + (0.015 * stopwatchSize) + textWidthMono("+00.000", stopwatchTextSize)
		// Height: 3 rows + delta text below third row
		h := 2*stopwatchRowSpacing + (0.022 * stopwatchSize) + stopwatchTextSize
		return timingX0, timingY0, timingX0 + w, timingY0 + h, true
	case "leaderboard":
		if !showLeaderboard {
			return 0, 0, 0, 0, false
		}
		w := leaderboardBaseWidth * leaderboardSize
		bottom := leaderboardBackgroundBottom
		if bottom <= leaderboardY {
			bottom = leaderboardY + leaderboardBaseHeight*leaderboardSize
		}
		return leaderboardX, leaderboardY, leaderboardX + w, bottom, true
	case "speedometer":
		if !showSpeedometer {
			return 0, 0, 0, 0, false
		}
		gearW := textWidthMono("0", speedometerBaseGearSize*speedometerSize)
		speedExtent := (0.025 * speedometerSize) + textWidthMono("(000)", speedometerBaseSpeedValueSize*speedometerSize)
		w := gearW
		if speedExtent > w {
			w = speedExtent
		}
		h := speedometerBaseGearSize * speedometerSize
		unitBottom := (0.059 * speedometerSize) + speedometerBaseUnitSize*speedometerSize
		if unitBottom > h {
			h = unitBottom
		}
		return speedometerX, speedometerY, speedometerX + w, speedometerY + h, true
	case "fuel":
		if !showFuelInfo {
			return 0, 0, 0, 0, false
		}
		w := textWidthMono("Fuel: 100.00 L", fuelTextSize)
		h := fuelRowSpacing + fuelTextSize
		return fuelX, fuelY, fuelX + w, fuelY + h, true
	case "suspension":
		if !showSuspension {
			return 0, 0, 0, 0, false
		}
		w := suspensionBaseWidth * suspensionSize
		h := suspensionBaseHeight * suspensionSize
		return suspensionX, suspensionY, suspensionX + w, suspensionY + h, true
	case "debug_panel":
		if !showDebugPanel {
			return 0, 0, 0, 0, false
		}
		textSize := float32(0.014) * debugPanelSize
		rowSpacing := float32(0.016) * debugPanelSize
		sectionSpacing := float32(0.006) * debugPanelSize
		// 8 section headers + 27 value lines + 8 section gaps
		numRows := 35
		w := textWidthMono("EventType:2 [EvInit]  BikeSession:5 [RunInit]___", textSize)
		h := float32(numRows)*rowSpacing + 8*sectionSpacing
		return debugPanelX, debugPanelY, debugPanelX + w, debugPanelY + h, true
	case "debug_console":
		if !showDebugConsole {
			return 0, 0, 0, 0, false
		}
		textSize := float32(0.012) * debugConsoleSize
		rowSpacing := float32(0.014) * debugConsoleSize
		lines := getVisibleDebugConsoleLines()
		longestLine := "--- DEBUG CONSOLE ---"
		for _, line := range lines {
			if len(line) > len(longestLine) {
				longestLine = line
			}
		}
		width := textWidthMono(longestLine, textSize)
		height := float32(len(lines)+1) * rowSpacing
		return debugConsoleX, debugConsoleY, debugConsoleX + width, debugConsoleY + height, true
	}
	return 0, 0, 0, 0, false
}

// getElementOrigin returns a pointer-pair to the X/Y origin of a UI element.
func getElementOrigin(id string) (px *float32, py *float32) {
	switch id {
	case "stopwatch":
		return &timingX0, &timingY0
	case "leaderboard":
		return &leaderboardX, &leaderboardY
	case "speedometer":
		return &speedometerX, &speedometerY
	case "fuel":
		return &fuelX, &fuelY
	case "suspension":
		return &suspensionX, &suspensionY
	case "debug_panel":
		return &debugPanelX, &debugPanelY
	case "debug_console":
		return &debugConsoleX, &debugConsoleY
	}
	return nil, nil
}

func toggleUIElementByIndex(i int) {
	if i < 0 || i >= len(uiElements) {
		return
	}
	e := &uiElements[i]
	*e.Enabled = !*e.Enabled
	writeLog("Toggled %s: %t", e.Name, *e.Enabled)
	saveIniFile()
}

// updateEditMode handles all edit mode logic - called each frame
func updateEditMode() {
	if !input.Available() {
		return
	}

	if len(uiElements) == 0 {
		initUIElements()
	}

	// Update focus detection with debounce.
	input.UpdateFocusState()

	// Auto-exit edit mode when the game loses focus (alt-tab).
	if editModeActive && !input.CursorEnabled() {
		editModeActive = false
		writeLog("Edit mode auto-disabled (game lost focus)")
	}

	// Tab toggles edit mode on rising edge (only when focused).
	tabDown := input.IsTabDown()
	if input.CursorEnabled() && tabDown && !wasTabDown {
		editModeActive = !editModeActive
		if editModeActive {
			writeLog("Edit mode enabled (Tab)")
		} else {
			writeLog("Edit mode disabled (Tab)")
		}
	}
	wasTabDown = tabDown

	if !editModeActive {
		return
	}

	// Mouse interactions use current pointer coordinates.
	mx, my, ok := input.NormalizedMousePosition()
	if !ok {
		return
	}
	leftDown := input.IsLeftMouseDown()
	rightDown := input.IsRightMouseDown()

	// Left-click on a row button toggles that setting or invokes an action.
	if leftDown && !wasLeftDown {
		clicked := false
		for i := range uiElements {
			bx0, by0, bx1, by1 := getButtonBounds(i)
			if isPointInRect(mx, my, bx0, by0, bx1, by1) {
				toggleUIElementByIndex(i)
				clicked = true
				break
			}
		}
		if !clicked && activeState.TestBikeEventType == 1 {
			rx0, ry0, rx1, ry1 := getResetButtonBounds()
			if isPointInRect(mx, my, rx0, ry0, rx1, ry1) {
				activeState.AccumulatedOnTrackTime = 0
				activeState.TelemetryTime = float32(-1)
				writeLog("Settings: Session time reset by user")
				clicked = true
			}
		}
		if !clicked && activeState.UpdateAvailable() && buttonClicked(mx, my, getUpdateButtonBounds()) {
			activeState.StartUpdate()
		}
	}
	wasLeftDown = leftDown

	// Right-click drag moves UI elements or the settings panel.
	if rightDown && !wasRightDown {
		draggingElement = ""
		// Check UI elements first (on-screen HUD widgets).
		for _, e := range uiElements {
			ex0, ey0, ex1, ey1, hit := getElementBounds(e.ID)
			if hit && isPointInRect(mx, my, ex0, ey0, ex1, ey1) {
				draggingElement = e.ID
				dragOffsetX = mx - ex0
				dragOffsetY = my - ey0
				break
			}
		}
		// Fall back to settings panel.
		if draggingElement == "" {
			px0, py0, px1, py1 := getPanelBounds()
			if isPointInRect(mx, my, px0, py0, px1, py1) {
				draggingElement = "panel"
				dragOffsetX = mx - panelX
				dragOffsetY = my - panelY
			}
		}
	} else if !rightDown && wasRightDown {
		// On release, save if we were dragging a UI element.
		if draggingElement != "" && draggingElement != "panel" {
			saveIniFile()
		}
		draggingElement = ""
	}
	if draggingElement == "panel" {
		windowLeft, windowTop, windowRight, windowBottom := input.WindowBounds()
		newX := mx - dragOffsetX
		newY := my - dragOffsetY
		_, _, bx1, by1 := getPanelBounds()
		pw := bx1 - panelX
		ph := by1 - panelY
		panelX = clampFloat32(newX, windowLeft, windowRight-pw)
		panelY = clampFloat32(newY, windowTop, windowBottom-ph)
	} else if draggingElement != "" {
		windowLeft, windowTop, windowRight, windowBottom := input.WindowBounds()
		px, py := getElementOrigin(draggingElement)
		if px != nil && py != nil {
			_, _, ex1, ey1, _ := getElementBounds(draggingElement)
			ew := ex1 - *px
			eh := ey1 - *py
			oldY := *py
			*px = clampFloat32(mx-dragOffsetX, windowLeft, windowRight-ew)
			*py = clampFloat32(my-dragOffsetY, windowTop, windowBottom-eh)
			// Keep leaderboard background bottom in sync with its origin.
			if draggingElement == "leaderboard" {
				leaderboardBackgroundBottom += *py - oldY
			}
		}
	}
	wasRightDown = rightDown
}

// drawArrowCursor draws the custom pointer sprite with its tip at (mx, my).
func drawArrowCursor(mx, my float32) {
	const (
		h          = float32(0.030)
		widthScale = float32(1.05)
	)
	tint := color(255, 255, 255, 255)
	w := h * input.CursorXScale() * widthScale

	sox := w * 0.10
	soy := h * 0.10

	// Clamp so the sprite + shadow stays within the visible window.
	windowLeft, windowTop, windowRight, windowBottom := input.WindowBounds()
	mx = clampFloat32(mx, windowLeft, windowRight-w-sox)
	my = clampFloat32(my, windowTop, windowBottom-h-soy)

	addSpriteQuadCCW(mx+sox, my+soy, mx+w+sox, my+h+soy, SpriteIndexPointerShadow, tint)
	addSpriteQuadCCW(mx, my, mx+w, my+h, SpriteIndexPointer, tint)
}

// drawEditModeOverlay draws the settings panel
func drawEditModeOverlay() {
	panelBgColor := color(32, 32, 32, 224)
	panelBorderColor := color(255, 255, 255, 255)
	headerBgColor := color(48, 48, 48, 255)
	textColor := color(255, 255, 255, 255)
	buttonOnColor := color(0, 170, 0, 255)
	buttonOffColor := color(170, 0, 0, 255)
	dimTextColor := color(170, 170, 170, 255)

	px0, py0, px1, py1 := getPanelBounds()

	addQuadCCW(px0, py0, px1, py1, panelBgColor)

	borderWidth := float32(0.003)
	addQuadCCW(px0, py0, px1, py0+borderWidth, panelBorderColor)
	addQuadCCW(px0, py1-borderWidth, px1, py1, panelBorderColor)
	addQuadCCW(px0, py0, px0+borderWidth, py1, panelBorderColor)
	addQuadCCW(px1-borderWidth, py0, px1, py1, panelBorderColor)

	addQuadCCW(px0+borderWidth, py0+borderWidth, px1-borderWidth, py0+panelHeaderHeight, headerBgColor)

	addText("PITBOARDZ SETTINGS", px0+panelWidth/2, py0+0.008, 0.022, 1, textColor)

	addQuadCCW(px0, py0+panelHeaderHeight, px1, py0+panelHeaderHeight+borderWidth, panelBorderColor)

	for i, e := range uiElements {
		rowY := py0 + panelHeaderHeight + panelPadding + float32(i)*panelRowHeight

		addText(e.Name, px0+panelPadding+0.01, rowY+0.005, 0.020, 0, textColor)

		bx0, by0, bx1, by1 := getButtonBounds(i)

		var btnColor Color
		if *e.Enabled {
			btnColor = buttonOnColor
		} else {
			btnColor = buttonOffColor
		}
		addQuadCCW(bx0, by0, bx1, by1, btnColor)

		btnBorder := float32(0.002)
		addQuadCCW(bx0, by0, bx1, by0+btnBorder, panelBorderColor)
		addQuadCCW(bx0, by1-btnBorder, bx1, by1, panelBorderColor)
		addQuadCCW(bx0, by0, bx0+btnBorder, by1, panelBorderColor)
		addQuadCCW(bx1-btnBorder, by0, bx1, by1, panelBorderColor)

		var btnText string
		if *e.Enabled {
			btnText = "ON"
		} else {
			btnText = "OFF"
		}
		addText(btnText, bx0+buttonWidth/2, by0+0.003, 0.016, 1, textColor)
	}

	// Reset Session Time button (practice mode only)
	if activeState.TestBikeEventType == 1 {
		resetBtnColor := color(0, 68, 136, 255)
		rx0, ry0, rx1, ry1 := getResetButtonBounds()
		addQuadCCW(rx0, ry0, rx1, ry1, resetBtnColor)
		btnBorderW := float32(0.002)
		addQuadCCW(rx0, ry0, rx1, ry0+btnBorderW, panelBorderColor)
		addQuadCCW(rx0, ry1-btnBorderW, rx1, ry1, panelBorderColor)
		addQuadCCW(rx0, ry0, rx0+btnBorderW, ry1, panelBorderColor)
		addQuadCCW(rx1-btnBorderW, ry0, rx1, ry1, panelBorderColor)
		addText("Reset Session Time", rx0+(rx1-rx0)/2, ry0+0.003, 0.016, 1, textColor)
	}

	if activeState.UpdateAvailable() {
		updateButtonColor := color(35, 185, 60, 255)
		updateButtonText := "Update"
		if activeState.UpdateInProgress() {
			updateButtonText = "Updating & Closing"
		}
		drawButton(updateButtonText, getUpdateButtonBounds(), 0.016, updateButtonColor, textColor)
	}

	instructY := py1 - panelFooterHeight + 0.008
	addText("LMB: Toggle | RMB: Drag | Tab: Close", px0+panelPadding, instructY, 0.012, 0, dimTextColor)

	if input.Available() {
		mx, my, ok := input.NormalizedMousePosition()
		if ok {
			drawArrowCursor(mx, my)
		}
	}
}

// saveIniFile writes current element positions and sizes to the INI file
func saveIniFile() {
	if configuration == nil || !configuration.Available() {
		writeLog("saveIniFile: No save path set")
		return
	}

	// Load existing INI to preserve comments and other settings
	cfg, err := configuration.Load()
	if err != nil {
		writeLog("saveIniFile: Failed to load INI: %v", err)
		cfg = ini.Empty()
	}

	// Update Stopwatch section
	cfg.Section("Stopwatch").Key("TopLeftX").SetValue(fmt.Sprintf("%.3f", timingX0))
	cfg.Section("Stopwatch").Key("TopLeftY").SetValue(fmt.Sprintf("%.3f", timingY0))
	cfg.Section("Stopwatch").Key("size").SetValue(fmt.Sprintf("%.2f", stopwatchSize))
	cfg.Section("Stopwatch").Key("enabled").SetValue(fmt.Sprintf("%t", showStopwatch))

	// Update Speedometer section
	cfg.Section("Speedometer").Key("topLeftX").SetValue(fmt.Sprintf("%.3f", speedometerX))
	cfg.Section("Speedometer").Key("topLeftY").SetValue(fmt.Sprintf("%.3f", speedometerY))
	cfg.Section("Speedometer").Key("size").SetValue(fmt.Sprintf("%.2f", speedometerSize))
	cfg.Section("Speedometer").Key("milesPerHour").SetValue(fmt.Sprintf("%t", isMilesPerHour))
	cfg.Section("Speedometer").Key("enabled").SetValue(fmt.Sprintf("%t", showSpeedometer))

	// Update Leaderboard section
	cfg.Section("Leaderboard").Key("topLeftX").SetValue(fmt.Sprintf("%.3f", leaderboardX))
	cfg.Section("Leaderboard").Key("topLeftY").SetValue(fmt.Sprintf("%.3f", leaderboardY))
	cfg.Section("Leaderboard").Key("size").SetValue(fmt.Sprintf("%.2f", leaderboardSize))
	cfg.Section("Leaderboard").Key("enabled").SetValue(fmt.Sprintf("%t", showLeaderboard))
	cfg.Section("Leaderboard").Key("background").SetValue(fmt.Sprintf("%t", showLeaderboardBackground))

	// Update Fuel Info section
	cfg.Section("Fuel Info").Key("topLeftX").SetValue(fmt.Sprintf("%.3f", fuelX))
	cfg.Section("Fuel Info").Key("topLeftY").SetValue(fmt.Sprintf("%.3f", fuelY))
	cfg.Section("Fuel Info").Key("size").SetValue(fmt.Sprintf("%.2f", fuelSize))
	cfg.Section("Fuel Info").Key("enabled").SetValue(fmt.Sprintf("%t", showFuelInfo))

	// Update Suspension section
	cfg.Section("Suspension").Key("topLeftX").SetValue(fmt.Sprintf("%.3f", suspensionX))
	cfg.Section("Suspension").Key("topLeftY").SetValue(fmt.Sprintf("%.3f", suspensionY))
	cfg.Section("Suspension").Key("size").SetValue(fmt.Sprintf("%.2f", suspensionSize))
	cfg.Section("Suspension").Key("enabled").SetValue(fmt.Sprintf("%t", showSuspension))

	// Update Debug Panel section
	cfg.Section("Debug Panel").Key("topLeftX").SetValue(fmt.Sprintf("%.3f", debugPanelX))
	cfg.Section("Debug Panel").Key("topLeftY").SetValue(fmt.Sprintf("%.3f", debugPanelY))
	cfg.Section("Debug Panel").Key("size").SetValue(fmt.Sprintf("%.2f", debugPanelSize))
	cfg.Section("Debug Panel").Key("enabled").SetValue(fmt.Sprintf("%t", showDebugPanel))

	// Update Debug Console section
	cfg.Section("Debug Console").Key("topLeftX").SetValue(fmt.Sprintf("%.3f", debugConsoleX))
	cfg.Section("Debug Console").Key("topLeftY").SetValue(fmt.Sprintf("%.3f", debugConsoleY))
	cfg.Section("Debug Console").Key("size").SetValue(fmt.Sprintf("%.2f", debugConsoleSize))
	cfg.Section("Debug Console").Key("enabled").SetValue(fmt.Sprintf("%t", showDebugConsole))

	// Save to file
	err = configuration.Save(cfg)
	if err != nil {
		writeLog("saveIniFile: Failed to save INI: %v", err)
	} else {
		writeLog("saveIniFile: Settings saved successfully")
	}
}
