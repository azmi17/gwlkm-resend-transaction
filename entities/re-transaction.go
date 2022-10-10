package entities

type TransHistory struct {
	Stan                        string
	Ref_Stan                    string
	Tgl_Trans_Str               string
	Bank_Code                   string
	Rek_Id                      string
	Mti                         string
	Processing_Code             string
	Biller_Code                 string
	Product_Code                string
	Subscriber_Id               string
	Dc                          string
	Response_Code               string
	Amount                      int
	Qty                         int
	Profit_Included             int
	Profit_Excluded             int
	Profit_Share_Biller         int
	Profit_Share_Aggregator     int
	Profit_Share_Bank           int
	Markup_Total                int
	Markup_Share_Aggregator     int
	Markup_Share_Bank           int
	Msg                         string
	Msg_Response                string
	Bit39_Bit48_Hulu            string
	Saldo_Before_Trans          int
	Keterangan                  string
	Ref                         string
	Synced_Ibs_Core             int
	Synced_Ibs_Core_Description string
	Bris_Original_Data          string
	Gateway_Id                  int
	Id_User                     int
	Id_Raw                      int
	Advice_Count                int
	Status_Id                   int
	Nohp_Notif                  string
	Score                       int
	No_Hp_Alternatif            string
	Inc_Notif_Status            int
	Fee_Rek_Induk               int
}

type MsgTransHistory struct {
	MTI              string // MTI
	ProcessingCode   string // bit 003
	Amount           int    // bit 004
	MarkupTotal      int    // bit 005
	ProfitEx         int    // bit 006
	TransmissionTime string // bit 007
	ProfitIn         int    // bit 008
	Stan             string // bit 011
	DateTime         string // bit 012
	Date             string // bit 013
	MerchantType     string // bit 018
	PosCC            string // bit 026 (not used)
	BankCode         string // bit 032
	NoReference      string // bit 037
	ResponseCode     string // bit 040 (not used)
	PhoneNumber      string // bit 041
	DeviceID         string // bit 042
	AccountNumber    string // bit 043
	Msg              string // bit 047
	Pin              string // bit 061
	BillerCode       string // bit 100
	SubscriberID     string // bit 103
	ProductCode      string // bit 104

	Ref_Stan string
	Ref      string
}

type CoreAddr struct {
	BankCode string
	IPaddr   string
	TCPPort  int
}

type TransHistoryRequest struct {
	Stan string `form:"stan"`
}

type TransHistoryResponse struct {
	ResponseCode    string `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	Ref_Stan        string `json:"receipt"`
}
