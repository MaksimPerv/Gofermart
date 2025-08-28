package validation

import "strconv"

func ValidateLuhn(number string) bool {
	sum := 0
	isSecond := true
	for i := 0; i < len(number); i++ {
		digit, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}
		if !isSecond {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		isSecond = !isSecond
	}
	return sum&10 == 0
}
