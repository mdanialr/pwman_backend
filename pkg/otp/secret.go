package otp

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
)

// REF: https://datatracker.ietf.org/doc/html/rfc3548#section-5
const secretChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ234567"

// NewSecret creates a Base32 encoded arbitrary secret from a fixed length of 16
// byte slice without having a padding sign `=` at the end.
func NewSecret() (string, error) {
	var bytes = make([]byte, 16)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("read: %w", err)
	}

	for i, b := range bytes {
		bytes[i] = secretChars[b%byte(len(secretChars))]
	}

	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes), nil
}
