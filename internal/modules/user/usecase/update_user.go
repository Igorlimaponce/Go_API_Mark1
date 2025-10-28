package usecase

import (
	"context"
	"errors"
	"rest-api_mark1/internal/modules/user/models"
	"rest-api_mark1/internal/modules/user/repository"
	"time"

	"github.com/google/uuid"
)

type UpdateUserUseCase struct {
	userRepo repository.UserRepository
}

type UpdateUserInput struct {
	Name     *string          `json:"name,omitempty"`
	Email    *string          `json:"email,omitempty"`
	Role     *models.UserRole `json:"role,omitempty"`
	IsActive *bool            `json:"is_active,omitempty"`
}

type UpdateUserOutput struct {
	ID        string
	Name      string
	Email     string
	Role      models.UserRole
	IsActive  bool
	CreatedAt string
	UpdatedAt string
}

func NewUpdateUserUseCase(useRepo repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo: useRepo,
	}
}

func (uc *UpdateUserUseCase) Exec(ctx context.Context, id uuid.UUID, input *UpdateUserInput) error {
	if err := uc.validateInput(input); err != nil {
		return err
	}

	user, err := uc.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if input.Name != nil {
		if err := models.validateName(*input.Name); err != nil {
			return err
		}
		user.Name = *input.Name
	}

	if input.Email != nil {
		if err := models.validateEmail(*input.Email); err != nil {
			return err
		}
		user.Email = *input.Email
	}

	if input.Role != nil {
		if err := models.validateRole(*input.Role); err != nil {
			return err
		}
		user.Role = *input.Role
	}

	if err := models.ValidateIsActive(input.IsActive); err != nil {
		return err
	}

	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, &user); err != nil {
		return err
	}
	return nil
}

func (uc *UpdateUserUseCase) validateInput(input *UpdateUserInput) error {
	if input.Name == nil && input.Email == nil && input.Role == nil && input.IsActive == nil {
		return errors.New("at least one field must be provided for update")
	}
	return nil
}
