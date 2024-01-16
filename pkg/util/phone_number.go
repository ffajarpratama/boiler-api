package util

import "strings"

func AddPhoneCode(phone string) string {
	if strings.HasPrefix(phone, "0") {
		return "62" + phone[1:]
	}

	return phone
}
