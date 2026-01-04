package repositories

import (
	"database/sql"
	"ssi-signin/backend/models"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (did, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	
	now := time.Now()
	err := r.db.QueryRow(query, user.DID, user.Phone, now, now).Scan(&user.ID)
	if err != nil {
		return err
	}
	
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *UserRepository) FindByDID(did string) (*models.User, error) {
	query := `
		SELECT id, did, phone, created_at, updated_at
		FROM users
		WHERE did = $1
	`
	
	user := &models.User{}
	err := r.db.QueryRow(query, did).Scan(
		&user.ID,
		&user.DID,
		&user.Phone,
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

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	query := `
		SELECT id, did, phone, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.DID,
		&user.Phone,
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

