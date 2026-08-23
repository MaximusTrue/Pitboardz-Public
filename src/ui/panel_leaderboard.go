package ui

import (
	"fmt"
	"strconv"

	"github.com/MaximusTrue/PitBoardzz/src/plugin"
)

// Leaderboard Coords
const (
	leaderboardBaseHeight     = 0.03
	leaderboardBaseRowSpacing = 0.02
	leaderboardBaseTextSize   = 0.021
	leaderboardBaseWidth      = 0.13
)

var (
	leaderboardX          float32 = 0.87
	leaderboardY          float32 = 0.45
	leaderboardSize       float32 = 1.0
	leaderboardRowSpacing float32 = leaderboardBaseRowSpacing
)

var (
	leaderboardRowY             = leaderboardY
	leaderboardBackgroundBottom = leaderboardY + leaderboardBaseHeight
)

type RacerInfo = plugin.RacerInfo
type RiderAddEntry = plugin.RiderAddEntry

func leaderboard(riders []RacerInfo, clientIndex int, racerLength int) []RacerInfo {
	n := racerLength
	if n == 0 {
		return nil
	}

	if clientIndex < 0 {
		clientIndex = 0
	} else if clientIndex >= n {
		clientIndex = n - 1
	}

	out := make([]RacerInfo, 0)
	used := make([]bool, n)

	out = append(out, riders[0])
	used[0] = true
	if n == 1 {
		return out
	}

	var start, end int
	switch clientIndex {
	case 0:
		start = 1
		end = min(n-1, 4)
	case n - 1:
		start = max(1, n-4)
		end = n - 1
	default:
		ridersAhead := clientIndex - 1
		ridersBehind := n - 1 - clientIndex
		if ridersAhead >= 2 && ridersBehind >= 2 {
			start = clientIndex - 2
			end = clientIndex + 2
		} else if ridersAhead < 2 {
			start = 1
			end = min(n-1, clientIndex+(4-ridersAhead))
		} else {
			start = max(1, clientIndex-(4-ridersBehind))
			end = n - 1
		}
	}

	for i := start; i <= end; i++ {
		if !used[i] {
			out = append(out, riders[i])
			used[i] = true
		}
	}
	return out
}

func practiceLeaderboardHeader() {
	leaderboardLabel("-------------------------", false, false)

	ms := int(activeState.TelemetryTime * 1000.0)
	if ms >= 0 {
		m := ms / 60000
		s := (ms / 1000) % 60
		leaderboardLabel(" Session Time: "+fmt.Sprintf("%02d:%02d", m, s), false, false)
	}
	leaderboardLabel(" Laps Completed: "+fmt.Sprintf("%d", activeState.RaceClassEntryLap), false, false)
}

func drawSessionDetails() {
	if activeState.SessionNumber < 5 {
		drawPracticeSessionDetails()
		return
	}

	if activeState.IsTimedPlusLaps == 1 {
		drawTimedPlusLapsSessionDetails()
		return
	}

	drawRaceSessionDetails()
}

func drawPracticeSessionDetails() {
	ms := int(activeState.TelemetryTime * 1000.0)
	if ms >= 0 {
		leaderboardLabel(" Session Time: "+formatClockMs(ms), false, false)
	} else {
		leaderboardLabel(" Session Time: --:--", false, false)
	}
	leaderboardLabel(" Laps Completed: "+fmt.Sprintf("%d", activeState.RaceClassEntryLap), false, false)
}

func drawTimedPlusLapsSessionDetails() {
	if activeState.RaceSessionClockMS > 0 {
		leaderboardLabel(" Time remaining: "+formatClockMs(activeState.RaceSessionClockMS), false, false)
		activeState.TimeExpired = 0
	} else {
		handleTimedPlusLapsExpiry()
		leaderboardLabel(" Time Expired", false, false)
	}

	if activeState.TimeExpired == 1 {
		drawTimedPlusLapsRemainingLaps()
	}
}

