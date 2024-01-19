package hashing

type Hasher interface {
	GenerateHashFromString(value string) (string, error)
}
