package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

type cryptoAES struct{}

var _ Encryptor = (*cryptoAES)(nil)

func NewCryptoAES() *cryptoAES {
	return &cryptoAES{}
}

func (ce *cryptoAES) Encrypt(opts EncryptOptions) (string, error) {
	c, err := aes.NewCipher([]byte(opts.Secret))
	if err != nil {
		return "", fmt.Errorf("create cipher: %w", err)
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", fmt.Errorf("create gcm: %w", err)
	}

	// fill nonce with a cryptographically secure
	// random sequence
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", fmt.Errorf("fill nonce: %w", err)
	}

	// encrypt data
	encryptedData := gcm.Seal(nonce, nonce, opts.Data, nil)

	// convert bytes to hex string
	encryptedDataHexString := hex.EncodeToString(encryptedData)

	return encryptedDataHexString, nil
}

func (ce *cryptoAES) Decrypt(opts DecryptOptions) ([]byte, error) {
	// convert hex string to bytes
	encryptedData, err := hex.DecodeString(opts.EncryptedData)
	if err != nil {
		return nil, fmt.Errorf("decode hex encrypted data: %w", err)
	}

	c, err := aes.NewCipher([]byte(opts.Secret))
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("invalid encrypted data size")
	}

	nonce, encryptedData := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// decrypt data
	decryptedData, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt data: %w", err)
	}

	return decryptedData, nil
}
