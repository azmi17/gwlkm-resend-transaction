package helper

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func ApiResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResp := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResp
}

var (
	AlreadyTransacted = "Transaksi sudah dilakukan"
	AlreadyReversed   = "44-Transaksi Sudah di reversal!"
)
