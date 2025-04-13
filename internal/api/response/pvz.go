package response

type PvzRequest struct {
	City string `json:"City"`
}

type PvzFullResponse struct {
	Pvz        PvzResponse        `json:"pvz"`
	Receptions []ReceptionWrapper `json:"receptions"`
}

type PvzResponse struct {
	Id               string `json:"id"`
	RegistrationDate string `json:"registrationDate"`
	City             string `json:"city"`
}

type ReceptionWrapper struct {
	Reception ReceptionResponse `json:"reception"`
	Products  []ProductResponse `json:"products"`
}
