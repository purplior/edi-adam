package response

type (
	JSONResponse struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data"`
	}
)

func (r JSONResponse) ErrorMessage() string {
	if r.IsOK() {
		return ""
	}

	return "일시적인 서버 오류가 발생 했습니다."
}

func (r JSONResponse) IsOK() bool {
	return r.Status >= 200 && r.Status < 300
}
