package response

type ProductRequest struct {
	Type  string `json:"Type"`
	PvzId string `json:"PvzId"`
}

type ProductResponse struct {
	Id          string `json:"Id"`
	DateTime    string `json:"DateTime"`
	Type        string `json:"Type"`
	ReceptionId string `json:"ReceptionId"`
}
