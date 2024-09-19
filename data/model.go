package data

import (
	"context"
	"database/sql"
	"log"
	"time"
)

// db timeout period
const dbTimeout = time.Second * 3

// data of sqlDB type here connections to DB will live
var db *sql.DB

type PostgresRepository struct {
	Conn *sql.DB
}

// new instance of the PostgresRepository struct
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		Conn: db,
	}
}

func (u *PostgresRepository) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `SELECT id, email, first_name, last_name, password, verified, updated_at, created_at FROM users`

	rows, err := u.Conn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Password,
			&user.Verified,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			log.Println("Error scanning", err)
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

func (u *PostgresRepository) Signup(payload SignupPayload) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// Create the insert statement to add a new user to the 'users' table
	query := `
		INSERT INTO users (first_name, last_name, email, phone, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, first_name, last_name, email, phone, password, verified, created_at, updated_at
	`

	// Prepare a variable to store the new user
	var user User

	// Execute the insert query and scan the returned row into the user struct
	err := u.Conn.QueryRowContext(ctx, query,
		payload.FirstName,
		payload.LastName,
		payload.Email,
		payload.Phone,
		payload.Password,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Phone,
		&user.Password,
		&user.Verified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Println("Error inserting user:", err)
		return nil, err
	}

	return &user, nil
}