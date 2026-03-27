package moph

import (
	"fmt"
	"log"
	"time"
)

func CreateSchedule(clinicID, roomID, targetDate string, count int) (string, error) {
	t, _ := time.Parse("2006-01-02", targetDate)
	body := map[string]any{
		"clinic_id": clinicID,
		"year":      t.Year(),
		"schedule_data": []map[string]any{
			{
				"id":             "",
				"clinic_room_id": roomID,
				"doctor_id":      nil,
				"date":           t.Day(),
				"month":          int(t.Month()),
				"slot":           count,
				"start_time":     "08:00",
				"end_time":       "16:00",
			},
		},
	}
	res, err := post("/open-api/schedule", body)
	if err != nil {
		return "", err
	}
	if status, _ := res["status"].(float64); status != 200 {
		return "", fmt.Errorf("create schedule failed: %v", res)
	}
	return findScheduleID(clinicID, t)
}

func findScheduleID(clinicID string, t time.Time) (string, error) {
	body := map[string]any{
		"sort": "created_at", "order": "desc",
		"offset": 0, "limit": 20,
		"filter": map[string]any{"clinic_id": clinicID},
	}
	// retry 3 ครั้ง เพราะ API อาจยังไม่ sync ข้อมูลทัน
	for attempt := 1; attempt <= 3; attempt++ {
		if attempt > 1 {
			log.Printf("retry %d/3 ...", attempt)
			time.Sleep(2 * time.Second)
		}
		res, err := post("/open-api/schedule/list", body)
		if err != nil {
			return "", err
		}
		rows, _ := res["rows"].([]any)
		for _, item := range rows {
			s, _ := item.(map[string]any)
			date, _ := s["date"].(float64)
			month, _ := s["month"].(float64)
			year, _ := s["year"].(float64)
			if int(date) == t.Day() && int(month) == int(t.Month()) && int(year) == t.Year() {
				id, _ := s["id"].(string)
				return id, nil
			}
		}
	}
	return "", fmt.Errorf("schedule not found for %s", t.Format("2006-01-02"))
}
