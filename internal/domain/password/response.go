package password

import (
	"strings"

	"github.com/mdanialr/pwman_backend/internal/entity"
	paginate "github.com/mdanialr/pwman_backend/pkg/pagination"
)

// ResponseCategory standard response object that may be used in password domain.
type ResponseCategory struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Icon  string `json:"icon"`
}

// NewResponseFromEntity transform given entity.Category to Response. Also
// prepend given prefix to both Image & Icon fields after cleaning the trailing
// slash.
func NewResponseFromEntity(cat entity.Category, prefix string) *ResponseCategory {
	pr := strings.TrimSuffix(prefix, "/") + "/"

	r := &ResponseCategory{
		ID:    cat.ID,
		Name:  cat.Name,
		Image: pr + cat.ImagePath,
		Icon:  pr + cat.IconPath,
	}
	return r
}

// IndexResponse response that's used in use case Index.
type IndexResponse struct {
	Data       []*Response `json:"-"`
	Pagination *paginate.M
}

// NewIndexResponseFromEntity create new pointer IndexResponse from given slices
// of entity.Category. Also prepend given prefix to both Image & Icon fields
// after cleaning the trailing slash.
func NewIndexResponseFromEntity(cats []*entity.Category, prefix string) *IndexResponse {
	var res []*Response

	for _, cat := range cats {
		res = append(res, NewResponseFromEntity(*cat, prefix))
	}

	return &IndexResponse{Data: res}
}
