package util

import (
	"github.com/cheekybits/is"
	"strconv"
	"testing"
)

func TestLuhn(t *testing.T) {
	is := is.New(t)
	validNumbers := []string{
		"378282246310005",
		"371449635398431",
		//mastercard
		"5555555555554444",
		"5105105105105100",
		//visa
		"4111111111111111",
		"4012888888881881",
		"4222222222222",
		"79927398713",
	}

	invalidNumbers := []string{
		// uncool (thanks wikipedia)
		"79927398710", "79927398711",
		"79927398712",
		"79927398714", "79927398715",
		"79927398716", "79927398717",
		"79927398718", "79927398719",
	}

	for _, str := range validNumbers {

		t.Logf("Check digit for %s should be good...", str)
		calculatedCheckDigit, err := Luhn(str[:len(str)-1])
		is.Nil(err)

		givenCheckDigit, err := strconv.Atoi(string([]byte{str[len(str)-1]}))
		is.Nil(err)

		is.Equal(givenCheckDigit, calculatedCheckDigit)
		if givenCheckDigit == calculatedCheckDigit {
			t.Logf("it was!")
		} else {
			t.Logf("it wasn't :-(")
		}
	}
	for _, str := range invalidNumbers {
		t.Logf("Check digit for %s should be bad...", str)
		calculatedCheckDigit, err := Luhn(str[:len(str)-1])
		is.Nil(err)

		givenCheckDigit, err := strconv.Atoi(string([]byte{str[len(str)-1]}))
		is.Nil(err)

		is.NotEqual(givenCheckDigit, calculatedCheckDigit)
		if givenCheckDigit != calculatedCheckDigit {
			t.Logf("it was!")
		} else {
			t.Logf("it wasn't :-(")
		}
	}
}

func TestCreditCardValidation(t *testing.T) {
	validCards := []string{
		//amex
		"378282246310005",
		"371449635398431",
		//mastercard
		"5555555555554444",
		"5105105105105100",
		//visa
		"4111111111111111",
		"4012888888881881",
		"4222222222222",
	}

	invalidCards := []string{
		// too short
		"4", "42", "422", "4222", "42222",
		"422222", "4222222", "42222222",
		"422222222", "4222222222", "42222222222",
		"422222222222",
		// invalid length (14, 15)
		"42222222222222", "4222222222222222",
		//too long (17)
		"42222222222222222",
	}

	is := is.New(t)
	for _, cc := range validCards {
		t.Logf("%s should be valid card\r\n", cc)

		is.True(validCC(cc))

	}
	for _, cc := range invalidCards {
		t.Logf("%s should be invalid card\r\n", cc)
		is.False(validCC(cc))
	}

}
