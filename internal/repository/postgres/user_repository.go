package postgres

import (
	"database/sql"
	"time"

	"github.com/user-api/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new PostgreSQL user repository
func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (email, password, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	return r.db.QueryRow(query, user.Email, user.Password, user.Name, now).Scan(&user.ID)
}

func (r *userRepository) GetByID(id uint) (*domain.User, error) {
	query := `
		SELECT id, email, password, name, created_at, updated_at
		FROM users
		WHERE id = $1`

	user := &domain.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, password = $2, name = $3, updated_at = $4
		WHERE id = $5`

	user.UpdatedAt = time.Now()
	result, err := r.db.Exec(query, user.Email, user.Password, user.Name, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *userRepository) Delete(id uint) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *userRepository) List(page, limit int) ([]*domain.User, error) {
	query := `
		SELECT id, email, password, name, created_at, updated_at
		FROM users
		ORDER BY id
		LIMIT $1 OFFSET $2`

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// GetByEmail finds a user by email address
func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, email, password, name, created_at, updated_at
		FROM users
		WHERE email = $1`

	user := &domain.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
