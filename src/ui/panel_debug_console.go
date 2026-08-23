package ui

import "unicode/utf8"

const (
	debugConsoleMaxVisibleLines = 24
	debugConsoleMaxLineLength   = 96
)

var (
	debugConsoleX    float32 = 0.02
	debugConsoleY    float32 = 0.02
	debugConsoleSize float32 = 1.0
)

func wrapDebugConsoleMessage(message string) []string {
	if message == "" {
		return []string{""}
	}

	lines := make([]string, 0, (len(message)/debugConsoleMaxLineLength)+1)
	for len(message) > debugConsoleMaxLineLength {
		splitIndex := debugConsoleMaxLineLength
		for splitIndex > 0 && !utf8.RuneStart(message[splitIndex]) {
			splitIndex--
		}
		lines = append(lines, message[:splitIndex])
		message = message[splitIndex:]
	}
	return append(lines, message)
}

func getVisibleDebugConsoleLines() []string {
	messages := getDebugConsoleMessages(debugConsoleMaxVisibleLines)
	lines := make([]string, 0, len(messages))
	for _, message := range messages {
		lines = append(lines, wrapDebugConsoleMessage(message)...)
	}
	if len(lines) > debugConsoleMaxVisibleLines {
		lines = lines[len(lines)-debugConsoleMaxVisibleLines:]
	}
	return lines
}

func drawDebugConsole() {
	textSize := float32(0.012) * debugConsoleSize
	rowSpacing := float32(0.014) * debugConsoleSize
	textColor := color(255, 255, 255, 255)
	lines := getVisibleDebugConsoleLines()

	addText("--- DEBUG CONSOLE ---", debugConsoleX, debugConsoleY, textSize, 0, textColor)
	for lineIndex, line := range lines {
		lineY := debugConsoleY + float32(lineIndex+1)*rowSpacing
		addText(line, debugConsoleX, lineY, textSize, 0, textColor)
	}
}
