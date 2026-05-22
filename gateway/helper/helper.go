package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	"io"
)

func Encrypt(payload []byte, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err

	}
	nounce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nounce); err != nil {
		return nil, err
	}
	return gcm.Seal(nounce, nounce, payload, nil), nil
}

func Decrypt(En_payload []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(En_payload) < nonceSize {
		return nil, errors.New("TO SHORT")
	}
	nounce, ENDATA := En_payload[:nonceSize], En_payload[nonceSize:]
	return gcm.Open(nil, nounce, ENDATA, nil)

}
