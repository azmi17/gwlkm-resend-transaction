package entities

import "time"

type TransHistory struct {
	Trans_id                    int
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

type IsoMessageBody struct {
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

type CoreAddrInfo struct {
	BankCode string
	IPaddr   string
	TCPPort  int
}

type StanReferences struct {
	Trans_id int
	Ref_Stan string
	Stan     string
}

type TransApx struct {
	Tabtrans_id      int
	Tgl_trans        time.Time
	No_rekening      string
	Kode_trans       string
	My_kode_trans    string
	Pokok            float64
	Kuitansi         string
	Userid           int
	Keterangan       string
	Verifikasi       string
	Tob              string
	Sandi_trans      string
	Posted_to_gl     string
	Kode_kantor      string
	Jam              string
	Tgl_real_trans   time.Time
	Pay_lkm_source   string
	Pay_lkm_norek    string
	Pay_idpel        string
	Pay_biller_code  string
	Pay_product_code string
}
