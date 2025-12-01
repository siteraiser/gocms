package sys

import (
	"fmt"
	"gocms/app/helpers"
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

			fmt.Printf("\rServing: %v Request%v - Total: %v - Allocated Memory: %v KB", l, helpers.Grammar.LowerIfPluralS(helpers.Grammar{}, l), Stats.TotalHits, memStats.Alloc/1024)

			//maybe add in the number of current requests as a stat too
		}
	}()

}
