package password_test

import (
	"testing"

	pwMock "github.com/mdanialr/pwman_backend/internal/domain/password/repository/mocks"
	strMock "github.com/mdanialr/pwman_backend/pkg/storage/mocks"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type (
	deps struct {
		config  *viper.Viper
		log     *zap.Logger
		storage *strMock.MockstoragePort
		repo    *pwMock.MockpasswordRepository
	}
	helperSetup struct {
		Dep deps
	}
)

func setupTestHelper(t *testing.T) *helperSetup {
	d := deps{
		config:  viper.New(),
		log:     zaptest.NewLogger(t),
		storage: new(strMock.MockstoragePort),
		repo:    new(pwMock.MockpasswordRepository),
	}

	return &helperSetup{
		Dep: d,
	}
}
