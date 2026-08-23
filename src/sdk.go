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
    int m_iState;         // 1 DNS; 3 retired; 4 DSQ
    int m_iReason;        // 0 jump; 1 offences; 2 director
    int m_iOffence;       // 1 jump; 2 cutting
    int m_iLap;
    int m_iStart;         // 1 before the start line
    int m_iType;          // 0 time penalty
    int m_iTime;          // ms (penalty)
} SPluginsRaceCommunication_t;

typedef struct {
    int m_iSession;
    int m_iSessionState;
    int m_iSessionTime;  // ms. SDK session clock; interpretation depends on session type
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
	"unsafe"

	"github.com/MaximusTrue/PitBoardzz/src/config"
	"github.com/MaximusTrue/PitBoardzz/src/logging"
	"github.com/MaximusTrue/PitBoardzz/src/plugin"
	"github.com/MaximusTrue/PitBoardzz/src/ui"
)

var pluginState = plugin.NewState(writeLog)
var loggingStore *logging.Store

// ---------------- C drawing primitives ----------------
const (
	maxQuads   = 160
	maxStrings = 128
)

var (
	cQuads      unsafe.Pointer
	cStrings    unsafe.Pointer
	quadCount   int
	strCount    int
	fontReady   = false
	cFontList   *C.char
	cSpriteList *C.char
)

const pitboardzDataDir = "pitboardz_data"

var (
	drawSpriteFiles = []string{
		pitboardzDataDir + "/pointer.tga",
		pitboardzDataDir + "/pointer_shadow.tga",
	}

	drawFontFiles = []string{
		pitboardzDataDir + "/FontFix_CqMono.fnt",
	}
)

func quadAt(i int) *C.SPluginQuad_t {
	return (*C.SPluginQuad_t)(unsafe.Add(cQuads, i*int(unsafe.Sizeof(C.SPluginQuad_t{}))))
}

func strAt(i int) *C.SPluginString_t {
	return (*C.SPluginString_t)(unsafe.Add(cStrings, i*int(unsafe.Sizeof(C.SPluginString_t{}))))
}

func clampCoord(v float32) C.float {
	return C.float(v)
}

func addQuadCCW(x0, y0, x1, y1 float32, color ui.Color) {
	if quadCount >= maxQuads {
		return
	}
	q := quadAt(quadCount)
	quadCount++
	q.m_aafPos[0][0] = clampCoord(x0)
	q.m_aafPos[0][1] = clampCoord(y0)
	q.m_aafPos[1][0] = clampCoord(x0)
	q.m_aafPos[1][1] = clampCoord(y1)
	q.m_aafPos[2][0] = clampCoord(x1)
	q.m_aafPos[2][1] = clampCoord(y1)
	q.m_aafPos[3][0] = clampCoord(x1)
	q.m_aafPos[3][1] = clampCoord(y0)
	q.m_iSprite = C.int(0)
	q.m_ulColor = C.ulong(color.ABGR())
}

