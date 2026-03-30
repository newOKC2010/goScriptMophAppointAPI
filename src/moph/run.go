package moph

import "log"

// Run เป็น main orchestration function สำหรับสร้าง schedule และ appointment ทั้งหมด
// รับรายชื่อผู้ป่วย, clinic ID, room ID และวันที่เป้าหมาย
func Run(patients []Patient, clinicID, roomID, targetDate string) {
	// ตรวจสอบว่ามีผู้ป่วยหรือไม่
	if len(patients) == 0 {
		log.Println("ไม่มีข้อมูลผู้ป่วย")
		return
	}

	// สร้าง schedule ด้วยจำนวน slot เท่ากับจำนวนผู้ป่วย
	log.Printf("สร้าง schedule วันที่ %s จำนวน %d slot", targetDate, len(patients))
	scheduleID, err := CreateSchedule(clinicID, roomID, targetDate, len(patients))
	if err != nil {
		log.Println("create schedule error:", err)
		return
	}
	log.Printf("schedule id: %s", scheduleID)

	// วนลูปสร้าง appointment สำหรับผู้ป่วยแต่ละคน
	for i, p := range patients {
		log.Printf("  [%d/%d] %s %s %s", i+1, len(patients), p.PName, p.FName, p.LName)
		if err := CreateAppointment(p, clinicID, roomID, scheduleID); err != nil {
			log.Printf("  appointment error: %v", err)
		}
	}
	// แสดงสรุปผลการทำงาน
	log.Printf("เสร็จสิ้น สร้าง appointment %d รายการ", len(patients))
}
