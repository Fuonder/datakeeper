package cipher

type Service interface {
	Encrypt(plainText []byte) (cipherText []byte, err error)
	Decrypt(cipherText []byte) (plainText []byte, err error)
}
