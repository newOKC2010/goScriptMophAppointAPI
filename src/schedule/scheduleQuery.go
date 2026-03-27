package schedule

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"go-script-moph-appoint/src/loadenv"
	"go-script-moph-appoint/src/moph"
)

type ClinicCount struct {
	Name  string
	Total int
}

func isHoliday(db *sql.DB, date string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM holiday WHERE holiday_date = $1", date).Scan(&count)
	if err != nil {
		log.Println("check holiday error:", err)
		return false
	}
	return count > 0
}

func queryClinicCount(db *sql.DB) {
	targetDate := loadenv.LoadDateCount()
	if targetDate == "" {
		targetDate = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}

	if isHoliday(db, targetDate) {
		log.Printf("⚠️  วันที่ %s ตรงกับวันหยุด ข้ามการประมวลผล", targetDate)
		return
	}

	rows, err := db.Query(`
		SELECT c.name, COUNT(o.vn) AS total
		FROM oapp o
		LEFT JOIN clinic c ON c.clinic = o.clinic
		WHERE o.nextdate = $1
		AND c.clinic IN('031','002','001','027')
		GROUP BY o.clinic, c.name
		ORDER BY total DESC`, targetDate)
	if err != nil {
		log.Println("query error:", err)
		return
	}
	defer rows.Close()

	var results []ClinicCount
	for rows.Next() {
		var c ClinicCount
		if err := rows.Scan(&c.Name, &c.Total); err != nil {
			log.Println("scan error:", err)
			continue
		}
		results = append(results, c)
	}

	grandTotal := 0
	for _, c := range results {
		grandTotal += c.Total
	}

	log.Println("┌─────────────────────────────────────────────────┐")
	log.Printf("│  รายงานนัดหมายวันที่ : %-26s│", targetDate)
	log.Printf("│  ประมวลผลเมื่อ       : %-26s│", time.Now().Format("2006-01-02 15:04:05"))
	log.Println("├────┬──────────────────────────────┬──────────────┤")
	log.Println("│ ลำ │ คลินิก                       │    จำนวน    │")
	log.Println("├────┼──────────────────────────────┼──────────────┤")
	for i, c := range results {
		log.Printf("│ %2d │ %-28s │ %10d   │", i+1, c.Name, c.Total)
	}
	log.Println("├────┴──────────────────────────────┼──────────────┤")
	log.Printf("│  รวมทั้งหมด                        │ %10d   │", grandTotal)
	log.Println("└───────────────────────────────────┴──────────────┘")

	patients := queryPatients(db, targetDate)
	clinicID := "69c4d4ecaf6b65ba386b4ce8"
	roomID := "69c4d62b94822eb80329d454"
	moph.Run(patients, clinicID, roomID, targetDate)
}

func queryPatients(db *sql.DB, date string) []moph.Patient {
	rows, err := db.Query(`
		SELECT p.cid, p.pname, p.fname, p.lname,
		       TO_CHAR(p.birthday, 'YYYY-MM-DD'), p.sex, COALESCE(p.informtel,'')
		FROM oapp o
		LEFT JOIN patient p ON p.hn = o.hn
		WHERE o.nextdate = $1
		AND o.clinic = '002'`, date)
	if err != nil {
		log.Println("query patients error:", err)
		return nil
	}
	defer rows.Close()

	var patients []moph.Patient
	for rows.Next() {
		var p moph.Patient
		if err := rows.Scan(&p.CID, &p.PName, &p.FName, &p.LName, &p.Birthday, &p.Sex, &p.Tel); err != nil {
			log.Println("scan patient error:", err)
			continue
		}
		patients = append(patients, p)
	}
	return patients
}

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
