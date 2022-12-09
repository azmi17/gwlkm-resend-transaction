package web

type StanFilter struct {
	Stan string `form:"stan"`
}

type RecreateApexRequest struct {
	KodeLKM string `form:"kode_lkm" binding:"required"`
	Stan    string `form:"stan" binding:"required"`
}

type KuitansiFilter struct {
	Kuitansi string `form:"kuitansi"`
}

type ChangeResponseCode struct {
	Stan string `form:"stan"`
	RC   string `form:"response_code"`
}

type UpdateIsoMsg struct {
	Stan    string `form:"stan"`
	Iso_Msg string `form:"iso_message"`
}
