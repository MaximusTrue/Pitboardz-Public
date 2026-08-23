package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// ---------------- logger ----------------
const maxDebugConsoleMessages = 200

var (
	logFile              *os.File
	logMutex             sync.Mutex
	debugConsoleMessages []string
)

func writeLog(format string, args ...any) {
	message := fmt.Sprintf(format, args...)

	logMutex.Lock()
	defer logMutex.Unlock()

	if logFile == nil {
		return
	}

	_, _ = fmt.Fprintln(logFile, message)
	debugConsoleMessages = append(debugConsoleMessages, message)
	if len(debugConsoleMessages) > maxDebugConsoleMessages {
		copy(debugConsoleMessages, debugConsoleMessages[len(debugConsoleMessages)-maxDebugConsoleMessages:])
		debugConsoleMessages = debugConsoleMessages[:maxDebugConsoleMessages]
	}
}

func getDebugConsoleMessages(maxMessages int) []string {
	logMutex.Lock()
	defer logMutex.Unlock()

	startIndex := 0
	if len(debugConsoleMessages) > maxMessages {
		startIndex = len(debugConsoleMessages) - maxMessages
	}
	messages := make([]string, len(debugConsoleMessages)-startIndex)
	copy(messages, debugConsoleMessages[startIndex:])
	return messages
}

func initializeLogging(savePath string) error {
	if savePath != "" {
		_ = os.MkdirAll(savePath, 0o755)
	}
	logPath := filepath.Join(savePath, "Pitboardz.log")
	if file, err := os.Create(logPath); err == nil {
		logMutex.Lock()
		logFile = file
		debugConsoleMessages = nil
		logMutex.Unlock()
		writeLog("Startup: version=%s savePath=%s", pluginVersion, savePath)
		return nil
	} else {
		return err
	}
}

func cleanupResources() {
	writeLog("Shutdown")
	logMutex.Lock()
	defer logMutex.Unlock()

	if logFile != nil {
		_ = logFile.Close()
		logFile = nil
	}
}
