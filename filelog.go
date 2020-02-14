// Copyright (C) 2010, Kyle Lemons <kyle@kylelemons.net>.  All rights reserved.

package logs

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// This log writer sends output to a file
type FileLogWriter struct {
	rec chan *LogRecord
	rot chan bool

	// The opened file
	filename string
	file     *os.File

	// The logging format
	format string

	// File header/trailer
	header, trailer string

	// Rotate at linecount
	maxLines          int
	maxLinesCurLines int

	// Rotate at size
	maxsize         int
	maxsizeCurSize int

	// Rotate daily
	daily          bool
	dailyOpenDate int

	// Keep old logfiles (.001, .002, etc)
	rotate    bool
	maxBackup int

	// Sanitize newlines to prevent log injection
	sanitize bool
}

// This is the FileLogWriter's output method
func (w *FileLogWriter) LogWrite(rec *LogRecord) {
	w.rec <- rec
}

func (w *FileLogWriter) Close() {
	close(w.rec)
	_ = w.file.Sync()
}

// NewFileLogWriter creates a new LogWriter which writes to the given file and
// has rotation enabled if rotate is true.
//
// If rotate is true, any time a new log file is opened, the old one is renamed
// with a .### extension to preserve it.  The various Set* methods can be used
// to configure log rotation based on lines, size, and daily.
//
// The standard log-line format is:
//   [%D %T] [%L] (%S) %M
func NewFileLogWriter(fileName string, rotate bool, daily bool) *FileLogWriter {
	w := &FileLogWriter{
		rec:       make(chan *LogRecord, LogBufferLength),
		rot:       make(chan bool),
		filename:  fileName,
		format:    "[%D %T] [%L] (%S) %M",
		daily:     daily,
		rotate:    rotate,
		maxBackup: 999,
		sanitize:  false, // set to false so as not to break compatibility
	}
	// open the file for the first time
	if err := w.intRotate(); err != nil {
		fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
		return nil
	}

	go func() {
		defer recoverPanic()
		defer func() {
			if w.file != nil {
				fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
				w.file.Close()
			}
		}()

		for {
			select {
			case <-w.rot:
				if err := w.intRotate(); err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
					return
				}
			case rec, ok := <-w.rec:
				if !ok {
					return
				}
				now := time.Now()
				if (w.maxLines > 0 && w.maxLinesCurLines >= w.maxLines) ||
					(w.maxsize > 0 && w.maxsizeCurSize >= w.maxsize) ||
					(w.daily && now.Day() != w.dailyOpenDate) {
					if err := w.intRotate(); err != nil {
						fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
						return
					}
				}

				// Sanitize newlines
				if w.sanitize {
					rec.Message = strings.Replace(rec.Message, "\n", "\\n", -1)
				}

				// Perform the write
				n, err := fmt.Fprint(w.file, FormatLogRecord(w.format, rec))
				if err != nil {
					fmt.Fprintf(os.Stderr, "FileLogWriter(%q): %s\n", w.filename, err)
					return
				}

				// Update the counts
				w.maxLinesCurLines++
				w.maxsizeCurSize += n
			}
		}
	}()

	return w
}

// Request that the logs rotate
func (w *FileLogWriter) Rotate() {
	w.rot <- true
}

