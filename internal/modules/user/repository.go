package user

import (
	"context"
	"database/sql"
	"fmt"
	"log" // MANTIDO: Apenas para exemplo, idealmente use um logger estruturado

	"github.com/google/uuid"
)

// Interface (sem mudanças, estava ótima)
type UserRepository interface {
	CreateUser(ctx context.Context, user *User) error
	GetAllUser(ctx context.Context) ([]User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type dbRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &dbRepository{
		db: db,
	}
}

func (r *dbRepository) CreateUser(ctx context.Context, user *User) error {
	// Lógica de hash de senha e validação de "Password nil".
	// O 'user' que chega aqui já deve vir pronto do 'UserService'.

	// Passando user.ID (o valor) em vez de uuid.UUID (o tipo).
	query := `
		INSERT INTO users (id, name, email, password, role, is_active, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Name,
		user.Email,
		user.Password, // Espera-se que esta senha JÁ ESTEJA HASHEADA
		user.Role,
		user.isActive,
		user.createdAt,
		user.updatedAt,
	)

	if err != nil {
		// Apenas logamos e retornamos o erro
		log.Printf("Error creating user: %v", err)
		return err
	}
	return nil
}

func (r *dbRepository) GetAllUser(ctx context.Context) ([]User, error) {
	// Seu código aqui estava perfeito, sem mudanças.
	rows, err := r.db.QueryContext(ctx, "SELECT id, name, email, password, role, is_active, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.isActive, &user.createdAt, &user.updatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *dbRepository) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	// Seu código aqui também estava perfeito.
	var user User
	query := `
		SELECT id, name, email, password, role, is_active, created_at, updated_at 
		FROM users WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email, &user.Password,
		&user.Role, &user.isActive, &user.createdAt, &user.updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("user with id %s not found", id)
		}
		return User{}, fmt.Errorf("error getting user: %w", err)
	}
	return user, nil
}

func (r *dbRepository) UpdateUser(ctx context.Context, user *User) error {
	// Lógica de hash e validação de senha.
	// log.Fatal(err)
	// Sintaxe SQL (usando vírgulas ,)

	query := `
		UPDATE users 
		SET name = $1, email = $2, password = $3, role = $4, is_active = $5, updated_at = $6 
		WHERE id = $7
	`
	// Nota: Esta query assume que você está atualizando a senha.
	// O UserService seria responsável por decidir se a senha antiga ou a nova hasheada é passada.
	res, err := r.db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.Password, // Espera-se que JÁ ESTEJA HASHEADA (se foi alterada)
		user.Role,
		user.isActive,
		user.updatedAt,
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	// Checagem correta de "não encontrado" para ExecContext
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %s not found for update", user.ID)
	}

	return nil
}

func (r *dbRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// Checagem correta de "não encontrado" para ExecContext
	res, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %s not found for deletion", id)
	}

	return nil
}
