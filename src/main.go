package main

/*
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

// ---------- MX Bikes SDK struct copies (must match exactly) ----------
typedef struct {
    char  m_szRiderName[100];
    char  m_szBikeID[100];
    char  m_szBikeName[100];
    int   m_iNumberOfGears;
    int   m_iMaxRPM;
    int   m_iLimiter;
    int   m_iShiftRPM;
    float m_fEngineOptTemperature;           // degrees Celsius
    float m_afEngineTemperatureAlarm[2];     // deg C. lower,upper
    float m_fMaxFuel;                        // liters
    float m_afSuspMaxTravel[2];              // meters
    float m_fSteerLock;                      // degrees
    char  m_szCategory[100];
    char  m_szTrackID[100];
    char  m_szTrackName[100];
    float m_fTrackLength;                    // meters
    int   m_iType;                           // 1=testing;2=race;4=straight rhythm
} SPluginsBikeEvent_t;

typedef struct {
    int   m_iSession;                        // testing: 0 wait;1 prog. Race: 0..6
    int   m_iConditions;                     // 0 sun;1 cloud;2 rain
    float m_fAirTemperature;                 // deg C
    char  m_szSetupFileName[100];
} SPluginsBikeSession_t;

typedef struct {
    int   m_iRPM;
    float m_fEngineTemperature;
    float m_fWaterTemperature;
    int   m_iGear;                           // 0=Neutral
    float m_fFuel;                           // liters
    float m_fSpeedometer;                    // m/s
    float m_fPosX, m_fPosY, m_fPosZ;
    float m_fVelocityX, m_fVelocityY, m_fVelocityZ;
    float m_fAccelerationX, m_fAccelerationY, m_fAccelerationZ; // G avg over 10ms
    float m_aafRot[3][3];
    float m_fYaw, m_fPitch, m_fRoll;         // deg -180..180
    float m_fYawVelocity, m_fPitchVelocity, m_fRollVelocity; // deg/s
    float m_afSuspLength[2];                 // meters
    float m_afSuspVelocity[2];               // m/s
    int   m_iCrashed;                        // 1 = rider off bike
    float m_fSteer;                          // deg (neg=right)
    float m_fThrottle;                       // 0..1
    float m_fFrontBrake;                     // 0..1
    float m_fRearBrake;                      // 0..1
    float m_fClutch;                         // 0..1 (0 = engaged)
    float m_afWheelSpeed[2];                 // m/s
    int   m_aiWheelMaterial[2];              // 0 = not in contact
    float m_afBrakePressure[2];              // kPa
    float m_fSteerTorque;                    // Nm
} SPluginsBikeData_t;

typedef struct {
    int m_iLapNum;                            // lap index
    int m_iInvalid;
    int m_iLapTime;                           // ms
    int m_iBest;                              // 1=best
} SPluginsBikeLap_t;

typedef struct {
    int m_iSplit;                             // split index
    int m_iSplitTime;                         // ms
    int m_iBestDiff;                          // ms (diff vs best)
} SPluginsBikeSplit_t;

typedef struct {
    int   m_iType;      // 0 straight; 1 curve
    float m_fLength;    // m
    float m_fRadius;    // m (<0 left; 0 straight)
    float m_fAngle;     // deg (0 north)
    float m_afStart[2]; // m
    float m_fHeight;    // m
} SPluginsTrackSegment_t;

// ---------------- race data ----------------
typedef struct {
    int   m_iType;        // 1 test;2 race;4 straight rhythm; -1 replay
    char  m_szName[100];
    char  m_szTrackName[100];
    float m_fTrackLength; // m
} SPluginsRaceEvent_t;

typedef struct {
    int   m_iRaceNum;
    char  m_szName[100];
    char  m_szBikeName[100];
    char  m_szBikeShortName[100];
    char  m_szCategory[100];
    int   m_iUnactive;
    int   m_iNumberOfGears;
    int   m_iMaxRPM;
} SPluginsRaceAddEntry_t;

typedef struct { int m_iRaceNum; } SPluginsRaceRemoveEntry_t;

typedef struct {
    int   m_iSession;          // see docs
    int   m_iSessionState;     // flags
    int   m_iSessionLength;    // ms (0 = none)
    int   m_iSessionNumLaps;
    int   m_iConditions;       // 0 sun;1 cloud;2 rain
    float m_fAirTemperature;   // deg C
} SPluginsRaceSession_t;

typedef struct {
    int m_iSession;
    int m_iSessionState;
    int m_iSessionLength;
} SPluginsRaceSessionState_t;

typedef struct {
    int m_iSession;
    int m_iRaceNum;
    int m_iLapNum;      // lap index
    int m_iInvalid;
    int m_iLapTime;     // ms
    int m_aiSplit[2];   // ms
    int m_iBest;        // 1 personal; 2 overall
} SPluginsRaceLap_t;

typedef struct {
    int m_iSession;
    int m_iRaceNum;
    int m_iLapNum;
    int m_iSplit;
    int m_iSplitTime;   // ms
} SPluginsRaceSplit_t;

typedef struct {
    int m_iSession;
    int m_iRaceNum;
    int m_iTime;
} SPluginsRaceHoleshot_t;

typedef struct {
    int m_iSession;
    int m_iRaceNum;
    int m_iCommunication; // 1 change state; 2 penalty
    int m_iState;         // 1 DNS; 2 retired; 3 DSQ
    int m_iReason;        // 0 jump; 1 offences; 2 director
    int m_iOffence;       // 1 jump; 2 cutting
    int m_iLap;
    int m_iType;          // 0 time penalty
    int m_iTime;          // ms (penalty)
} SPluginsRaceCommunication_t;

typedef struct {
    int m_iSession;
    int m_iSessionState;
    int m_iSessionTime;  // ms (current)
    int m_iNumEntries;
} SPluginsRaceClassification_t;

typedef struct {
    int m_iRaceNum;
    int m_iState;        // 1 DNS; 2 retired; 3 DSQ
    int m_iBestLap;      // ms
    int m_iBestLapNum;
    int m_iNumLaps;      // laps completed
    int m_iGap;          // ms
    int m_iGapLaps;
    int m_iPenalty;      // ms
    int m_iPit;          // 0 track;1 pits
} SPluginsRaceClassificationEntry_t;

typedef struct {
    int   m_iRaceNum;
    float m_fPosX, m_fPosY, m_fPosZ;
    float m_fYaw;        // deg from north
    float m_fTrackPos;   // 0..1
    int   m_iCrashed;
} SPluginsRaceTrackPosition_t;

typedef struct {
    int   m_iRaceNum;
    int   m_iActive;
    int   m_iRPM;
    int   m_iGear;
    float m_fSpeedometer;  // m/s
    float m_fThrottle;
    float m_fFrontBrake;
    float m_fLean;         // deg (neg=left)
} SPluginsRaceVehicleData_t;

// ---------------- draw ----------------
typedef struct {
    float m_aafPos[4][2];        // TL,BL,BR,TR (0..1)
    int   m_iSprite;             // 0 = solid color
    unsigned long m_ulColor;     // ABGR
} SPluginQuad_t;

typedef struct {
    char  m_szString[100];
    float m_afPos[2];            // 0..1
    int   m_iFont;               // 1-based
    float m_fSize;
    int   m_iJustify;            // 0 L,1 C,2 R
    unsigned long m_ulColor;     // ABGR
} SPluginString_t;

// helper to return static "mxbikes" C string for GetModID
static const char* static_modid() { return "mxbikes"; }
*/
import "C"

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"unsafe"

	"gopkg.in/ini.v1"
)

// ---------------- logger ----------------
var logFile *os.File

func writeLog(format string, a ...any) {
	if logFile == nil {
		return
	}
	fmt.Fprintf(logFile, format+"\n", a...)
	logFile.Sync()
}

func initializeLogging(savePath string) error {
	if savePath != "" {
		_ = os.MkdirAll(savePath, 0o755)
	}
	logPath := filepath.Join(savePath, "Pitboardz.log")
	if f, err := os.Create(logPath); err == nil {
		logFile = f
		writeLog("Startup: savePath=%s", savePath)
		return nil
	} else {
		return err
	}
}

func cleanupResources() {
	writeLog("Shutdown")
	if logFile != nil {
		_ = logFile.Close()
		logFile = nil
	}
}

