package usecase

import (
	"context"
	"fmt"
	"regexp"
	"rest-api_mark1/internal/modules/user/models"
	"rest-api_mark1/internal/modules/user/repository"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type CreateUserUseCase struct {
	userRepo repository.UserRepository
}

// NewCreateUserUseCase cria uma nova instância do caso de uso de criação de usuário.
// Recebe um UserRepository para interagir com a camada de dados.
func NewCreateUserUseCase(userRepo repository.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo: userRepo,
	}
}

// CreateUserInput representa os dados necessários para criar um usuário.
type CreateUserInput struct {
	Name     string          `json:"name" validate:"required,min=3,max=100"`
	Email    string          `json:"email" validate:"required,email"`
	Password string          `json:"password" validate:"required,min=8"`
	Role     models.UserRole `json:"role" validate:"required"`
}

// CreateUserOutput representa os dados retornados após a criação do usuário.
type CreateUserOutput struct {
	ID        uuid.UUID       `json:"id"`
	Name      string          `json:"name"`
	Email     string          `json:"email"`
	Role      models.UserRole `json:"role"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// Esse regex simples valida o formato básico de um email.
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// Exec executa a lógica de criação de um novo usuário.
// Valida os dados de entrada, criptografa a senha e persiste no repositório.
func (uc *CreateUserUseCase) Exec(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
	// Validações
	if err := uc.validate(input); err != nil {
		return nil, err
	}

	// Verifica se o email já existe
	// Esse GetUserByEmail deve ser implementado no UserRepository
	existingUser, err := uc.userRepo.GetUserByEmail(ctx, input.Email)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already in use")
	}

	// Criptografa a senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Cria o modelo do usuário
	user := &models.User{
		ID:        uuid.New(),
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Role:      input.Role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Persiste no banco
	if err := uc.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Retorna o output (sem a senha)
	return &CreateUserOutput{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// validate valida os dados de entrada
func (uc *CreateUserUseCase) validate(input CreateUserInput) error {
	if input.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(input.Name) < 3 {
		return fmt.Errorf("name must be at least 3 characters")
	}
	if input.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !emailRegex.MatchString(input.Email) {
		return fmt.Errorf("invalid email format")
	}
	if input.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(input.Password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if input.Role == "" {
		return fmt.Errorf("role is required")
	}
	return nil
}
