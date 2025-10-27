package user

// internal/modules/user/models.go define os modelos de dados relacionados aos usuarios
// e encapsula a logica de negocio associada a esses modelos.

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	AdminRole  UserRole = "admin"
	CommonRole UserRole = "common"
)

type User struct {
	ID        uuid.UUID
	Name      string
	Email     string
	Password  string
	Role      UserRole
	isActive  bool
	createdAt time.Time
	updatedAt time.Time
}
