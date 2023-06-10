// Package otp
//
// TOTP: https://en.wikipedia.org/wiki/One-time_password
//
//	https://datatracker.ietf.org/doc/html/rfc6238
//
// HOTP: https://en.wikipedia.org/wiki/HMAC-based_one-time_password
//
//	https://datatracker.ietf.org/doc/html/rfc4226
//
// Copyright 2021 Inanzzz.com
package otp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"time"
)

const (
	// length defines the OTP code in character length.
	length = 6
	// period defines the TTL of a TOTP code in seconds.
	period = 30
	// issuer the issuer name.
	issuer = "Password Manager API"
	// account the account for this otp, currently only one account, so it's
	// safe to has a placeholder.
	account = "admin"
)

type OTP struct {
	// issuer represents the service provider. It is you! e.g. your service,
	// your application, your organisation so on.
	issuer string
	// account represents the service user. It is the user! e.g. username, email
	// address so on.
	account string
	// secret is an arbitrary key value encoded in Base32 and belongs to the
	// service user.
	secret string
	// window is used for time (TOTP) and counter (HOTP) synchronization. Given
	// that the possible time and counter drifts between client and server, this
	// parameter helps overcome such issue. TOTP uses backward and forward time
	// window whereas HOTP uses look-ahead counter window that depends on the
	// counter parameter.
	// Resynchronisation is an official recommended practise, however the
	// lower, the better.
	// 0 = not recommended as synchronization is disabled
	//   TOTP: current time
	//   HOTP: current counter
	// 1 = recommended option
	//   TOTP: previous - current - next
	//   HOTP: current counter - next counter
	// 2 = being overcautious
	//   TOTP: previous,previous - current - next,next
	//   HOTP: current counter - next counter - next counter
	// * = Higher numbers may cause denial-of-service attacks.
	// REF: https://datatracker.ietf.org/doc/html/rfc6238#page-7
	// REF: https://datatracker.ietf.org/doc/html/rfc4226#page-11
	window int
	// counter is required for HOTP only and used for provisioning the code. Set
	// it to 0 if you with to use TOTP. Start from 1 for HOTP then fetch and use
	// the one in the persistent storage. The server counter is incremented only
	// after a successful code verification, however the counter on the code is
	// incremented every time a new code is requested by the user which causes
	// counters being out of sync. For that reason, time-synchronization should
	// be enabled.
	// REF: https://datatracker.ietf.org/doc/html/rfc4226#page-11
	counter int
}

// NewTOTP return new otp that's use Time-based One-Time Password. Good when
// used with mobile apps such as Microsoft Authenticator or Google
// Authenticator etc.
func NewTOTP(secret string) *OTP {
	return newOTP(secret, 0)
}

// NewHOTP return new otp that's use HMAC-based One-Time Password. Good when
// sending otp via sms or email. Not recommended yet for this app since this
// app does not support sending any email yet.
func NewHOTP(secret string) *OTP {
	return newOTP(secret, 1)
}

// newOTP base function that will create an OTP based on given secret and
// counter (for HOTP only).
func newOTP(secret string, counter int) *OTP {
	return &OTP{
		issuer:  issuer,
		account: account,
		secret:  secret,
		window:  0,
		counter: counter,
	}
}

// CreateURI builds the authentication URI which is used to create a QR code.
// If the counter is set to 0, the algorithm is assumed to be TOTP, otherwise
// HOTP.
// REF: https://github.com/google/google-authenticator/wiki/Key-Uri-Format
func (o *OTP) CreateURI() string {
	algorithm := "totp"
	counter := ""
	if o.counter != 0 {
		algorithm = "hotp"
		counter = fmt.Sprintf("&counter=%d", o.counter)
	}

	return fmt.Sprintf("otpauth://%s/%s:%s?secret=%s&issuer=%s%s",
		algorithm,
		o.issuer,
		o.account,
		o.secret,
		o.issuer,
		counter,
	)
}

// CreateHOTPCode creates a new HOTP with a specific counter. This method is
// ideal if you are planning to send manually created code via email, SMS etc.
// The user should not be present a QR code for this option otherwise there is
// a high possibility that the client and server counters will be out of sync,
// unless the user will be forced to rescan a newly generated QR with
// up-to-date counter value.
func (o *OTP) CreateHOTPCode(counter int) (string, error) {
	val, err := o.createCode(counter)
	if err != nil {
		return "", fmt.Errorf("create code: %w", err)
	}

	o.counter = counter
	return val, nil
}

// VerifyCode talks to an algorithm specific validator to verify the integrity
// of the code. If the counter is set to 0, the algorithm is assumed to be
// TOTP, otherwise HOTP.
func (o *OTP) VerifyCode(code string) (bool, error) {
	if len(code) != length {
		return false, fmt.Errorf("invalid length")
	}

	if o.counter != 0 {
		ok, err := o.verifyHOTP(code)
		if err != nil {
			return false, fmt.Errorf("verify HOTP: %w", err)
		}
		if !ok {
			return false, nil
		}
		return true, nil
	}

	ok, err := o.verifyTOTP(code)
	if err != nil {
		return false, fmt.Errorf("verify TOTP: %w", err)
	}
	if !ok {
		return false, nil
	}

	return true, nil
}

// verifyTOTP depending on the given windows size, we handle clock
// resynchronisation. If the window size is set to 0, resynchronisation is
// disabled, and we just use the current time. Otherwise, backward and forward
// window is taken into account as well.
func (o *OTP) verifyTOTP(code string) (bool, error) {
	curr := int(time.Now().UTC().Unix() / period)
	back := curr
	forw := curr
	if o.window != 0 {
		back -= o.window
		forw += o.window
	}

	for i := back; i <= forw; i++ {
		val, err := o.createCode(i)
		if err != nil {
			return false, fmt.Errorf("create code: %w", err)
		}
		if val == code {
			return true, nil
		}
	}

	return false, nil
}

// verifyHOTP depending on the given windows size, we handle counter
// resynchronisation. If the window size is set to 0, resynchronisation is
// disabled, and we just use the current counter. Otherwise, look-ahead counter
// window is used. When the look-ahead window is used, we calculate the next
// codes and determine if there is a match by utilising counter
// resynchronisation.
func (o *OTP) verifyHOTP(code string) (bool, error) {
	size := 0
	if o.window != 0 {
		size = o.window
	}

	for i := 0; i <= size; i++ {
		val, err := o.createCode(o.counter + i)
		if err != nil {
			return false, fmt.Errorf("create code: %w", err)
		}
		if val == code {
			o.counter += i + 1
			return true, nil
		}
	}

	o.counter++
	return false, nil
}

// createCode creates a new OTP code based on either a time or counter interval.
// The time is used for TOTP and the counter is used for HOTP algorithm.
func (o *OTP) createCode(interval int) (string, error) {
	sec, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(o.secret)
	if err != nil {
		return "", fmt.Errorf("decode string: %w", err)
	}

	hash := hmac.New(sha1.New, sec)
	if err := binary.Write(hash, binary.BigEndian, int64(interval)); err != nil {
		return "", fmt.Errorf("binary write: %w", err)
	}
	sign := hash.Sum(nil)

	offset := sign[19] & 15
	trunc := binary.BigEndian.Uint32(sign[offset : offset+4])

	return fmt.Sprintf("%0*d", length, (trunc&0x7fffffff)%1000000), nil
}
