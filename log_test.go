package log

import (
	"io/ioutil"
	stdLog "log"
	"os"
	"regexp"
	"testing"
)

func getTempFileName(t *testing.T) string {
	f, err := ioutil.TempFile("", "logTest")
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	return f.Name()
}

func getTempFileData(tempFile string, t *testing.T) []byte {
	data, err := ioutil.ReadFile(tempFile)
	if err != nil {
		t.Error(err)
	}
	os.Remove(tempFile)
	return data
}

func initTempLogFile(t *testing.T) string {
	tempFile := getTempFileName(t)
	InitFile(tempFile, L_debug)
	SetConsoleLevel(L_debug)
	SetMaxMessagesPerSecond(0)
	return tempFile
}

func TestCaptureStdLog(t *testing.T) {
	tempFile := initTempLogFile(t)
	for x := 0; x < 4; x++ {
		stdLog.Println("std")
	}
	stdLog.Println("std")

	data := getTempFileData(tempFile, t)
	check := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[Notice\] log_test.go:\d+ - std\n` +
		`...\n\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[Notice\] log_test.go:\d+ - std`)
	if !check.Match(data) {
		t.Error("Wrong output above, should be something like:\n2020-02-21 12:34:56 [Notice] log_test.go:19 - std\n" +
			"...\n2020-02-21 12:34:56 [Notice] log_test.go:20 - std")
	}
	RestoreStdLog()
}

func TestSetMaxMessagesPerSecond(t *testing.T) {
	tempFile := initTempLogFile(t)
	SetMaxMessagesPerSecond(2)
	Info(1)
	Info(2)
	Info(3)

	data := getTempFileData(tempFile, t)
	check := regexp.MustCompile(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[Info\] log_test.go:\d+ - 1\n` +
		`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[Info\] log_test.go:\d+ - 2`)
	if !check.Match(data) {
		t.Error("Wrong output above, should be something like:\n" +
			"2020-02-21 12:50:38 [Info] log_test.go:53 - 1\n" +
			"2020-02-21 12:50:38 [Info] log_test.go:54 - 2")
	}
}

func TestSyslog(t *testing.T) {
	err := InitSyslog("testLogApp", L_warning)
	if err != nil {
		t.Error(err)
	}
	SetSyslogLevel(L_error)
}

func TestInterfaces(t *testing.T) {
	tempFile := initTempLogFile(t)
	Debugf("%s", "test")
	Debug("test")
	Infof("%s", "test")
	Info("test")
	Noticef("%s", "test")
	Notice("test")
	Warningf("%s", "test")
	Warning("test")
	Errorf("%s", "test")
	Error("test")
	Criticalf("%s", "test")
	Critical("test")
	Alertf("%s", "test")
	Alert("test")

	data := getTempFileData(tempFile, t)
	check := regexp.MustCompile(`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[Debug\] log_test.go:\d+ - test\n){2}` +
		`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[Info\] log_test.go:\d+ - test\n){2}` +
		`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[Notice\] log_test.go:\d+ - test\n){2}` +
		`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[Warning\] log_test.go:\d+ - test\n){2}` +
		`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[Error\] log_test.go:\d+ - test\n){2}` +
		`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[Critical\] log_test.go:\d+ - test\n){2}` +
		`(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2} \[Alert\] log_test.go:\d+ - test\n){2}`)
	if !check.Match(data) {
		t.Error("Wrong output above, should be something like:\n" +
			"2020-02-21 13:26:53 [Debug] log_test.go:81 - test\n" +
			"2020-02-21 13:26:53 [Debug] log_test.go:82 - test\n" +
			"2020-02-21 13:26:53 [Info] log_test.go:83 - test\n" +
			"2020-02-21 13:26:53 [Info] log_test.go:84 - test\n" +
			"2020-02-21 13:26:53 [Notice] log_test.go:85 - test\n" +
			"2020-02-21 13:26:53 [Notice] log_test.go:86 - test\n" +
			"2020-02-21 13:26:53 [Warning] log_test.go:87 - test\n" +
			"2020-02-21 13:26:53 [Warning] log_test.go:88 - test\n" +
			"2020-02-21 13:26:53 [Error] log_test.go:89 - test\n" +
			"2020-02-21 13:26:53 [Error] log_test.go:90 - test\n" +
			"2020-02-21 13:26:53 [Critical] log_test.go:91 - test\n" +
			"2020-02-21 13:26:53 [Critical] log_test.go:92 - test\n" +
			"2020-02-21 13:26:53 [Alert] log_test.go:93 - test\n" +
			"2020-02-21 13:26:53 [Alert] log_test.go:94 - test\n")
	}
}
