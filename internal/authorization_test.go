package internal

import (
	"errors"
	"testing"
)

func TestAuthenticate(t *testing.T) {
	InitLogger("", true)

	authStruct := GetAuthorizationService()
	accounts := map[string]EncryptedPin{}

	accounts["jc5678"] = setEncryptedPin("5678", t)
	accounts["jc0123"] = setEncryptedPin("1111", t)
	accounts["jc0001"] = setEncryptedPin("1234", t)

	authStruct.SetAuthData(accounts)

	//goland:noinspection SpellCheckingInspection
	testCases := []struct {
		description   string
		accountNumber string
		pin           string
		expected      bool
		err           error
	}{
		{"empty value", "jc0001", "", false, &InvalidInputError{"invalid pin: \"\""}},
		{"non-numeric", "jc0001", "xxxx", false, &InvalidInputError{"pin must be numeric"}},
		{"too short", "jc0001", "123", false, &InvalidInputError{"the pin must be a 4-digit number"}},
		{"too long", "jc0001", "123456", false, &InvalidInputError{"the pin must be a 4-digit number"}},
		{"invalid pin", "jc5678", "1234", false, nil},
		{"valid pin", "jc0123", "1111", true, nil},
	}

	for _, tc := range testCases {

		result, err := authStruct.Authenticate(tc.accountNumber, tc.pin)
		if result != tc.expected {
			t.Errorf("test \"%s\" failed - returned %t; expected %t\n", tc.description, result, tc.expected)
		}
		if !errors.Is(err, tc.err) {
			var tcMsg string
			if tc.err == nil {
				tcMsg = "nil"
			} else {
				tcMsg = tc.err.Error()
			}
			var errMsg string
			if err == nil {
				errMsg = "nil"
			} else {
				errMsg = err.Error()
			}
			t.Errorf("test \"%s\" failed - expected error to be %s, but was %s\n", tc.description, tcMsg, errMsg)
		}
	}
}

func setEncryptedPin(pin string, t *testing.T) EncryptedPin {
	data, err := EncryptPin(pin)
	if err != nil {
		t.Fatal("failed to encrypt pin")
	}
	return data
}
