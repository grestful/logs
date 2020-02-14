package logs

import (
	"fmt"
	"os"
	"strings"

	"github.com/toolkits/file"
)

type ConsoleConfig struct {
	Enable  bool   `json:"enable"`
	Level   string `json:"level"`
	Pattern string `json:"pattern"`
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
	maxLines string `json:"maxLines"` //\d+[KMG]? Suffixes are in terms of thousands
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
		fmt.Fprintf(os.Stderr, "LoadJsonConfiguration: Error: Required level <%s> for filter has unknown value: %s\n", "level", l)
		os.Exit(1)
	}
	return lvl
}

func jsonToConsoleLogWriter(filename string, cf *ConsoleConfig) (*ConsoleLogWriter, bool) {
	format := "[%D %T] [%C] [%L] (%S) %M"

	if len(cf.Pattern) > 0 {
		format = strings.Trim(cf.Pattern, " \r\n")
	}

	if !cf.Enable {
		return nil, true
	}

	clw := NewConsoleLogWriter()
	clw.SetFormat(format)

	return clw, true
}

func jsonToFileLogWriter(filename string, ff *FileConfig) (*FileLogWriter, bool) {
	file := "app.log"
	format := "[%D %T] [%C] [%L] (%S) %M"
	maxLines := 0
	maxsize := 0
	daily := false
	rotate := false
	sanitize := false

	if len(ff.Filename) > 0 {
		file = ff.Filename
	}
	if len(ff.Pattern) > 0 {
		format = strings.Trim(ff.Pattern, " \r\n")
	}
	if len(ff.maxLines) > 0 {
		maxLines = strToNumSuffix(strings.Trim(ff.maxLines, " \r\n"), 1000)
	}
	if len(ff.Maxsize) > 0 {
		maxsize = strToNumSuffix(strings.Trim(ff.Maxsize, " \r\n"), 1024)
	}
	daily = ff.Daily
	rotate = ff.Rotate
	sanitize = ff.Sanitize

	if !ff.Enable {
		return nil, true
	}

	flw := NewFileLogWriter(file, rotate, daily)
	flw.SetFormat(format)
	flw.SetRotateLines(maxLines)
	flw.SetRotateSize(maxsize)
	flw.SetSanitize(sanitize)
	return flw, true
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
