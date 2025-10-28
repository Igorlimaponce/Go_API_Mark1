package handler

import (
	"rest-api_mark1/internal/modules/user/usecase"
)

type UserHandler struct {
	createUserUseCase *usecase.CreateUserUseCase
}
