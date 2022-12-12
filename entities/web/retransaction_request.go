package web

type StanFilter struct {
	Stan string `form:"stan" binding:"required"`
}

type RecreateApexRequest struct {
	KodeLKM string `form:"kode_lkm" binding:"required"`
	Stan    string `form:"stan" binding:"required"`
}

type KuitansiFilter struct {
	Kuitansi string `form:"kuitansi"`
}

type ChangeResponseCode struct {
	Stan string `form:"stan" binding:"required"`
	RC   string `form:"response_code" binding:"required"`
}

type UpdateIsoMsg struct {
	Stan    string `form:"stan" binding:"required"`
	Iso_Msg string `form:"iso_message" binding:"required"`
}
