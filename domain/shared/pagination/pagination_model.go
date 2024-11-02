package pagination

type (
	PaginationMeta struct {
		Page      int `json:"page"`
		Size      int `json:"size"`
		TotalPage int `json:"totalPage"`
	}
)
