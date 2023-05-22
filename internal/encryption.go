package internal

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

const (
	saltSize   = 16     // Salt size in bytes
	iterations = 100000 // Number of iterations for the hash function
)

// Generate a random salt value
func generateSalt() ([]byte, error) {
	salt := make([]byte, saltSize)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// Hash the pin using the salt
func hashPin(pin string, salt []byte) string {
	// Concatenate the PIN and salt
	data := append([]byte(pin), salt...)

	// Perform the hash function with iterations
	for i := 0; i < iterations; i++ {
		hash := sha256.Sum256(data)
		data = hash[:]
	}

	// Convert the hashed pin to a hex-encoded string
	hashedPin := hex.EncodeToString(data)
	return hashedPin
}

// Compare a provided pin with the stored hashed pin
func comparePins(providedPin string, storedHashedPin string, salt []byte) bool {
	// Hash the provided PIN with the same salt
	hashedProvidedPin := hashPin(providedPin, salt)

	// Compare the hashed PINs for equality
	return hashedProvidedPin == storedHashedPin
}

// EncryptPin encrypts the pin using the hashing algorithm
func EncryptPin(pin string) (EncryptedPin, error) {
	salt, err := generateSalt()
	if err != nil {
		return EncryptedPin{}, err
	}
	return EncryptedPin{encryptedPin: hashPin(pin, salt), salt: salt}, nil
}
