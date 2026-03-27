package moph

import "log"

func Run(patients []Patient, clinicID, roomID, targetDate string) {
	if len(patients) == 0 {
		log.Println("ไม่มีข้อมูลผู้ป่วย")
		return
	}

	log.Printf("สร้าง schedule วันที่ %s จำนวน %d slot", targetDate, len(patients))
	scheduleID, err := CreateSchedule(clinicID, roomID, targetDate, len(patients))
	if err != nil {
		log.Println("create schedule error:", err)
		return
	}
	log.Printf("schedule id: %s", scheduleID)

	for i, p := range patients {
		log.Printf("  [%d/%d] %s %s %s", i+1, len(patients), p.PName, p.FName, p.LName)
		if err := CreateAppointment(p, clinicID, roomID, scheduleID); err != nil {
			log.Printf("  appointment error: %v", err)
		}
	}
	log.Printf("เสร็จสิ้น สร้าง appointment %d รายการ", len(patients))
}
