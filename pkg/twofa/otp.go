package twofa

import (
	"errors"
	"log"
	"strings"

	conf "github.com/mdanialr/pwman_backend/pkg/config"
	"github.com/mdanialr/pwman_backend/pkg/otp"

	"github.com/spf13/viper"
)

// GenerateQR generate QR code in bytes.
func GenerateQR() []byte {
	otpObj, err := initOTP(nil)
	if err != nil {
		log.Fatalln("failed to init otp:", err)
	}
	qr, err := otp.NewQR(otpObj.CreateURI())
	if err != nil {
		log.Fatalln("failed to generate QR code:", err)
	}
	return qr
}

// Verify do verify given code either using TOTP ot HOTP based on the config.
func Verify(code string) bool {
	otpObj, err := initOTP(nil)
	if err != nil {
		log.Fatalln("failed to init otp:", err)
	}
	// verify
	valid, err := otpObj.VerifyCode(code)
	if err != nil {
		log.Fatalln("failed to verify given code:", err)
	}
	return valid
}

// InitOTPWithConfig init new otp with given viper instance to retrieve the
// credential and the otp type whether it's hotp or totp.
func InitOTPWithConfig(v *viper.Viper) (*otp.OTP, error) {
	return initOTP(v)
}

// initOTP return pointer to otp.OTP which already initialized either using
// TOTP or HOTP based on the config.
func initOTP(v *viper.Viper) (*otp.OTP, error) {
	// init new viper config
	newV, err := conf.InitConfigYml()
	if err != nil {
		return nil, err
	}
	// if the provided viper is not nil then use that instead
	if v != nil {
		newV = v
	}
	// retrieve the secret
	secret := newV.GetString("cred.secret")
	// then the otp type
	otpType := newV.GetString("cred.type")
	// decide the otp type
	var otpObj *otp.OTP
	switch strings.ToLower(otpType) {
	case "hotp":
		otpObj = otp.NewHOTP(secret)
	case "totp":
		otpObj = otp.NewTOTP(secret)
	}
	// throw error if it's unsupported otp type
	if otpObj == nil {
		return nil, errors.New("unsupported otp type. should be either totp or hotp")
	}
	return otpObj, nil
}
