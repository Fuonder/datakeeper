package cipher

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPKCS7Pad(t *testing.T) {
	// Test cases for padding
	tests := map[string]struct {
		data      []byte
		blockSize int
		expected  []byte
	}{
		"validPaddingForDataShorterThanBlockSize": {
			data:      []byte("Hello"), // 5 bytes
			blockSize: 16,
			expected:  []byte("Hello\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B"), // Padded to 16 bytes
		},
		"validPaddingForDataEqualToBlockSize": {
			data:      []byte("1234567890123456"), // Exactly 16 bytes
			blockSize: 16,
			expected:  []byte("1234567890123456\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10\x10"), // Full block of padding
		},
		"validPaddingForDataGreaterThanBlockSize": {
			data:      []byte("This is longer than the block size"),
			blockSize: 16,
			expected:  []byte("This is longer than the block size\x0E\x0E\x0E\x0E\x0E\x0E\x0E\x0E\x0E\x0E\x0E\x0E\x0E\x0E"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			padded, err := pkcs7Pad(tc.data, tc.blockSize)
			require.NoError(t, err)              // No error should occur
			assert.Equal(t, tc.expected, padded) // Check if padding is correct
		})
	}
}

func TestPKCS7Unpad(t *testing.T) {
	// Test cases for unpadding
	tests := map[string]struct {
		data      []byte
		blockSize int
		expected  []byte
		expectErr bool
	}{
		"validUnpad": {
			data:      []byte("Hello\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B"),
			blockSize: 16,
			expected:  []byte("Hello"), // Should remove padding
			expectErr: false,
		},
		"invalidUnpad": {
			data:      []byte("Hello\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0B\x0C"), // Incorrect padding
			blockSize: 16,
			expected:  nil,
			expectErr: true,
		},
		"emptyData": {
			data:      []byte(""),
			blockSize: 16,
			expected:  nil,
			expectErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			unpadded, err := pkcs7Unpad(tc.data, tc.blockSize)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, unpadded)
			}
		})
	}
}

func TestEncryptDecrypt(t *testing.T) {
	// Define a valid key for AES256
	key := []byte("32-byte-long-encryption-key-1234")

	// Create the cipher service
	cipherService, err := NewAES256Cipher(key)
	require.NoError(t, err)

	tests := map[string]struct {
		plainText []byte
		expectErr bool
	}{
		"validEncryptionAndDecryption": {
			plainText: []byte("This is a secret message."),
			expectErr: false,
		},
		"emptyStringEncryption": {
			plainText: []byte(""),
			expectErr: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			// Encrypt
			cipherText, err := cipherService.Encrypt(tc.plainText)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Decrypt
			decryptedText, err := cipherService.Decrypt(cipherText)
			require.NoError(t, err)
			assert.Equal(t, tc.plainText, decryptedText)
		})
	}
}

func TestWrongKeyLength(t *testing.T) {
	// Test encryption with wrong key size
	invalidKey := []byte("shortkey")

	_, err := NewAES256Cipher(invalidKey)
	require.Error(t, err)
	require.Equal(t, errors.New("invalid key size, must be 32 bytes for AES256"), err)
}
