package ui

import (
	"math"
	"time"
)

// ---------------- INI layout config ----------------

var (
	iniLoadErrorLogged  bool
	iniValueErrorLogged = make(map[string]bool)
	iniLastLoadTime     time.Time
)

const (
	iniReloadInterval = 2 * time.Second
	defaultIniContent = `# Pitboardz Configuration
#
# UI EDIT MODE:
#   Press Tab to toggle edit mode while in-game
#   - Left-click to toggle elements on/off
#   - Right-click + drag to move settings panel or HUD elements
#   Changes are saved automatically

[Stopwatch]
TopLeftX = 0.02 # default: 0.02
TopLeftY = 0.55 # default: 0.55
size = 1.0 # size multiplier for entire stopwatch (1.0 = default size)
enabled = true # show/hide stopwatch display

[Speedometer]
topLeftX = 0.895 # default: 0.895
topLeftY = 0.881 # default: 0.881
size = 1.0 # size multiplier for entire speedometer (1.0 = default size)
milesPerHour = true
enabled = true # show/hide speedometer display

[Leaderboard]
topLeftX = 0.87 # default: 0.87
topLeftY = 0.45 # default: 0.45
size = 1.0 # size multiplier for entire leaderboard (1.0 = default size)
enabled = true # show/hide leaderboard display
background = true # show/hide leaderboard background

[Fuel Info]
topLeftX = 0.75 # default: 0.75
topLeftY = 0.02 # default: 0.02
size = 1.0 # size multiplier for entire fuel info (1.0 = default size)
enabled = true # show/hide fuel info display

[Suspension]
topLeftX = 0.75 # default: 0.75
topLeftY = 0.09 # default: 0.09
size = 1.0 # size multiplier for entire suspension panel (1.0 = default size)
enabled = true # show/hide suspension panel

[Debug Panel]
topLeftX = 0.30 # default: 0.30
topLeftY = 0.01 # default: 0.01
size = 1.0 # size multiplier for entire debug panel (1.0 = default size)
enabled = false # show/hide debug panel (shows all plugin state variables with sources)

[Debug Console]
topLeftX = 0.02 # default: 0.02
topLeftY = 0.02 # default: 0.02
size = 1.0 # size multiplier for debug console text (1.0 = default size)
enabled = false # show/hide the transparent on-screen log console

`
)

func initLayoutConfig() {
	if configuration == nil || !configuration.Available() {
		return
	}
	created, err := configuration.EnsureDefault(defaultIniContent)
	if err != nil {
		writeLog("Failed to create default ini file: %v", err)
	} else if created {
		writeLog("Created default ini file at: %s", configuration.Path())
	}
	ensureIniKeysExist()
}

func ensureIniKeysExist() {
	if configuration == nil || !configuration.Available() {
		return
	}
	cfg, err := configuration.Load()
	if err != nil {
		writeLog("ensureIniKeysExist: Failed to load INI: %v", err)
		return
	}

	modified := false

	stopwatchSection := cfg.Section("Stopwatch")
	if !stopwatchSection.HasKey("TopLeftX") {
		stopwatchSection.Key("TopLeftX").SetValue("0.02")
		modified = true
	}
	if !stopwatchSection.HasKey("TopLeftY") {
		stopwatchSection.Key("TopLeftY").SetValue("0.55")
		modified = true
	}
	if !stopwatchSection.HasKey("size") {
		stopwatchSection.Key("size").SetValue("1.0")
		modified = true
	}
	if !stopwatchSection.HasKey("enabled") {
		stopwatchSection.Key("enabled").SetValue("true")
		modified = true
	}

	speedometerSection := cfg.Section("Speedometer")
	if !speedometerSection.HasKey("topLeftX") {
		speedometerSection.Key("topLeftX").SetValue("0.895")
		modified = true
	}
	if !speedometerSection.HasKey("topLeftY") {
		speedometerSection.Key("topLeftY").SetValue("0.881")
		modified = true
	}
	if !speedometerSection.HasKey("size") {
		speedometerSection.Key("size").SetValue("1.0")
		modified = true
	}
	if !speedometerSection.HasKey("milesPerHour") {
		speedometerSection.Key("milesPerHour").SetValue("true")
		modified = true
	}
	if !speedometerSection.HasKey("enabled") {
		speedometerSection.Key("enabled").SetValue("true")
		modified = true
	}

	leaderboardSection := cfg.Section("Leaderboard")
	if !leaderboardSection.HasKey("topLeftX") {
		leaderboardSection.Key("topLeftX").SetValue("0.87")
		modified = true
	}
	if !leaderboardSection.HasKey("topLeftY") {
		leaderboardSection.Key("topLeftY").SetValue("0.45")
		modified = true
	}
	if !leaderboardSection.HasKey("size") {
		leaderboardSection.Key("size").SetValue("1.0")
		modified = true
	}
	if !leaderboardSection.HasKey("enabled") {
		leaderboardSection.Key("enabled").SetValue("true")
		modified = true
	}
	if !leaderboardSection.HasKey("background") {
		leaderboardSection.Key("background").SetValue("true")
		modified = true
	}

	fuelSection := cfg.Section("Fuel Info")
	if !fuelSection.HasKey("topLeftX") {
		fuelSection.Key("topLeftX").SetValue("0.75")
		modified = true
	}
	if !fuelSection.HasKey("topLeftY") {
		fuelSection.Key("topLeftY").SetValue("0.02")
		modified = true
	}
	if !fuelSection.HasKey("size") {
		fuelSection.Key("size").SetValue("1.0")
		modified = true
	}
	if !fuelSection.HasKey("enabled") {
		fuelSection.Key("enabled").SetValue("true")
		modified = true
	}

	suspensionSection := cfg.Section("Suspension")
	if !suspensionSection.HasKey("topLeftX") {
		suspensionSection.Key("topLeftX").SetValue("0.75")
		modified = true
	}
	if !suspensionSection.HasKey("topLeftY") {
		suspensionSection.Key("topLeftY").SetValue("0.09")
		modified = true
	}
	if !suspensionSection.HasKey("size") {
		suspensionSection.Key("size").SetValue("1.0")
		modified = true
	}
	if !suspensionSection.HasKey("enabled") {
		suspensionSection.Key("enabled").SetValue("true")
		modified = true
	}

	debugPanelSection := cfg.Section("Debug Panel")
	if !debugPanelSection.HasKey("topLeftX") {
		debugPanelSection.Key("topLeftX").SetValue("0.30")
		modified = true
	}
	if !debugPanelSection.HasKey("topLeftY") {
		debugPanelSection.Key("topLeftY").SetValue("0.01")
		modified = true
	}
	if !debugPanelSection.HasKey("size") {
		debugPanelSection.Key("size").SetValue("1.0")
		modified = true
	}
	if !debugPanelSection.HasKey("enabled") {
		debugPanelSection.Key("enabled").SetValue("false")
		modified = true
	}

	debugConsoleSection := cfg.Section("Debug Console")
	if !debugConsoleSection.HasKey("topLeftX") {
		debugConsoleSection.Key("topLeftX").SetValue("0.02")
		modified = true
	}
	if !debugConsoleSection.HasKey("topLeftY") {
		debugConsoleSection.Key("topLeftY").SetValue("0.02")
		modified = true
	}
	if !debugConsoleSection.HasKey("size") {
		debugConsoleSection.Key("size").SetValue("1.0")
		modified = true
	}
	if !debugConsoleSection.HasKey("enabled") {
		debugConsoleSection.Key("enabled").SetValue("false")
		modified = true
	}

	if modified {
		err = configuration.Save(cfg)
		if err != nil {
			writeLog("ensureIniKeysExist: Failed to save INI: %v", err)
		} else {
			writeLog("ensureIniKeysExist: Added missing keys to INI file")
		}
	}
}

