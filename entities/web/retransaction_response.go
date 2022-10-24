package web

type RetransResponse struct {
	ResponseCode    string `json:"response_code"`
	ResponseMessage string `json:"response_message"`
}

type RetransWithNewStanResponse struct {
	ResponseCode    string `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	NewStan         string `json:"new_stan"`
}

type RetransTxInfo struct {
	Trans_Id      int    `json:"trans_id"`
	Idpel         string `json:"idpel"`
	Stan          string `json:"stan"`
	Ref_Stan      string `json:"ref_stan"`
	Tgl_Trans_Str string `json:"tgl_trans"`
	Response_Code string `json:"response_code"`
	BankCode      string `json:"kode_lkm"`
	RekID         string `json:"rek_id"`
	Amount        string `json:"amount"`
	Kuitansi      string `json:"kuitansi"`
	Iso_Msg_Resp  string `json:"iso_message_response"`
}

type TabtransInfoApx struct {
	Tgl_Trans    string  `json:"tgl_trans"`
	Kode_LKM     string  `json:"kode_lkm"`
	Kode_Trans   string  `json:"kode_trans"`
	Idpel        string  `json:"subscriber_id"`
	Kuitansi     string  `json:"kuitansi"`
	Pokok        float64 `json:"pokok"`
	Keterangan   string  `json:"keterangan"`
	Product_Code string  `json:"product_code"`
	Biller_Code  string  `json:"biller_code"`
	User_Id      int     `json:"user_id"`
	// Tgl_real_trans string  `json:"tgl_real_trans"`
}