// If this is called in a threaded context, it MUST be synchronized
func (w *FileLogWriter) intRotate() error {
	// Close any log file that may be open
	if w.file != nil {
		fmt.Fprint(w.file, FormatLogRecord(w.trailer, &LogRecord{Created: time.Now()}))
		w.file.Close()
	}
	// If we are keeping log files, move it to the next available number
	if w.rotate {
		info, err := os.Stat(w.filename)
		// _, err = os.Lstat(w.filename)

		if err == nil { // file exists
			// Find the next available number
			modifiedtime := info.ModTime()
			w.dailyOpenDate = modifiedtime.Day()
			num := 1
			fileName := ""
			if w.daily && time.Now().Day() != w.dailyOpenDate {
				modifiedDate := modifiedtime.Format("2006-01-02")
				// for ; err == nil && num <= w.maxBackup; num++ {
				// 	fileName = w.filename + fmt.Sprintf(".%s.%03d", yesterday, num)
				// 	_, err = os.Lstat(fileName)
				// }
				// if err == nil {
				// 	return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", w.filename)
				// }
				fileName = w.filename + fmt.Sprintf(".%s", modifiedDate)
				w.file.Close()
				// Rename the file to its newfound home
				err = os.Rename(w.filename, fileName)
				if err != nil {
					return fmt.Errorf("Rotate: %s\n", err)
				}
			} else if !w.daily {
				num = w.maxBackup - 1
				for ; num >= 1; num-- {
					fileName = w.filename + fmt.Sprintf(".%d", num)
					nfileName := w.filename + fmt.Sprintf(".%d", num+1)
					_, err = os.Lstat(fileName)
					if err == nil {
						os.Rename(fileName, nfileName)
					}
				}
				w.file.Close()
				// Rename the file to its newfound home
				err = os.Rename(w.filename, fileName)
				// return error if the last file checked still existed
				if err != nil {
					return fmt.Errorf("Rotate: %s\n", err)
				}
			}

		}
	}

	// Open the log file
	fd, err := os.OpenFile(w.filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	w.file = fd

	now := time.Now()
	fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: now}))

	// Set the daily open date to the current date
	w.dailyOpenDate = now.Day()

	// initialize rotation values
	w.maxLinesCurLines = 0
	w.maxsizeCurSize = 0

	return nil
}

// Set the logging format (chainable).  Must be called before the first log
// message is written.
func (w *FileLogWriter) SetFormat(format string) *FileLogWriter {
	w.format = format
	return w
}

// Set the logfile header and footer (chainable).  Must be called before the first log
// message is written.  These are formatted similar to the FormatLogRecord (e.g.
// you can use %D and %T in your header/footer for date and time).
func (w *FileLogWriter) SetHeadFoot(head, foot string) *FileLogWriter {
	w.header, w.trailer = head, foot
	if w.maxLinesCurLines == 0 {
		fmt.Fprint(w.file, FormatLogRecord(w.header, &LogRecord{Created: time.Now()}))
	}
	return w
}

// Set rotate at linecount (chainable). Must be called before the first log
// message is written.
func (w *FileLogWriter) SetRotateLines(maxLines int) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateLines: %v\n", maxLines)
	w.maxLines = maxLines
	return w
}

// Set rotate at size (chainable). Must be called before the first log message
// is written.
func (w *FileLogWriter) SetRotateSize(maxsize int) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateSize: %v\n", maxsize)
	w.maxsize = maxsize
	return w
}

// Set rotate daily (chainable). Must be called before the first log message is
// written.
func (w *FileLogWriter) SetRotateDaily(daily bool) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotateDaily: %v\n", daily)
	w.daily = daily
	return w
}

// Set max backup files. Must be called before the first log message
// is written.
func (w *FileLogWriter) SetRotatemaxBackup(maxBackup int) *FileLogWriter {
	w.maxBackup = maxBackup
	return w
}

// SetRotate changes whether or not the old logs are kept. (chainable) Must be
// called before the first log message is written.  If rotate is false, the
// files are overwritten; otherwise, they are rotated to another file before the
// new log is opened.
func (w *FileLogWriter) SetRotate(rotate bool) *FileLogWriter {
	//fmt.Fprintf(os.Stderr, "FileLogWriter.SetRotate: %v\n", rotate)
	w.rotate = rotate
	return w
}

// SetSanitize changes whether or not the sanitization of newline characters takes
// place. This is to prevent log injection, although at some point the sanitization
// of other non-printable characters might be valueable just to prevent binary
// data from mucking up the logs.
func (w *FileLogWriter) SetSanitize(sanitize bool) *FileLogWriter {
	w.sanitize = sanitize
	return w
}

// NewXMLLogWriter is a utility method for creating a FileLogWriter set up to
// output XML record log messages instead of line-based ones.
func NewXMLLogWriter(fileName string, rotate bool, daily bool) *FileLogWriter {
	return NewFileLogWriter(fileName, rotate, daily).SetFormat(
		`	<record level="%L">
		<timestamp>%D %T</timestamp>
		<source>%S</source>
		<message>%M</message>
	</record>`).SetHeadFoot("<log created=\"%D %T\">", "</log>")
}
