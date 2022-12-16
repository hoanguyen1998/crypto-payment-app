package postgresRepo

import (
	"context"
	"time"

	"github.com/hoanguyen1998/crypto-payment-system/internal/models"
)

func (m *postgresDBRepo) CreateApp(app models.Application) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into applications (name, user_id)
			values ($1, $2)`

	row, err := m.DB.ExecContext(ctx, stmt,
		app.Name,
		app.UserId,
	)

	if err != nil {
		return 0, err
	}

	appId, err := row.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(appId), nil
}

func (m *postgresDBRepo) GetAppById(appId int) (*models.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select name, app_index, user_id from applications where id = $1`

	row := m.DB.QueryRowContext(ctx, query, appId)

	var app models.Application

	err := row.Scan(
		&app.Name,
		&app.AppIndex,
		&app.UserId,
	)

	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, helpers.ErrAppNotFound
		// }
		return nil, err
	}

	return &app, nil
}

func (m *postgresDBRepo) GetAppIndexByAppId(appId int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select number_of_apps from applications where id = $1`

	row := m.DB.QueryRowContext(ctx, query, appId)

	var appIndex int

	err := row.Scan(
		&appIndex,
	)

	if err != nil {
		return 0, err
	}

	return appIndex, nil
}

func (m *postgresDBRepo) GetNumberOfAppsByUser(userId int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select count(id) from applications where user_id = $1`

	row := m.DB.QueryRowContext(ctx, query, userId)

	var count int

	err := row.Scan(
		&count,
	)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *postgresDBRepo) GetApps(skip int) ([]*models.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name, app_index, user_id from applications where id > $1 limit 10`

	rows, err := m.DB.QueryContext(ctx, query, skip)

	var apps []*models.Application

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var app models.Application
		rows.Scan(
			&app.Id,
			&app.Name,
			&app.AppIndex,
			&app.UserId,
		)

		apps = append(apps, &app)
	}

	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (m *postgresDBRepo) GetAppsByUserId(userId int) ([]*models.Application, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name, app_index from applications where user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, userId)

	var apps []*models.Application

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var app models.Application
		rows.Scan(
			&app.Id,
			&app.Name,
			&app.AppIndex,
		)

		apps = append(apps, &app)
	}

	if err != nil {
		return nil, err
	}

	return apps, nil
}
