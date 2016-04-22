package util

import (
	"strconv"
	"strings"
	"time"
)

func validAlways(string) bool {
	return true
}

func validEmptyOr(otherFn func(string) bool) func(string) bool {
	return func(s string) bool {
		return len(s) == 0 || otherFn(s)
	}
}

func validNonEmpty(s string) bool {
	return len(s) > 1
}

func validName(s string) bool {
	if strings.ContainsAny(s, " \t\r\n") {
		return false
	} else if len(s) >= 3 {
		return true
	}
	return false
}

func validPostcode(s string) bool {
	return len(s) > 1
}

func validNumber(s string) bool {
	for _, c := range s {
		if c < '0' && c > '9' {
			return false
		}
	}
	return true
}

func validCC(cc string) bool {
	if !(len(cc) == 13 || len(cc) == 16 || len(cc) == 15) {

		return false
	}
	if !validNumber(cc) {
		return false
	}
	calculatedCheckDigit, err := Luhn(cc[:len(cc)-1])
	if err != nil {
		return false
	}
	givenCheckDigit, err := strconv.Atoi(cc[len(cc)-1:])
	if err != nil {
		return false
	}
	return calculatedCheckDigit == givenCheckDigit
}

func validCVV(cvv string) bool {
	return validNumber(cvv) && (len(cvv) == 3 || len(cvv) == 4)
}

func validISOCountry(code string) bool {
	return len(code) == 2
}

func validExpiry(exp string) bool {
	if len(exp) != 4 {
		return false
	}
	if !validNumber(exp) {
		return false
	}
	mo, _ := strconv.ParseInt(exp[0:1], 10, 8)
	if mo < 1 || mo > 12 {
		return false
	}
	yr, _ := strconv.ParseInt(exp[2:3], 10, 8)

	// this doesn't handle century boundaries well.
	// but if this code is still in use at the end of the 2000s let me know and I'll
	// cook and eat various head-toppers.
	thisYear := time.Now().Year() % 100
	if int(yr) < thisYear || int(yr) > (thisYear+10) {
		return false
	}
	return true

}

// Luhn calculates the Luhn checksum for the given number
func Luhn(number string) (int, error) {
	sum := 0
	for i, dStr := range strings.Split(number, "") {
		d, err := strconv.Atoi(dStr)
		newDigit := d
		if err != nil {
			return -1, err
		}
		if len(number)%2 == 1 {
			if i%2 == 0 {
				newDigit = d * 2
			}
		} else {
			if i%2 == 1 {
				newDigit = d * 2
			}
		}
		if newDigit >= 10 {
			newDigit = newDigit - 9
		}
		sum += newDigit
	}
	return (10 - (sum % 10)) % 10, nil

}
