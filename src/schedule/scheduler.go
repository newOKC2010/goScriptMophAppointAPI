package schedule

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Start เริ่ม scheduler ที่จะรันทุกวันตามเวลาที่กำหนด
func Start(db *sql.DB, scheduleTime string) {
	nextRun := nextRunTime(scheduleTime)
	log.Printf("เวลาที่เริ่มประมวลผลของทุกวัน %s", nextRun.Format("15:04:05"))

	go func() {
		for {
			time.Sleep(time.Until(nextRun))
			queryClinicCount(db)
			nextRun = nextRun.Add(24 * time.Hour)
			log.Printf("waiting for next run...")
		}
	}()
}

// nextRunTime คำนวณเวลาที่จะรันครั้งถัดไป
func nextRunTime(t string) time.Time {
	now := time.Now()
	var h, m int
	fmt.Sscanf(t, "%d:%d", &h, &m)
	next := time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, now.Location())
	if next.Before(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}
