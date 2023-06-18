package cons

import "errors"

var (
	ErrUsedOTP        = errors.New("otp has been used")
	ErrInvalidOTP     = errors.New("invalid otp")
	ErrInternalServer = errors.New("something wasn't right")
	ErrSigningToken   = errors.New("signing token error")
	ErrAlreadyExist   = errors.New("data is already exist")
	ErrNotFound       = errors.New("data not found")
	ErrDataInUse      = errors.New("data still in use")
)
