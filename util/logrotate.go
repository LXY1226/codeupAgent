package util

import (
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
	"unsafe"
)

type RotatingLog struct {
	path, format string
	curFile *os.File
}

func NewRotateLog(path, format string, Rotation time.Duration) *RotatingLog {
	info, err := os.Stat(path)
	if err != nil {
		log.Panicln("Can't stat log dir", err)
	}
	if !info.IsDir() {
		log.Panicln("Log dir isn't dir")
	}
	l := RotatingLog{
		path: path,
		format: format,
	}
	l.tidyLog()
	if l.curFile == nil {
		log.Panicln("no log file available")
	}
	go l.daemonFunc(Rotation)
	return &l
}
// Do this at start of day
func (l *RotatingLog) daemonFunc(Rotation time.Duration) {
	for {
		t := time.Now()
		time.Sleep(time.Until(t.Round(Rotation)))
		l.tidyLog()
	}
}

func (l *RotatingLog) Write(p []byte) (n int, err error) {
	return l.curFile.Write(p)
}

func (l *RotatingLog) tidyLog() {
	fName := filepath.Join(l.path, time.Now().Format(l.format))
	f, err := os.OpenFile(fName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Println("log file open failed", err)
		return
	}
	oldF := l.curFile
	// switch log File
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(l.curFile)), unsafe.Pointer(f))
	oldF.Close()
}