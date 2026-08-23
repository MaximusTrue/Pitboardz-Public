package ui

const (
	CHAR_ASPECT = 0.22
)

// Bounds defines the rectangular area occupied by a UI control.
type Bounds struct {
	X0 float32
	Y0 float32
	X1 float32
	Y1 float32
}

func textWidthMono(s string, size float32) float32 {
	return float32(len(s)) * (size * CHAR_ASPECT)
}

// buttonClicked reports whether a click position is inside a button's bounds.
func buttonClicked(mouseX, mouseY float32, bounds Bounds) bool {
	return isPointInRect(mouseX, mouseY, bounds.X0, bounds.Y0, bounds.X1, bounds.Y1)
}

// drawButton draws a filled box with centered text.
func drawButton(text string, bounds Bounds, textSize float32, buttonColor, textColor Color) {
	addQuadCCW(bounds.X0, bounds.Y0, bounds.X1, bounds.Y1, buttonColor)
	textX := bounds.X0 + (bounds.X1-bounds.X0)/2
	textY := bounds.Y0 + (bounds.Y1-bounds.Y0-textSize)/2
	addText(text, textX, textY, textSize, 1, textColor)
}
