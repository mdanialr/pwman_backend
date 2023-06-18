package password

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	cons "github.com/mdanialr/pwman_backend/internal/constant"
	"github.com/mdanialr/pwman_backend/internal/domain/password"
	pw "github.com/mdanialr/pwman_backend/internal/domain/password/repository"
	"github.com/mdanialr/pwman_backend/internal/entity"
	stderr "github.com/mdanialr/pwman_backend/internal/err"
	repo "github.com/mdanialr/pwman_backend/internal/repository"
	help "github.com/mdanialr/pwman_backend/pkg/helper"
	"github.com/mdanialr/pwman_backend/pkg/storage"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// NewUseCase return concrete implementation of UseCase in password domain.
func NewUseCase(conf *viper.Viper, log *zap.Logger, st storage.Port, repo pw.Repository) UseCase {
	return &useCase{conf: conf, log: log, st: st, repo: repo}
}

type useCase struct {
	conf *viper.Viper
	log  *zap.Logger
	st   storage.Port
	repo pw.Repository
}

func (u *useCase) IndexCategory(ctx context.Context, req password.RequestCategory) (*password.IndexResponse[password.ResponseCategory], error) {
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
	resp := password.NewIndexResponseCategoryFromEntity(cats, u.conf.GetString("storage.url"))
	resp.Pagination = &req.M
	resp.Pagination.Paginate()

	return resp, nil
}

func (u *useCase) SaveCategory(ctx context.Context, req password.RequestCategory) (*password.ResponseCategory, error) {
	// make sure given category name not used yet in data store
	cond := "name = '" + req.Name + "'"
	c, _ := u.repo.GetCategoryByID(ctx, 0, repo.Cols("id"), repo.Cons(cond))
	// return error if already exist
	if c.ID != 0 {
		return nil, stderr.NewUCErr(cons.InvalidPayload, cons.ErrAlreadyExist)
	}

	// save multipart icon to Storage
	ico, err := u.SaveFile(req.Icon)
	if err != nil {
		u.log.Error(help.Pad("failed to save icon:", req.Icon.Filename, "with err:", err.Error()))
		return nil, stderr.NewUCErr(cons.DepsErr, cons.ErrInternalServer)
	}
	// save multipart image to Storage
	img, err := u.SaveFile(req.Image)
	if err != nil {
		u.log.Error(help.Pad("failed to save image:", req.Image.Filename, "with err:", err.Error()))
		return nil, stderr.NewUCErr(cons.DepsErr, cons.ErrInternalServer)
	}

	// save the category to data store
	obj := entity.Category{
		Name:      req.Name,
		IconPath:  ico,
		ImagePath: img,
	}
	newObj, err := u.repo.CreateCategory(ctx, obj)
	if err != nil {
		u.log.Error(help.Pad("failed to create new category:", err.Error()))
		return nil, stderr.NewUCErr(cons.DepsErr, cons.ErrInternalServer)
	}

	// adapt to appropriate response
	resp := password.NewResponseCategoryFromEntity(*newObj, u.conf.GetString("storage.url"))

	return resp, nil
}

func (u *useCase) UpdateCategory(ctx context.Context, id uint, req password.RequestCategory) error {
	// retrieve category from repo using given id
	c, _ := u.repo.GetCategoryByID(ctx, id)
	// do additional validation if the name from request and from repo is different
	if c.Name != req.Name {
		// make sure it's unique and not taken yet
		cond := "name = '" + req.Name + "'"
		oldC, _ := u.repo.GetCategoryByID(ctx, 0, repo.Cols("id"), repo.Cons(cond))
		// return error if already exist
		if oldC.ID != 0 {
			return stderr.NewUCErr(cons.InvalidPayload, cons.ErrAlreadyExist)
		}
	}

	// record the updated fields
	updatedFields := []string{"name"}
	newCategory := entity.Category{Name: req.Name}

	// update Image if provided
	if req.Image != nil {
		updatedFields = append(updatedFields, "image_path")
		// save new image
		img, err := u.SaveFile(req.Image)
		if err != nil {
			u.log.Error(help.Pad("failed to save image:", req.Image.Filename, "with err:", err.Error()))
			return stderr.NewUCErr(cons.DepsErr, cons.ErrInternalServer)
		}
		newCategory.ImagePath = img
	}
	// update Icon if provided
	if req.Icon != nil {
		updatedFields = append(updatedFields, "icon_path")
		// save new icon
		icon, err := u.SaveFile(req.Icon)
		if err != nil {
			u.log.Error(help.Pad("failed to save image:", req.Icon.Filename, "with err:", err.Error()))
			return stderr.NewUCErr(cons.DepsErr, cons.ErrInternalServer)
		}
		newCategory.IconPath = icon
	}

	if _, err := u.repo.UpdateCategory(ctx, c.ID, newCategory, repo.Cols(updatedFields...)); err != nil {
		u.log.Error(help.Pad("failed to update existing category with id:", strconv.Itoa(int(c.ID)), "and err:", err.Error()))
		return stderr.NewUCErr(cons.DepsErr, cons.ErrInternalServer)
	}

	// lastly remove the old image & icon
	go u.removeOldMedia(*c, updatedFields...)

	return nil
}

func (u *useCase) SaveFile(f *multipart.FileHeader, prefix ...string) (string, error) {
	fl, err := f.Open()
	if err != nil {
		return "", err
	}
	defer fl.Close()

	// set up the target path
	pt := strings.TrimSuffix(u.conf.GetString("storage.path"), "/") + "/"
	// generate random name then append it with the file extension
	fn := uuid.NewString() + filepath.Ext(f.Filename)

	// optionally append given prefix
	if len(prefix) > 0 {
		fn = strings.TrimSuffix(prefix[0], "/") + "/" + fn
	}

	// save in separate goroutine
	go u.st.Store(fl, pt+fn)

	return fn, nil
}

func (u *useCase) RemoveFile(fn string) {
	pt := strings.TrimSuffix(u.conf.GetString("storage.path"), "/") + "/"
	u.st.Remove(pt + fn)
}

// removeOldMedia do call RemoveFile method for each media's field name that
// exist in given fields.
func (u *useCase) removeOldMedia(c entity.Category, fields ...string) {
	for _, field := range fields {
		switch field {
		case "image_path":
			u.RemoveFile(c.ImagePath)
		case "icon_path":
			u.RemoveFile(c.IconPath)
		}
	}
}
