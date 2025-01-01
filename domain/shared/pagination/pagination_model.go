package pagination

type (
	PaginationRequest struct {
		Page int `json:"page"`
		Size int `json:"size"`

		// 총 count에 대한 정보가 있을 때는 쿼리하지 않는다. (쿼리 비용을 아끼기 위해)
		TotalCount int `json:"totalCount"`
	}

	PaginationMeta struct {
		Page int `json:"page"`
		Size int `json:"size"`

		TotalPage  int `json:"totalPage"`
		TotalCount int `json:"totalCount"`
	}
)
