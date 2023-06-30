package password_test

import (
	"context"
	"errors"
	"testing"

	"github.com/mdanialr/pwman_backend/internal/domain/password/repository/mocks"
	password "github.com/mdanialr/pwman_backend/internal/domain/password/usecase"
	"github.com/mdanialr/pwman_backend/internal/entity"
	stderr "github.com/mdanialr/pwman_backend/internal/err"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUseCase_DeletePassword(t *testing.T) {
	testCases := []struct {
		name       string
		setup      func(repo *mocks.MockpasswordRepository)
		sample     uint
		expectCode string
		expectMsg  string
		wantErr    bool
	}{
		{
			name: "Given id 201 that does not exist in deps repository should return UC instance, " +
				"INVALID_PAYLOAD as code and data not found as message",
			setup: func(repo *mocks.MockpasswordRepository) {
				repo.EXPECT().
					GetPasswordByID(mock.Anything, uint(201), mock.Anything).
					Return(nil, errors.New("error")).
					Once()
			},
			sample:     201,
			expectCode: "INVALID_PAYLOAD",
			expectMsg:  "data not found",
			wantErr:    true,
		},
		{
			name: "Given id 12 that does exist in deps repository but deps somehow failed to delete" +
				" the record should return UC instance, DEPS_ERROR as code and something wasn't " +
				"right as message",
			setup: func(repo *mocks.MockpasswordRepository) {
				obj := entity.Password{ID: 12}
				repo.EXPECT().
					GetPasswordByID(mock.Anything, obj.ID, mock.Anything).
					Return(&obj, nil).
					Once()
				repo.EXPECT().
					DeletePassword(mock.Anything, obj.ID).
					Return(errors.New("error")).
					Once()
			},
			sample:     12,
			expectCode: "DEPS_ERROR",
			expectMsg:  "something wasn't right",
			wantErr:    true,
		},
		{
			name: "Given id 5 that does exist in deps repository but deps somehow failed to delete" +
				" the record should return UC instance, DEPS_ERROR as code and something wasn't " +
				"right as message",
			setup: func(repo *mocks.MockpasswordRepository) {
				obj := entity.Password{ID: 5}
				repo.EXPECT().
					GetPasswordByID(mock.Anything, obj.ID, mock.Anything).
					Return(&obj, nil).
					Once()
				repo.EXPECT().
					DeletePassword(mock.Anything, obj.ID).
					Return(nil).
					Once()
			},
			sample: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			h := setupTestHelper(t)
			tc.setup(h.Dep.repo)

			newUC := password.NewUseCase(h.Dep.config, h.Dep.log, h.Dep.storage, h.Dep.repo)
			err := newUC.DeletePassword(context.Background(), tc.sample)

			if tc.wantErr {
				assert.Error(t, err)
				// assert error instance
				assert.IsType(t, &stderr.UC{}, err)
				// assert Code and Message
				assert.NotPanics(t, func() {
					stdErrUC := err.(*stderr.UC)
					assert.Equal(t, tc.expectCode, stdErrUC.Code)
					assert.Equal(t, tc.expectMsg, stdErrUC.Msg)
				})
				return
			}

			assert.NoError(t, err)
			assert.Panics(t, func() {
				_ = err.(*stderr.UC)
			})
		})
	}
}
