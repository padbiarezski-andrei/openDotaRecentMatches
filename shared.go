package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"
	"time"
)

func _ErrLogExit0(err error) {
	if err != nil {
		log.Panicln(err)
		os.Exit(deInit())
	}
}

const mscFormat = "D2006-01-02T15:04:05|MSC"

func timeStamp() string {
	return time.Now().Format(mscFormat)
}

var f = func() func(string) {
	i := int64(0)
	return func(s string) {
		atomic.AddInt64(&i, 1)
		//fmt.Printf("[%s]#%s\t\t%s\n", WPrintf("", timeStamp()), WPrintf("", i), CPrintf("", s))
		fmt.Printf("[%s]#%d\t\t%s\n", timeStamp(), i, s)
	}
}()

func openInBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	_ErrLogExit0(err)
}

// Output color: black
func BkPrintf(f string, a ...interface{}) string {
	return fmt.Sprintf("\033[0;30m"+f+"\033[0m", a...)
}

// Output color: red
func RPrintf(f string, a ...interface{}) string {
	return fmt.Sprintf("\033[0;31m"+f+"\033[0m", a...)
}

// Output color: green
func GPrintf(f string, a ...interface{}) string {
	return fmt.Sprintf("\033[0;32m"+f+"\033[0m", a...)
}

// Output color: yellow
func YPrintf(f string, a ...interface{}) string {
	return fmt.Sprintf("\033[0;33m"+f+"\033[0m", a...)
}

// Output color: blue
func BlPrintf(f string, a ...interface{}) string {
	return fmt.Sprintf("\033[0;34m"+f+"\033[0m", a...)
}

// Output color: magenta
func MPrintf(f string, a ...interface{}) string {
	return fmt.Sprintf("\033[0;35m"+f+"\033[0m", a...)
}

// Output color: cyan
func CPrintf(f string, a ...interface{}) string {
	return fmt.Sprintf("\033[0;36m"+f+"\033[0m", a...)
}

// Output color: white
func WPrintf(f string, a ...interface{}) string {
	return fmt.Sprintf("\033[0;37m"+f+"\033[0m", a...)
}
