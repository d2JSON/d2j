package encryption

// Encryptor represents an interface for working with encryption.
type Encryptor interface {
	Encrypt(opts EncryptOptions) (string, error)
	Decrypt(opts DecryptOptions) ([]byte, error)
}

// EncryptOptions represents an input options for Encrypt method.
type EncryptOptions struct {
	Data   []byte
	Secret string
}

// DecryptOptions represents an input options for Decrypt method.
type DecryptOptions struct {
	EncryptedData string
	Secret        string
}
