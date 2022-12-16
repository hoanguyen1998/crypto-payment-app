package postgresRepo

import (
	"context"
	"time"

	"github.com/hoanguyen1998/crypto-payment-system/internal/models"
)

func (m *postgresDBRepo) CreateMasterPublicKey(masterKey models.MasterPublicKey) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into master_public_keys (public_key, user_id, payment_method_id)
			values ($1, $2, $3)`

	_, err := m.DB.ExecContext(ctx, stmt,
		masterKey.PublicKey,
		masterKey.UserId,
		masterKey.PaymentMethodId,
	)

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (m *postgresDBRepo) GetMasterKeyByUserIdAndPaymentMethod(userId int, paymentMethodId int) (*models.MasterPublicKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select public_key from master_public_keys where user_id = $1 and payment_method_id = $2
	`

	row := m.DB.QueryRowContext(ctx, query, userId, paymentMethodId)

	var masterKey models.MasterPublicKey

	err := row.Scan(
		&masterKey.PublicKey,
	)

	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return nil, helpers.ErrMasterKeyNotFound
		// }
		return nil, err
	}

	return &masterKey, nil
}

func (m *postgresDBRepo) GetMasterKeys(skip int) ([]*models.MasterPublicKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, public_key, user_id, payment_method_id from master_public_keys where id > $1 limit 10`

	rows, err := m.DB.QueryContext(ctx, query, skip)

	var masterKeys []*models.MasterPublicKey

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var masterKey models.MasterPublicKey
		rows.Scan(
			&masterKey.Id,
			&masterKey.PublicKey,
			&masterKey.UserId,
			&masterKey.PaymentMethodId,
		)

		masterKeys = append(masterKeys, &masterKey)
	}

	if err != nil {
		return nil, err
	}

	return masterKeys, nil
}

func (m *postgresDBRepo) GetMasterKeysWithUserId(userId int) ([]*models.MasterPublicKey, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, public_key, user_id, payment_method_id from master_public_keys where user_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, userId)

	var masterKeys []*models.MasterPublicKey

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var masterKey models.MasterPublicKey
		rows.Scan(
			&masterKey.Id,
			&masterKey.PublicKey,
			&masterKey.UserId,
			&masterKey.PaymentMethodId,
		)

		masterKeys = append(masterKeys, &masterKey)
	}

	if err != nil {
		return nil, err
	}

	return masterKeys, nil
}
