package util

func ValidatePhone(phone string) bool {
	if len(phone) < 10 || len(phone) > 12 {
		return false
	}

	switch len(phone) {
	case 10:
		if phone[0:1] != "0" {
			return false
		}
	case 11:
		if phone[0:2] != "84" {
			return false
		}
	case 12:
		if phone[0:3] != "+84" {
			return false
		}
	}

	return true
}
