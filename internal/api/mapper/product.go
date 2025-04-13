package mapper

import (
	"pvz/internal/api/response"
	"pvz/internal/repository/model"
)

func ToProduct(req response.ProductRequest) model.Product {
	return model.Product{
		Type: req.Type,
	}
}

func ToProductResponse(product model.Product) response.ProductResponse {
	return response.ProductResponse{
		Id:          product.Id.String(),
		DateTime:    product.DateTime.Format("2006-01-02 15:04:05"),
		Type:        product.Type,
		ReceptionId: product.ReceptionId.String(),
	}
}
