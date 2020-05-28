/*
Package hashutil provides utility functions to generate and validate hash.

 */
package hashutil

import (
	"crypto/sha1"
	"fmt"
)

// HashSHA1 generates and returns an SHA1 hash.
func HashSHA1(value string) (string, error) {
	sh := sha1.New()
	_, err := sh.Write([]byte(value))
	if err != nil {
		return "", err
	}

	hash := fmt.Sprintf("%x", sh.Sum(nil))
	return hash, nil
}

// VerifySHA1 checks if the plaintext value matches the given sha1 hash.
func VerifySHA1(hashedValue, value string) (bool, error) {
	hash, err := HashSHA1(value)
	if err != nil {
		return false, err
	}

	return hashedValue == hash, nil
}
