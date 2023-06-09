package paginate

import (
	"math"

	"gorm.io/gorm"
)

// M standard pagination object that should be embedded to request/response
// object that need pagination feature.
type M struct {
	// Limit set the number of data that should be retrieved from DB. This may
	// be filled from query param in request.
	Limit int `json:"per_page,omitempty" query:"limit" example:"1"`
	// Page based on given Limit which page of data we want to retrieve.
	Page int `json:"current_page,omitempty" query:"page" example:"1"`
	// Next the next page number from current page.
	Next int `json:"next_page,omitempty"`
	// Prev the previous page number from current page.
	Prev int `json:"previous_page,omitempty"`
	// TotalPage based on given Limit how many page that can be divided from
	// total available data.
	TotalPage int `json:"total_page,omitempty"`
	// count just a placeholder to count the total retrieved number of data.
	count int64
}

// Paginate setup current, previous and next page based on the count & Limit
// fields. This should be called after retrieving the actual data from DB.
func (m *M) Paginate() {
	m.TotalPage = int(math.Ceil(float64(m.count) / float64(m.Limit)))
	if m.TotalPage > 0 {
		switch {
		case m.Page == 1 && m.TotalPage > 1: // 1st page
			m.Next = m.Page + 1
		case m.Page > 1 && m.Page < m.TotalPage: // 2nd page until before last page
			m.Next = m.Page + 1
			m.Prev = m.Page - 1
		case m.Page == m.TotalPage && m.TotalPage > 1: // last page
			m.Prev = m.Page - 1
		}
	}
	if m.TotalPage <= 0 {
		m.TotalPage = 1
	}
}

// Set implementation of repository.IOptions.
func (m *M) Set(db *gorm.DB) *gorm.DB {
	db.Count(&m.count)

	if m.Page == 0 { // set to first page if only has one page
		m.Page = 1
	}
	if m.Limit != 0 {
		db = db.Limit(m.Limit)
	}
	db = db.Offset((m.Page - 1) * m.Limit)

	return db
}
