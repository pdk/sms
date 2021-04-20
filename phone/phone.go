package phone

import (
	"fmt"
	"unicode"
)

// IsValid checks if it looks like a phone number. If it's valid, return the
// canonical phone number and true. If not valid, return empty string and false.
func IsValid(number string) (string, bool) {

	phoneNumber, err := FormatNumber(number)
	if err == nil {
		return phoneNumber, true
	}

	return "", false
}

// FormatNumber formats a given number into E.164. Currently only handles
// U.S. numbers. E.g.: "(832) 981-1702" => "+1 832-981-1702" Should upgrade to a
// real package to handle this someday, e.g.
// https://github.com/ttacon/libphonenumber
func FormatNumber(phoneNumber string) (string, error) {
	//      "twilioPhoneNumber": "+1 832-981-1702"
	var digits []rune
	for _, char := range phoneNumber {
		if unicode.IsDigit(char) {
			digits = append(digits, char)
		}
	}

	if len(digits) == 11 {
		if digits[0] != '1' {
			return phoneNumber, fmt.Errorf("Failed to convert phone number to E.164 format, wrong length (%d): %s", len(digits), phoneNumber)
		}
		digits = digits[1:]
	}

	if len(digits) != 10 {
		return phoneNumber, fmt.Errorf("Failed to convert phone number to E.164 format, wrong length (%d): %s", len(digits), phoneNumber)
	}

	return fmt.Sprintf("+1 %s-%s-%s", string(digits[0:3]), string(digits[3:6]), string(digits[6:])), nil
}
