package plugin

import (
	"os"
	"sort"
	"sync/atomic"
)

// Logger is the logging surface used by the plugin core.
type Logger func(format string, args ...any)

// TrackSegment is one segment of the MX Bikes track centerline.
type TrackSegment struct {
	Type   int
	Length float32
	Radius float32
	Angle  float32
	StartX float32
	StartY float32
	Height float32
}

// PositionTimePoint records the elapsed lap time at a normalized track position.
type PositionTimePoint struct {
	Position float32
	TimeMS   int
}

// RacerInfo contains the classification data needed by the leaderboard.
type RacerInfo struct {
	RaceNum        int
	Position       int
	GapSeconds     float32
	PenaltySeconds float32
	Name           string
	BikeName       string
	NumLaps        int
}

// RiderAddEntry contains the stable rider information supplied at race entry.
type RiderAddEntry struct {
	RaceNum  int
	Name     string
	BikeName string
}

// State contains all mutable game and telemetry state shared by SDK callbacks and the HUD.
type State struct {
	RaceSessionClockMS     int
	TelemetryTime          float32
	AccumulatedOnTrackTime float32
	HaveTelemetry          int
	TotalLaps              int
	LapIndex               int
	MyRaceNum              int
	PenaltyMS              int
	SessionLengthMS        int
	SessionNumLaps         int
	IsTimedPlusLaps        int
	TimeExpired            int
	ExpiryLapStart         int
	LapsAfterExpiry        int
	SessionFormat          string
	FirstZeroSeen          int
	LapAtZero              int
	RaceClassEntryLap      int

	ClientRaceNum    int
	ClientGapMS      int
	ClientClassIndex int

	BikeSpeed    float32
	BikeGear     int
	BikeRPM      int
	BikeMaxRPM   int
	BikeShiftRPM int

	CurrentTrackPosition       float32
	TrackCenterlineNumSegments int
	TrackCenterlineSegments    []TrackSegment
	TrackCenterlineRaceData    [4]float32

	TestBikeEventType           int
	SessionNumber               int
	TestRaceEventType           int
	TestRaceSessionSession      int
	TestRaceSessionStateSession int
	TestRaceLapSession          int
	TestRaceCommSession         int
	TestRaceClassSession        int

	RaceAddEntry      []int
	RaceClassEntry    []int
	RaceNumToRider    map[int]RiderAddEntry
	LeaderboardRacers []RacerInfo

	ClientBestLapTimeMS int
	ClientLastLapTimeMS int
	StopwatchMS         int
	LapStartTime        int
	IsOnTrack           bool
	HasStartedLap       bool
	ClientInPits        bool
	LastTrackPosition   float32
	HasSeenPosition     bool

	BestLapPositionData    []PositionTimePoint
	CurrentLapPositionData []PositionTimePoint
	DeltaTimeMS            int
	HasBestLapData         bool

	CurrentFuel      float32
	LapStartFuel     float32
	LastLapFuelDelta float32
	HasFuelData      bool
	MaxFuel          float32
	HasMaxFuel       bool
	PrevStopwatchMS  int
	PrevClientInPits bool

	SessionStateRaceClassification int

	SuspensionLength     [2]float32
	SuspensionVelocity   [2]float32
	SuspensionBrakePress [2]float32
	SuspensionWheelMat   [2]int
	SuspensionMaxTravel  [2]float32
	HasSuspensionData    bool
	HasSuspensionMaxData bool

	updateRelease atomic.Pointer[UpdateRelease]
	updateStarted atomic.Bool
	logger        Logger
}

// NewState creates state with the same sentinel defaults expected by the SDK lifecycle.
func NewState(logger Logger) *State {
	return &State{
		RaceSessionClockMS:             -1,
		TelemetryTime:                  -1,
		TotalLaps:                      -1,
		LapIndex:                       -1,
		MyRaceNum:                      -1,
		ExpiryLapStart:                 -1,
		LapAtZero:                      -1,
		RaceClassEntryLap:              -1,
		ClientRaceNum:                  -1,
		ClientClassIndex:               -1,
		TestBikeEventType:              -1,
		SessionNumber:                  -1,
		TestRaceEventType:              -1,
		TestRaceSessionSession:         -1,
		TestRaceSessionStateSession:    -1,
		TestRaceLapSession:             -1,
		TestRaceCommSession:            -1,
		TestRaceClassSession:           -1,
		LapStartTime:                   -1,
		SessionStateRaceClassification: -1,
		RaceNumToRider:                 make(map[int]RiderAddEntry),
		logger:                         logger,
	}
}