// ---------------- state ----------------
var (
	sessionTimeMs      = -1
	telemetryTime      = float32(-1)
	haveTelemetry      = 0
	totalLaps          = -1
	lapIndex           = -1
	myRaceNum          = -1
	penaltyMs          = 0
	sessionLengthMs    = 0
	sessionNumLaps     = 0 // 0 if unlimited
	isTimedPlusLaps    = 0
	zeroStableTicks    = 0
	timeExpired        = 0
	expiryLapStart     = -1
	lapsAfterExpiry    = 0
	sessionFormat      = ""
	firstZeroSeen      = 0
	lapAtZero          = -1
	raceClassEntryLap  = -1
	lapsAfterExpired   = -1
	firstExpirySeen    = true
	removedRaceNumbers []int
)

var (
	clientRaceNum    = -1
	clientGapMs      = 0
	clientClassIndex = -1
)

var (
	bikeSpeed    = float32(0)
	bikeGear     = 0
	bikeRPM      = 0
	bikeMaxRPM   = 0
	bikeLimiter  = 0
	bikeShiftRPM = 0
)

// Track position display
var (
	currentTrackPosition = float32(0.0) // Current raw track position (0.0 - 1.0)
)

// Test variables for struct fields
var (
	testBikeEventType           = -1 // m_iType from SPluginsBikeEvent_t
	testBikeSessionSession      = -1 // m_iSession from SPluginsBikeSession_t
	testRaceEventType           = -1 // m_iType from SPluginsRaceEvent_t
	testRaceSessionSession      = -1 // m_iSession from SPluginsRaceSession_t
	testRaceSessionStateSession = -1 // m_iSession from SPluginsRaceSessionState_t
	testRaceLapSession          = -1 // m_iSession from SPluginsRaceLap_t
	testRaceCommSession         = -1 // m_iSession from SPluginsRaceCommunication_t
	testRaceClassSession        = -1 // m_iSession from SPluginsRaceClassification_t
)

// Leaderboard
var (
	raceAddEntry     []int
	raceClassEntry   []int
	raceCommEntry    []int
	raceAddEntryName []RiderAddEntry
)

var leaderboardRacers []RacerInfo

// Lap Time
var (
	clientBestLapTimeMS       int
	clientLastLapTimeMS       int
	ridersFastestLapTimeMS    int
	riderFastestLapRaceNumber int
	stopwatchMS               int
	lapStartTime              int     // Time when current lap started (ms)
	lastRaceClassEntryLap     int     // Previous lap number to detect lap changes
	isOnTrack                 bool    // Whether client is currently on track (not in pits)
	hasStartedLap             bool    // Whether we've started timing a lap
	clientInPits              bool    // Whether client is currently in pits
	lastTrackPosition         float32 // Previous track position (0-1) to detect finish line crossing
	hasSeenPosition           bool    // Whether we've received position data yet
)

// Delta Timing - Track position vs time for best lap comparison
var (
	bestLapPositionData    []PositionTimePoint // Best lap position/time data points
	currentLapPositionData []PositionTimePoint // Current lap position/time data points
	deltaTimeMS            int                 // Current delta vs best lap (+ = slower, - = faster)
	hasBestLapData         bool                // Whether we have valid best lap data
	isRecordingBestLap     bool                // Whether current lap might become new best
)

// Fuel tracking
var (
	currentFuel      float32 // Current fuel level in liters
	lapStartFuel     float32 // Fuel level at start of current lap
	lastLapFuelDelta float32 // Fuel consumed in last completed lap
	hasFuelData      bool    // Whether we have valid fuel data
	maxFuel          float32 // Maximum fuel capacity (detected from first track entry)
	hasMaxFuel       bool    // Whether we have detected max fuel capacity
	prevStopwatchMS  int     // Previous stopwatch value to detect resets
	prevClientInPits bool    // Previous pit status to detect pit exit transitions
)

// ---------------- types ----------------
type RacerInfo struct {
	raceNum         int
	position        int
	gapSeconds      float32
	penalitySeconds float32
	name            string
	bikeName        string
}

type RiderAddEntry struct {
	raceNum  int
	name     string
	bikeName string
}

// Delta timing data point
type PositionTimePoint struct {
	position float32 // Track position (0.0 - 1.0)
	timeMS   int     // Time from lap start in milliseconds
}

// ---------------- hooks (low-boilerplate feature system) ----------------
type Ctx struct{}

func newCtx() *Ctx { return &Ctx{} }

func (c *Ctx) StopwatchMS() int       { return stopwatchMS }
func (c *Ctx) BestLapTimeMS() int     { return clientBestLapTimeMS }
func (c *Ctx) LastLapTimeMS() int     { return clientLastLapTimeMS }
func (c *Ctx) DeltaTimeMS() int       { return deltaTimeMS }
func (c *Ctx) HasBestLapData() bool   { return hasBestLapData }
func (c *Ctx) TrackPosition() float32 { return currentTrackPosition }
func (c *Ctx) IsOnTrack() bool        { return isOnTrack }
func (c *Ctx) HasStartedLap() bool    { return hasStartedLap }

func (c *Ctx) BikeSpeedMS() float32 { return bikeSpeed }
func (c *Ctx) BikeRPM() int         { return bikeRPM }
func (c *Ctx) BikeMaxRPM() int      { return bikeMaxRPM }
func (c *Ctx) BikeShiftRPM() int    { return bikeShiftRPM }
func (c *Ctx) BikeGear() int        { return bikeGear }

func (c *Ctx) Fuel() (current float32, lastLapDelta float32, max float32, has bool) {
	return currentFuel, lastLapFuelDelta, maxFuel, hasFuelData
}

func (c *Ctx) Leaderboard() []RacerInfo { return leaderboardRacers }
func (c *Ctx) ClientIndex() int         { return clientClassIndex }

type Painter struct{}

func NewPainter() *Painter { return &Painter{} }

func (p *Painter) Text(s string, x, y, size float32, justify int, abgr uint32) {
	addText(s, x, y, size, justify, abgr)
}
func (p *Painter) Rect(x0, y0, x1, y1 float32, abgr uint32) { addQuadCCW(x0, y0, x1, y1, abgr) }
func (p *Painter) Triangle(x, y, h float32, abgr uint32)    { addTriangle(x, y, h, abgr) }

type (
	DrawHook               interface{ OnDraw(p *Painter, ctx *Ctx) }
	TelemetryHook          interface{ OnTelemetry(ctx *Ctx) }
	RaceClassificationHook interface{ OnRaceClassification(ctx *Ctx) }
)

var registeredFeatures []any

func RegisterFeature(f any) { registeredFeatures = append(registeredFeatures, f) }

func dispatchDraw(p *Painter, ctx *Ctx) {
	for i := range registeredFeatures {
		if h, ok := registeredFeatures[i].(DrawHook); ok {
			h.OnDraw(p, ctx)
		}
	}
}

func dispatchTelemetry(ctx *Ctx) {
	for i := range registeredFeatures {
		if h, ok := registeredFeatures[i].(TelemetryHook); ok {
			h.OnTelemetry(ctx)
		}
	}
}

func dispatchRaceClassification(ctx *Ctx) {
	for i := range registeredFeatures {
		if h, ok := registeredFeatures[i].(RaceClassificationHook); ok {
			h.OnRaceClassification(ctx)
		}
	}
}

// ---------------- logic ----------------
func formatMinutes(ms int) string {
	m := ms / 60000
	return fmt.Sprintf("%dm", m)
}

