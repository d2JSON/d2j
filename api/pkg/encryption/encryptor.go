package encryption

type Encryptor interface {
	Encrypt(opts EncryptOptions) (string, error)
	Decrypt(opts DecryptOptions) ([]byte, error)
}

type EncryptOptions struct {
	Data   []byte
	Secret string
}

type DecryptOptions struct {
	EncryptedData string
	Secret        string
}