func handleTimedPlusLapsExpiry() {
	if activeState.FirstZeroSeen == 0 {
		activeState.FirstZeroSeen = 1
		activeState.LapAtZero = activeState.LapIndex
	}
	if activeState.TimeExpired != 0 || activeState.SessionStateRaceClassification != 16 {
		return
	}

	activeState.TimeExpired = 1
	if activeState.LapAtZero >= 0 {
		activeState.ExpiryLapStart = activeState.LapAtZero
	} else {
		activeState.ExpiryLapStart = activeState.LapIndex
	}
	activeState.LapsAfterExpiry = 0
}

func drawTimedPlusLapsRemainingLaps() {
	lapsRemaining := activeState.SessionNumLaps - activeState.LapsAfterExpiry + 1
	if lapsRemaining < 0 {
		lapsRemaining = 0
	}
	if lapsRemaining != 0 {
		leaderboardLabel(" Laps remaining: "+fmt.Sprintf("%d", lapsRemaining), false, false)
	} else {
		leaderboardLabel(" Laps remaining: Finished", false, false)
	}
}

func drawRaceSessionDetails() {
	tbuf := "00:00"
	if activeState.SessionLengthMS-activeState.RaceSessionClockMS > 0 {
		tbuf = formatClockMs(activeState.RaceSessionClockMS)
	}
	if activeState.SessionNumber != 7 {
		leaderboardLabel(" Time remaining:"+tbuf, false, false)
	}

	if activeState.TotalLaps > 0 {
		drawRaceLapsRemaining()
	}
}

func drawRaceLapsRemaining() {
	left := activeState.SessionNumLaps - activeState.RaceClassEntryLap
	if left > 0 {
		leaderboardLabel(" Laps remaining: "+fmt.Sprintf("%d", left), false, false)
	} else {
		leaderboardLabel(" Lap remaining: Finished", false, false)
	}
}

func formatClockMs(ms int) string {
	m := ms / 60000
	s := (ms / 1000) % 60
	return fmt.Sprintf("%02d:%02d", m, s)
}

func drawLeaderboardHeaderInfo() {
	const penaltyIconBaseOffsetX = 0.112
	penaltyIconOffsetX := penaltyIconBaseOffsetX * leaderboardSize
	penaltyTriangleHeight := 0.018 * leaderboardSize
	penaltyExclamationHeight := 0.016 * leaderboardSize
	penaltyExclamationThickness := 0.001 * leaderboardSize
	drawPenaltyIcon(leaderboardX+penaltyIconOffsetX, leaderboardRowY, penaltyTriangleHeight, penaltyExclamationHeight, penaltyExclamationThickness)
	leaderboardLabel(fmt.Sprintf("%19s", "Gap"), false, false)
}

func setNames(racers []RacerInfo, byRaceNum map[int]RiderAddEntry) {
	for i := range racers {
		if entry, ok := byRaceNum[racers[i].RaceNum]; ok {
			racers[i].Name = entry.Name
			racers[i].BikeName = entry.BikeName
		}
	}
}

func leaderboardLabel(val string, isClient, isFastestRider bool) {
	size := leaderboardBaseTextSize * leaderboardSize
	addText(val, leaderboardX, leaderboardRowY, size, 0, getColorLeaderboard(isClient, isFastestRider))
	leaderboardRowY += leaderboardRowSpacing
}

func getColorLeaderboard(isClient, isFastestRider bool) Color {
	if isClient {
		return color(255, 102, 0, 230)
	}
	if isFastestRider {
		return color(0, 204, 0, 230)
	}
	return color(255, 255, 255, 255)
}

func drawPenaltyIcon(x, y, triangleHeight, exclamationHeight, exclamationThickness float32) {
	addPyramidTriangle(x, y, triangleHeight, color(255, 0, 0, 255))
	exclamationOffsetY := 0.002 * leaderboardSize
	addExclamation(x+(triangleHeight*0.8*0.5), y+exclamationOffsetY, exclamationHeight, exclamationThickness, color(255, 255, 255, 255))
}

