package models

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	AdminRole  UserRole = "admin"
	CommonRole UserRole = "common"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey"`
	Name      string    `json:"name" gorm:"type:varchar(100);not null"`
	Email     string    `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"type:varchar(255);not null"` // Nunca retorna no JSON
	Role      UserRole  `json:"role" gorm:"type:varchar(20);not null;default:'common'"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName define o nome da tabela no banco de dados
func (User) TableName() string {
	return "users"
}

// NewUser cria um novo usuário com valores padrão
func NewUser(name, email, password string, role UserRole) (*User, error) {
	if err := ValidateEmail(email); err != nil {
		return nil, err
	}

	if err := ValidateName(name); err != nil {
		return nil, err
	}

	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	if err := ValidateRole(role); err != nil {
		return nil, err
	}

	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		Password:  password, // Deve ser hash na camada de usecase
		Role:      role,
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Validações

func ValidateName(name string) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if len(name) < 3 || len(name) > 100 {
		return errors.New("name must be between 3 and 100 characters")
	}
	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password cannot be empty")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	return nil
}

func ValidateRole(role UserRole) error {
	if role != AdminRole && role != CommonRole {
		return errors.New("invalid role: must be 'admin' or 'common'")
	}
	return nil
}

func (u *User) Validate() error {
	if err := ValidateName(u.Name); err != nil {
		return err
	}
	if err := ValidateEmail(u.Email); err != nil {
		return err
	}
	if err := ValidateRole(u.Role); err != nil {
		return err
	}
	return nil
}

// IsAdmin verifica se o usuário é admin
func (u *User) IsAdmin() bool {
	return u.Role == AdminRole
}

// Activate ativa o usuário
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Deactivate desativa o usuário
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}
