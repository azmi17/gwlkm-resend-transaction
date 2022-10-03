package entities

type TransHisotry struct {
	MTI              string // MTI
	ProcessingCode   string // bit 003
	Amount           int    // bit 004
	MarkupTotal      int    // bit 005
	ProfitEx         int    // bit 006
	TransmissionTime string // bit 007 => Now
	ProfitIn         int    // bit 008
	// Qty                    int       // bit 009
	Stan         string // bit 011
	DateTime     string // bit 012
	Date         string // bit 013
	MerchantType string // bit 018
	PosCC        string // bit 026 (not used)
	BankCode     string // bit 032
	NoReference  string // bit 037
	// ServiceRestrictionCode string    // bit 039
	ResponseCode  string // bit 040 (not used)
	PhoneNumber   string // bit 041
	DeviceID      string // bit 042
	AccountNumber string // bit 043
	Msg           string // bit 047
	Pin           string // bit 061
	BillerCode    string // bit 100
	SubscriberID  string // bit 103
	ProductCode   string // bit 104

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
}
