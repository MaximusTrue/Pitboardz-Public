package ui

import "fmt"

// Fuel Info Coords
const (
	fuelBaseTextSize   = 0.025
	fuelBaseRowSpacing = 0.03
)

var (
	fuelX          = float32(0.75)
	fuelY          = float32(0.02)
	fuelSize       = float32(1.0)
	fuelTextSize   = float32(fuelBaseTextSize)
	fuelRowSpacing = float32(fuelBaseRowSpacing)
)

func drawFuelPanel() {
	var fuelText string
	fuelColor := color(255, 255, 255, 255)
	if !activeState.HasFuelData {
		fuelText = "Fuel: ---"
	} else {
		fuelText = fmt.Sprintf("Fuel: %.2f L", activeState.CurrentFuel)
		if activeState.CurrentFuel < 1.0 {
			fuelColor = color(255, 0, 0, 255)
		} else if activeState.CurrentFuel < 2.0 {
			fuelColor = color(255, 255, 0, 255)
		}
	}
	addText(fuelText, fuelX, fuelY, fuelTextSize, 0, fuelColor)
	var deltaText string
	if activeState.LastLapFuelDelta == 0.0 {
		deltaText = "Last: ---"
	} else {
		deltaText = fmt.Sprintf("Last: %.2f L", activeState.LastLapFuelDelta)
	}
	addText(deltaText, fuelX, fuelY+fuelRowSpacing, fuelTextSize, 0, color(255, 255, 255, 255))
}
