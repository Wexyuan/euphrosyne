package base

const (
	// defaultPage is the default page number.
	defaultPage = 1
	// minPageSize is the minimum page size.
	minPageSize = 10
	// maxPageSize is the maximum page size.
	maxPageSize = 100
)

// PageQuery holds the pagination query.
type PageQuery struct {
	Page     int `json:"page"`      // page number
	PageSize int `json:"page_size"` // page size
}

// clamp limits the page query to valid bounds.
func (q PageQuery) clamp() PageQuery {
	if q.Page < defaultPage {
		q.Page = defaultPage
	}
	if q.PageSize < minPageSize {
		q.PageSize = minPageSize
	}
	if q.PageSize > maxPageSize {
		q.PageSize = maxPageSize
	}
	return q
}

// PageResult holds the pagination result.
type PageResult[T any] struct {
	Total    int64 `json:"total"`     // total row count
	Page     int   `json:"page"`      // current page number
	PageSize int   `json:"page_size"` // rows per page
	List     []*T  `json:"list"`      // current page rows
}
