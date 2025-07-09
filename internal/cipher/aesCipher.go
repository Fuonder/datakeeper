package cipher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
)

type aes256Cipher struct {
	key   []byte
	block cipher.Block
}

const blockSize = aes.BlockSize

// NewAES256Cipher creates a new AES256 cipher service with the provided key.
func NewAES256Cipher(key []byte) (Encryptor, error) {
	if len(key) != 32 {
		return nil, errors.New("invalid key size, must be 32 bytes for AES256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	return &aes256Cipher{
		key:   key,
		block: block,
	}, nil
}

// Encrypt encrypts the given plaintext using AES256-CBC with a random IV.
func (c *aes256Cipher) Encrypt(plainText []byte) (cipherText []byte, err error) {
	plainText, err = pkcs7Pad(plainText, blockSize)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, blockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(c.block, iv)

	cipherText = make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)

	return append(iv, cipherText...), nil
}

// Decrypt decrypts the given ciphertext using AES256-CBC.
func (c *aes256Cipher) Decrypt(cipherText []byte) (plainText []byte, err error) {
	if len(cipherText) < blockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := cipherText[:blockSize]
	cipherText = cipherText[blockSize:]

	mode := cipher.NewCBCDecrypter(c.block, iv)

	plainText = make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)

	return pkcs7Unpad(plainText, blockSize)
}

// pkcs7Pad adds PKCS#7 padding to the given data.
func pkcs7Pad(data []byte, blockSize int) ([]byte, error) {
	padding := blockSize - len(data)%blockSize
	if padding == 0 {
		padding = blockSize
	}

	pad := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, pad...), nil
}

// pkcs7Unpad removes PKCS#7 padding from the given data.
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("invalid padding size")
	}

	padding := int(data[len(data)-1])

	if padding == 0 || padding > blockSize {
		return nil, errors.New("invalid padding value")
	}

	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, errors.New("invalid padding bytes")
		}
	}

	return data[:len(data)-padding], nil
}

// IsValidKey checks if the provided key is a valid AES256 key (32 bytes).
func IsValidKey(key string) bool {
	return len(key) == 32
}
