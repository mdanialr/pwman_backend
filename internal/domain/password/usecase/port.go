package password

import (
	"context"

	pw "github.com/mdanialr/pwman_backend/internal/domain/password"
)

// UseCase signature that's used in password domain for use case layer.
type UseCase interface {
	// IndexCategory retrieve all category information including the url to both image
	// and icon.
	IndexCategory(ctx context.Context, req pw.Request) (*pw.IndexResponse, error)
}
