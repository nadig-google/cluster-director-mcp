package genericCore

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const maxLogFiles = 100

var logFile *os.File

// Hack - fix this later
var filePathMap = map[string]string{
	"/usr/local/google/home/nadig/cluster-director-mcp/pkg/tools/cluster/cluster.go":     "pkg/tools/cluster/cluster.go",
	"/usr/local/google/home/nadig/cluster-director-mcp/pkg/tools/cluster/clusterCore.go": "pkg/tools/cluster/clusterCore.go",
	"/usr/local/google/home/nadig/cluster-director-mcp/pkg/config/config.go":             "pkg/config.go",
}

func deleteOldLogFiles(logNameRoot string) {
	for i := 0; i < maxLogFiles; i++ {
		os.Remove(fmt.Sprintf("%s.%d", logNameRoot, i))
	}
}

/*
// findProjectRoot searches upwards from the given directory to find a go.mod file.
// The directory containing go.mod is considered the project root.
func findProjectRoot(startDir string) (string, error) {
	dir := startDir
	for {
		// Check if go.mod exists in the current directory
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			// Found it!
			return dir, nil
		}

		// Move up to the parent directory
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// Reached the filesystem root (e.g., "/") without finding it.
			return "", fmt.Errorf("go.mod not found in any parent directory")
		}
		dir = parentDir
	}
}
*/

func WriteToLog(message string) {
	msg := ""

	if logFile == nil {
		logFile = CreateUniqueFilePath("logs/log.cluster-director-mcp")
	}

	// Compute caller's package, file and line number
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		fmt.Println(message)
		msg = fmt.Sprintf("<UNKNOWN> : %s\n", message)
	} else {
		relativePath := filePathMap[file]
		msg = fmt.Sprintf("%s:%d: %s\n", relativePath, line, message)
	}
	if logFile != nil {
		logFile.WriteString(msg)
	} else {
		log.Println(msg)
	}
}

func getUniqueLogFileName(logNameRoot string) string {
	for i := 0; i < maxLogFiles; i++ {
		_, err := os.Stat(fmt.Sprintf("%s.%d", logNameRoot, i))
		if err != nil && !os.IsNotExist(err) {
			return fmt.Sprintf("%s.%d", logNameRoot, i)
		}
	}

	// There were too many logfiles, delete the old ones
	deleteOldLogFiles(logNameRoot)

	return fmt.Sprintf("%s.%d", logNameRoot, 0)
}

func CreateUniqueFilePath(logNameRoot string) *os.File {
	// Make the directory if it does not exist, fail silently
	os.MkdirAll(filepath.Dir(logNameRoot), 0755)
	logFile, err := os.OpenFile(getUniqueLogFileName(logNameRoot), os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		// If we can't open the log file, it's a fatal error, so we exit.
		return nil
	}
	return logFile
}

func QueryURLAndGetResult(authToken string, url string) (string, bool) {
	WriteToLog("URL : " + url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	authHeader := fmt.Sprintf("Bearer %s", authToken)
	req.Header.Set("Authorization", authHeader)

	// --- Printing the Request Object ---
	WriteToLog("\n--- Request Details ---")
	WriteToLog("Method: " + req.Method + "\n")
	WriteToLog("Headers:")
	for key, values := range req.Header {
		// Dont write the Authorization token
		if key != "Authorization" {
			WriteToLog(fmt.Sprintf("  %s: %s\n", key, strings.Join(values, ", ")))
		}
	}
	WriteToLog("-----------------------")

	client := &http.Client{
		Timeout: 30 * time.Second, // Set a reasonable timeout.
	}

	WriteToLog("\nSending GET request to: " + url + "\n")
	resp, err := client.Do(req)
	if err != nil {
		WriteToLog("Error making HTTP request")
		return "", false
	}
	// Defer the closing of the response body.
	// This is important to free up network resources.
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		WriteToLog("http.Get() did NOT return StatusOK")
		return "", false
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		WriteToLog("io.ReadAll(body) returned error. Returning ERROR")
		return "", false
	}

	bodyString := string(body)
	return bodyString, true
}
