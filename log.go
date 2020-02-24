// Package log provides simple logging to console, file and syslog
package log

import (
	"bytes"
	"fmt"
	stdLog "log"
	"log/syslog"
	"os"
	"regexp"
	"runtime"
	"time"
)

type Level uint8
type logger struct{}

const (
	L_fatal Level = iota
	L_alert
	L_critical
	L_error
	L_warning
	L_notice
	L_info
	L_debug
)

var sys *syslog.Writer
var maxConsoleLevel = L_debug
var maxSyslogLevel = L_debug
var maxFileLevel = L_debug
var maxMessagesPerSecond int
var thisSecond int64
var messagesThisSecond int
var logFile string
var myLogger *logger
var msgLevelCheck = regexp.MustCompile(`(?i)^([^\p{L}\p{M}\d]*` + // Any register, any non-letters in beginning
	`(fatal|alert|critical|crit|error|err|e|warning|warn|w|notice|info|debug)s?` + // List of levels
	`[^\p{L}\p{M}\d_/-]+)`) // At least one non-letter in the end
var lastMessage string
var isRepeated = false

// String returns string representation of Level
func (l Level) String() string {
	switch l {
	default:
		fallthrough
	case L_fatal:
		return "Fatal"
	case L_alert:
		return "Alert"
	case L_critical:
		return "Critical"
	case L_error:
		return "Error"
	case L_warning:
		return "Warning"
	case L_notice:
		return "Notice"
	case L_info:
		return "Info"
	case L_debug:
		return "Debug"
	}
}

// getText turns message arguments into single string
func getText(args []interface{}, pattern string, level Level, file string, line int) string {
	var msg string
	if pattern == "" {
		msg = fmt.Sprint(args...)
	} else {
		msg = fmt.Sprintf(pattern, args...)
	}
	return fmt.Sprintf("[%s] %s:%d - %s\n", level, file, line, msg)
}

// checkMessagesLimit returns true if we hit messages limit
func checkMessagesLimit() bool {
	if maxMessagesPerSecond < 1 {
		return false
	}
	if thisSecond != time.Now().Unix() {
		thisSecond = time.Now().Unix()
		messagesThisSecond = 1
		return false
	}
	if messagesThisSecond < maxMessagesPerSecond {
		messagesThisSecond++
		return false
	}
	return true
}

// outputMessage outputs formatted mesage into stdout, syslog and file if message level is high enough
func outputMessage(pattern string, args []interface{}, level Level, depth int) error {
	if level > maxConsoleLevel && level > maxSyslogLevel && level > maxFileLevel {
		return nil // silently ignore message
	}
	if checkMessagesLimit() {
		return nil
	}
	_, file, line, _ := runtime.Caller(depth)
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	text := getText(args, pattern, level, short, line)

	var err error
	msg := "."
	if text != lastMessage {
		if isRepeated {
			msg = time.Now().Format("\n2006-01-02 15:04:05 ") + text
			isRepeated = false
		} else {
			msg = time.Now().Format("2006-01-02 15:04:05 ") + text
		}
		lastMessage = text
	} else {
		isRepeated = true
	}

	if level <= maxConsoleLevel {
		if level <= L_warning {
			_, err = fmt.Fprint(os.Stderr, msg)
		} else {
			_, err = fmt.Fprint(os.Stdout, msg)
		}
		if err != nil {
			return err
		}
	}
	if (sys != nil) && (level <= maxSyslogLevel) {
		switch level {
		default:
			fallthrough
		case L_fatal:
			err = sys.Emerg(text)
		case L_alert:
			err = sys.Alert(text)
		case L_critical:
			err = sys.Crit(text)
		case L_error:
			err = sys.Err(text)
		case L_warning:
			err = sys.Warning(text)
		case L_notice:
			err = sys.Notice(text)
		case L_info:
			err = sys.Info(text)
		case L_debug:
			err = sys.Debug(text)
		}
		if err != nil {
			return err
		}
	}
	if (logFile != "") && (level <= maxFileLevel) {
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			// Don't try to write in this file anymore
			logFile = ""
			return err
		} else {
			if _, err = f.WriteString(msg); err != nil {
				logFile = ""
				return err
			}
			return f.Close()
		}
	}
	return nil
}

// Write is standard writer for this log library. It used when standard log is captured to route all log messages into
// this library. Regexp is used to extract message level from text
func (t *logger) Write(buf []byte) (ln int, err error) {
	level := L_notice
	m := msgLevelCheck.FindSubmatch(buf)
	if len(m) == 3 {
		buf = bytes.Replace(buf, m[1], []byte(""), 1)
		letter := bytes.ToLower(m[2][0:1])
		switch letter[0] {
		case 'f':
			level = L_fatal
		case 'a':
			level = L_alert
		case 'c':
			level = L_critical
		case 'e':
			level = L_error
		case 'w':
			level = L_warning
		case 'n':
			level = L_notice
		case 'i':
			level = L_info
		case 'd':
			level = L_debug
		default:
		}
	}
	ln = len(buf)
	i := []interface{}{string(bytes.TrimSpace(buf))}
	err = outputMessage("", i, level, 4)
	if level == L_fatal {
		os.Exit(-1)
	}
	return
}

// CaptureStdLog does what the name says. It means that any message sent to standard logger will be parsed via Write
func CaptureStdLog() {
	stdLog.SetPrefix("")
	stdLog.SetFlags(0)
	stdLog.SetOutput(myLogger)
}

// RestoreStdLog restores standard logger. Any message to it will be again sent to Stdout
func RestoreStdLog() {
	stdLog.SetPrefix("")
	stdLog.SetFlags(stdLog.LstdFlags)
	stdLog.SetOutput(os.Stdout)
}

// InitSyslog should be called before sending messages to Syslog
func InitSyslog(tag string, level Level) error {
	maxSyslogLevel = level
	var err error
	sys, err = syslog.New(syslog.LOG_DEBUG|syslog.LOG_USER, tag)
	return err
}

// SetSyslogLevel can be used to change Syslog messages level filter
func SetSyslogLevel(level Level) {
	maxConsoleLevel = level
}

// InitFile should be called before writing messages to a local file
func InitFile(file string, level Level) {
	logFile = file
	maxFileLevel = level
}

// SetConsoleLevel sets log level which will be written into stdout
func SetConsoleLevel(level Level) {
	maxConsoleLevel = level
}

// SetMaxMessagesPerSecond sets limit for messages per seconds. All messages above the limit will be discarded
// Set m = 0 to turn off the limit (default)
func SetMaxMessagesPerSecond(m int) {
	maxMessagesPerSecond = m
}

// Lib captures StdLog upon init
func init() {
	CaptureStdLog()
}
