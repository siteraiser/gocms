package sys

import (
	"fmt"
	"runtime"
	"time"
)

type SystemStats struct {
	SysStats func()
}

var Stats = SystemStats{}

func SysStats() {
	go func() {
		for {
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			fmt.Printf("\rAllocated Memory: %v KB", memStats.Alloc/1024)
			time.Sleep(500 * time.Millisecond)
		}
	}()

}
