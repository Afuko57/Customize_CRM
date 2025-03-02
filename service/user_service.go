package service

import (
	"context"
	"database/sql"
	"errors"

	"customize_crm/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *pgxpool.Pool
}

func NewUserService(db *pgxpool.Pool) *UserService {
	return &UserService{db: db}
}

// GetByID
func (s *UserService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, 
			   role_id, department, created_at, updated_at, is_active
		FROM users
		WHERE id = $1
	`

	var user model.User
	var department sql.NullString

	err := s.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.RoleID, &department,
		&user.CreatedAt, &user.UpdatedAt, &user.IsActive,
	)
	if err != nil {
		return nil, err
	}

	if department.Valid {
		dept := department.String
		user.Department = &dept
	}

	return &user, nil
}

// GetByUsername
func (s *UserService) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, 
			   role_id, department, created_at, updated_at, is_active
		FROM users
		WHERE username = $1
	`

	var user model.User
	var department sql.NullString

	err := s.db.QueryRow(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.RoleID, &department,
		&user.CreatedAt, &user.UpdatedAt, &user.IsActive,
	)
	if err != nil {
		return nil, err
	}

	if department.Valid {
		dept := department.String
		user.Department = &dept
	}

	return &user, nil
}

// GetAll
func (s *UserService) GetAll(ctx context.Context) ([]*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, first_name, last_name, 
			   role_id, department, created_at, updated_at, is_active
		FROM users
	`

	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User

	for rows.Next() {
		var user model.User
		var department sql.NullString

		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash,
			&user.FirstName, &user.LastName, &user.RoleID, &department,
			&user.CreatedAt, &user.UpdatedAt, &user.IsActive,
		)
		if err != nil {
			return nil, err
		}

		if department.Valid {
			dept := department.String
			user.Department = &dept
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Create
func (s *UserService) Create(ctx context.Context, user *model.User, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)

	query := `
		INSERT INTO users (username, email, password_hash, first_name, last_name, role_id, department, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	err = s.db.QueryRow(ctx, query,
		user.Username, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.RoleID, user.Department, user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

// Update
func (s *UserService) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, department = $3, role_id = $4, is_active = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING updated_at
	`

	err := s.db.QueryRow(ctx, query,
		user.FirstName, user.LastName, user.Department, user.RoleID, user.IsActive, user.ID,
	).Scan(&user.UpdatedAt)

	return err
}

// UpdatePassword
func (s *UserService) UpdatePassword(ctx context.Context, id uuid.UUID, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET password_hash = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err = s.db.Exec(ctx, query, string(hashedPassword), id)
	return err
}

// Delete
func (s *UserService) Delete(ctx context.Context, ids []uuid.UUID) error {
	query := `DELETE FROM users WHERE id = ANY($1)`
	_, err := s.db.Exec(ctx, query, ids)
	return err
}

// Authenticate
func (s *UserService) Authenticate(ctx context.Context, username, password string) (*model.User, error) {
	user, err := s.GetByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("user account is disabled")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *UserService) GetRoleByID(ctx context.Context, id uuid.UUID) (*model.Role, error) {
	query := `
		SELECT id, name, description, permissions, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	var role model.Role
	err := s.db.QueryRow(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.Permissions,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &role, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return s.GetByID(ctx, id)
}
