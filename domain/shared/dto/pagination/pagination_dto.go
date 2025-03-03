package pagination

type (
	PaginationRequest struct {
		Size int `query:"ps"`

		// pagination 방식
		Page int `query:"p"`

		// 커서 방식
		Cursor string `query:"cursor"`
	}

	PaginationMeta struct {
		Size int `json:"size"`

		// pagination 방식
		Page      int `json:"page"`
		TotalPage int `json:"totalPage"`

		// cursor 방식
		NextCursor string `json:"nextCursor"`
	}

	PaginationQuery[Q any] struct {
		QueryOption Q
		PageRequest PaginationRequest
		Order       string
	}
)
