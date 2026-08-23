package ui

import "fmt"

// Suspension Panel Coords
const (
	suspensionBaseTextSize   = 0.016
	suspensionBaseRowSpacing = 0.023
	suspensionBaseWidth      = 0.205
	suspensionBaseHeight     = 0.105
)

var (
	suspensionX          = float32(0.75)
	suspensionY          = float32(0.09)
	suspensionSize       = float32(1.0)
	suspensionTextSize   = float32(suspensionBaseTextSize)
	suspensionRowSpacing = float32(suspensionBaseRowSpacing)
)

func suspensionCompression(index int) float32 {
	if !activeState.HasSuspensionData || !activeState.HasSuspensionMaxData || index < 0 || index >= 2 {
		return 0
	}
	maxTravel := activeState.SuspensionMaxTravel[index]
	if maxTravel <= 0 {
		return 0
	}
	return clampFloat32((maxTravel-activeState.SuspensionLength[index])/maxTravel, 0, 1)
}

func suspensionContactText(index int) string {
	if activeState.SuspensionWheelMat[index] == 0 {
		return "AIR"
	}
	return "GRIP"
}

func suspensionBarColor(compression float32) Color {
	switch {
	case compression >= 0.90:
		return color(220, 40, 35, 220)
	case compression >= 0.75:
		return color(235, 160, 35, 220)
	default:
		return color(45, 180, 110, 220)
	}
}

func drawSuspensionRow(label string, index int, x, y, width, textSize float32) {
	labelColor := color(235, 235, 235, 255)
	valueColor := color(210, 210, 210, 255)
	barBgColor := color(42, 42, 42, 210)
	barBorderColor := color(110, 110, 110, 220)

	compression := suspensionCompression(index)
	barX := x + 0.055*suspensionSize
	barY := y + 0.006*suspensionSize
	barW := width - 0.065*suspensionSize
	barH := 0.008 * suspensionSize
	fillW := barW * compression

	addText(label, x, y, textSize, 0, labelColor)
	addQuadCCW(barX, barY, barX+barW, barY+barH, barBgColor)
	addQuadCCW(barX, barY, barX+barW, barY+0.0015*suspensionSize, barBorderColor)
	addQuadCCW(barX, barY+barH-0.0015*suspensionSize, barX+barW, barY+barH, barBorderColor)
	if fillW > 0 {
		addQuadCCW(barX, barY, barX+fillW, barY+barH, suspensionBarColor(compression))
	}

	y += 0.012 * suspensionSize
	addText(fmt.Sprintf("%3.0f%%  %+05.2fm/s  %s", compression*100, activeState.SuspensionVelocity[index], suspensionContactText(index)), barX, y, textSize*0.78, 0, valueColor)
}

func drawSuspensionPanel() {
	x := suspensionX
	y := suspensionY
	width := suspensionBaseWidth * suspensionSize
	height := suspensionBaseHeight * suspensionSize
	textSize := suspensionTextSize

	addQuadCCW(x, y, x+width, y+height, color(0, 0, 0, 155))
	addText("SUSPENSION", x+0.008*suspensionSize, y+0.006*suspensionSize, textSize*0.9, 0, color(255, 255, 255, 255))

	if !activeState.HasSuspensionData {
		addText("No telemetry", x+0.008*suspensionSize, y+0.039*suspensionSize, textSize, 0, color(190, 190, 190, 255))
		return
	}

	rowY := y + 0.032*suspensionSize
	drawSuspensionRow("F", 0, x+0.008*suspensionSize, rowY, width-0.016*suspensionSize, textSize)
	rowY += suspensionRowSpacing
	drawSuspensionRow("R", 1, x+0.008*suspensionSize, rowY, width-0.016*suspensionSize, textSize)
}