func addSpriteQuadCCW(x0, y0, x1, y1 float32, spriteIndex int, color ui.Color) {
	if quadCount >= maxQuads || spriteIndex <= 0 {
		return
	}
	q := quadAt(quadCount)
	quadCount++
	q.m_aafPos[0][0] = clampCoord(x0)
	q.m_aafPos[0][1] = clampCoord(y0)
	q.m_aafPos[1][0] = clampCoord(x0)
	q.m_aafPos[1][1] = clampCoord(y1)
	q.m_aafPos[2][0] = clampCoord(x1)
	q.m_aafPos[2][1] = clampCoord(y1)
	q.m_aafPos[3][0] = clampCoord(x1)
	q.m_aafPos[3][1] = clampCoord(y0)
	q.m_iSprite = C.int(spriteIndex)
	q.m_ulColor = C.ulong(color.ABGR())
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

func addText(s string, x, y, size float32, justify int, color ui.Color) {
	if strCount >= maxStrings {
		return
	}
	t := strAt(strCount)
	strCount++
	setCStringFixed100((*[100]C.char)(unsafe.Pointer(&t.m_szString[0])), s)
	t.m_afPos[0] = clampCoord(x)
	t.m_afPos[1] = clampCoord(y)
	t.m_iFont = C.int(1)
	t.m_fSize = C.float(size)
	t.m_iJustify = C.int(justify)
	t.m_ulColor = C.ulong(color.ABGR())
}

func addTriangle(x, y, height float32, color ui.Color) {
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
	q.m_aafPos[0][0] = clampCoord(x1)
	q.m_aafPos[0][1] = clampCoord(y1)
	q.m_aafPos[1][0] = clampCoord(x2)
	q.m_aafPos[1][1] = clampCoord(y2)
	q.m_aafPos[2][0] = clampCoord(x3)
	q.m_aafPos[2][1] = clampCoord(y3)
	q.m_aafPos[3][0] = clampCoord(x1)
	q.m_aafPos[3][1] = clampCoord(y1)
	q.m_iSprite = C.int(0)
	q.m_ulColor = C.ulong(color.ABGR())
}

func addPyramidTriangle(x, y, height float32, color ui.Color) {
	if quadCount >= maxQuads {
		return
	}
	width := height * 0.8
	ax, ay := x+width*0.5, y
	lx, ly := x, y+height
	rx, ry := x+width, y+height
	q := quadAt(quadCount)
	quadCount++
	q.m_aafPos[0][0] = clampCoord(ax)
	q.m_aafPos[0][1] = clampCoord(ay)
	q.m_aafPos[1][0] = clampCoord(lx)
	q.m_aafPos[1][1] = clampCoord(ly)
	q.m_aafPos[2][0] = clampCoord(rx)
	q.m_aafPos[2][1] = clampCoord(ry)
	q.m_aafPos[3][0] = clampCoord(ax)
	q.m_aafPos[3][1] = clampCoord(ay)
	q.m_iSprite = C.int(0)
	q.m_ulColor = C.ulong(color.ABGR())
}

func addExclamation(x, y, height, thickness float32, color ui.Color) {
	if quadCount+2 >= maxQuads {
		return
	}
	lineHeight := height * 0.6
	squareSize := thickness * 1.5
	x0 := x - thickness/2
	x1 := x + thickness/2
	y0 := y
	y1 := y + lineHeight
	qLine := quadAt(quadCount)
	quadCount++
	qLine.m_aafPos[0][0], qLine.m_aafPos[0][1] = clampCoord(x0), clampCoord(y0)
	qLine.m_aafPos[1][0], qLine.m_aafPos[1][1] = clampCoord(x0), clampCoord(y1)
	qLine.m_aafPos[2][0], qLine.m_aafPos[2][1] = clampCoord(x1), clampCoord(y1)
	qLine.m_aafPos[3][0], qLine.m_aafPos[3][1] = clampCoord(x1), clampCoord(y0)
	qLine.m_iSprite = C.int(0)
	qLine.m_ulColor = C.ulong(color.ABGR())
	ySquareTop := y + lineHeight + (thickness * 2)
	ySquareBottom := ySquareTop + (squareSize * 2)
	xSquareLeft := x - squareSize/2
	xSquareRight := x + squareSize/2
	qSquare := quadAt(quadCount)
	quadCount++
	qSquare.m_aafPos[0][0], qSquare.m_aafPos[0][1] = clampCoord(xSquareLeft), clampCoord(ySquareTop)
	qSquare.m_aafPos[1][0], qSquare.m_aafPos[1][1] = clampCoord(xSquareLeft), clampCoord(ySquareBottom)
	qSquare.m_aafPos[2][0], qSquare.m_aafPos[2][1] = clampCoord(xSquareRight), clampCoord(ySquareBottom)
	qSquare.m_aafPos[3][0], qSquare.m_aafPos[3][1] = clampCoord(xSquareRight), clampCoord(ySquareTop)
	qSquare.m_iSprite = C.int(0)
	qSquare.m_ulColor = C.ulong(color.ABGR())
}

type sdkRenderer struct{}

func (sdkRenderer) Quad(x0, y0, x1, y1 float32, color ui.Color) {
	addQuadCCW(x0, y0, x1, y1, color)
}

func (sdkRenderer) Sprite(x0, y0, x1, y1 float32, spriteIndex int, color ui.Color) {
	addSpriteQuadCCW(x0, y0, x1, y1, spriteIndex, color)
}

func (sdkRenderer) Text(text string, x, y, size float32, justify int, color ui.Color) {
	addText(text, x, y, size, justify, color)
}

func (sdkRenderer) Triangle(x, y, height float32, color ui.Color) {
	addTriangle(x, y, height, color)
}

func (sdkRenderer) PyramidTriangle(x, y, height float32, color ui.Color) {
	addPyramidTriangle(x, y, height, color)
}

func (sdkRenderer) Exclamation(x, y, height, thickness float32, color ui.Color) {
	addExclamation(x, y, height, thickness, color)
}

func allocCStringList(entries []string) *C.char {
	if len(entries) == 0 {
		return nil
	}
	total := 1
	for _, entry := range entries {
		total += len(entry) + 1
	}
	b := make([]byte, 0, total)
	for _, entry := range entries {
		b = append(b, []byte(entry)...)
		b = append(b, 0)
	}
	b = append(b, 0)

	ptr := C.malloc(C.size_t(len(b)))
	if ptr == nil {
		return nil
	}
	C.memcpy(ptr, unsafe.Pointer(&b[0]), C.size_t(len(b)))
	return (*C.char)(ptr)
}

// ---------------- SDK exports ----------------

//export GetModID
func GetModID() *C.char { return (*C.char)(C.static_modid()) }

//export GetModDataVersion
func GetModDataVersion() C.int { return 8 }

//export GetInterfaceVersion
func GetInterfaceVersion() C.int { return 9 }

//export Startup
func Startup(szSavePath *C.char) C.int {
	savePath := ""
	if szSavePath != nil {
		savePath = C.GoString(szSavePath)
	}
	logFileErr := initializeLogging(savePath)
	store, databaseErr := logging.Open(savePath)
	if databaseErr != nil {
		writeLog("Startup: Failed to open logging database: %v", databaseErr)
	} else {
		loggingStore = store
		writeLog("Startup: Logging database opened at %s", loggingStore.Path())
	}
	ui.Initialize(pluginState, config.NewStore(savePath), writeLog, getDebugConsoleMessages)
	plugin.WriteUpdateState(savePath, pluginVersion, writeLog)
	if logFileErr == nil {
		go func() {
			updateRelease := plugin.LogLatestVersionComparison(pluginVersion, writeLog)
			pluginState.SetUpdateRelease(updateRelease)
		}()
	}
	return 0
}

//export Shutdown
func Shutdown() {
	if loggingStore != nil {
		if err := loggingStore.Close(); err != nil {
			writeLog("Shutdown: Failed to close logging database: %v", err)
		}
		loggingStore = nil
	}
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
	if cSpriteList != nil {
		C.free(unsafe.Pointer(cSpriteList))
		cSpriteList = nil
	}
}

//export EventInit
func EventInit(p unsafe.Pointer, size C.int) {
	pluginState.BeginBikeEvent()
	data := (*C.SPluginsBikeEvent_t)(p)
	pluginState.BikeMaxRPM = int(data.m_iMaxRPM)
	pluginState.BikeShiftRPM = int(data.m_iShiftRPM)
	pluginState.TestBikeEventType = int(data.m_iType)
	pluginState.SetSuspensionMaxTravel([2]float32{
		float32(data.m_afSuspMaxTravel[0]),
		float32(data.m_afSuspMaxTravel[1]),
	})
	writeLog("EventInit: Event initialized, fuel tracking state reset")
}

//export EventDeinit
func EventDeinit() {
	pluginState.EndBikeEvent()
	writeLog("EventDeinit: Event closed, fuel tracking state reset")
}

//export RunInit
func RunInit(p unsafe.Pointer, size C.int) {
	if p != nil && int(size) >= int(unsafe.Sizeof(C.SPluginsBikeSession_t{})) {
		data := (*C.SPluginsBikeSession_t)(p)
		newSessionNumber := int(data.m_iSession)
		if newSessionNumber != pluginState.SessionNumber {
			writeLog("RunInit: Session type changed from %d to %d", pluginState.SessionNumber, newSessionNumber)
			pluginState.SessionNumber = newSessionNumber
		}
	}
	pluginState.BeginRun()
	writeLog("RunInit: Client entered track, pluginState.IsOnTrack=%t, pluginState.HasStartedLap=%t", pluginState.IsOnTrack, pluginState.HasStartedLap)
}

//export RunDeinit
func RunDeinit() {
	pluginState.EndRun()
	writeLog("RunDeinit: Client left track, timing stopped")
}

//export RunStart
func RunStart() { writeLog("RunStart: Simulation started/resumed") }

//export RunStop
func RunStop() { writeLog("RunStop: Simulation paused") }

//export RunLap
func RunLap(p unsafe.Pointer, size C.int) {
	lap := (*C.SPluginsBikeLap_t)(p)
	pluginState.LapIndex = int(lap.m_iLapNum)
	if pluginState.TimeExpired != 0 && pluginState.ExpiryLapStart >= 0 {
		pluginState.LapsAfterExpiry = pluginState.LapIndex - pluginState.ExpiryLapStart
		if pluginState.LapsAfterExpiry < 0 {
			pluginState.LapsAfterExpiry = 0
		}
	}
	writeLog("RunLap: Fuel calc check - pluginState.HasFuelData=%t, pluginState.LapStartFuel=%.2f, pluginState.CurrentFuel=%.2f, pluginState.LastLapFuelDelta=%.2f", pluginState.HasFuelData, pluginState.LapStartFuel, pluginState.CurrentFuel, pluginState.LastLapFuelDelta)
	if pluginState.HasFuelData && pluginState.LapStartFuel > 0.0 {
		pluginState.LastLapFuelDelta = pluginState.LapStartFuel - pluginState.CurrentFuel
		writeLog("RunLap: Fuel consumption calculated - Start: %.2f, End: %.2f, Delta: %.2f", pluginState.LapStartFuel, pluginState.CurrentFuel, pluginState.LastLapFuelDelta)
	} else {
		writeLog("RunLap: No fuel consumption calculated - pluginState.HasFuelData=%t, pluginState.LapStartFuel=%.2f, pluginState.CurrentFuel=%.2f", pluginState.HasFuelData, pluginState.LapStartFuel, pluginState.CurrentFuel)
	}
	newLapTime := int(lap.m_iLapTime)
	if newLapTime > 0 && (pluginState.ClientBestLapTimeMS == 0 || newLapTime < pluginState.ClientBestLapTimeMS) {
		pluginState.UpdateBestLapData()
		writeLog("RunLap: New best lap time %d ms, updated delta data", newLapTime)
	}
	pluginState.ClientLastLapTimeMS = newLapTime
	writeLog("RunLap: Called with lap %d, pluginState.IsOnTrack=%t, pluginState.TestBikeEventType=%d", pluginState.LapIndex, pluginState.IsOnTrack, pluginState.TestBikeEventType)
	var currentTime int
	if pluginState.HaveTelemetry != 0 && pluginState.TelemetryTime >= 0 {
		currentTime = int(pluginState.TelemetryTime * 1000.0)
	} else if pluginState.RaceSessionClockMS >= 0 {
		currentTime = pluginState.RaceSessionClockMS
	} else {
		currentTime = 0
	}
	pluginState.LapStartTime = currentTime
	pluginState.HasStartedLap = true
	pluginState.ResetCurrentLapData()
	if pluginState.HasFuelData {
		pluginState.LapStartFuel = pluginState.CurrentFuel
		writeLog("RunLap: Lap start fuel captured: %.2f liters", pluginState.LapStartFuel)
	}
	writeLog("RunLap: Finish line crossed in testing mode, lap %d, stopwatch reset at time %d", pluginState.LapIndex, currentTime)
}

//export RunSplit
func RunSplit(p unsafe.Pointer, size C.int) {}

//export TrackCenterline
func TrackCenterline(n C.int, seg *C.SPluginsTrackSegment_t, raceData unsafe.Pointer) {
	pluginState.TrackCenterlineNumSegments = int(n)
	pluginState.TrackCenterlineSegments = pluginState.TrackCenterlineSegments[:0]
	if n > 0 && seg != nil {
		segments := unsafe.Slice(seg, int(n))
		if cap(pluginState.TrackCenterlineSegments) < len(segments) {
			pluginState.TrackCenterlineSegments = make([]plugin.TrackSegment, 0, len(segments))
		}
		for _, s := range segments {
			pluginState.TrackCenterlineSegments = append(pluginState.TrackCenterlineSegments, plugin.TrackSegment{
				Type:   int(s.m_iType),
				Length: float32(s.m_fLength),
				Radius: float32(s.m_fRadius),
				Angle:  float32(s.m_fAngle),
				StartX: float32(s.m_afStart[0]),
				StartY: float32(s.m_afStart[1]),
				Height: float32(s.m_fHeight),
			})
		}
	}
	if raceData != nil {
		raceFloats := unsafe.Slice((*C.float)(raceData), 4)
		for i := range 4 {
			pluginState.TrackCenterlineRaceData[i] = float32(raceFloats[i])
		}
	} else {
		pluginState.TrackCenterlineRaceData = [4]float32{}
	}
}

//export RaceEvent
func RaceEvent(p unsafe.Pointer, size C.int) {
	pluginState.BeginRaceEvent()
	if p != nil && int(size) >= int(unsafe.Sizeof(C.SPluginsRaceEvent_t{})) {
		data := (*C.SPluginsRaceEvent_t)(p)
		pluginState.TestRaceEventType = int(data.m_iType)
	}
	writeLog("RaceEvent: Race event initialized, fuel tracking state reset")
}

//export RaceDeinit
func RaceDeinit() {
	pluginState.EndRaceEvent()
	writeLog("RaceDeinit: Race event closed, fuel tracking state reset")
}

//export RaceAddEntry
func RaceAddEntry(p unsafe.Pointer, size C.int) {
	add := (*C.SPluginsRaceAddEntry_t)(p)
	pluginState.RaceAddEntry = append(pluginState.RaceAddEntry, int(add.m_iRaceNum))
	var addEntry plugin.RiderAddEntry
	addEntry.RaceNum = int(add.m_iRaceNum)
	sub := ""
	fullName := C.GoString(&add.m_szName[0])
	if len(fullName) >= 3 {
		sub = fullName[:3]
	} else {
		sub = fullName
	}
	addEntry.Name = sub
	addEntry.BikeName = C.GoString(&add.m_szBikeName[0])
	if pluginState.RaceNumToRider == nil {
		pluginState.RaceNumToRider = make(map[int]plugin.RiderAddEntry)
	}
	pluginState.RaceNumToRider[addEntry.RaceNum] = addEntry
}

//export RaceRemoveEntry
func RaceRemoveEntry(p unsafe.Pointer, size C.int) {
}

//export RaceClassification
func RaceClassification(p unsafe.Pointer, size C.int, arr unsafe.Pointer, elemSize C.int) {
	if p == nil || int(size) < int(unsafe.Sizeof(C.SPluginsRaceClassification_t{})) {
		return
	}
	rc := (*C.SPluginsRaceClassification_t)(p)
	pluginState.RaceSessionClockMS = int(rc.m_iSessionTime)
	pluginState.TestRaceClassSession = int(rc.m_iSession)
	pluginState.SessionStateRaceClassification = int(rc.m_iSessionState)

	pluginState.BeginRaceClassification()

	n := int(rc.m_iNumEntries)
	for i := range n {
		entries := (*C.SPluginsRaceClassificationEntry_t)(unsafe.Add(arr, i*int(unsafe.Sizeof(C.SPluginsRaceClassificationEntry_t{}))))
		pluginState.RaceClassEntry = append(pluginState.RaceClassEntry, int(entries.m_iRaceNum))
		var racer plugin.RacerInfo
		racer.RaceNum = int(entries.m_iRaceNum)
		racer.GapSeconds = float32(float32(entries.m_iGap) / 1000.0)
		racer.PenaltySeconds = float32(float32(entries.m_iPenalty) / 1000.0)
		racer.Position = i + 1
		racer.NumLaps = int(entries.m_iNumLaps)
		pluginState.LeaderboardRacers = append(pluginState.LeaderboardRacers, racer)
		if racer.RaceNum == pluginState.ClientRaceNum {
			pluginState.ClientClassIndex = i
		}
		if int(entries.m_iRaceNum) == pluginState.ClientRaceNum {
			pluginState.ClientGapMS = int(entries.m_iGap)
			oldBestTime := pluginState.ClientBestLapTimeMS
			pluginState.ClientBestLapTimeMS = int(entries.m_iBestLap)
			newClientInPits := (int(entries.m_iPit) == 1)
			if pluginState.PrevClientInPits && !newClientInPits {
				writeLog("RaceClassification: Pit exit detected - resetting max fuel detection for new bike setup")
				pluginState.HasMaxFuel = false
				pluginState.MaxFuel = 0.0
			}
			pluginState.PrevClientInPits = pluginState.ClientInPits
			pluginState.ClientInPits = newClientInPits
			if oldBestTime != pluginState.ClientBestLapTimeMS && pluginState.ClientBestLapTimeMS > 0 {
				writeLog("RaceClassification: Best lap time updated from %d to %d ms, pluginState.HasBestLapData=%t, dataPoints=%d", oldBestTime, pluginState.ClientBestLapTimeMS, pluginState.HasBestLapData, len(pluginState.BestLapPositionData))
			}
		}
		if int(entries.m_iRaceNum) == pluginState.MyRaceNum {
			pluginState.PenaltyMS = int(entries.m_iPenalty)
			pluginState.RaceClassEntryLap = int(entries.m_iNumLaps)
		}
	}
}

//export RaceCommunication
func RaceCommunication(p unsafe.Pointer, size C.int) {
	rc := (*C.SPluginsRaceCommunication_t)(p)
	if int(rc.m_iCommunication) == 2 && int(rc.m_iType) == 0 && int(rc.m_iRaceNum) == pluginState.MyRaceNum {
		pluginState.PenaltyMS += int(rc.m_iTime)
	}
	pluginState.TestRaceCommSession = int(rc.m_iSession)
}

//export RunTelemetry
func RunTelemetry(p unsafe.Pointer, size C.int, fTime C.float, fPos C.float) {
	if pluginState.TestBikeEventType == 1 {
		pluginState.TelemetryTime = float32(fTime) + pluginState.AccumulatedOnTrackTime
	} else {
		pluginState.TelemetryTime = float32(fTime)
	}
	pluginState.HaveTelemetry = 1
	currentPos := float32(fPos)
	var rawFuel float32
	validFuel := false
	var rawSuspensionLength [2]float32
	var rawSuspensionVelocity [2]float32
	var rawBrakePressure [2]float32
	var rawWheelMaterial [2]int
	validSuspension := false
	if p != nil && int(size) >= int(unsafe.Sizeof(C.SPluginsBikeData_t{})) {
		data := (*C.SPluginsBikeData_t)(p)
		rawFuel = float32(data.m_fFuel)
		validFuel = true

		for i := range rawSuspensionLength {
			rawSuspensionLength[i] = float32(data.m_afSuspLength[i])
			rawSuspensionVelocity[i] = float32(data.m_afSuspVelocity[i])
			rawBrakePressure[i] = float32(data.m_afBrakePressure[i])
			rawWheelMaterial[i] = int(data.m_aiWheelMaterial[i])
		}
		validSuspension = true
	}
	pluginState.UpdateFuel(rawFuel, validFuel)
	pluginState.UpdateSuspension(rawSuspensionLength, rawSuspensionVelocity, rawBrakePressure, rawWheelMaterial, validSuspension)
	if pluginState.HasSeenPosition {
		if pluginState.LastTrackPosition > 0.9 && currentPos < 0.1 {
			writeLog("RunTelemetry: FINISH LINE CROSSED! pos %.3f->%.3f, pluginState.HasStartedLap=%t", pluginState.LastTrackPosition, currentPos, pluginState.HasStartedLap)
			currentTime := int(pluginState.TelemetryTime * 1000.0)
			shouldStartTiming := (pluginState.TestBikeEventType == 1 && pluginState.IsOnTrack) || (pluginState.TestBikeEventType == 2)
			if shouldStartTiming {
				pluginState.LapStartTime = currentTime
				pluginState.HasStartedLap = true
				writeLog("RunTelemetry: Lap timing started at time %d", currentTime)
			}
		}
	}
	var currentLapTime int
	if pluginState.HasStartedLap && pluginState.LapStartTime >= 0 {
		var currentTime int
		if pluginState.HaveTelemetry != 0 && pluginState.TelemetryTime >= 0 {
			currentTime = int(pluginState.TelemetryTime * 1000.0)
		} else if pluginState.RaceSessionClockMS >= 0 {
			currentTime = pluginState.RaceSessionClockMS
		}
		if currentTime > 0 && currentTime >= pluginState.LapStartTime {
			currentLapTime = currentTime - pluginState.LapStartTime
		}
	}
	if pluginState.HasStartedLap && currentLapTime > 0 {
		pluginState.RecordPositionData(currentPos, currentLapTime)
		pluginState.DeltaTimeMS = pluginState.CalculateDelta(currentPos, currentLapTime)
	}
	pluginState.LastTrackPosition = currentPos
	pluginState.HasSeenPosition = true
	pluginState.CurrentTrackPosition = currentPos
}

//export RaceLap
func RaceLap(p unsafe.Pointer, size C.int) {
	data := (*C.SPluginsRaceLap_t)(p)
	n := int(data.m_iLapNum)
	isClientLap := (int(data.m_iRaceNum) == pluginState.ClientRaceNum) || (int(data.m_iRaceNum) == pluginState.MyRaceNum)
	writeLog("RaceLap: Called with race num %d, lap %d, pluginState.ClientRaceNum=%d, pluginState.MyRaceNum=%d, pluginState.ClientInPits=%t, isClientLap=%t", int(data.m_iRaceNum), n, pluginState.ClientRaceNum, pluginState.MyRaceNum, pluginState.ClientInPits, isClientLap)
	if isClientLap {
		if n > pluginState.LapIndex {
			pluginState.LapIndex = n
		}
		if pluginState.TimeExpired != 0 && pluginState.ExpiryLapStart >= 0 {
			pluginState.LapsAfterExpiry = pluginState.LapIndex - pluginState.ExpiryLapStart
			if pluginState.LapsAfterExpiry < 0 {
				pluginState.LapsAfterExpiry = 0
			}
		}
	}
	if isClientLap && !pluginState.ClientInPits {
		writeLog("RaceLap: Skipping fuel calculation to avoid overwriting RunLap result. Current pluginState.LastLapFuelDelta=%.2f", pluginState.LastLapFuelDelta)
		newLapTime := int(data.m_iLapTime)
		if len(pluginState.CurrentLapPositionData) > 0 {
			if !pluginState.HasBestLapData || (newLapTime > 0 && (pluginState.ClientBestLapTimeMS == 0 || newLapTime < pluginState.ClientBestLapTimeMS)) {
				pluginState.UpdateBestLapData()
			}
		}
		var currentTime int
		if pluginState.HaveTelemetry != 0 && pluginState.TelemetryTime >= 0 {
			currentTime = int(pluginState.TelemetryTime * 1000.0)
		} else if pluginState.RaceSessionClockMS >= 0 {
			currentTime = pluginState.RaceSessionClockMS
		} else {
			currentTime = 0
		}
		pluginState.LapStartTime = currentTime
		pluginState.HasStartedLap = true
		pluginState.ResetCurrentLapData()
		writeLog("RaceLap: Skipping fuel capture to avoid overwriting RunLap values")
		writeLog("RaceLap: Finish line crossed in race mode, race num %d, lap %d, stopwatch reset at time %d", int(data.m_iRaceNum), n, currentTime)
	}
	pluginState.TestRaceLapSession = int(data.m_iSession)
}

//export RaceSplit
func RaceSplit(p unsafe.Pointer, size C.int) {}

//export RaceHoleshot
func RaceHoleshot(p unsafe.Pointer, size C.int) {}

func formatMinutes(ms int) string {
	m := ms / 60000
	return fmt.Sprintf("%dm", m)
}

//export RaceSession
func RaceSession(p unsafe.Pointer, size C.int) {
	pluginState.BeginRaceSession()
	rs := (*C.SPluginsRaceSession_t)(p)
	pluginState.TotalLaps = int(rs.m_iSessionNumLaps)
	pluginState.SessionNumLaps = int(rs.m_iSessionNumLaps)
	pluginState.SessionLengthMS = int(rs.m_iSessionLength)
	pluginState.IsTimedPlusLaps = 0
	if pluginState.SessionLengthMS > 0 && pluginState.SessionNumLaps > 0 {
		pluginState.IsTimedPlusLaps = 1
	}
	pluginState.TestRaceSessionSession = int(rs.m_iSession)
	if pluginState.SessionLengthMS > 0 && pluginState.SessionNumLaps > 0 {
		pluginState.SessionFormat = fmt.Sprintf("%s + %d laps", formatMinutes(pluginState.SessionLengthMS), pluginState.SessionNumLaps)
	} else if pluginState.SessionLengthMS > 0 {
		pluginState.SessionFormat = formatMinutes(pluginState.SessionLengthMS)
	} else if pluginState.SessionNumLaps > 0 {
		pluginState.SessionFormat = fmt.Sprintf("%d laps", pluginState.SessionNumLaps)
	} else {
		pluginState.SessionFormat = "Practice"
	}
	writeLog("RaceSession: Starting race session, format=%s, pluginState.TestBikeEventType=%d", pluginState.SessionFormat, pluginState.TestBikeEventType)
}

//export RaceSessionState
func RaceSessionState(p unsafe.Pointer, size C.int) {
	if p != nil && int(size) >= int(unsafe.Sizeof(C.SPluginsRaceSessionState_t{})) {
		data := (*C.SPluginsRaceSessionState_t)(p)
		pluginState.TestRaceSessionStateSession = int(data.m_iSession)
	}
}

//export RaceVehicleData
func RaceVehicleData(p unsafe.Pointer, size C.int) {
	if p == nil || int(size) < int(unsafe.Sizeof(C.SPluginsRaceVehicleData_t{})) {
		return
	}
	data := (*C.SPluginsRaceVehicleData_t)(p)
	if data.m_iActive == 0 {
		return
	}
	pluginState.UpdateRaceVehicleData(int(data.m_iRaceNum), int(data.m_iActive), float32(data.m_fSpeedometer), int(data.m_iGear), int(data.m_iRPM))
}

//export RaceTrackPosition
func RaceTrackPosition(n C.int, arr unsafe.Pointer, elemSize C.int) {
}

//export DrawInit
func DrawInit(nSprites *C.int, spriteNames **C.char, nFonts *C.int, fontNames **C.char) C.int {
	if cQuads == nil {
		cQuads = C.malloc(C.size_t(maxQuads) * C.size_t(unsafe.Sizeof(C.SPluginQuad_t{})))
	}
	if cStrings == nil {
		cStrings = C.malloc(C.size_t(maxStrings) * C.size_t(unsafe.Sizeof(C.SPluginString_t{})))
	}
	if nSprites != nil && spriteNames != nil {
		if cSpriteList == nil {
			cSpriteList = allocCStringList(drawSpriteFiles)
		}
		*nSprites = C.int(len(drawSpriteFiles))
		*spriteNames = cSpriteList
	} else {
		if nSprites != nil {
			*nSprites = 0
		}
		if spriteNames != nil {
			*spriteNames = nil
		}
	}
	if nFonts != nil && fontNames != nil {
		if cFontList == nil {
			cFontList = allocCStringList(drawFontFiles)
		}
		*nFonts = C.int(len(drawFontFiles))
		*fontNames = cFontList
		fontReady = true
	}
	return 0
}

//export Draw
func Draw(state C.int, nQuads *C.int, ppQuads *unsafe.Pointer, nStrings *C.int, ppStrings *unsafe.Pointer) {
	ui.Periodic()
	quadCount = 0
	strCount = 0
	ui.Draw(sdkRenderer{})

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