// SetUpdateRelease records the newer release offered by the edit menu.
func (state *State) SetUpdateRelease(release *UpdateRelease) {
	state.updateRelease.Store(release)
}

// UpdateAvailable reports whether the edit menu should offer an update.
func (state *State) UpdateAvailable() bool {
	return state.updateRelease.Load() != nil
}

// UpdateInProgress reports whether the updater is downloading or closing the game.
func (state *State) UpdateInProgress() bool {
	return state.updateStarted.Load()
}

// StartUpdate downloads and launches the update installer once per release check.
func (state *State) StartUpdate() {
	release := state.updateRelease.Load()
	if release == nil || !state.updateStarted.CompareAndSwap(false, true) {
		return
	}
	state.log("Update: Downloading PitBoardzz %s", release.Version)

	go func() {
		if err := DownloadAndLaunchUpdate(release.InstallerURL); err != nil {
			state.updateStarted.Store(false)
			state.log("Update: Failed to start update: %v", err)
			return
		}
		state.log("Update: Installer downloaded; closing MX Bikes to install update")
		// The detached installer now owns the update and will continue after the
		// host process exits, so terminate MX Bikes without waiting for shutdown.
		os.Exit(0)
	}()
}

func (state *State) log(format string, args ...any) {
	if state.logger != nil {
		state.logger(format, args...)
	}
}

// BeginBikeEvent clears state owned by the local rider's bike event.
func (state *State) BeginBikeEvent() {
	state.resetBikeEventState()
}

// EndBikeEvent clears local rider state when the bike event closes.
func (state *State) EndBikeEvent() {
	state.resetBikeEventState()
}

func (state *State) resetBikeEventState() {
	state.TestBikeEventType = -1
	state.SessionNumber = -1
	state.BikeMaxRPM = 0
	state.BikeShiftRPM = 0
	state.HasSuspensionMaxData = false
	state.SuspensionMaxTravel = [2]float32{}
	state.IsOnTrack = false
	state.LapIndex = -1
	state.ClientBestLapTimeMS = 0
	state.ClientLastLapTimeMS = 0
	state.resetBestLapState()
	state.ResetRunInitTrackState()
	state.ResetFuelState()
	state.ResetSuspensionData()
	state.ResetTelemetrySessionState()
	state.PrevClientInPits = false
}

// BeginRun clears state that belongs to one period on track.
func (state *State) BeginRun() {
	state.IsOnTrack = true
	state.ResetRunInitTrackState()
	state.ResetFuelReadings()
	state.ResetSuspensionData()
	state.TimeExpired = 0
}

// EndRun preserves accumulated testing time and clears on-track state.
func (state *State) EndRun() {
	if state.TestBikeEventType == 1 && state.TelemetryTime > 0 {
		state.AccumulatedOnTrackTime = state.TelemetryTime
	}
	state.IsOnTrack = false
	state.ResetRunTrackState()
	state.ResetFuelReadings()
	state.ResetSuspensionData()
	state.TimeExpired = 0
}

// BeginRaceEvent clears state owned by the wider race event.
func (state *State) BeginRaceEvent() {
	state.resetRaceEventState()
}

// EndRaceEvent clears race-owned data when the event closes.
func (state *State) EndRaceEvent() {
	state.resetRaceEventState()
}

func (state *State) resetRaceEventState() {
	state.TestRaceEventType = -1
	state.TestRaceSessionSession = -1
	state.TestRaceSessionStateSession = -1
	state.TestRaceLapSession = -1
	state.TestRaceCommSession = -1
	state.TestRaceClassSession = -1
	state.RaceSessionClockMS = -1
	state.TotalLaps = -1
	state.LapIndex = -1
	state.PenaltyMS = 0
	state.SessionLengthMS = 0
	state.SessionNumLaps = 0
	state.IsTimedPlusLaps = 0
	state.SessionFormat = ""
	state.RaceClassEntryLap = -1
	state.ClientRaceNum = -1
	state.MyRaceNum = -1
	state.ClientGapMS = 0
	state.ClientClassIndex = -1
	state.ClientBestLapTimeMS = 0
	state.ClientLastLapTimeMS = 0
	state.ClientInPits = false
	state.PrevClientInPits = false
	state.SessionStateRaceClassification = -1
	state.RaceAddEntry = state.RaceAddEntry[:0]
	state.resetClassificationState()
	if state.RaceNumToRider == nil {
		state.RaceNumToRider = make(map[int]RiderAddEntry)
	} else {
		clear(state.RaceNumToRider)
	}
	state.resetBestLapState()
	state.ResetTimedSessionState()
	state.ResetRunInitTrackState()
	state.ResetFuelState()
	state.ResetSuspensionData()
	state.IsOnTrack = false
}

