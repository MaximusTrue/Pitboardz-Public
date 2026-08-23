package ui

import (
	"fmt"
)

// Stopwatch Coords
const (
	stopwatchBaseTextSize   = 0.021
	stopwatchBaseRowSpacing = 0.031
)

var (
	timingX0            float32 = 0.02
	timingY0            float32 = 0.55
	stopwatchSize       float32 = 1.0
	stopwatchTextSize   float32 = stopwatchBaseTextSize
	stopwatchRowSpacing float32 = stopwatchBaseRowSpacing
	timingRowY                  = timingY0
)

func formatDelta(deltaMS int) string {
	if deltaMS == 0 {
		return "---"
	}
	absMS := deltaMS
	if absMS < 0 {
		absMS = -absMS
	}
	seconds := absMS / 1000
	ms := absMS % 1000
	var sign string
	if deltaMS > 0 {
		sign = "+"
	} else {
		sign = "-"
	}
	if seconds > 0 {
		return fmt.Sprintf("%s%d.%03d", sign, seconds, ms)
	}
	return fmt.Sprintf("%s0.%03d", sign, ms)
}

func formatLapTime(laptimeMS int) string {
	minutes := laptimeMS / 60000
	seconds := (laptimeMS % 60000) / 1000
	milliseconds := laptimeMS % 1000
	return fmt.Sprintf("%02d:%02d.%03d", minutes, seconds, milliseconds)
}

func getDeltaColor() Color {
	if !activeState.HasBestLapData || !activeState.HasStartedLap || activeState.DeltaTimeMS == 0 {
		return color(255, 255, 255, 255)
	}
	if activeState.DeltaTimeMS < 0 {
		return color(0, 255, 0, 255)
	}
	return color(255, 0, 0, 255)
}

func drawTimingPanel() {
	timingRowY = timingY0
	bestText := "Best: " + formatLapTime(activeState.ClientBestLapTimeMS)
	addText(bestText, timingX0, timingRowY, stopwatchTextSize, 0, color(255, 255, 255, 255))
	timingRowY += stopwatchRowSpacing
	lastText := "Last: " + formatLapTime(activeState.ClientLastLapTimeMS)
	addText(lastText, timingX0, timingRowY, stopwatchTextSize, 0, color(255, 255, 255, 255))
	if activeState.ClientLastLapTimeMS > 0 && activeState.ClientBestLapTimeMS > 0 && activeState.ClientLastLapTimeMS != activeState.ClientBestLapTimeMS {
		lastDelta := activeState.ClientLastLapTimeMS - activeState.ClientBestLapTimeMS
		lastDeltaText := formatDelta(lastDelta)
		lastTextWidth := textWidthMono(lastText, stopwatchTextSize)
		deltaX := timingX0 + lastTextWidth + (0.015 * stopwatchSize)
		addText(lastDeltaText, deltaX, timingRowY, stopwatchTextSize, 0, color(255, 0, 0, 255))
	}
	timingRowY += stopwatchRowSpacing
	currentText := formatLapTime(activeState.StopwatchMS)
	prefixWidth := textWidthMono("Best: ", stopwatchTextSize)
	currentX := timingX0 + prefixWidth
	addText(currentText, currentX, timingRowY, stopwatchTextSize, 0, color(255, 255, 255, 255))
	deltaText := formatDelta(activeState.DeltaTimeMS)
	currentWidth := textWidthMono(currentText, stopwatchTextSize)
	deltaWidth := textWidthMono(deltaText, stopwatchTextSize*0.9)
	deltaY := timingRowY + (0.022 * stopwatchSize)
	deltaX := currentX + currentWidth - deltaWidth
	addText(deltaText, deltaX, deltaY, stopwatchTextSize*0.9, 0, getDeltaColor())
}
