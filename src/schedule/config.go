package schedule

// ClinicCount เก็บข้อมูลจำนวนนัดหมายของแต่ละคลินิก
type ClinicCount struct {
	Code  string // รหัสคลินิก
	Name  string // ชื่อคลินิก
	Total int    // จำนวนนัดหมาย
}

// clinicConfig เก็บ ID ของคลินิกและห้องตรวจใน MOPH API
type clinicConfig struct {
	ClinicID string // ID ของคลินิกใน MOPH
	RoomID   string // ID ของห้องตรวจใน MOPH
}

// clinicMap แมพรหัสคลินิกกับ ID ใน MOPH API
// ⚠️  วิธีเพิ่ม/ลด clinic:
//     1. เพิ่ม/ลด entry ใน map นี้ โดยใช้รูปแบบ: "รหัสคลินิก": {ClinicID, RoomID}
//     2. แก้ SQL query ใน scheduleQuery.go (ฟังก์ชัน queryClinicCount) ส่วน IN(...) ให้ตรงกัน
//     3. ต้องมี Clinic ID และ Room ID จาก MOPH API ก่อน
//
// ตัวอย่าง mapping:
//   031 = 00 คลินิกโรคทั่วไป       → clinic: 69c4d4ac3bb5b8fa699af212, room: 69c4d5bd77c2f73814b7861e
//   002 = 01 คลินิกความดันโลหิตสูง → clinic: 69c4d4c9af6b65ba386b4caf, room: 69c4d5eeaf6b65ba386b4e8a
//   001 = 02 คลินิกเบาหวาน        → clinic: 69c4d4d6b41b0eb311732747, room: 69c4d60b52f2e01a5cc139c3
//   027 = 03 คลินิกไตวายเรื้อรัง   → clinic: 69c4d4ecaf6b65ba386b4ce8, room: 69c4d62b94822eb80329d454
//   041 = 04 จิตเวชและยาเสพติด      → clinic: 69e59a57271ceed1cb284fb3, room: 69e59aa1271ceed1cb2850e3
var clinicMap = map[string]clinicConfig{
	"031": {"69c4d4ac3bb5b8fa699af212", "69c4d5bd77c2f73814b7861e"}, // 00 คลินิกโรคทั่วไป
	"002": {"69c4d4c9af6b65ba386b4caf", "69c4d5eeaf6b65ba386b4e8a"}, // 01 คลินิกความดันโลหิตสูง
	"001": {"69c4d4d6b41b0eb311732747", "69c4d60b52f2e01a5cc139c3"}, // 02 คลินิกเบาหวาน
	"027": {"69c4d4ecaf6b65ba386b4ce8", "69c4d62b94822eb80329d454"}, // 03 คลินิกไตวายเรื้อรัง
	"041": {"69e59a57271ceed1cb284fb3", "69e59aa1271ceed1cb2850e3"}, // 04 จิตเวชและยาเสพติด
}