// BeginRaceSession clears state that must not leak between race sessions.
func (state *State) BeginRaceSession() {
	state.RaceSessionClockMS = -1
	state.LapIndex = -1
	state.PenaltyMS = 0
	state.RaceClassEntryLap = -1
	state.ClientGapMS = 0
	state.ClientClassIndex = -1
	state.ClientBestLapTimeMS = 0
	state.ClientLastLapTimeMS = 0
	state.ClientInPits = false
	state.PrevClientInPits = false
	state.SessionStateRaceClassification = -1
	state.resetClassificationState()
	state.resetBestLapState()
	state.ResetTimedSessionState()
	state.ResetRunInitTrackState()
	state.ResetTelemetrySessionState()
	state.ResetFuelReadings()
	state.ResetSuspensionData()
	state.IsOnTrack = true
}

// BeginRaceClassification clears values rebuilt by one classification callback.
func (state *State) BeginRaceClassification() {
	state.ClientClassIndex = -1
	state.ClientGapMS = 0
	state.RaceClassEntryLap = -1
	state.resetClassificationState()
}

func (state *State) resetClassificationState() {
	state.RaceClassEntry = state.RaceClassEntry[:0]
	state.LeaderboardRacers = state.LeaderboardRacers[:0]
}

func (state *State) resetBestLapState() {
	state.BestLapPositionData = state.BestLapPositionData[:0]
	state.ResetCurrentLapData()
	state.HasBestLapData = false
}

func (state *State) ResetFuelState() {
	state.CurrentFuel = 0
	state.LapStartFuel = 0
	state.LastLapFuelDelta = 0
	state.HasFuelData = false
	state.MaxFuel = 0
	state.HasMaxFuel = false
}

func (state *State) ResetFuelReadings() {
	state.CurrentFuel = 0
	state.LapStartFuel = 0
	state.HasFuelData = false
}

func (state *State) ResetRunTrackState() {
	state.HasStartedLap = false
	state.StopwatchMS = 0
	state.PrevStopwatchMS = 0
	state.ResetCurrentLapData()
}

func (state *State) ResetRunInitTrackState() {
	state.ResetRunTrackState()
	state.LapStartTime = -1
	state.LastTrackPosition = 0
	state.HasSeenPosition = false
	state.CurrentTrackPosition = 0
}

func (state *State) ResetTimedSessionState() {
	state.TimeExpired = 0
	state.ExpiryLapStart = -1
	state.LapsAfterExpiry = 0
	state.FirstZeroSeen = 0
	state.LapAtZero = -1
}

func (state *State) ResetTelemetrySessionState() {
	state.AccumulatedOnTrackTime = 0
	state.TelemetryTime = -1
	state.HaveTelemetry = 0
}

func (state *State) UpdateFuel(rawFuel float32, valid bool) {
	if !valid {
		state.HasFuelData = false
		return
	}
	if rawFuel >= 0 && rawFuel <= 50 {
		state.CurrentFuel = rawFuel
		state.HasFuelData = true
		if !state.HasMaxFuel && rawFuel > 0 {
			state.MaxFuel = rawFuel
			state.HasMaxFuel = true
			state.log("updateFuel: Max fuel capacity detected: %.2f liters (new bike setup)", state.MaxFuel)
		}
		if state.LapStartFuel == 0 || (!state.HasStartedLap && rawFuel > 0) {
			state.LapStartFuel = rawFuel
		}
		return
	}
	if rawFuel < 0 {
		state.log("Warning: Negative fuel reading: %.2f", rawFuel)
	}
	state.HasFuelData = false
}

func (state *State) UpdateSuspension(length, velocity, brakePress [2]float32, wheelMat [2]int, valid bool) {
	if !valid {
		state.HasSuspensionData = false
		return
	}
	state.SuspensionLength = length
	state.SuspensionVelocity = velocity
	state.SuspensionBrakePress = brakePress
	state.SuspensionWheelMat = wheelMat
	state.HasSuspensionData = true
}

func (state *State) SetSuspensionMaxTravel(maxTravel [2]float32) {
	state.SuspensionMaxTravel = maxTravel
	state.HasSuspensionMaxData = maxTravel[0] > 0 || maxTravel[1] > 0
}

func (state *State) ResetSuspensionData() {
	state.SuspensionLength = [2]float32{}
	state.SuspensionVelocity = [2]float32{}
	state.SuspensionBrakePress = [2]float32{}
	state.SuspensionWheelMat = [2]int{}
	state.HasSuspensionData = false
}

