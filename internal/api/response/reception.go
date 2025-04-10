package response

type ReceptionRequest struct {
	PvzId string `json:"pvzId"`
}

type ReceptionResponse struct {
	Id       string `json:"Id"`
	DateTime string `json:"DateTime"`
	PvzId    string `json:"PvzId"`
	Status   string `json:"Status"`
}
