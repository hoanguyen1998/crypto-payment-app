package postgresRepo

import (
	"context"
	"time"

	"github.com/hoanguyen1998/crypto-payment-system/internal/models"
)

func (m *postgresDBRepo) CreateAppKey(appKey models.ApplicationKey, masterPublicKeyId int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.Begin()

	if err != nil {
		return 0, err
	}

	stmt1 := `insert into application_keys (public_key, address_generated, app_id, payment_method_id)
			values ($1, $2, $3, $4)`

	row, err := m.DB.ExecContext(ctx, stmt1,
		appKey.PublicKey,
		appKey.AddressGenerated,
		appKey.AppId,
		appKey.MasterPublicKeyId,
	)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	appKeyId, err := row.LastInsertId()

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	stmt2 := `update master_public_keys set number_of_apps = number_of_apps + 1 where id = $1`

	_, err = m.DB.ExecContext(ctx, stmt2,
		masterPublicKeyId,
	)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return int(appKeyId), nil
}

func (m *postgresDBRepo) CreateAppAndKey(app models.Application, appKey models.ApplicationKey) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.Begin()

	if err != nil {
		return 0, 0, err
	}

	stmt := `insert into applications (name, user_id)
		values ($1, $2)`

	row, err := m.DB.ExecContext(ctx, stmt,
		app.Name,
		app.UserId,
	)

	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	appId, err := row.LastInsertId()

	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	stmt1 := `insert into application_keys (public_key, address_generated, app_id, payment_method_id)
			values ($1, $2, $3, $4)`

	row, err = m.DB.ExecContext(ctx, stmt1,
		appKey.PublicKey,
		appKey.AddressGenerated,
		appId,
		appKey.MasterPublicKeyId,
	)

	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	appKeyId, err := row.LastInsertId()

	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	stmt2 := `update master_public_keys set number_of_apps = number_of_apps + 1 where id = $1`

	_, err = m.DB.ExecContext(ctx, stmt2,
		appKey.MasterPublicKeyId,
	)

	if err != nil {
		tx.Rollback()
		return 0, 0, err
	}

	tx.Commit()

	return int(appId), int(appKeyId), nil
}

func (m *postgresDBRepo) GetAppKeyByAppIdAndPaymentMethod(appId int, paymentMethodId int) (*models.ApplicationKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, public_key, address_generated from application_keys where app_id = $1 and payment_method_id = $2
	`

	row := m.DB.QueryRowContext(ctx, query, appId, paymentMethodId)

	var appKey models.ApplicationKey

	err := row.Scan(
		&appKey.Id,
		&appKey.PublicKey,
		&appKey.AddressGenerated,
	)

	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, helpers.ErrAppKeyNotFound
		// }
		return nil, err
	}

	return &appKey, nil
}

func (m *postgresDBRepo) GetUserIdOfApp(appId int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select user_id from applications where id = $1
	`

	row := m.DB.QueryRowContext(ctx, query, appId)

	var userId int

	err := row.Scan(&userId)

	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, helpers.ErrAppKeyNotFound
		// }
		return 0, err
	}

	return userId, nil
}

// func (m *postgresDBRepo) GetAppKeys(skip int) ([]*models.ApplicationKey, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	query := `select id, public_key, address_generated, app_id, payment_method_id from application_keys where id > $1 limit 10`

// 	rows, err := m.DB.QueryContext(ctx, query, skip)

// 	var appKeys []*models.ApplicationKey

// 	if err != nil {
// 		return nil, err
// 	}

// 	for rows.Next() {
// 		var appKey models.ApplicationKey
// 		rows.Scan(
// 			&appKey.Id,
// 			&appKey.PublicKey,
// 			&appKey.AddressGenerated,
// 			&appKey.AppId,
// 			&appKey.PaymentMethodId,
// 		)

// 		appKeys = append(appKeys, &appKey)
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	return appKeys, nil
// }

// func (m *postgresDBRepo) GetAppKeysWithUserId(skip int, userId int) ([]*models.ApplicationKey, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	query := `select id, public_key, address_generated, app_id, payment_method_id from application_keys where user_id = $1` //id > $1 limit 10`

// 	rows, err := m.DB.QueryContext(ctx, query, userId)

// 	var appKeys []*models.ApplicationKey

// 	if err != nil {
// 		return nil, err
// 	}

// 	for rows.Next() {
// 		var appKey models.ApplicationKey
// 		rows.Scan(
// 			&appKey.Id,
// 			&appKey.PublicKey,
// 			&appKey.AddressGenerated,
// 			&appKey.AppId,
// 			&appKey.PaymentMethodId,
// 		)

// 		appKeys = append(appKeys, &appKey)
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	return appKeys, nil
// }

// func (m *postgresDBRepo) GetAppKeyByAppId(skip int, appId int) ([]*models.ApplicationKey, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	query := `select id, public_key, address_generated, payment_method_id from application_keys where app_id = $1`

// 	rows, err := m.DB.QueryContext(ctx, query, appId)

// 	var appKeys []*models.ApplicationKey

// 	if err != nil {
// 		return nil, err
// 	}

// 	for rows.Next() {
// 		var appKey models.ApplicationKey
// 		rows.Scan(
// 			&appKey.Id,
// 			&appKey.PublicKey,
// 			&appKey.AddressGenerated,
// 			&appKey.PaymentMethodId,
// 		)

// 		appKeys = append(appKeys, &appKey)
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	return appKeys, nil
// }

// func (m *postgresDBRepo) GetPaymentMethodsOfApp(appId int) ([]int, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	query := `select payment_method_id from application_keys where app_id = $1`

// 	rows, err := m.DB.QueryContext(ctx, query, appId)

// 	var paymentMethods []int

// 	if err != nil {
// 		return nil, err
// 	}

// 	for rows.Next() {
// 		var paymentMethod int
// 		rows.Scan(
// 			&paymentMethod,
// 		)

// 		paymentMethods = append(paymentMethods, paymentMethod)
// 	}

// 	return paymentMethods, nil
// }
