package ui

import (
	"fmt"
	"strconv"
)

// Speedometer Coords
const (
	speedometerBaseGearSize       = 0.1
	speedometerBaseSpeedValueSize = 0.05
	speedometerBaseUnitSize       = 0.02
)

var (
	speedometerX    float32 = 0.895
	speedometerY    float32 = 0.881
	speedometerSize float32 = 1.0
	isMilesPerHour          = true
)

func getBikeSpeed(isMPH bool) int {
	bikeMPH := int(activeState.BikeSpeed * 2.237)
	if !isMPH {
		bikeMPH = int((float32(bikeMPH) * 1.609))
	}
	return bikeMPH
}

func getBikeGear(gear int) string {
	var outputGear string = strconv.Itoa(gear)
	if gear == 0 {
		outputGear = "N"
	}
	return outputGear
}

func getSpeedometerColor(rpm int, shiftprm int, maxrpm int) Color {
	if rpm >= maxrpm {
		return color(200, 0, 0, 159)
	}
	if rpm >= shiftprm {
		return color(200, 100, 0, 159)
	}
	return color(0, 0, 0, 159)
}

func drawSpeedometer() {
	gearSize := speedometerBaseGearSize * speedometerSize
	speedValueSize := speedometerBaseSpeedValueSize * speedometerSize
	unitSize := speedometerBaseUnitSize * speedometerSize
	speedValueOffsetX := 0.025 * speedometerSize
	speedValueOffsetY := 0.019 * speedometerSize
	unitOffsetX := 0.035 * speedometerSize
	unitOffsetY := 0.059 * speedometerSize

	baseX := speedometerX
	baseY := speedometerY
	speedColor := getSpeedometerColor(activeState.BikeRPM, activeState.BikeShiftRPM, activeState.BikeMaxRPM)

	addText(getBikeGear(activeState.BikeGear), baseX, baseY, gearSize, 0, speedColor)
	addText(fmt.Sprintf("(%d)", getBikeSpeed(isMilesPerHour)), baseX+speedValueOffsetX, baseY+speedValueOffsetY, speedValueSize, 0, speedColor)

	var speedText string
	if isMilesPerHour {
		speedText = "MPH"
	} else {
		speedText = "KPH"
	}
	addText(speedText, baseX+unitOffsetX, baseY+unitOffsetY, unitSize, 0, speedColor)
}
