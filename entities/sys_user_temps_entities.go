package entities

import "time"

type SysDaftarUser struct {
	User_Id                int
	User_Name              string
	User_Password          string
	Nama_Lengkap           string
	Penerimaan             float32
	Pengeluaran            float32
	Unit_Kerja             string
	Jabatan                string
	User_Code              string
	Tgl_Expired            time.Time
	Flag                   int
	Status_Aktif           int
	User_Web_Password      string
	User_Web_Password_Hash string
}

type KodeLKMFilter struct {
	KodeLkm string `form:"kode_lkm" binding:"required"`
}

type ResetApexPwdResponse struct {
	KodeLkm       string `json:"kode_lkm"`
	Password_Smec string `json:"new_apex_password"`
}

type ResetPasswrodResponse struct {
	Response_Code string                `json:"response_code"`
	Response_Msg  string                `json:"response_message"`
	Data          *ResetApexPwdResponse `json:"data,omitempty"`
}
