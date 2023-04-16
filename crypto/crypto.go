package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
)

// GenerateSalt generates a new 16-byte random salt.
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// GenerateKey generates a new encryption key using a password and salt.
// The key length should be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
func GenerateKey(password, salt []byte, keyLength int) ([]byte, error) {
	if keyLength != 16 && keyLength != 24 && keyLength != 32 {
		return nil, errors.New("key length must be 16, 24, or 32 bytes")
	}
	key := make([]byte, keyLength)
	derivedKey, err := scrypt.Key(password, salt, 16384, 8, 1, keyLength)
	if err != nil {
		return nil, err
	}
	copy(key, derivedKey)
	return key, nil
}

// Encrypt encrypts the plaintext using the given key with AES-256.
func Encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// Decrypt decrypts the ciphertext using the given key with AES-256.
func Decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// HexEncode encodes a byte array as a hexadecimal string.
func HexEncode(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

// HexDecode decodes a hexadecimal string into a byte array.
func HexDecode(str string) ([]byte, error) {
	return hex.DecodeString(str)
}
