package mapper

import (
	"pvz/internal/api/response"
	"pvz/internal/repository/model"
)

func ToUser(req response.RegisterPostRequest) model.User {
	return model.User{
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}
}

func ToRegisterResponse(user model.User) response.RegisterResponse {
	return response.RegisterResponse{
		Id:    user.Id.String(),
		Email: user.Email,
		Role:  user.Role,
	}
}
