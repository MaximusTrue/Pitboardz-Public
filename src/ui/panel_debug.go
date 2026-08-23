package ui

import "fmt"

// Debug Panel
var (
	debugPanelX    float32 = 0.30
	debugPanelY    float32 = 0.01
	debugPanelSize float32 = 1.0
)

func boolStr(b bool) string {
	if b {
		return "T"
	}
	return "F"
}

func drawDebugPanel() {
	textSize := float32(0.014) * debugPanelSize
	rowSpacing := float32(0.016) * debugPanelSize
	sectionSpacing := float32(0.006) * debugPanelSize
	headerColor := color(255, 204, 0, 255)
	valueColor := color(204, 204, 204, 255)
	y := debugPanelY
	x := debugPanelX

	addText("--- STATE DEBUG ---", x, y, textSize*1.1, 0, color(255, 255, 255, 255))
	y += rowSpacing + sectionSpacing

	// Session - sources: EventInit, RunInit, RaceEvent, RaceSession, RaceSessionState, RaceClassification
	addText("[Session]", x, y, textSize, 0, headerColor)
	y += rowSpacing
	addText(fmt.Sprintf("EventType:%d [EvInit]  BikeSession:%d [RunInit]", activeState.TestBikeEventType, activeState.SessionNumber), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("RaceEvType:%d [RaceEv]  RaceSes:%d [RaceSes]", activeState.TestRaceEventType, activeState.TestRaceSessionSession), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("SesState:%d [RaceSeSt]  ClassSes:%d [RaceClass]", activeState.TestRaceSessionStateSession, activeState.TestRaceClassSession), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("LapSes:%d [RaceLap]  CommSes:%d [RaceComm]", activeState.TestRaceLapSession, activeState.TestRaceCommSession), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("Format:%s  TimedPlusLaps:%d", activeState.SessionFormat, activeState.IsTimedPlusLaps), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("RaceClock:%d  SesLen:%d  NumLaps:%d", activeState.RaceSessionClockMS, activeState.SessionLengthMS, activeState.SessionNumLaps), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("SesStateRC:%d [RaceClass]", activeState.SessionStateRaceClassification), x, y, textSize, 0, valueColor)
	y += rowSpacing + sectionSpacing

	// Track - sources: RunInit, RunTelemetry, RunLap
	addText("[Track]", x, y, textSize, 0, headerColor)
	y += rowSpacing
	addText(fmt.Sprintf("OnTrack:%s  StartedLap:%s  InPits:%s", boolStr(activeState.IsOnTrack), boolStr(activeState.HasStartedLap), boolStr(activeState.ClientInPits)), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("TrackPos:%.3f  Lap:%d  TotalLaps:%d", activeState.CurrentTrackPosition, activeState.LapIndex, activeState.TotalLaps), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("Stopwatch:%d  LapStart:%d [RunLap]", activeState.StopwatchMS, activeState.LapStartTime), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("TelTime:%.3f  HaveTel:%d [RunTel]", activeState.TelemetryTime, activeState.HaveTelemetry), x, y, textSize, 0, valueColor)
	y += rowSpacing + sectionSpacing

	// Timer Expiry - sources: RaceClassification, RunLap
	addText("[Timer]", x, y, textSize, 0, headerColor)
	y += rowSpacing
	addText(fmt.Sprintf("TimeExpired:%d  ExpiryStart:%d  LapsAfter:%d", activeState.TimeExpired, activeState.ExpiryLapStart, activeState.LapsAfterExpiry), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("FirstZero:%d  LapAtZero:%d  ClassLap:%d", activeState.FirstZeroSeen, activeState.LapAtZero, activeState.RaceClassEntryLap), x, y, textSize, 0, valueColor)
	y += rowSpacing + sectionSpacing

	// Identity - source: RaceVehicleData, RaceClassification
	addText("[Identity]", x, y, textSize, 0, headerColor)
	y += rowSpacing
	addText(fmt.Sprintf("Client#:%d  My#:%d [RaceVD]", activeState.ClientRaceNum, activeState.MyRaceNum), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("ClassIdx:%d  Gap:%dms  Penalty:%dms", activeState.ClientClassIndex, activeState.ClientGapMS, activeState.PenaltyMS), x, y, textSize, 0, valueColor)
	y += rowSpacing + sectionSpacing

	// Timing - sources: RunLap, RaceLap, RaceClassification
	addText("[Timing]", x, y, textSize, 0, headerColor)
	y += rowSpacing
	addText(fmt.Sprintf("Best:%s  Last:%s", formatLapTime(activeState.ClientBestLapTimeMS), formatLapTime(activeState.ClientLastLapTimeMS)), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("Delta:%d  HasBest:%s  Pts:%d/%d", activeState.DeltaTimeMS, boolStr(activeState.HasBestLapData), len(activeState.BestLapPositionData), len(activeState.CurrentLapPositionData)), x, y, textSize, 0, valueColor)
	y += rowSpacing + sectionSpacing

	// Fuel - sources: RunTelemetry, RunLap
	addText("[Fuel]", x, y, textSize, 0, headerColor)
	y += rowSpacing
	addText(fmt.Sprintf("Fuel:%.2f  Max:%.2f  Start:%.2f [RunTel]", activeState.CurrentFuel, activeState.MaxFuel, activeState.LapStartFuel), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("Delta:%.2f  HasData:%s  HasMax:%s", activeState.LastLapFuelDelta, boolStr(activeState.HasFuelData), boolStr(activeState.HasMaxFuel)), x, y, textSize, 0, valueColor)
	y += rowSpacing + sectionSpacing

	// Bike - sources: EventInit, RunTelemetry, RaceVehicleData
	addText("[Bike]", x, y, textSize, 0, headerColor)
	y += rowSpacing
	addText(fmt.Sprintf("Speed:%.1f  Gear:%d  RPM:%d/%d  Shift:%d", activeState.BikeSpeed, activeState.BikeGear, activeState.BikeRPM, activeState.BikeMaxRPM, activeState.BikeShiftRPM), x, y, textSize, 0, valueColor)
	y += rowSpacing + sectionSpacing

	// Entries - sources: RaceAddEntry, RaceClassification
	addText("[Entries]", x, y, textSize, 0, headerColor)
	y += rowSpacing
	addText(fmt.Sprintf("Add:%d  Class:%d  LB:%d", len(activeState.RaceAddEntry), len(activeState.RaceClassEntry), len(activeState.LeaderboardRacers)), x, y, textSize, 0, valueColor)
	y += rowSpacing + sectionSpacing

	// Centerline - source: TrackCenterline
	addText("[Centerline]", x, y, textSize, 0, headerColor)
	y += rowSpacing
	addText(fmt.Sprintf("Segs:%d [TrackCL]", activeState.TrackCenterlineNumSegments), x, y, textSize, 0, valueColor)
	y += rowSpacing
	addText(fmt.Sprintf("SF:%.1f  S1:%.1f  S2:%.1f  HS:%.1f [RaceData]", activeState.TrackCenterlineRaceData[0], activeState.TrackCenterlineRaceData[1], activeState.TrackCenterlineRaceData[2], activeState.TrackCenterlineRaceData[3]), x, y, textSize, 0, valueColor)
	y += rowSpacing
	var totalLen float32
	var straightCount, curveCount int
	for _, s := range activeState.TrackCenterlineSegments {
		totalLen += s.Length
		if s.Type == 0 {
			straightCount++
		} else {
			curveCount++
		}
	}
	addText(fmt.Sprintf("TotLen:%.1fm  Str:%d  Crv:%d", totalLen, straightCount, curveCount), x, y, textSize, 0, valueColor)
	y += rowSpacing
	const maxCenterlineSegRows = 8
	shown := min(len(activeState.TrackCenterlineSegments), maxCenterlineSegRows)
	for i := range shown {
		s := activeState.TrackCenterlineSegments[i]
		addText(fmt.Sprintf("S%d: T:%d L:%.1f R:%.1f A:%.1f X:%.1f Y:%.1f H:%.1f", i, s.Type, s.Length, s.Radius, s.Angle, s.StartX, s.StartY, s.Height), x, y, textSize, 0, valueColor)
		y += rowSpacing
	}
	if len(activeState.TrackCenterlineSegments) > shown {
		addText(fmt.Sprintf("...+%d more segs stored", len(activeState.TrackCenterlineSegments)-shown), x, y, textSize, 0, valueColor)
	}
}
