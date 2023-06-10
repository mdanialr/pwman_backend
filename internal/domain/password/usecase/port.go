package password

import (
	"context"

	pw "github.com/mdanialr/pwman_backend/internal/domain/password"
)

// UseCase signature that's used in password domain for use case layer.
type UseCase interface {
	// Index retrieve all category information including the url to both image
	// and icon.
	Index(ctx context.Context, req pw.Request) (*pw.IndexResponse, error)
}
