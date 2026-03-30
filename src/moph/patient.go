package moph

// Patient เก็บข้อมูลผู้ป่วยสำหรับสร้าง appointment
type Patient struct {
	CID      string // เลขบัตรประชาชน
	PName    string // คำนำหน้าชื่อ
	FName    string // ชื่อ
	LName    string // นามสกุล
	Birthday string // วันเกิด (YYYY-MM-DD)
	Sex      string // เพศ (1=ชาย, 2=หญิง)
	Tel      string // เบอร์โทรศัพท์
}

// gender แปลงรหัสเพศจาก database เป็นรูปแบบที่ MOPH API ต้องการ
func (p Patient) gender() int {
	if p.Sex == "1" {
		return 1
	}
	return 2
}

// cleanTel ลบเครื่องหมาย - และช่องว่างออกจากเบอร์โทรศัพท์
func cleanTel(tel string) string {
	result := make([]byte, 0, len(tel))
	for i := 0; i < len(tel); i++ {
		if tel[i] != '-' && tel[i] != ' ' {
			result = append(result, tel[i])
		}
	}
	return string(result)
}
