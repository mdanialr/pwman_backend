package password

import (
	"strings"

	paginate "github.com/mdanialr/pwman_backend/pkg/pagination"
)

// pagination standard object that could be reused in password domain for every
// request that need pagination capability.
type pagination struct {
	paginate.M
	// Order the field name to query Order. Default to id.
	Order string `json:"-" query:"order"`
	// Sort to query Order. Should be filled with either asc or desc. Default
	// to asc.
	Sort string `json:"-" query:"sort"`
	// Search do search for category name and or username from given string.
	Search string `json:"-" query:"search"`
}

// SetQuery do setup Order and Sort.
func (p *pagination) SetQuery() {
	if p.Order == "" {
		p.Order = "id" // set default to id
	}
	// sanitize Sort
	p.Sort = p.sanitizeQuerySort()
	if p.Sort == "" {
		p.Sort = "asc" // set default to asc
	}
	// make sure the Sort is upper-cased
	p.Sort = strings.ToUpper(p.Sort)
}

// sanitizeQuerySort make sure Sort has the expected value.
func (p *pagination) sanitizeQuerySort() string {
	switch strings.ToLower(p.Sort) {
	case "asc", "desc":
		return p.Sort
	}
	return ""
}
