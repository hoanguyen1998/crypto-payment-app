package postgresRepo

import (
	"context"
	"time"

	"github.com/hoanguyen1998/crypto-payment-system/internal/models"
)

func (m *postgresDBRepo) CreateUser(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into users (email, username, password_hash)
			values ($1, $2, $3)`

	row, err := m.DB.ExecContext(ctx, stmt,
		user.Email,
		user.Username,
		user.PasswordHash,
	)

	if err != nil {
		return 0, err
	}

	userId, err := row.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(userId), nil
}

func (m *postgresDBRepo) GetUserById(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select username, password_hash, email from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var user models.User

	err := row.Scan(
		&user.Username,
		&user.PasswordHash,
		&user.Email,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *postgresDBRepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select username, password_hash from users where email = $1`

	row := m.DB.QueryRowContext(ctx, query, email)

	var user models.User

	err := row.Scan(
		&user.Username,
		&user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *postgresDBRepo) GetUsers(skip int) ([]*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, username, password_hash, email from users where id > $1 limit 10`

	rows, err := m.DB.QueryContext(ctx, query, skip)

	var users []*models.User

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user models.User
		rows.Scan(
			&user.Id,
			&user.Username,
			&user.PasswordHash,
			&user.Email,
		)

		users = append(users, &user)
	}

	if err != nil {
		return nil, err
	}

	return users, nil
}
