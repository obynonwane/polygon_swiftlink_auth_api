package data

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/obynonwane/polygon_swiftlink_auth_api/util"
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

	//convert password to hash
	hashPassword, err := util.HashPassword(payload.Password)
	if err != nil {
		log.Println("Error inserting user:", err)
		return nil, errors.New("error hashing user password")
	}

	// Create the insert statement to add a new user to the 'users' table
	query := `
		INSERT INTO users (first_name, last_name, email, phone, password)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, first_name, last_name, email, phone, password, verified, created_at, updated_at
	`

	// Prepare a variable to store the new user
	var user User

	// Execute the insert query and scan the returned row into the user struct
	err = u.Conn.QueryRowContext(ctx, query,
		payload.FirstName,
		payload.LastName,
		payload.Email,
		payload.Phone,
		hashPassword,
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

func (u *PostgresRepository) GetUserWithEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `SELECT id, first_name, last_name, email, phone, password, updated_at, created_at FROM users WHERE email = $1 LIMIT 1;`
	row := u.Conn.QueryRowContext(ctx, stmt, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Password,
		&i.UpdatedAt,
		&i.CreatedAt,
	)

	return &i, err
}
