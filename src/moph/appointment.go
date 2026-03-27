package moph

import "log"

func CreateAppointment(p Patient, clinicID, roomID, slotID string) error {
	tel := cleanTel(p.Tel)
	if tel == "" {
		tel = "0000000000"
	}
	body := map[string]any{
		"cid":                p.CID,
		"title_name":         p.PName,
		"first_name":         p.FName,
		"last_name":          p.LName,
		"birthdate":          p.Birthday,
		"mobile_no":          tel,
		"gender":             p.gender(),
		"line_user_id":       "",
		"clinic_id":          clinicID,
		"room_id":            roomID,
		"slot_id":            slotID,
		"reason":             "sent by moph appointment script",
		"remark":             "-",
		"appointment_mode":   "walkin",
		"appointment_status": "อนุมัติ",
	}
	res, err := post("/open-api/appointment/create", body)
	if err != nil {
		return err
	}
	log.Printf("  appointment %s %s %s → %v", p.PName, p.FName, p.LName, res["message"])
	return nil
}