// leaderboard returns up to 6 riders, always including P1.
// It then shows a window around the client with robust edge handling.
func leaderboard(riders []RacerInfo, clientIndex int, racerLength int) []RacerInfo {
	n := racerLength
	if n == 0 {
		return nil
	}

	// Clamp client index
	if clientIndex < 0 {
		clientIndex = 0
	} else if clientIndex >= n {
		clientIndex = n - 1
	}

	out := make([]RacerInfo, 0)
	used := make([]bool, n)

	// Always include first place
	out = append(out, riders[0])
	used[0] = true
	if n == 1 {
		return out
	}

	// Determine the window around the client
	var start, end int
	switch clientIndex {
	case 0:
		start = 1
		end = min(n-1, 4)
	case n - 1:
		start = max(1, n-4)
		end = n - 1
	default:
		ridersAhead := clientIndex - 1 // excluding first place
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

	ms := int(telemetryTime * 1000.0)
	if ms >= 0 {
		m := ms / 60000
		s := (ms / 1000) % 60
		leaderboardLabel(" Session Time: "+fmt.Sprintf("%02d:%02d", m, s), false, false)
	}
	leaderboardLabel(" Laps Completed: "+fmt.Sprintf("%d", raceClassEntryLap), false, false)
}

func leaderboardHeader() {
	remainingMs := sessionLengthMs - sessionTimeMs

	if testBikeSessionSession >= 5 {
		if isTimedPlusLaps == 1 {
			if sessionTimeMs > 0 {
				zeroStableTicks = 0
				m := sessionTimeMs / 60000
				s := (sessionTimeMs / 1000) % 60
				leaderboardLabel(" Time remaining: "+fmt.Sprintf("%02d:%02d", m, s), false, false)
				timeExpired = 0
			} else {
				if firstZeroSeen == 0 {
					firstZeroSeen = 1
					lapAtZero = lapIndex
				}
				if timeExpired == 0 && sessionStateRaceClassification == 16 {
					timeExpired = 1
					if lapAtZero >= 0 {
						expiryLapStart = lapAtZero
					} else {
						expiryLapStart = lapIndex
					}
					lapsAfterExpiry = 0
				}
				leaderboardLabel(" Time Expired", false, false)
			}

			// Laps Display
			if timeExpired == 1 {
				lapsRemaining := sessionNumLaps - lapsAfterExpiry + 1
				if lapsRemaining < 0 {
					lapsRemaining = 0
				}
				if lapsRemaining != 0 {
					leaderboardLabel(" Laps remaining: "+fmt.Sprintf("%d", lapsRemaining), false, false)
				} else {
					leaderboardLabel(" Laps remaining: Finished", false, false)
				}
			}
		} else {
			var tbuf string
			if remainingMs > 0 {
				m := sessionTimeMs / 60000
				s := (sessionTimeMs / 1000) % 60
				tbuf = fmt.Sprintf("%02d:%02d", m, s)
			} else {
				tbuf = "00:00"
			}
			if testBikeSessionSession != 7 {
				leaderboardLabel(" Time remaining:"+tbuf, false, false)
			}

			// Laps Display
			if totalLaps > 0 {
				left := sessionNumLaps - raceClassEntryLap
				if left > 0 {
					leaderboardLabel(" Laps remaining: "+fmt.Sprintf("%d", left), false, false)
				} else {
					leaderboardLabel(" Lap remaining: Finished", false, false)
				}
			}
		}
	} else {
		ms := int(telemetryTime * 1000.0)
		if ms >= 0 {
			m := ms / 60000
			s := (ms / 1000) % 60
			leaderboardLabel(" Session Time: "+fmt.Sprintf("%02d:%02d", m, s), false, false)
		} else {
			leaderboardLabel(" Session Time: --:--", false, false)
		}
		leaderboardLabel(" Laps Completed: "+fmt.Sprintf("%d", raceClassEntryLap), false, false)
	}
}

func drawLeaderboardDetails() {
	const penaltyIconBaseOffsetX = 0.112
	penaltyIconOffsetX := penaltyIconBaseOffsetX * leaderboardSize
	penaltyTriangleHeight := 0.018 * leaderboardSize
	penaltyExclamationHeight := 0.016 * leaderboardSize
	penaltyExclamationThickness := 0.001 * leaderboardSize
	drawPenaltyIcon(leaderboardX+penaltyIconOffsetX, leaderboardRowY, penaltyTriangleHeight, penaltyExclamationHeight, penaltyExclamationThickness)
	leaderboardLabel(fmt.Sprintf("%19s", "Gap"), false, false)
}

func setNames(racers []RacerInfo, names []RiderAddEntry) {
	for i := range racers {
		for j := range names {
			if racers[i].raceNum == names[j].raceNum {
				racers[i].name = names[j].name
				racers[i].bikeName = names[j].bikeName
				continue
			}
		}
	}
}

// ---------------- delta timing logic ----------------
func recordPositionData(position float32, timeMS int) {
	if !hasStartedLap || timeMS <= 0 {
		return
	}
	point := PositionTimePoint{position: position, timeMS: timeMS}
	currentLapPositionData = append(currentLapPositionData, point)
}

func calculateDelta(currentPosition float32, currentTimeMS int) int {
	if !hasBestLapData || len(bestLapPositionData) == 0 || currentTimeMS <= 0 {
		return 0
	}
	bestTimeAtPosition := interpolateBestLapTime(currentPosition)
	if bestTimeAtPosition <= 0 {
		return 0
	}
	return currentTimeMS - bestTimeAtPosition
}

func interpolateBestLapTime(targetPosition float32) int {
	if len(bestLapPositionData) == 0 {
		return 0
	}
	if targetPosition <= bestLapPositionData[0].position {
		return bestLapPositionData[0].timeMS
	}
	if targetPosition >= bestLapPositionData[len(bestLapPositionData)-1].position {
		return bestLapPositionData[len(bestLapPositionData)-1].timeMS
	}
	for i := 0; i < len(bestLapPositionData)-1; i++ {
		p1 := bestLapPositionData[i]
		p2 := bestLapPositionData[i+1]
		if targetPosition >= p1.position && targetPosition <= p2.position {
			if p2.position == p1.position {
				return p1.timeMS
			}
			ratio := (targetPosition - p1.position) / (p2.position - p1.position)
			interpolatedTime := float32(p1.timeMS) + ratio*float32(p2.timeMS-p1.timeMS)
			return int(interpolatedTime)
		}
	}
	return 0
}

func updateBestLapData() {
	if len(currentLapPositionData) == 0 {
		return
	}
	bestLapPositionData = make([]PositionTimePoint, len(currentLapPositionData))
	copy(bestLapPositionData, currentLapPositionData)
	for i := 0; i < len(bestLapPositionData)-1; i++ {
		for j := 0; j < len(bestLapPositionData)-1-i; j++ {
			if bestLapPositionData[j].position > bestLapPositionData[j+1].position {
				temp := bestLapPositionData[j]
				bestLapPositionData[j] = bestLapPositionData[j+1]
				bestLapPositionData[j+1] = temp
			}
		}
	}
	hasBestLapData = true
}

func resetCurrentLapData() {
	currentLapPositionData = currentLapPositionData[:0]
	deltaTimeMS = 0
}

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

// ---------------- drawing ----------------
const (
	maxQuads   = 48
	maxStrings = 32
)

var (
	cQuads    unsafe.Pointer
	cStrings  unsafe.Pointer
	quadCount int
	strCount  int
	fontReady = false
	cFontList *C.char
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
	leaderboardSize       float32 = 1.0 // size multiplier for entire leaderboard
	leaderboardRowSpacing float32 = leaderboardBaseRowSpacing
)

var (
	leaderboardRowY             = leaderboardY
	leaderboardBackgroundBottom = leaderboardY + leaderboardBaseHeight
)

const (
	CHAR_ASPECT = 0.22
)

func textWidthMono(s string, size float32) float32 { return float32(len(s)) * (size * CHAR_ASPECT) }

func quadAt(i int) *C.SPluginQuad_t {
	return (*C.SPluginQuad_t)(unsafe.Add(cQuads, i*int(unsafe.Sizeof(C.SPluginQuad_t{}))))
}

func strAt(i int) *C.SPluginString_t {
	return (*C.SPluginString_t)(unsafe.Add(cStrings, i*int(unsafe.Sizeof(C.SPluginString_t{}))))
}

func addQuadCCW(x0, y0, x1, y1 float32, abgr uint32) {
	if quadCount >= maxQuads {
		return
	}
	q := quadAt(quadCount)
	quadCount++
	q.m_aafPos[0][0] = C.float(x0)
	q.m_aafPos[0][1] = C.float(y0)
	q.m_aafPos[1][0] = C.float(x0)
	q.m_aafPos[1][1] = C.float(y1)
	q.m_aafPos[2][0] = C.float(x1)
	q.m_aafPos[2][1] = C.float(y1)
	q.m_aafPos[3][0] = C.float(x1)
	q.m_aafPos[3][1] = C.float(y0)
	q.m_iSprite = C.int(0)
	q.m_ulColor = C.ulong(abgr)
}

func setCStringFixed100(dst *[100]C.char, s string) {
	C.memset(unsafe.Pointer(dst), 0, C.size_t(100))
	if len(s) == 0 {
		return
	}
	if len(s) > 99 {
		s = s[:99]
	}
	b := []byte(s)
	C.memcpy(unsafe.Pointer(dst), unsafe.Pointer(&b[0]), C.size_t(len(b)))
}

func addText(s string, x, y, size float32, justify int, abgr uint32) {
	if strCount >= maxStrings {
		return
	}
	t := strAt(strCount)
	strCount++
	setCStringFixed100((*[100]C.char)(unsafe.Pointer(&t.m_szString[0])), s)
	t.m_afPos[0] = C.float(x)
	t.m_afPos[1] = C.float(y)
	t.m_iFont = C.int(1)
	t.m_fSize = C.float(size)
	t.m_iJustify = C.int(justify)
	t.m_ulColor = C.ulong(abgr)
}

func leaderboardLabel(val string, isClient, isFastestRider bool) {
	size := leaderboardBaseTextSize * leaderboardSize
	addText(val, leaderboardX, leaderboardRowY, size, 0, getColorLeaderboard(isClient, isFastestRider))
	leaderboardRowY += leaderboardRowSpacing
}

func getColorLeaderboard(isClient, isFastestRider bool) uint32 {
	if isClient {
		return 0xE60066FF
	}
	if isFastestRider {
		return 0xE600CC00
	}
	return 0xFFFFFFFF
}

func getBikeManufacturerColor(bikeName string) uint32 {
	name := bikeName
	if len(name) > 0 {
		name = bikeName[:min(len(bikeName), 20)]
	}
	switch {
	case containsIgnoreCase(name, "Fantic"):
		return 0xFF800020
	case containsIgnoreCase(name, "GasGas"):
		return 0xFFFFFFFF
	case containsIgnoreCase(name, "Honda"):
		return 0xFF0000FF
	case containsIgnoreCase(name, "Husqvarna"):
		return 0xFFE0E0E0
	case containsIgnoreCase(name, "Kawasaki"):
		return 0xFF00FF00
	case containsIgnoreCase(name, "KTM"):
		return 0xFF0080FF
	case containsIgnoreCase(name, "Suzuki"):
		return 0xFF00FFFF
	case containsIgnoreCase(name, "TM MX"):
		return 0xFFFFB000
	case containsIgnoreCase(name, "Triumph"):
		return 0xFF3FE3FF
	case containsIgnoreCase(name, "Yamaha"):
		return 0xFFFF0000
	case containsIgnoreCase(name, "Beta"):
		return 0xFF0000C0
	default:
		return 0xFF808080
	}
}

func containsIgnoreCase(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			c1 := s[i+j]
			c2 := substr[j]
			if c1 >= 'A' && c1 <= 'Z' {
				c1 = c1 + 32
			}
			if c2 >= 'A' && c2 <= 'Z' {
				c2 = c2 + 32
			}
			if c1 != c2 {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

func addTriangle(x, y, height float32, color uint32) {
	if quadCount >= maxQuads {
		return
	}
	halfHeight := height / 2
	width := height * 0.5
	x1, y1 := x+width, y+halfHeight
	x2, y2 := x, y
	x3, y3 := x, y+height
	q := quadAt(quadCount)
	quadCount++
	q.m_aafPos[0][0] = C.float(x1)
	q.m_aafPos[0][1] = C.float(y1)
	q.m_aafPos[1][0] = C.float(x2)
	q.m_aafPos[1][1] = C.float(y2)
	q.m_aafPos[2][0] = C.float(x3)
	q.m_aafPos[2][1] = C.float(y3)
	q.m_aafPos[3][0] = C.float(x1)
	q.m_aafPos[3][1] = C.float(y1)
	q.m_iSprite = C.int(0)
	q.m_ulColor = C.ulong(color)
}

func addPyramidTriangle(x, y, height float32, color uint32) {
	if quadCount >= maxQuads {
		return
	}

	// Pick a base width. For an isosceles look, width = height.
	// For an *equilateral* look, use: width := 2*height/float32(math.Sqrt(3))
	width := height * 0.8

	// Vertices for an upward-pointing triangle:
	// apex at the middle of the top edge, base along the bottom
	ax, ay := x+width*0.5, y    // apex (top-center)
	lx, ly := x, y+height       // left base
	rx, ry := x+width, y+height // right base

	q := quadAt(quadCount)
	quadCount++

	// Fill the “quad” slot with 3 triangle verts + repeat the first to close
	q.m_aafPos[0][0] = C.float(ax)
	q.m_aafPos[0][1] = C.float(ay)
	q.m_aafPos[1][0] = C.float(lx)
	q.m_aafPos[1][1] = C.float(ly)
	q.m_aafPos[2][0] = C.float(rx)
	q.m_aafPos[2][1] = C.float(ry)
	q.m_aafPos[3][0] = C.float(ax)
	q.m_aafPos[3][1] = C.float(ay)

	q.m_iSprite = C.int(0)
	q.m_ulColor = C.ulong(color)
}

func addExclamation(x, y, height, thickness float32, color uint32) {
	if quadCount+2 >= maxQuads {
		return
	}

	// Define proportions
	lineHeight := height * 0.6    // top line
	squareSize := thickness * 1.5 // bottom block (dot)

	// --- Draw vertical line ---
	// Center it on x
	x0 := x - thickness/2
	x1 := x + thickness/2
	y0 := y
	y1 := y + lineHeight

	qLine := quadAt(quadCount)
	quadCount++
	qLine.m_aafPos[0][0], qLine.m_aafPos[0][1] = C.float(x0), C.float(y0)
	qLine.m_aafPos[1][0], qLine.m_aafPos[1][1] = C.float(x0), C.float(y1)
	qLine.m_aafPos[2][0], qLine.m_aafPos[2][1] = C.float(x1), C.float(y1)
	qLine.m_aafPos[3][0], qLine.m_aafPos[3][1] = C.float(x1), C.float(y0)
	qLine.m_iSprite = C.int(0)
	qLine.m_ulColor = C.ulong(color)

	// --- Draw bottom square (the “dot”) ---
	ySquareTop := y + lineHeight + (thickness * 2) // small gap
	ySquareBottom := ySquareTop + (squareSize * 2)
	xSquareLeft := x - squareSize/2
	xSquareRight := x + squareSize/2

	qSquare := quadAt(quadCount)
	quadCount++
	qSquare.m_aafPos[0][0], qSquare.m_aafPos[0][1] = C.float(xSquareLeft), C.float(ySquareTop)
	qSquare.m_aafPos[1][0], qSquare.m_aafPos[1][1] = C.float(xSquareLeft), C.float(ySquareBottom)
	qSquare.m_aafPos[2][0], qSquare.m_aafPos[2][1] = C.float(xSquareRight), C.float(ySquareBottom)
	qSquare.m_aafPos[3][0], qSquare.m_aafPos[3][1] = C.float(xSquareRight), C.float(ySquareTop)
	qSquare.m_iSprite = C.int(0)
	qSquare.m_ulColor = C.ulong(color)
}

func drawPenaltyIcon(x, y, triangleHeight, exclamationHeight, exclamationThickness float32) {
	addPyramidTriangle(x, y, triangleHeight, 0xFF0000FF)
	// Scale the vertical offset for exclamation mark positioning
	exclamationOffsetY := 0.002 * leaderboardSize
	addExclamation(x+(triangleHeight*0.8*0.5), y+exclamationOffsetY, exclamationHeight, exclamationThickness, 0xFFFFFFFF)
}

func getBikeSpeed(isMPH bool) int {
	bikeMPH := int(bikeSpeed * 2.237)
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

func getSpeedometerColor(rpm int, shiftprm int, maxrpm int) uint32 {
	if rpm >= maxrpm {
		return 0x9F0000C8
	}
	if rpm >= shiftprm {
		return 0x9F0064C8
	}
	return 0x9F000000
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

func formatLapTime(laptimeMS int) string {
	minutes := laptimeMS / 60000
	seconds := (laptimeMS % 60000) / 1000
	milliseconds := laptimeMS % 1000
	return fmt.Sprintf("%02d:%02d.%03d", minutes, seconds, milliseconds)
}

func allocCString(s string) *C.char {
	b := append([]byte(s), 0)
	ptr := C.malloc(C.size_t(len(b)))
	if ptr == nil {
		return nil
	}
	C.memcpy(ptr, unsafe.Pointer(&b[0]), C.size_t(len(b)))
	return (*C.char)(ptr)
}

func drawTest() {
	practiceLeaderboardHeader()
}

func drawRace() {
	var lbR []RacerInfo = leaderboard(leaderboardRacers, clientClassIndex, len(raceClassEntry))
	setNames(lbR, raceAddEntryName)
	leaderboardHeader()
	drawLeaderboardDetails() // Leaderboard Info Above Riders
	for i := range lbR {
		racer := lbR[i]
		if i == 1 && (racer.position != 2) {
			leaderboardLabel(fmt.Sprintf("%3s ------------------", " "), false, false)
		}
		isClient := racer.position-1 == clientClassIndex
		leaderboardText := fmt.Sprintf("%3s %3d %4s %7s %3s", fmt.Sprintf("%d.", racer.position), racer.raceNum, racer.name, getGapOutput(racer.gapSeconds), getPenaltyOutput(racer.penalitySeconds))
		size := leaderboardBaseTextSize * leaderboardSize
		addText(leaderboardText, leaderboardX, leaderboardRowY, size, 0, getColorLeaderboard(isClient, false))
		positionText := fmt.Sprintf("%3s", fmt.Sprintf("%d.", racer.position))
		raceNumText := fmt.Sprintf(" %3d", racer.raceNum)
		textBeforeTriangle := positionText + raceNumText
		triangleX := leaderboardX + textWidthMono(textBeforeTriangle, size) + (0.005 * leaderboardSize)
		triangleHeight := size * 0.59375
		triangleY := leaderboardRowY + (0.01 * leaderboardSize) - (triangleHeight * 0.125) - (0.005 * leaderboardSize)
		manufacturerColor := getBikeManufacturerColor(racer.bikeName)
		addTriangle(triangleX, triangleY, triangleHeight, manufacturerColor)
		leaderboardRowY += leaderboardRowSpacing
	}
}

func getDeltaColor() uint32 {
	if !hasBestLapData || !hasStartedLap || deltaTimeMS == 0 {
		return 0xFFFFFFFF
	}
	if deltaTimeMS < 0 {
		return 0xFF00FF00
	}
	return 0xFF0000FF
}

// Stopwatch Coords
const (
	stopwatchBaseTextSize   = 0.021
	stopwatchBaseRowSpacing = 0.031
)

var (
	timingX0            float32 = 0.02
	timingY0            float32 = 0.55
	stopwatchSize       float32 = 1.0 // size multiplier for entire stopwatch
	stopwatchTextSize   float32 = stopwatchBaseTextSize
	stopwatchRowSpacing float32 = stopwatchBaseRowSpacing
	timingRowY                  = timingY0
)

func drawTimingPanel() {
	timingRowY = timingY0
	bestText := "Best: " + formatLapTime(clientBestLapTimeMS)
	addText(bestText, timingX0, timingRowY, stopwatchTextSize, 0, 0xFFFFFFFF)
	timingRowY += stopwatchRowSpacing
	lastText := "Last: " + formatLapTime(clientLastLapTimeMS)
	addText(lastText, timingX0, timingRowY, stopwatchTextSize, 0, 0xFFFFFFFF)
	if clientLastLapTimeMS > 0 && clientBestLapTimeMS > 0 && clientLastLapTimeMS != clientBestLapTimeMS {
		lastDelta := clientLastLapTimeMS - clientBestLapTimeMS
		lastDeltaText := formatDelta(lastDelta)
		lastTextWidth := textWidthMono(lastText, stopwatchTextSize)
		deltaX := timingX0 + lastTextWidth + (0.015 * stopwatchSize)
		addText(lastDeltaText, deltaX, timingRowY, stopwatchTextSize, 0, 0xFF0000FF)
	}
	timingRowY += stopwatchRowSpacing
	currentText := formatLapTime(stopwatchMS)
	prefixWidth := textWidthMono("Best: ", stopwatchTextSize)
	currentX := timingX0 + prefixWidth
	addText(currentText, currentX, timingRowY, stopwatchTextSize, 0, 0xFFFFFFFF)
	deltaText := formatDelta(deltaTimeMS)
	currentWidth := textWidthMono(currentText, stopwatchTextSize)
	deltaWidth := textWidthMono(deltaText, stopwatchTextSize*0.9)
	deltaY := timingRowY + (0.022 * stopwatchSize)
	deltaX := currentX + currentWidth - deltaWidth
	addText(deltaText, deltaX, deltaY, stopwatchTextSize*0.9, 0, getDeltaColor())
}

// Fuel Info Coords
const (
	fuelBaseTextSize   = 0.025
	fuelBaseRowSpacing = 0.03
)

var (
	fuelX          = float32(0.75)
	fuelY          = float32(0.02)
	fuelSize       = float32(1.0) // size multiplier for fuel info
	fuelTextSize   = float32(fuelBaseTextSize)
	fuelRowSpacing = float32(fuelBaseRowSpacing)
)

func gasTank() {
	var fuelText string
	var fuelColor uint32 = 0xFFFFFFFF
	if !hasFuelData {
		fuelText = "Fuel: ---"
	} else {
		fuelText = fmt.Sprintf("Fuel: %.2f L", currentFuel)
		if currentFuel < 1.0 {
			fuelColor = 0xFF0000FF
		} else if currentFuel < 2.0 {
			fuelColor = 0xFF00FFFF
		}
	}
	addText(fuelText, fuelX, fuelY, fuelTextSize, 0, fuelColor)
	var deltaText string
	if lastLapFuelDelta == 0.0 {
		deltaText = "Last: ---"
	} else {
		deltaText = fmt.Sprintf("Last: %.2f L", lastLapFuelDelta)
	}
	addText(deltaText, fuelX, fuelY+fuelRowSpacing, fuelTextSize, 0, 0xFFFFFFFF)
}

func drawSessionFormat() {
	leaderboardLabel(" Pitboardz", false, false)
	if sessionFormat == "" {
		leaderboardLabel(" Format: --", false, false)
	} else {
		leaderboardLabel(" Format: "+sessionFormat, false, false)
	}
}

// Speedometer Coords
const (
	speedometerBaseGearSize       = 0.1
	speedometerBaseSpeedValueSize = 0.05
	speedometerBaseUnitSize       = 0.02
)

var (
	speedometerX    float32 = 0.895
	speedometerY    float32 = 0.881
	speedometerSize float32 = 1.0 // size multiplier for speedometer
	isMilesPerHour          = true
)

// Display toggles
var (
	showStopwatch   bool = true
	showLeaderboard bool = true
	showSpeedometer bool = true
	showFuelInfo    bool = true
)

// Session state indicators
var (
	sessionStateRaceClassification int = -1 // from SPluginsRaceClassification_t
)

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
	speedColor := getSpeedometerColor(bikeRPM, bikeShiftRPM, bikeMaxRPM)

	addText(getBikeGear(bikeGear), baseX, baseY, gearSize, 0, speedColor)
	addText(fmt.Sprintf("(%d)", getBikeSpeed(isMilesPerHour)), baseX+speedValueOffsetX, baseY+speedValueOffsetY, speedValueSize, 0, speedColor)

	var speedText string
	if isMilesPerHour {
		speedText = "MPH"
	} else {
		speedText = "KPH"
	}
	addText(speedText, baseX+unitOffsetX, baseY+unitOffsetY, unitSize, 0, speedColor)
}

func drawCommon() {
	// addPyramidTriangle(0.1, 0.15, 0.2, 0xFF0000FF)
	if showSpeedometer {
		drawSpeedometer()
	}
	// drawTrackPosition()
	if showStopwatch {
		drawTimingPanel()
	}
	if showFuelInfo {
		gasTank()
	}
	// Hooks draw after core widgets
	dispatchDraw(NewPainter(), newCtx())
}

// ---------------- exports ----------------
//
//export GetModID
func GetModID() *C.char { return (*C.char)(C.static_modid()) }

//export GetModDataVersion
func GetModDataVersion() C.int { return 8 }

//export GetInterfaceVersion
func GetInterfaceVersion() C.int { return 9 }

var (
	iniSavePath         string
	iniLoadErrorLogged  bool
	iniValueErrorLogged = make(map[string]bool)
)

//export Startup
func Startup(szSavePath *C.char) C.int {
	savePath := ""
	if szSavePath != nil {
		savePath = C.GoString(szSavePath)
	}
	iniSavePath = C.GoString(szSavePath)
	_ = initializeLogging(savePath)
	return 0
}

//export Shutdown
func Shutdown() {
	cleanupResources()
	if cQuads != nil {
		C.free(cQuads)
		cQuads = nil
	}
	if cStrings != nil {
		C.free(cStrings)
		cStrings = nil
	}
	if cFontList != nil {
		C.free(unsafe.Pointer(cFontList))
		cFontList = nil
	}
}

//export EventInit
func EventInit(p unsafe.Pointer, size C.int) {
	data := (*C.SPluginsBikeEvent_t)(p)
	bikeMaxRPM = int(data.m_iMaxRPM)
	bikeLimiter = int(data.m_iLimiter)
	bikeShiftRPM = int(data.m_iShiftRPM)
	testBikeEventType = int(data.m_iType)
	currentFuel = 0.0
	lapStartFuel = 0.0
	lastLapFuelDelta = 0.0
	hasFuelData = false
	maxFuel = 0.0
	hasMaxFuel = false
	prevStopwatchMS = 0
	prevClientInPits = false
	writeLog("EventInit: Event initialized, fuel tracking state reset")
}

//export EventDeinit
func EventDeinit() {
	currentFuel = 0.0
	lapStartFuel = 0.0
	lastLapFuelDelta = 0.0
	hasFuelData = false
	maxFuel = 0.0
	hasMaxFuel = false
	prevStopwatchMS = 0
	prevClientInPits = false
	writeLog("EventDeinit: Event closed, fuel tracking state reset")
}

//export RunInit
func RunInit(p unsafe.Pointer, size C.int) {
	if p != nil && int(size) >= int(unsafe.Sizeof(C.SPluginsBikeSession_t{})) {
		data := (*C.SPluginsBikeSession_t)(p)
		testBikeSessionSession = int(data.m_iSession)
	}
	isOnTrack = true
	hasStartedLap = false
	stopwatchMS = 0
	prevStopwatchMS = 0
	lapStartTime = -1
	lastTrackPosition = 0.0
	hasSeenPosition = false
	currentTrackPosition = 0.0
	resetCurrentLapData()
	currentFuel = 0.0
	lapStartFuel = 0.0
	hasFuelData = false
	timeExpired = 0
	writeLog("RunInit: Client entered track, isOnTrack=%t, hasStartedLap=%t", isOnTrack, hasStartedLap)
}

//export RunDeinit
func RunDeinit() {
	isOnTrack = false
	hasStartedLap = false
	stopwatchMS = 0
	prevStopwatchMS = 0
	resetCurrentLapData()
	currentFuel = 0.0
	lapStartFuel = 0.0
	hasFuelData = false
	timeExpired = 0
	writeLog("RunDeinit: Client left track, timing stopped")
}

//export RunStart
func RunStart() { writeLog("RunStart: Simulation started/resumed") }

//export RunStop
func RunStop() { writeLog("RunStop: Simulation paused") }

//export RunLap
func RunLap(p unsafe.Pointer, size C.int) {
	lap := (*C.SPluginsBikeLap_t)(p)
	lapIndex = int(lap.m_iLapNum)
	if timeExpired != 0 && expiryLapStart >= 0 {
		lapsAfterExpiry = lapIndex - expiryLapStart
		if lapsAfterExpiry < 0 {
			lapsAfterExpiry = 0
		}
	}
	if !firstExpirySeen {
		lapsAfterExpired++
	}
	writeLog("RunLap: Fuel calc check - hasFuelData=%t, lapStartFuel=%.2f, currentFuel=%.2f, lastLapFuelDelta=%.2f", hasFuelData, lapStartFuel, currentFuel, lastLapFuelDelta)
	if hasFuelData && lapStartFuel > 0.0 {
		lastLapFuelDelta = lapStartFuel - currentFuel
		writeLog("RunLap: Fuel consumption calculated - Start: %.2f, End: %.2f, Delta: %.2f", lapStartFuel, currentFuel, lastLapFuelDelta)
	} else {
		writeLog("RunLap: No fuel consumption calculated - hasFuelData=%t, lapStartFuel=%.2f, currentFuel=%.2f", hasFuelData, lapStartFuel, currentFuel)
	}
	newLapTime := int(lap.m_iLapTime)
	if newLapTime > 0 && (clientBestLapTimeMS == 0 || newLapTime < clientBestLapTimeMS) {
		updateBestLapData()
		writeLog("RunLap: New best lap time %d ms, updated delta data", newLapTime)
	}
	clientLastLapTimeMS = newLapTime
	writeLog("RunLap: Called with lap %d, isOnTrack=%t, testBikeEventType=%d", lapIndex, isOnTrack, testBikeEventType)
	var currentTime int
	if haveTelemetry != 0 && telemetryTime >= 0 {
		currentTime = int(telemetryTime * 1000.0)
	} else if sessionTimeMs >= 0 {
		currentTime = sessionTimeMs
	} else {
		currentTime = 0
	}
	lapStartTime = currentTime
	hasStartedLap = true
	resetCurrentLapData()
	if hasFuelData {
		lapStartFuel = currentFuel
		writeLog("RunLap: Lap start fuel captured: %.2f liters", lapStartFuel)
	}
	writeLog("RunLap: Finish line crossed in testing mode, lap %d, stopwatch reset at time %d", lapIndex, currentTime)
}

//export RunSplit
func RunSplit(p unsafe.Pointer, size C.int) {}

//export TrackCenterline
func TrackCenterline(n C.int, seg *C.SPluginsTrackSegment_t, raceData unsafe.Pointer) {}

//export RaceEvent
func RaceEvent(p unsafe.Pointer, size C.int) {
	if p != nil && int(size) >= int(unsafe.Sizeof(C.SPluginsRaceEvent_t{})) {
		data := (*C.SPluginsRaceEvent_t)(p)
		testRaceEventType = int(data.m_iType)
	}
	currentFuel = 0.0
	lapStartFuel = 0.0
	lastLapFuelDelta = 0.0
	hasFuelData = false
	maxFuel = 0.0
	hasMaxFuel = false
	prevStopwatchMS = 0
	prevClientInPits = false
	timeExpired = 0
	writeLog("RaceEvent: Race event initialized, fuel tracking state reset")
}

//export RaceDeinit
func RaceDeinit() {
	currentFuel = 0.0
	lapStartFuel = 0.0
	lastLapFuelDelta = 0.0
	hasFuelData = false
	maxFuel = 0.0
	hasMaxFuel = false
	prevStopwatchMS = 0
	prevClientInPits = false
	clientLastLapTimeMS = 0
	resetCurrentLapData()
	hasBestLapData = false
	timeExpired = 0
	writeLog("RaceDeinit: Race event closed, fuel tracking state reset")
}

//export RaceAddEntry
func RaceAddEntry(p unsafe.Pointer, size C.int) {
	add := (*C.SPluginsRaceAddEntry_t)(p)
	raceAddEntry = append(raceAddEntry, int(add.m_iRaceNum))
	var addEntry RiderAddEntry
	addEntry.raceNum = int(add.m_iRaceNum)
	sub := ""
	fullName := C.GoString(&add.m_szName[0])
	if len(fullName) >= 3 {
		sub = fullName[:3]
	} else {
		sub = fullName
	}
	addEntry.name = sub
	addEntry.bikeName = C.GoString(&add.m_szBikeName[0])
	raceAddEntryName = append(raceAddEntryName, addEntry)
}

//export RaceRemoveEntry
func RaceRemoveEntry(p unsafe.Pointer, size C.int) {
	rem := (*C.SPluginsRaceRemoveEntry_t)(p)
	num := int(rem.m_iRaceNum)
	removedRaceNumbers = append(removedRaceNumbers, num)
}

//export RaceClassification
func RaceClassification(p unsafe.Pointer, size C.int, arr unsafe.Pointer, elemSize C.int) {
	if p == nil || int(size) < int(unsafe.Sizeof(C.SPluginsRaceClassification_t{})) {
		return
	}
	rc := (*C.SPluginsRaceClassification_t)(p)
	sessionTimeMs = int(rc.m_iSessionTime)
	testRaceClassSession = int(rc.m_iSession)
	sessionStateRaceClassification = int(rc.m_iSessionState)
	if timeExpired == 1 && firstExpirySeen {
		firstExpirySeen = false
	}
	if len(raceClassEntry) > 0 {
		raceClassEntry = raceClassEntry[:0]
	}
	if len(leaderboardRacers) > 100 {
		leaderboardRacers = leaderboardRacers[:0]
	}
	n := int(rc.m_iNumEntries)
	for i := range n {
		entries := (*C.SPluginsRaceClassificationEntry_t)(unsafe.Add(arr, i*int(unsafe.Sizeof(C.SPluginsRaceClassificationEntry_t{}))))
		raceClassEntry = append(raceClassEntry, int(entries.m_iRaceNum))
		var racer RacerInfo
		racer.raceNum = int(entries.m_iRaceNum)
		racer.gapSeconds = float32(float32(entries.m_iGap) / 1000.0)
		racer.penalitySeconds = float32(float32(entries.m_iPenalty) / 1000.0)
		racer.position = i + 1
		leaderboardRacers = append(leaderboardRacers, racer)
		if racer.raceNum == clientRaceNum {
			clientClassIndex = i
		}
		if int(entries.m_iRaceNum) == clientRaceNum {
			clientGapMs = int(entries.m_iGap)
			oldBestTime := clientBestLapTimeMS
			clientBestLapTimeMS = int(entries.m_iBestLap)
			newClientInPits := (int(entries.m_iPit) == 1)
			if prevClientInPits && !newClientInPits {
				writeLog("RaceClassification: Pit exit detected - resetting max fuel detection for new bike setup")
				hasMaxFuel = false
				maxFuel = 0.0
			}
			prevClientInPits = clientInPits
			clientInPits = newClientInPits
			if oldBestTime != clientBestLapTimeMS && clientBestLapTimeMS > 0 {
				writeLog("RaceClassification: Best lap time updated from %d to %d ms, hasBestLapData=%t, dataPoints=%d", oldBestTime, clientBestLapTimeMS, hasBestLapData, len(bestLapPositionData))
			}
		}
		if int(entries.m_iRaceNum) == myRaceNum {
			penaltyMs = int(entries.m_iPenalty)
			raceClassEntryLap = int(entries.m_iNumLaps)
		}
		if int(entries.m_iBestLap) < ridersFastestLapTimeMS {
			ridersFastestLapTimeMS = int(entries.m_iBestLap)
			riderFastestLapRaceNumber = int(entries.m_iRaceNum)
		}
	}
	// Hooks after leaderboard snapshot update
	dispatchRaceClassification(newCtx())
}

//export RaceCommunication
func RaceCommunication(p unsafe.Pointer, size C.int) {
	rc := (*C.SPluginsRaceCommunication_t)(p)
	if int(rc.m_iCommunication) == 2 && int(rc.m_iType) == 0 && int(rc.m_iRaceNum) == myRaceNum {
		penaltyMs += int(rc.m_iTime)
	}
	testRaceCommSession = int(rc.m_iSession)
	raceCommEntry = append(raceCommEntry, int(rc.m_iRaceNum))
}

//export RunTelemetry
func RunTelemetry(p unsafe.Pointer, size C.int, fTime C.float, fPos C.float) {
	telemetryTime = float32(fTime)
	haveTelemetry = 1
	currentPos := float32(fPos)
	if p != nil && int(size) >= int(unsafe.Sizeof(C.SPluginsBikeData_t{})) {
		data := (*C.SPluginsBikeData_t)(p)
		rawFuel := float32(data.m_fFuel)
		if rawFuel >= 0.0 && rawFuel <= 50.0 {
			currentFuel = rawFuel
			hasFuelData = true
			if !hasMaxFuel && rawFuel > 0.0 {
				maxFuel = rawFuel
				hasMaxFuel = true
				writeLog("RunTelemetry: Max fuel capacity detected: %.2f liters (new bike setup)", maxFuel)
			}
			if lapStartFuel == 0.0 || (!hasStartedLap && rawFuel > 0.0) {
				lapStartFuel = rawFuel
				writeLog("RunTelemetry: Initial fuel reading captured: %.2f liters (hasStartedLap=%t)", lapStartFuel, hasStartedLap)
			}
		} else {
			if rawFuel < 0.0 {
				writeLog("Warning: Negative fuel reading: %.2f", rawFuel)
			} else if rawFuel > 50.0 {
				writeLog("Warning: Unusually high fuel reading: %.2f", rawFuel)
			}
			hasFuelData = false
		}
	} else {
		hasFuelData = false
	}
	if hasSeenPosition {
		if lastTrackPosition > 0.9 && currentPos < 0.1 {
			writeLog("RunTelemetry: FINISH LINE CROSSED! pos %.3f->%.3f, hasStartedLap=%t", lastTrackPosition, currentPos, hasStartedLap)
			currentTime := int(telemetryTime * 1000.0)
			shouldStartTiming := (testBikeEventType == 1 && isOnTrack) || (testBikeEventType == 2)
			if shouldStartTiming {
				lapStartTime = currentTime
				hasStartedLap = true
				writeLog("RunTelemetry: Lap timing started at time %d", currentTime)
			}
		}
	}
	var currentLapTime int
	if hasStartedLap && lapStartTime >= 0 {
		var currentTime int
		if haveTelemetry != 0 && telemetryTime >= 0 {
			currentTime = int(telemetryTime * 1000.0)
		} else if sessionTimeMs >= 0 {
			currentTime = sessionTimeMs
		}
		if currentTime > 0 && currentTime >= lapStartTime {
			currentLapTime = currentTime - lapStartTime
		}
	}
	if hasStartedLap && currentLapTime > 0 {
		recordPositionData(currentPos, currentLapTime)
		deltaTimeMS = calculateDelta(currentPos, currentLapTime)
	}
	if int(telemetryTime*10)%50 == 0 {
		writeLog("RunTelemetry: Track position %.3f, hasSeenPosition=%t", currentPos, hasSeenPosition)
	}
	lastTrackPosition = currentPos
	hasSeenPosition = true
	currentTrackPosition = currentPos
	// Hooks after state updates
	dispatchTelemetry(newCtx())
}

//export RaceLap
func RaceLap(p unsafe.Pointer, size C.int) {
	data := (*C.SPluginsRaceLap_t)(p)
	n := int(data.m_iLapNum)
	if n > lapIndex {
		lapIndex = n
	}
	if timeExpired != 0 && expiryLapStart >= 0 {
		lapsAfterExpiry = lapIndex - expiryLapStart
		if lapsAfterExpiry < 0 {
			lapsAfterExpiry = 0
		}
	}
	isClientLap := (int(data.m_iRaceNum) == clientRaceNum) || (int(data.m_iRaceNum) == myRaceNum)
	writeLog("RaceLap: Called with race num %d, lap %d, clientRaceNum=%d, myRaceNum=%d, clientInPits=%t, isClientLap=%t", int(data.m_iRaceNum), n, clientRaceNum, myRaceNum, clientInPits, isClientLap)
	if isClientLap && !clientInPits {
		writeLog("RaceLap: Skipping fuel calculation to avoid overwriting RunLap result. Current lastLapFuelDelta=%.2f", lastLapFuelDelta)
		newLapTime := int(data.m_iLapTime)
		if len(currentLapPositionData) > 0 {
			if !hasBestLapData || (newLapTime > 0 && (clientBestLapTimeMS == 0 || newLapTime < clientBestLapTimeMS)) {
				updateBestLapData()
			}
		}
		var currentTime int
		if haveTelemetry != 0 && telemetryTime >= 0 {
			currentTime = int(telemetryTime * 1000.0)
		} else if sessionTimeMs >= 0 {
			currentTime = sessionTimeMs
		} else {
			currentTime = 0
		}
		lapStartTime = currentTime
		hasStartedLap = true
		resetCurrentLapData()
		writeLog("RaceLap: Skipping fuel capture to avoid overwriting RunLap values")
		writeLog("RaceLap: Finish line crossed in race mode, race num %d, lap %d, stopwatch reset at time %d", int(data.m_iRaceNum), n, currentTime)
	}
	testRaceLapSession = int(data.m_iSession)
}

//export RaceSplit
func RaceSplit(p unsafe.Pointer, size C.int) {}

//export RaceHoleshot
func RaceHoleshot(p unsafe.Pointer, size C.int) {}

//export RaceSession
func RaceSession(p unsafe.Pointer, size C.int) {
	rs := (*C.SPluginsRaceSession_t)(p)
	totalLaps = int(rs.m_iSessionNumLaps)
	sessionNumLaps = int(rs.m_iSessionNumLaps)
	sessionLengthMs = int(rs.m_iSessionLength)
	isTimedPlusLaps = 0
	if sessionLengthMs > 0 && sessionNumLaps > 0 {
		isTimedPlusLaps = 1
	}
	testRaceSessionSession = int(rs.m_iSession)
	if sessionLengthMs > 0 && sessionNumLaps > 0 {
		sessionFormat = fmt.Sprintf("%s + %d laps", formatMinutes(sessionLengthMs), sessionNumLaps)
	} else if sessionLengthMs > 0 {
		sessionFormat = formatMinutes(sessionLengthMs)
	} else if sessionNumLaps > 0 {
		sessionFormat = fmt.Sprintf("%d laps", sessionNumLaps)
	} else {
		sessionFormat = "Practice"
	}
	writeLog("RaceSession: Starting race session, format=%s, testBikeEventType=%d", sessionFormat, testBikeEventType)
	lapIndex = -1
	penaltyMs = 0
	zeroStableTicks = 0
	timeExpired = 0
	expiryLapStart = -1
	lapsAfterExpiry = 0
	firstZeroSeen = 0
	lapAtZero = -1
	lapStartTime = -1
	lastRaceClassEntryLap = -1
	stopwatchMS = 0
	prevStopwatchMS = 0
	isOnTrack = true
	hasStartedLap = false
	clientInPits = false
	prevClientInPits = false
	lastTrackPosition = 0.0
	hasSeenPosition = false
	currentTrackPosition = 0.0
	currentLapPositionData = currentLapPositionData[:0]
	deltaTimeMS = 0
	isRecordingBestLap = false
	currentFuel = 0.0
	lapStartFuel = 0.0
	hasFuelData = false
}

//export RaceSessionState
func RaceSessionState(p unsafe.Pointer, size C.int) {
	if p != nil && int(size) >= int(unsafe.Sizeof(C.SPluginsRaceSessionState_t{})) {
		data := (*C.SPluginsRaceSessionState_t)(p)
		testRaceSessionStateSession = int(data.m_iSession)
	}
}

//export RaceVehicleData
func RaceVehicleData(p unsafe.Pointer, size C.int) {
	v := (*C.SPluginsRaceVehicleData_t)(p)
	oldClientRaceNum := clientRaceNum
	oldMyRaceNum := myRaceNum
	if int(v.m_iActive) != 0 {
		myRaceNum = int(v.m_iRaceNum)
	}
	clientRaceNum = int(v.m_iRaceNum)
	if (oldClientRaceNum != clientRaceNum && clientRaceNum != -1) || (oldMyRaceNum != myRaceNum && myRaceNum != -1) {
		writeLog("RaceVehicleData: Race numbers updated - clientRaceNum=%d, myRaceNum=%d, active=%d", clientRaceNum, myRaceNum, int(v.m_iActive))
	}
	bikeSpeed = float32(v.m_fSpeedometer)
	bikeGear = int(v.m_iGear)
	bikeRPM = int(v.m_iRPM)
	if stopwatchMS%1000 == 0 && stopwatchMS > 0 {
		writeLog("RaceVehicleData: hasStartedLap=%t, lapStartTime=%d, stopwatchMS=%d, speed=%.2f", hasStartedLap, lapStartTime, stopwatchMS, bikeSpeed)
	}
	if hasStartedLap && lapStartTime >= 0 {
		var currentTime int
		if haveTelemetry != 0 && telemetryTime >= 0 {
			currentTime = int(telemetryTime * 1000.0)
		} else if sessionTimeMs >= 0 {
			currentTime = sessionTimeMs
		} else {
			return
		}
		if currentTime >= lapStartTime {
			newStopwatchMS := currentTime - lapStartTime
			if prevStopwatchMS > 10000 && newStopwatchMS < prevStopwatchMS/2 {
				writeLog("RaceVehicleData: STOPWATCH RESET DETECTED! prev=%d, new=%d - LAP COMPLETED (RunLap handles fuel calc)", prevStopwatchMS, newStopwatchMS)
			}
			prevStopwatchMS = stopwatchMS
			stopwatchMS = newStopwatchMS
		} else {
			// do nothing
		}
	} else {
		stopwatchMS = 0
		prevStopwatchMS = 0
	}
}

//export RaceTrackPosition
func RaceTrackPosition(n C.int, arr unsafe.Pointer, elemSize C.int) {}

//export DrawInit
func DrawInit(nSprites *C.int, spriteNames **C.char, nFonts *C.int, fontNames **C.char) C.int {
	if cQuads == nil {
		cQuads = C.malloc(C.size_t(maxQuads) * C.size_t(unsafe.Sizeof(C.SPluginQuad_t{})))
	}
	if cStrings == nil {
		cStrings = C.malloc(C.size_t(maxStrings) * C.size_t(unsafe.Sizeof(C.SPluginString_t{})))
	}
	if nSprites != nil {
		*nSprites = 0
	}
	if spriteNames != nil {
		*spriteNames = nil
	}
	if nFonts != nil && fontNames != nil {
		if cFontList == nil {
			cFontList = allocCString("FontFix_CqMono.fnt")
		}
		*nFonts = 1
		*fontNames = cFontList
		fontReady = true
	}
	return 0
}

func periodic() {
	if iniSavePath == "" {
		return
	}
	iniPath := filepath.Join(iniSavePath, "pitboardz.ini")
	cfg, _ := ini.Load(iniPath)

	// Fuel Info
	fuelX = float32(cfg.Section("Fuel Info").Key("topLeftX").MustFloat64(float64(fuelX)))
	fuelY = float32(cfg.Section("Fuel Info").Key("topLeftY").MustFloat64(float64(fuelY)))
	fuelSize = float32(cfg.Section("Fuel Info").Key("size").MustFloat64(1.0))
	fuelTextSize = fuelBaseTextSize * fuelSize
	fuelRowSpacing = fuelBaseRowSpacing * fuelSize
	showFuelInfo = cfg.Section("Fuel Info").Key("enabled").MustBool(true)

	// Leaderboard
	prevLeaderboardY := leaderboardY
	prevLeaderboardHeight := leaderboardBackgroundBottom - prevLeaderboardY
	leaderboardX = float32(cfg.Section("Leaderboard").Key("topLeftX").MustFloat64(float64(leaderboardX)))
	newLeaderboardY := float32(cfg.Section("Leaderboard").Key("topLeftY").MustFloat64(float64(leaderboardY)))
	leaderboardSize = float32(cfg.Section("Leaderboard").Key("size").MustFloat64(1.0))
	leaderboardRowSpacing = leaderboardBaseRowSpacing * leaderboardSize
	showLeaderboard = cfg.Section("Leaderboard").Key("enabled").MustBool(true)

	// idk, weird math stuff, should prob look into it
	leaderboardY = newLeaderboardY
	scaledDefaultHeight := leaderboardBaseHeight * leaderboardSize
	if prevLeaderboardHeight <= 0 {
		prevLeaderboardHeight = scaledDefaultHeight
	}
	leaderboardBackgroundBottom = leaderboardY + prevLeaderboardHeight
	minHeight := float32(math.Max(float64(scaledDefaultHeight), float64(leaderboardRowSpacing)))
	if leaderboardBackgroundBottom < leaderboardY+minHeight {
		leaderboardBackgroundBottom = leaderboardY + minHeight
	}

	// Speedometer
	speedometerX = float32(cfg.Section("Speedometer").Key("topLeftX").MustFloat64(float64(speedometerX)))
	speedometerY = float32(cfg.Section("Speedometer").Key("topLeftY").MustFloat64(float64(speedometerY)))
	speedometerSize = float32(cfg.Section("Speedometer").Key("size").MustFloat64(1.0))
	isMilesPerHour = cfg.Section("Speedometer").Key("milesPerHour").MustBool(true)
	showSpeedometer = cfg.Section("Speedometer").Key("enabled").MustBool(true)

	// Stopwatch
	timingX0 = float32(cfg.Section("Stopwatch").Key("TopLeftX").MustFloat64(float64(timingX0)))
	timingY0 = float32(cfg.Section("Stopwatch").Key("TopLeftY").MustFloat64(float64(timingY0)))
	stopwatchSize = float32(cfg.Section("Stopwatch").Key("size").MustFloat64(1.0))
	stopwatchTextSize = stopwatchBaseTextSize * stopwatchSize
	stopwatchRowSpacing = stopwatchBaseRowSpacing * stopwatchSize
	showStopwatch = cfg.Section("Stopwatch").Key("enabled").MustBool(true)
}

//export Draw
func Draw(state C.int, nQuads *C.int, ppQuads *unsafe.Pointer, nStrings *C.int, ppStrings *unsafe.Pointer) {
	periodic()
	quadCount = 0
	strCount = 0
	drawCommon()
	if showLeaderboard {
		// Save initial background height (from previous frame) to draw background first
		initialBackgroundBottom := leaderboardBackgroundBottom
		if initialBackgroundBottom <= leaderboardY {
			initialBackgroundBottom = leaderboardY + (leaderboardBaseHeight * leaderboardSize)
		}
		// Draw background FIRST so it renders behind triangles and text
		scaledWidth := leaderboardBaseWidth * leaderboardSize
		addQuadCCW(leaderboardX, leaderboardY, leaderboardX+scaledWidth, initialBackgroundBottom, 0xA0000000)

		// Now draw all content (text, triangles, etc.)
		leaderboardRowY = leaderboardY
		drawSessionFormat()
		if testBikeEventType == 1 {
			drawTest()
		}
		if testBikeEventType == 2 {
			drawRace()
		}
		// Update background bottom for next frame based on final content height
		leaderboardBackgroundBottom = leaderboardRowY
		if leaderboardBackgroundBottom <= leaderboardY {
			leaderboardBackgroundBottom = leaderboardY + (leaderboardBaseHeight * leaderboardSize)
		}
	}
	if nQuads != nil {
		*nQuads = C.int(quadCount)
	}
	if ppQuads != nil {
		*ppQuads = cQuads
	}
	if !fontReady {
		if nStrings != nil {
			*nStrings = 0
		}
		return
	}
	if nStrings != nil {
		*nStrings = C.int(strCount)
	}
	if ppStrings != nil {
		*ppStrings = cStrings
	}
}

// Dummy main for -buildmode=c-shared
func main() {}

// ---------------- small helpers ----------------
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
