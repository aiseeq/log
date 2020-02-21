package log

import "os"

// Fatal outputs message and terminates application
func Fatal(args ...interface{}) {
	outputMessage("", args, L_fatal, 2)
	os.Exit(-1)
}

// Fatalf outputs formatted message and terminates application
func Fatalf(pattern string, args ...interface{}) {
	outputMessage(pattern, args, L_fatal, 2)
	os.Exit(-1)
}

func Alert(args ...interface{}) {
	outputMessage("", args, L_alert, 2)
}

func Alertf(pattern string, args ...interface{}) {
	outputMessage(pattern, args, L_alert, 2)
}

func Critical(args ...interface{}) {
	outputMessage("", args, L_critical, 2)
}

func Criticalf(pattern string, args ...interface{}) {
	outputMessage(pattern, args, L_critical, 2)
}

func Error(args ...interface{}) {
	outputMessage("", args, L_error, 2)
}

func Errorf(pattern string, args ...interface{}) {
	outputMessage(pattern, args, L_error, 2)
}

func Warning(args ...interface{}) {
	outputMessage("", args, L_warning, 2)
}

func Warningf(pattern string, args ...interface{}) {
	outputMessage(pattern, args, L_warning, 2)
}

func Notice(args ...interface{}) {
	outputMessage("", args, L_notice, 2)
}

func Noticef(pattern string, args ...interface{}) {
	outputMessage(pattern, args, L_notice, 2)
}

func Info(args ...interface{}) {
	outputMessage("", args, L_info, 2)
}

func Infof(pattern string, args ...interface{}) {
	outputMessage(pattern, args, L_info, 2)
}

func Debug(args ...interface{}) {
	outputMessage("", args, L_debug, 2)
}

func Debugf(pattern string, args ...interface{}) {
	outputMessage(pattern, args, L_debug, 2)
}