func periodic() {
	updateEditMode()

	if configuration == nil || !configuration.Available() {
		return
	}

	if editModeActive {
		return
	}

	now := time.Now()
	if !iniLastLoadTime.IsZero() && now.Sub(iniLastLoadTime) < iniReloadInterval {
		return
	}
	iniLastLoadTime = now

	cfg, err := configuration.Load()
	if err != nil {
		if !iniLoadErrorLogged {
			writeLog("periodic: Failed to load INI: %v", err)
			iniLoadErrorLogged = true
		}
		return
	}
	iniLoadErrorLogged = false

	// Fuel Info
	fuelX = float32(cfg.Section("Fuel Info").Key("topLeftX").MustFloat64(float64(fuelX)))
	fuelY = float32(cfg.Section("Fuel Info").Key("topLeftY").MustFloat64(float64(fuelY)))
	fuelSize = float32(cfg.Section("Fuel Info").Key("size").MustFloat64(1.0))
	fuelTextSize = fuelBaseTextSize * fuelSize
	fuelRowSpacing = fuelBaseRowSpacing * fuelSize
	showFuelInfo = cfg.Section("Fuel Info").Key("enabled").MustBool(true)

	// Suspension
	suspensionX = float32(cfg.Section("Suspension").Key("topLeftX").MustFloat64(float64(suspensionX)))
	suspensionY = float32(cfg.Section("Suspension").Key("topLeftY").MustFloat64(float64(suspensionY)))
	suspensionSize = float32(cfg.Section("Suspension").Key("size").MustFloat64(1.0))
	suspensionTextSize = suspensionBaseTextSize * suspensionSize
	suspensionRowSpacing = suspensionBaseRowSpacing * suspensionSize
	showSuspension = cfg.Section("Suspension").Key("enabled").MustBool(true)

	// Leaderboard
	prevLeaderboardY := leaderboardY
	prevLeaderboardHeight := leaderboardBackgroundBottom - prevLeaderboardY
	leaderboardX = float32(cfg.Section("Leaderboard").Key("topLeftX").MustFloat64(float64(leaderboardX)))
	newLeaderboardY := float32(cfg.Section("Leaderboard").Key("topLeftY").MustFloat64(float64(leaderboardY)))
	leaderboardSize = float32(cfg.Section("Leaderboard").Key("size").MustFloat64(1.0))
	leaderboardRowSpacing = leaderboardBaseRowSpacing * leaderboardSize
	showLeaderboard = cfg.Section("Leaderboard").Key("enabled").MustBool(true)
	showLeaderboardBackground = cfg.Section("Leaderboard").Key("background").MustBool(true)

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

	// Debug Panel
	debugPanelX = float32(cfg.Section("Debug Panel").Key("topLeftX").MustFloat64(float64(debugPanelX)))
	debugPanelY = float32(cfg.Section("Debug Panel").Key("topLeftY").MustFloat64(float64(debugPanelY)))
	debugPanelSize = float32(cfg.Section("Debug Panel").Key("size").MustFloat64(1.0))
	showDebugPanel = cfg.Section("Debug Panel").Key("enabled").MustBool(false)

	// Debug Console
	debugConsoleX = float32(cfg.Section("Debug Console").Key("topLeftX").MustFloat64(float64(debugConsoleX)))
	debugConsoleY = float32(cfg.Section("Debug Console").Key("topLeftY").MustFloat64(float64(debugConsoleY)))
	debugConsoleSize = float32(cfg.Section("Debug Console").Key("size").MustFloat64(1.0))
	showDebugConsole = cfg.Section("Debug Console").Key("enabled").MustBool(false)
}
