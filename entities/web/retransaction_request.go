package web

type StanFilter struct {
	Stan string `form:"stan"`
}

type ChangeResponseCode struct {
	Stan string `form:"stan"`
	RC   string `form:"response_code"`
}

type UpdateIsoMsg struct {
	Stan    string `form:"stan"`
	Iso_Msg string `form:"iso_message"`
}