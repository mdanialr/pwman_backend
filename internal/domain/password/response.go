package password

import (
	"strings"

	"github.com/mdanialr/pwman_backend/internal/entity"
	paginate "github.com/mdanialr/pwman_backend/pkg/pagination"
)

// responseAble generic type that holds all standard Response that can be
// transformed from entity to IndexResponse.
type responseAble interface {
	ResponseCategory
}

// ResponseCategory standard response object that may be used in password domain.
type ResponseCategory struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
	Icon  string `json:"icon"`
}

// NewResponseCategoryFromEntity transform given entity.Category to
// ResponseCategory. Also prepend given prefix to both Image & Icon fields
// after cleaning the trailing slash.
func NewResponseCategoryFromEntity(cat entity.Category, prefix string) *ResponseCategory {
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
type IndexResponse[T responseAble] struct {
	Data       []*T `json:"-"`
	Pagination *paginate.M
}

// NewIndexResponseFromEntity create new pointer IndexResponse from given slices
// of entity.Category. Also prepend given prefix to both Image & Icon fields
// after cleaning the trailing slash.
func NewIndexResponseCategoryFromEntity(cats []*entity.Category, prefix string) *IndexResponse[ResponseCategory] {
	var res []*ResponseCategory

	for _, cat := range cats {
		res = append(res, NewResponseCategoryFromEntity(*cat, prefix))
	}

	return &IndexResponse[ResponseCategory]{Data: res}
}
