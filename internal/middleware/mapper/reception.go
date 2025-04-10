package mapper

import (
	"github.com/google/uuid"
	"pvz/internal/api/response"
	"pvz/internal/logger"
	"pvz/internal/repository/model"
)

func ToReception(req response.ReceptionRequest) model.Reception {
	pvzId, err := uuid.Parse(req.PvzId)
	if err != nil {
		logger.SugaredLogger.Errorw("Invalid PvzId format", "pvzId", req.PvzId, "error", err)
		return model.Reception{}
	}

	return model.Reception{
		PvzId: pvzId,
	}
}

func ToReceptionResponse(reception model.Reception) response.ReceptionResponse {
	return response.ReceptionResponse{
		Id:       reception.Id.String(),
		DateTime: reception.DateTime.Format("2006-01-02 15:04:05"),
		PvzId:    reception.PvzId.String(),
		Status:   reception.Status,
	}
}
