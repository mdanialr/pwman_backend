package password

import (
	"context"
	"mime/multipart"

	pw "github.com/mdanialr/pwman_backend/internal/domain/password"
)

// UseCase signature that's used in password domain for use case layer.
type UseCase interface {
	// IndexCategory retrieve all category information including the url to
	// both image and icon.
	IndexCategory(ctx context.Context, req pw.RequestCategory) (*pw.IndexResponse[pw.ResponseCategory], error)
	// SaveCategory create new category from given request including the binary
	// files for both image and icon fields.
	SaveCategory(ctx context.Context, req pw.RequestCategory) (*pw.ResponseCategory, error)
	// UpdateCategory update existing Category that match given id. Optionally
	// replace either or both Image & Icon fields if provided.
	UpdateCategory(ctx context.Context, id uint, req pw.RequestCategory) error
	// SaveFile store given multipart to storage.Port then return filename of
	// the stored file that's ready to be saved. Optionally append given
	// prefix path too.
	SaveFile(f *multipart.FileHeader, prefix ...string) (string, error)
	// RemoveFile remove given filename using storage.Port and just log if
	// there is any error.
	RemoveFile(fn string)
}
