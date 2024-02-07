package hashing

// Hasher represents an interface for hashing values.
type Hasher interface {
	GenerateHashFromString(value string) (string, error)
}
