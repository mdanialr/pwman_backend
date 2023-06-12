package password

import (
	"context"

	cons "github.com/mdanialr/pwman_backend/internal/constant"
	"github.com/mdanialr/pwman_backend/internal/domain/password"
	pw "github.com/mdanialr/pwman_backend/internal/domain/password/repository"
	stderr "github.com/mdanialr/pwman_backend/internal/err"
	repo "github.com/mdanialr/pwman_backend/internal/repository"
	help "github.com/mdanialr/pwman_backend/pkg/helper"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewUseCase return concrete implementation of UseCase in password domain.
func NewUseCase(conf *viper.Viper, log *zap.Logger, repo pw.Repository) UseCase {
	return &useCase{conf: conf, log: log, repo: repo}
}

type useCase struct {
	conf *viper.Viper
	log  *zap.Logger
	repo pw.Repository
}

func (u *useCase) IndexCategory(ctx context.Context, req password.Request) (*password.IndexResponse, error) {
	// set up repo options
	opts := []repo.Options{repo.Paginate(&req.M), repo.Order(req.Order + " " + req.Sort)}
	// additionally add search option
	if req.Search != "" {
		q := "name ILIKE '%" + req.Search + "%'"
		opts = append(opts, repo.Cons(q))
	}

	// search for all categories that matched given conditions
	cats, err := u.repo.FindCategories(ctx, opts...)
	if err != nil {
		u.log.Error(help.Pad("failed to retrieve categories:", err.Error()))
		return nil, stderr.NewUCErr(cons.DepsErr, cons.ErrInternalServer)
	}

	// prepare the response to contain the actual data and the pagination info
	resp := password.NewIndexResponseFromEntity(cats, u.conf.GetString("storage.url"))
	resp.Pagination = &req.M
	resp.Pagination.Paginate()

	return resp, nil
}
