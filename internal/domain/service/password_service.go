package service

// PasswordHasher defines the interface for password hashing
type PasswordHasher interface {
	// Hash hashes a plain text password
	Hash(password string) (string, error)

	// Compare compares a plain text password with a hash
	Compare(password, hash string) error
}
