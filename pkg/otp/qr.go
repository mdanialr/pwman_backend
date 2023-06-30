package otp

import (
	"fmt"

	"rsc.io/qr"
)

// NewQR creates a new QR PNG from an OTP URI.
func NewQR(uri string) ([]byte, error) {
	code, err := qr.Encode(uri, qr.Q)
	if err != nil {
		return nil, fmt.Errorf("failed to encode: %s", err)
	}

	return code.PNG(), nil
}
