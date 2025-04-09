package response

type PvzRequest struct {
	City string
}

type PvzResponse struct {
	Id               string
	RegistrationDate string
	City             string
}
