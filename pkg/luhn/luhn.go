package luhn

import "strconv"

func IsValid(digits string) bool {
	var sum int
	for i := 0; i < len(digits); i++ {
		cardNum, err := strconv.Atoi(string(digits[i]))
		if err != nil {
			return false
		}
		if (len(digits)-i)%2 == 0 {
			cardNum = cardNum * 2
			if cardNum > 9 {
				cardNum = cardNum - 9
			}
		}
		sum += cardNum
	}
	return sum%10 == 0
}