func getGapOutput(gapSeconds float32) string {
	if gapSeconds < 0 {
		return "---"
	}
	if gapSeconds == 0 {
		return strconv.FormatFloat(float64(gapSeconds), 'f', 3, 32)
	}
	return "+" + strconv.FormatFloat(float64(gapSeconds), 'f', 3, 32)
}

func getPenaltyOutput(seconds float32) string {
	if seconds < 0.1 {
		return "-"
	}
	return "+" + fmt.Sprintf("%d", int(seconds))
}

func drawLeaderboardPanel() {
	initialBackgroundBottom := leaderboardBackgroundBottom
	if initialBackgroundBottom <= leaderboardY {
		initialBackgroundBottom = leaderboardY + (leaderboardBaseHeight * leaderboardSize)
	}
	scaledWidth := leaderboardBaseWidth * leaderboardSize
	if showLeaderboardBackground {
		addQuadCCW(leaderboardX, leaderboardY, leaderboardX+scaledWidth, initialBackgroundBottom, color(0, 0, 0, 160))
	}

	leaderboardRowY = leaderboardY
	drawSessionFormat()
	if activeState.TestBikeEventType == 1 {
		drawTest()
	}
	if activeState.TestBikeEventType == 2 {
		drawRace()
	}
	leaderboardBackgroundBottom = leaderboardRowY
	if leaderboardBackgroundBottom <= leaderboardY {
		leaderboardBackgroundBottom = leaderboardY + (leaderboardBaseHeight * leaderboardSize)
	}
}

func drawSessionFormat() {
	leaderboardLabel(" Pitboardz", false, false)
	if activeState.SessionFormat == "" {
		leaderboardLabel(" Format: --", false, false)
	} else {
		leaderboardLabel(" Format: "+activeState.SessionFormat, false, false)
	}
}

func drawLeaderboard() {
	var lbR []RacerInfo = leaderboard(activeState.LeaderboardRacers, activeState.ClientClassIndex, len(activeState.RaceClassEntry))
	setNames(lbR, activeState.RaceNumToRider)
	for i := range lbR {
		racer := lbR[i]
		if i == 1 && (racer.Position != 2) {
			leaderboardLabel(fmt.Sprintf("%3s ------------------", " "), false, false)
		}
		isClient := racer.Position-1 == activeState.ClientClassIndex
		leaderboardText := fmt.Sprintf("%3s %3d %4s %7s %3s", fmt.Sprintf("%d.", racer.Position), racer.RaceNum, racer.Name, getGapOutput(racer.GapSeconds), getPenaltyOutput(racer.PenaltySeconds))
		size := leaderboardBaseTextSize * leaderboardSize
		addText(leaderboardText, leaderboardX, leaderboardRowY, size, 0, getColorLeaderboard(isClient, false))
		positionText := fmt.Sprintf("%3s", fmt.Sprintf("%d.", racer.Position))
		raceNumText := fmt.Sprintf(" %3d", racer.RaceNum)
		textBeforeTriangle := positionText + raceNumText
		triangleX := leaderboardX + textWidthMono(textBeforeTriangle, size) + (0.005 * leaderboardSize)
		triangleHeight := size * 0.59375
		triangleY := leaderboardRowY + (0.01 * leaderboardSize) - (triangleHeight * 0.125) - (0.005 * leaderboardSize)
		manufacturerColor := getBikeManufacturerColor(racer.BikeName)
		addTriangle(triangleX, triangleY, triangleHeight, manufacturerColor)
		leaderboardRowY += leaderboardRowSpacing
	}
}

func drawTest() {
	practiceLeaderboardHeader()
}

func drawRace() {
	drawSessionDetails()
	drawLeaderboardHeaderInfo()
	drawLeaderboard()
}
