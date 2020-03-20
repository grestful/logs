package logs

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// LOGGER get the log Filter by category
func GetLogger(category string) *Filter {
	f, ok := Global[category]
	if !ok {
		f = &Filter{TRACE, NewConsoleLogWriter(), "DEFAULT"}
	} else {
		f.Category = category
	}
	return f
}

// Send a formatted log message internally
func (f *Filter) intLogf(lvl Level, format string, args ...interface{}) {
	skip := true

	// Determine if any logging will be done
	if lvl >= f.Level {
		skip = false
	}
	if skip {
		return
	}

	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	msg := format
	if len(args) > 0 {
		msg = fmt.Sprintf(format, args...)
	}

	// Make the log record
	rec := &LogRecord{
		Level:    lvl,
		Created:  time.Now(),
		Source:   src,
		Message:  msg,
		Category: f.Category,
	}

	// Dispatch the logs
	/*for _, filt := range log {
		if lvl < filt.Level {
			continue
		}
		filt.LogWrite(rec)
	}
	*/
	defaultFilter := Global["stdout"]

	if defaultFilter != nil && lvl > defaultFilter.Level {
		defaultFilter.LogWrite(rec)
	}

	if f.Category != "DEFAULT" && f.Category != "stdout" {
		f.LogWrite(rec)
	}

}

// Send a closure log message internally
func (f *Filter) intLogc(lvl Level, closure func() string) {
	skip := true

	// Determine if any logging will be done
	if lvl >= f.Level {
		skip = false
	}
	if skip {
		return
	}

	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(2)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}

	// Make the log record
	rec := &LogRecord{
		Level:    lvl,
		Created:  time.Now(),
		Source:   src,
		Message:  closure(),
		Category: f.Category,
	}

	default_filter := Global["stdout"]

	if default_filter != nil && lvl > default_filter.Level {
		default_filter.LogWrite(rec)
	}

	if f.Category != "DEFAULT" && f.Category != "stdout" {
		f.LogWrite(rec)
	}
}

// Send a log message with manual level, source, and message.
func (f *Filter) Log(lvl Level, source, message string) {
	skip := true

	// Determine if any logging will be done
	if lvl >= f.Level {
		skip = false
	}
	if skip {
		return
	}

	// Make the log record
	rec := &LogRecord{
		Level:    lvl,
		Created:  time.Now(),
		Source:   source,
		Message:  message,
		Category: f.Category,
	}

	defaultFilter := Global["stdout"]

	if defaultFilter != nil && lvl > defaultFilter.Level {
		defaultFilter.LogWrite(rec)
	}

	if f.Category != "DEFAULT" && f.Category != "stdout" {
		f.LogWrite(rec)
	}
}

// Logf logs a formatted log message at the given log level, using the caller as
// its source.
func (f *Filter) Logf(lvl Level, format string, args ...interface{}) {
	f.intLogf(lvl, format, args...)
}

// Logc logs a string returned by the closure at the given log level, using the caller as
// its source.  If no log message would be written, the closure is never called.
func (f *Filter) Logc(lvl Level, closure func() string) {
	f.intLogc(lvl, closure)
}

// Debug is a utility method for debug f messages.
// The behavior of Debug depends on the first argument:
// - arg0 is a string
//   When given a string as the first argument, this behaves like ff but with
//   the DEBUG f level: the first argument is interpreted as a format for the
//   latter arguments.
// - arg0 is a func()string
//   When given a closure of type func()string, this fs the string returned by
//   the closure iff it will be fged.  The closure runs at most one time.
// - arg0 is interface{}
//   When given anything else, the f message will be each of the arguments
//   formatted with %v and separated by spaces (ala Sprint).
func (f *Filter) Debug(arg0 interface{}, args ...interface{}) {
	f.intLogf(DEBUG, f.getMsg(arg0, args))
}

// Trace fs a message at the trace f level.
// See Debug for an explanation of the arguments.
func (f *Filter) Trace(arg0 interface{}, args ...interface{}) {
	f.intLogf(TRACE, f.getMsg(arg0, args))
}

// Info fs a message at the info f level.
// See Debug for an explanation of the arguments.
func (f *Filter) Info(arg0 interface{}, args ...interface{}) {
	f.intLogf(INFO, f.getMsg(arg0, args))
}

// Warn fs a message at the warning f level and returns the formatted error.
// At the warning level and higher, there is no performance benefit if the
// message is not actually fged, because all formats are processed and all
// closures are executed to format the error message.
// See Debug for further explanation of the arguments.
func (f *Filter) Warn(arg0 interface{}, args ...interface{}) {
	f.intLogf(WARN, f.getMsg(arg0, args))
}

// Error fs a message at the error f level and returns the formatted error,
// See Warn for an explanation of the performance and Debug for an explanation
// of the parameters.
func (f *Filter) Error(arg0 interface{}, args ...interface{}) {
	f.intLogf(ERROR, f.getMsg(arg0, args))
}

// Fatal fs a message at the error f level and returns the formatted error,
// See Fatal for an explanation of the performance and Debug for an explanation
// of the parameters.
func (f *Filter) Fatal(arg0 interface{}, args ...interface{}) {
	f.intLogf(FATAL, f.getMsg(arg0, args))
}

func (f *Filter) getMsg(arg0 interface{}, args ...interface{}) string {
	var msg string
	switch first := arg0.(type) {
	case string:
		// Use the string as a format string
		msg = fmt.Sprintf(first, args...)
	case func() string:
		// f the closure (no other arguments used)
		msg = first()
	default:
		// Build a format string so that it will be similar to Sprint
		msg = fmt.Sprintf(fmt.Sprint(first)+strings.Repeat(" %v", len(args)), args...)
	}
	return msg
}
