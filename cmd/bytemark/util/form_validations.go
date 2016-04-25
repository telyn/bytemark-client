package util

import (
	"bytemark.co.uk/client/util/log"
	"strconv"
	"strings"
	"time"
)

func validAlways(string) (string, bool) {
	return "", true
}

func validEmptyOr(otherFn func(string) (string, bool)) func(string) (string, bool) {
	return func(s string) (string, bool) {
		if len(s) == 0 {
			return "", true
		} else {
			return otherFn(s)
		}
	}
}

func validNonEmpty(s string) (string, bool) {
	if len(s) < 1 {
		return "is empty", false
	}
	return "", true
}

func validPassword(s string) (string, bool) {
	if len(s) < 12 {
		return "is not long enough", false
	}
	return "", true
}

func validName(s string) (string, bool) {
	if strings.ContainsAny(s, " \t\r\n") {
		return "contains spaces", false
	} else if len(s) >= 3 {
		return "", true
	}
	return "is less than three characters long", false
}

func validPostcode(s string) (string, bool) {
	if len(s) > 1 {
		return "", true
	} else {
		return "is less than two characters long", false
	}
}

func validNumber(s string) (string, bool) {
	if p, ok := validNonEmpty(s); !ok {
		return p, ok
	}
	for _, c := range s {
		if c < '0' && c > '9' {
			return "is not a number", false
		}
	}
	return "", true
}

func validCC(cc string) (string, bool) {
	if !(len(cc) == 13 || len(cc) == 16 || len(cc) == 15) {

		return "is not 13, 15 or 16 characters long", false
	}
	if p, ok := validNumber(cc); !ok {
		return p, ok
	}
	calculatedCheckDigit, err := Luhn(cc[:len(cc)-1])
	if err != nil {
		return "error calculating Luhn checksum. Please report this as a bug", false
	}
	givenCheckDigit, err := strconv.Atoi(cc[len(cc)-1:])
	if err != nil {
		return "error getting Luhn check digit. Please report this as a bug", false
	}
	if calculatedCheckDigit != givenCheckDigit {
		return "is invalid - check your card number carefully", false
	}
	return "", true
}

func validCVV(cvv string) (string, bool) {
	if p, ok := validNumber(cvv); !ok {
		return p, ok
	}
	if len(cvv) < 3 {
		return "too short to be a CVV", false
	} else if len(cvv) > 4 {
		return "too long to be a CVV", false
	}
	return "", true
}

func validISOCountry(code string) (string, bool) {
	if len(code) != 2 {
		return "is not 2-digit country code (ISO Alpha-2)", false
	}
	return "", true
}

func validExpiry(exp string) (string, bool) {
	if len(exp) != 5 {
		return "not in MM/YY format - wrong number of characters", false
	}
	if exp[2] != '/' {
		return "not in MM/YY format - no / character", false
	}
	mo, err := strconv.ParseInt(exp[0:2], 10, 8)
	if err != nil {
		return "couldn't parse month - " + err.Error(), false
	}
	if mo < 1 || mo > 12 {
		return "month not between 01 and 12", false
	}
	yr, err := strconv.ParseInt(exp[3:5], 10, 8)
	if err != nil {
		return "couldn't parse year - " + err.Error(), false
	}

	// this doesn't handle century boundaries well.
	// but if this code is still in use at the end of the 2000s let me know and I'll
	// cook and eat various head-toppers.
	thisYear := time.Now().Year() % 100
	log.Debugf(7, "exp[3:4]: %s\r\nyr: %d, thisYear: %d\r\n", exp[3:4], yr, thisYear)
	if int(yr) < thisYear {
		return "expiry in the past", false
	} else if int(yr) > (thisYear + 10) {
		return "expiry too far in the future", false
	}
	return "", true

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
