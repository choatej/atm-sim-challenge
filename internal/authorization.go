package internal

import (
	"fmt"
	"strconv"
	"time"
)

type SessionGetter interface {
	GetSession() *UserSession
}

// UserSession stores the state of the session
type UserSession struct {
	IsAuthenticated  bool
	AccountId        string
	LastActivityTime time.Time
}

// the shared session object
var session = &UserSession{
	IsAuthenticated: false,
}

func GetSession() *UserSession {
	return session
}

// EncryptedPin is the encrypted pin data and its salt
// the salt is needed so that we can hash
// a provided pin using the same salt for comparison
type EncryptedPin struct {
	encryptedPin string
	salt         []byte
}

type Authorization struct {
	accounts map[string]EncryptedPin
}

// the shared authorization object
var authorization = &Authorization{}

func GetAuthorizationService() *Authorization {
	return authorization
}

// SetAuthData sets the Authorization with a map of account id to the encrypted pin data
func (auth *Authorization) SetAuthData(authData map[string]EncryptedPin) {
	auth.accounts = authData
}

// Authenticate the provided pin against the hashed pin for a given account id
func (auth *Authorization) Authenticate(accountId string, pin string) (bool, error) {

	if pin == "" {
		return false, &InvalidInputError{fmt.Sprintf("invalid pin: \"%s\"", pin)}
	}
	_, err := strconv.ParseInt(pin, 0, 64)
	if err != nil {
		return false, &InvalidInputError{"pin must be numeric"}
	}
	if len(pin) != 4 {
		return false, &InvalidInputError{"the pin must be a 4-digit number"}
	}

	encryptedPinData := auth.accounts[accountId]
	if comparePins(pin, encryptedPinData.encryptedPin, encryptedPinData.salt) {
		return true, nil
	} else {
		return false, nil
	}
}
