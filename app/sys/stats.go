package sys

import (
	"fmt"
	"gocms/app/models"
	"runtime"
	"time"
)

type SystemStats struct {
	SysStats  func()
	ReqRef    map[string]*models.Request
	TotalHits int
}

var Stats = SystemStats{}

func SysStats() {
	go func() {
		for {
			time.Sleep(500 * time.Millisecond)
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			l := len(Stats.ReqRef)
			s := func(l int) string {
				if l != 1 {
					return "s"
				}
				return ""
			}(l)
			fmt.Printf("\rServing: %v Request%v - Total: %v - Allocated Memory: %v KB", l, s, Stats.TotalHits, memStats.Alloc/1024)

			//maybe add in the number of current requests as a stat too
		}
	}()

}
