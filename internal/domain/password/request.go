package password

import (
	"strings"

	paginate "github.com/mdanialr/pwman_backend/pkg/pagination"
)

// Request standard request object that may be used in password domain.
type Request struct {
	paginate.M
	// Order the field name to query Order. Default to id.
	Order string `json:"-" query:"order"`
	// Sort to query Order. Should be filled with either asc or desc. Default
	// to asc.
	Sort string `json:"-" query:"sort"`
	// Search do search for category name from given string.
	Search string `json:"-" query:"search"`
}

// SetQuery do setup Order and Sort.
func (r *Request) SetQuery() {
	if r.Order == "" {
		r.Order = "id" // set default to id
	}
	// sanitize Sort
	r.Sort = r.sanitizeQuerySort()
	if r.Sort == "" {
		r.Sort = "asc" // set default to asc
	}
	// make sure the Sort is upper-cased
	r.Sort = strings.ToUpper(r.Sort)
}

// sanitizeQuerySort make sure Sort has the expected value.
func (r *Request) sanitizeQuerySort() string {
	switch strings.ToLower(r.Sort) {
	case "asc", "desc":
		return r.Sort
	}
	return ""
}
