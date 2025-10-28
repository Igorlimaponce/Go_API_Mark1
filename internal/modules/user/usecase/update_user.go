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
	Name     string
	Email    string
	Role     models.UserRole
	IsActive bool
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

func (uu *UpdateUserUseCase) Exec(ctx context.Context, id uuid.UUID, input *UpdateUserInput) error {
	if err := uu.validate(input); err != nil {
		return err
	}

	user, err := uu.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	if input.Name != "" {
		if err := uu.validateName(input.Name); err != nil {
			return err
		}
		user.Name = input.Name
	}

	if input.Email != "" {
		if err := uu.validateEmail(input.Email); err != nil {
			return err
		}
		user.Email = input.Email
	}

	if input.Role != "" {
		if err := uu.validateRole(input.Role); err != nil {
			return err
		}
		user.Role = input.Role
	}

	if err := uu.validateIsActive(input.IsActive); err != nil {
		return err
	}
	user.IsActive = input.IsActive

	user.UpdatedAt = time.Now()

	if err := uu.userRepo.Update(ctx, &user); err != nil {
		return err
	}
	return nil
}

func (uu *UpdateUserUseCase) validate(input *UpdateUserInput) error {
	// Implementar validações necessárias para atualização de usuário
	if input.Name == "" {
		return errors.New("name is required")
	}
	if input.Email == "" {
		return errors.New("email is required")
	}
	return nil
}

func (uu *UpdateUserUseCase) validateName(name string) error {
	if len(name) < 3 || len(name) > 100 {
		return errors.New("name must be between 3 and 100 characters")
	}
	if name == "" {
		return errors.New("name cannot be empty")
	}
	return nil
}

func (uu *UpdateUserUseCase) validateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func (uu *UpdateUserUseCase) validateRole(role models.UserRole) error {
	if role != models.Admin && role != models.User {
		return errors.New("invalid role")
	}
	return nil
}

func (uu *UpdateUserUseCase) validateIsActive(isActive bool) error {
	if isActive != true && isActive != false {
		return errors.New("isActive must be a boolean value")
	}
	return nil
}
