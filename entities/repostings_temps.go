package entities

// ===================================================================================================================================================
// ==================================================================TEMPORARY FUNCTIONS==============================================================
// ===================================================================================================================================================
type RepostingData struct {
	KodeLKM     string
	TotalDebet  float64
	TotalKredit float64
}

type CalculateSaldoResult struct {
	KodeLKM    string
	SaldoAkhir float64
}

type LKMlist struct {
	KodeLKM string
}

type SchedulerResponse struct {
	ResponseCode    string `json:"response_code"`
	ResponseMessage string `json:"response_message"`
}
