package mapper

import (
	"pvz/internal/api/response"
	"pvz/internal/repository/model"
)

func ToPvz(pvz response.PvzRequest) model.Pvz {
	return model.Pvz{
		City: pvz.City,
	}
}

func ToPvzResponse(pvz model.Pvz) response.PvzResponse {
	return response.PvzResponse{
		Id:               pvz.Id.String(),
		RegistrationDate: pvz.RegistrationDate.Format("2006-01-02 15:04:05"),
		City:             pvz.City,
	}
}
