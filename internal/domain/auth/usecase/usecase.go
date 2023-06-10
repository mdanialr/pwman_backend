package auth

import (
	"context"
	"time"

	cons "github.com/mdanialr/pwman_backend/internal/constant"
	"github.com/mdanialr/pwman_backend/internal/domain/auth"
	authRepo "github.com/mdanialr/pwman_backend/internal/domain/auth/repository"
	stderr "github.com/mdanialr/pwman_backend/internal/err"
	help "github.com/mdanialr/pwman_backend/pkg/helper"
	"github.com/mdanialr/pwman_backend/pkg/twofa"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewUseCase return concrete implementation of UseCase in auth domain.
func NewUseCase(conf *viper.Viper, zap *zap.Logger, repo authRepo.Repository) UseCase {
	return &useCase{conf: conf, zap: zap, repo: repo}
}

type useCase struct {
	conf *viper.Viper
	zap  *zap.Logger
	repo authRepo.Repository
}

func (u *useCase) ValidateOTP(ctx context.Context, req auth.Request) (*auth.Response, error) {
	// init totp from pkg using given config
	ot, err := twofa.InitOTPWithConfig(u.conf)
	if err != nil {
		u.zap.Error(help.Pad("failed to init otp with config from app:", err.Error()))
		return nil, stderr.NewUC(cons.DepsErr, cons.ErrInternalServer.Error())
	}

	// verify the validity of given totp code from request
	valid, err := ot.VerifyCode(req.Code)
	if err != nil {
		u.zap.Error(help.Pad("failed to verify otp with code", req.Code, "and error:", err.Error()))
	}

	// make sure otp never used before
	if valid {
		if ro, _ := u.repo.GetByCode(ctx, req.Code); ro != nil {
			// return false if it's exist in db
			if ro.ID != 0 {
				return nil, stderr.NewUC(cons.UsedOTP, cons.ErrUsedOTP.Error())
			}
			// delete all past records
			if err = u.repo.DeleteAll(ctx); err != nil {
				u.zap.Error(help.Pad("failed to delete all records of RegisteredCode:", err.Error()))
				return nil, stderr.NewUC(cons.DepsErr, cons.ErrInternalServer.Error())
			}
			// then save the recent one
			if _, err = u.repo.Create(ctx, req.Code); err != nil {
				u.zap.Error(help.Pad("failed to save new RegisteredCode:", err.Error()))
				return nil, stderr.NewUC(cons.DepsErr, cons.ErrInternalServer.Error())
			}
			// create new jwt
			return u.CreateJWT(ctx)
		}
	}

	return nil, stderr.NewUC(cons.InvalidOTP, cons.ErrInvalidOTP.Error())
}

func (u *useCase) CreateJWT(_ context.Context) (*auth.Response, error) {
	// count the token's expiry time
	dur, _ := time.ParseDuration(u.conf.GetString("jwt.duration") + "m")
	exp := time.Now().Add(dur)

	// prepare the claims
	claims := jwt.MapClaims{
		"exp": exp.Unix(),
	}

	// generate and sign the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	at, err := token.SignedString([]byte(u.conf.GetString("jwt.secret")))
	if err != nil {
		return nil, stderr.NewUC(cons.DepsErr, cons.ErrSigningToken.Error())
	}

	// give back the constructed response
	resp := &auth.Response{
		AccessToken: at,
		ExpiredAt:   exp,
	}
	return resp, nil
}
