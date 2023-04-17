package crypto

import (
	"crypto/rand"
	"golang.org/x/crypto/argon2"
)

const (
	saltLength = 16
	keyLength  = 32
)

func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func DeriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, keyLength)
}

func EncryptWithPassword(plaintext string, password string, salt []byte) ([]byte, error) {
	key := DeriveKey(password, salt)
	return Encrypt([]byte(plaintext), key)
}

func DecryptWithPassword(ciphertext []byte, password string, salt []byte) (string, error) {
	key := DeriveKey(password, salt)
	plaintextBytes, err := Decrypt(ciphertext, key)
	if err != nil {
		return "", err
	}
	return string(plaintextBytes), nil
}
