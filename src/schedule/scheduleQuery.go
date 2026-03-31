package schedule

import (
	"database/sql"
	"log"
	"time"

	"go-script-moph-appoint/src/loadenv"
	"go-script-moph-appoint/src/moph"
)

// isHoliday ตรวจสอบว่าวันที่กำหนดเป็นวันหยุดหรือไม่
func isHoliday(db *sql.DB, date string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM holiday WHERE holiday_date = $1", date).Scan(&count)
	if err != nil {
		log.Println("check holiday error:", err)
		return false
	}
	return count > 0
}

// queryClinicCount ดึงข้อมูลจำนวนนัดหมายของทุกคลินิก แล้ววนลูปสร้าง schedule และ appointment ในแต่ละคลินิก
func queryClinicCount(db *sql.DB) {
	// อ่านวันที่เป้าหมายจาก env หรือใช้วันพรุ่งนี้
	targetDate := loadenv.LoadDateCount()
	if targetDate == "" {
		targetDate = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	}

	// ตรวจสอบว่าเป็นวันหยุดหรือไม่
	if isHoliday(db, targetDate) {
		log.Printf("⚠️  วันที่ %s ตรงกับวันหยุด ข้ามการประมวลผล", targetDate)
		return
	}

	// ⚠️  หากเพิ่ม/ลด clinic ต้องแก้ส่วน IN(...) ด้านล่างให้ตรงกับ clinicMap ด้านบน
	rows, err := db.Query(`
		SELECT o.clinic, c.name, COUNT(o.vn) AS total
		FROM oapp o
		LEFT JOIN clinic c ON c.clinic = o.clinic
		INNER JOIN patient p ON p.hn = o.hn
		WHERE o.nextdate = $1
		AND o.clinic IN('031','002','001','027')
		AND p.nationality = '99'
		AND p.citizenship = '99'
		GROUP BY o.clinic, c.name
		ORDER BY total DESC`, targetDate)
	if err != nil {
		log.Println("query error:", err)
		return
	}
	defer rows.Close()

	// เก็บผลลัพธ์ใน slice
	var results []ClinicCount
	for rows.Next() {
		var c ClinicCount
		if err := rows.Scan(&c.Code, &c.Name, &c.Total); err != nil {
			log.Println("scan error:", err)
			continue
		}
		results = append(results, c)
	}

	// คำนวณยอดรวมทั้งหมด
	grandTotal := 0
	for _, c := range results {
		grandTotal += c.Total
	}

	// แสดงรายงานสรุปในรูปแบบตาราง
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

	// วนลูปแต่ละคลินิก เพื่อสร้าง schedule และ appointment
	for _, c := range results {
		// หา config ของคลินิกจาก map
		cfg, ok := clinicMap[c.Code]
		if !ok {
			log.Printf("ไม่พบ config สำหรับ clinic %s", c.Code)
			continue
		}
		// แสดง separator และชื่อคลินิก
		log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		log.Printf("คลินิก: %s  จำนวน: %d", c.Name, c.Total)
		log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		// ดึงข้อมูลผู้ป่วยของคลินิกนี้
		patients := queryPatients(db, targetDate, c.Code)
		// เรียก moph.Run เพื่อสร้าง schedule และ appointment
		moph.Run(patients, cfg.ClinicID, cfg.RoomID, targetDate)
	}
}

// queryPatients ดึงข้อมูลผู้ป่วยที่มีนัดหมายในวันและคลินิกที่กำหนด (เฉพาะคนไทย)
func queryPatients(db *sql.DB, date, clinicCode string) []moph.Patient {
	rows, err := db.Query(`
		SELECT p.cid, p.pname, p.fname, p.lname,
		       COALESCE(TO_CHAR(p.birthday, 'YYYY-MM-DD'), ''), p.sex, COALESCE(p.informtel,'')
		FROM oapp o
		INNER JOIN patient p ON p.hn = o.hn
		WHERE o.nextdate = $1
		AND o.clinic = $2
		AND p.nationality = '99'
		AND p.citizenship = '99'`, date, clinicCode)
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
