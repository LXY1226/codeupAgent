package main

import (
	"github.com/LXY1226/codeupAgent/util"
	"log"
	"time"
)

func init() {
	log.SetFlags(log.Ltime)
	log.SetOutput(util.NewRotateLog("logs", "06-01-02_15-04-05.log", 24*time.Hour))
}

func main() {

}