func (state *State) UpdateBikeData(speedometer float32, gear, rpm int) {
	state.BikeSpeed = speedometer
	state.BikeGear = gear
	state.BikeRPM = rpm
}

func (state *State) UpdateRaceVehicleData(raceNum, active int, speedometer float32, gear, rpm int) {
	if active == 0 {
		return
	}

	oldClientRaceNum := state.ClientRaceNum
	oldMyRaceNum := state.MyRaceNum
	state.MyRaceNum = raceNum
	state.ClientRaceNum = raceNum
	if (oldClientRaceNum != state.ClientRaceNum && state.ClientRaceNum != -1) || (oldMyRaceNum != state.MyRaceNum && state.MyRaceNum != -1) {
		state.log("RaceVehicleData: Race numbers updated - clientRaceNum=%d, myRaceNum=%d, active=%d", state.ClientRaceNum, state.MyRaceNum, active)
	}
	state.UpdateBikeData(speedometer, gear, rpm)
	state.UpdateStopwatch()
}

func (state *State) RecordPositionData(position float32, timeMS int) {
	if !state.HasStartedLap || timeMS <= 0 {
		return
	}
	state.CurrentLapPositionData = append(state.CurrentLapPositionData, PositionTimePoint{Position: position, TimeMS: timeMS})
}

func (state *State) CalculateDelta(currentPosition float32, currentTimeMS int) int {
	if !state.HasBestLapData || len(state.BestLapPositionData) == 0 || currentTimeMS <= 0 {
		return 0
	}
	bestTimeAtPosition := state.interpolateBestLapTime(currentPosition)
	if bestTimeAtPosition <= 0 {
		return 0
	}
	return currentTimeMS - bestTimeAtPosition
}

func (state *State) interpolateBestLapTime(targetPosition float32) int {
	if len(state.BestLapPositionData) == 0 {
		return 0
	}
	if targetPosition <= state.BestLapPositionData[0].Position {
		return state.BestLapPositionData[0].TimeMS
	}
	lastPoint := state.BestLapPositionData[len(state.BestLapPositionData)-1]
	if targetPosition >= lastPoint.Position {
		return lastPoint.TimeMS
	}
	for pointIndex := 0; pointIndex < len(state.BestLapPositionData)-1; pointIndex++ {
		firstPoint := state.BestLapPositionData[pointIndex]
		secondPoint := state.BestLapPositionData[pointIndex+1]
		if targetPosition >= firstPoint.Position && targetPosition <= secondPoint.Position {
			if secondPoint.Position == firstPoint.Position {
				return firstPoint.TimeMS
			}
			ratio := (targetPosition - firstPoint.Position) / (secondPoint.Position - firstPoint.Position)
			return int(float32(firstPoint.TimeMS) + ratio*float32(secondPoint.TimeMS-firstPoint.TimeMS))
		}
	}
	return 0
}

func (state *State) UpdateBestLapData() {
	if len(state.CurrentLapPositionData) == 0 {
		return
	}
	state.BestLapPositionData = append(state.BestLapPositionData[:0], state.CurrentLapPositionData...)
	sort.Slice(state.BestLapPositionData, func(firstIndex, secondIndex int) bool {
		return state.BestLapPositionData[firstIndex].Position < state.BestLapPositionData[secondIndex].Position
	})
	state.HasBestLapData = true
}

func (state *State) UpdateStopwatch() {
	if state.HasStartedLap && state.LapStartTime >= 0 {
		currentTime := 0
		if state.HaveTelemetry != 0 && state.TelemetryTime >= 0 {
			currentTime = int(state.TelemetryTime * 1000)
		} else if state.RaceSessionClockMS >= 0 {
			currentTime = state.RaceSessionClockMS
		}
		if currentTime >= state.LapStartTime {
			newStopwatchMS := currentTime - state.LapStartTime
			if state.PrevStopwatchMS > 10000 && newStopwatchMS < state.PrevStopwatchMS/2 {
				state.log("RaceVehicleData: STOPWATCH RESET DETECTED! prev=%d, new=%d - LAP COMPLETED (RunLap handles fuel calc)", state.PrevStopwatchMS, newStopwatchMS)
			}
			state.PrevStopwatchMS = state.StopwatchMS
			state.StopwatchMS = newStopwatchMS
		}
		return
	}
	state.StopwatchMS = 0
	state.PrevStopwatchMS = 0
}

func (state *State) ResetCurrentLapData() {
	state.CurrentLapPositionData = state.CurrentLapPositionData[:0]
	state.DeltaTimeMS = 0
}
