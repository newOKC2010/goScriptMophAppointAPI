package moph

type Patient struct {
	CID      string
	PName    string
	FName    string
	LName    string
	Birthday string
	Sex      string
	Tel      string
}

func (p Patient) gender() int {
	if p.Sex == "1" {
		return 1
	}
	return 2
}

func cleanTel(tel string) string {
	result := make([]byte, 0, len(tel))
	for i := 0; i < len(tel); i++ {
		if tel[i] != '-' && tel[i] != ' ' {
			result = append(result, tel[i])
		}
	}
	return string(result)
}
