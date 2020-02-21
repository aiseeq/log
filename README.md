Log library
---

Simple fast and low resource consuming library for writing logs into console, file and syslog

Usage example
---

```go
package main

import "github.com/aiseeq/log"

func main() {
	log.SetConsoleLevel(log.L_debug)
	log.InitFile("log.txt", log.L_info)
	log.InitSyslog("AppName", log.L_warning)
	log.Info("AppName v1.0.0") 
}
```