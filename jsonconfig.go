package logs

import (
	"fmt"
	"github.com/toolkits/file"
	"os"
)

type ConsoleConfig struct {
	Enable bool   `json:"enable"`
	Level  string `json:"level"`
}

type FileConfig struct {
	Enable   bool   `json:"enable"`
	Category string `json:"category"`
	Level    string `json:"level"`
	Filename string `json:"filename"`

	// %T - Time (15:04:05 MST)
	// %t - Time (15:04)
	// %D - Date (2006/01/02)
	// %d - Date (01/02/06)
	// %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
	// %S - Source
	// %M - Message
	// %C - Category
	// It ignores unknown format strings (and removes them)
	// Recommended: "[%D %T] [%C] [%L] (%S) %M"//
	Pattern string `json:"pattern"`

	Rotate   bool   `json:"rotate"`
	Maxsize  string `json:"maxsize"`  // \d+[KMG]? Suffixes are in terms of 2**10
	MaxLines string `json:"MaxLines"` //\d+[KMG]? Suffixes are in terms of thousands
	Daily    bool   `json:"daily"`    //Automatically rotates by day
	Sanitize bool   `json:"sanitize"` //Sanitize newlines to prevent log injection
}

type SocketConfig struct {
	Enable   bool   `json:"enable"`
	Category string `json:"category"`
	Level    string `json:"level"`
	Pattern  string `json:"pattern"`

	Addr     string `json:"addr"`
	Protocol string `json:"protocol"`
}

// LogConfig presents json log config struct
type LogConfig struct {
	Console *ConsoleConfig  `json:"console"`
	Files   []*FileConfig   `json:"files"`
	Sockets []*SocketConfig `json:"sockets"`
}

func getLogLevel(l string) Level {
	var lvl Level
	switch l {
	case "DEBUG":
		lvl = DEBUG
	case "TRACE":
		lvl = TRACE
	case "INFO":
		lvl = INFO
	case "WARNING":
		lvl = WARN
	case "ERROR":
		lvl = ERROR
	case "FATAL":
		lvl = FATAL
	default:
		_, _ = fmt.Fprintf(os.Stderr, "LoadJsonConfiguration: Error: Required level <%s> for filter has unknown value: %s\n", "level", l)
		os.Exit(1)
	}
	return lvl
}

func ReadFile(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("[%s] path empty", path)
	}

	if !file.IsExist(path) {
		return "", fmt.Errorf("config file %s is nonexistent", path)
	}

	configContent, err := file.ToTrimString(path)
	if err != nil {
		return "", fmt.Errorf("read file %s fail %s", path, err)
	}

	return configContent, nil
}
