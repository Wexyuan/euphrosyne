package base

const (
	// defaultPage is the default page number.
	defaultPage = 1
	// minPageSize is the smallest page size.
	minPageSize = 10
	// maxPageSize is the largest page size.
	maxPageSize = 100
)

// PageQuery holds the pagination input.
type PageQuery struct {
	// Page is the requested page number.
	Page int `json:"page"`
	// PageSize is the requested page size.
	PageSize int `json:"page_size"`
}

// resolve clamps the query to a valid page number and page size.
func (q PageQuery) resolve() PageQuery {
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

// PageResult holds the pagination output.
type PageResult[T any] struct {
	// Total is the total number of entities.
	Total int64 `json:"total"`
	// Page is the current page number.
	Page int `json:"page"`
	// PageSize is the rows per page.
	PageSize int `json:"page_size"`
	// List holds the entities of the current page.
	List []*T `json:"list"`
}
