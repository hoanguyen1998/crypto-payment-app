package postgresRepo

import (
	"context"
	"time"

	"github.com/hoanguyen1998/crypto-payment-system/internal/models"
)

func (m *postgresDBRepo) CreateOrder(addressNum, appKeyId int, order models.Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.Begin()

	if err != nil {
		return 0, err
	}

	stmt1 := `insert into orders (application_order_id, amount, received_address, status, path, application_key_id)
			values ($1, $2, $3, $4, $5, $6)`

	row, err := m.DB.ExecContext(ctx, stmt1,
		order.ApplicationOrderId,
		order.Amount,
		order.ReceivedAddress,
		order.Status,
		order.Path,
		order.ApplicationKeyId,
	)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	orderId, err := row.LastInsertId()

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	stmt2 := `update application_keys set address_generated = $1 where id = $2`

	_, err = m.DB.ExecContext(ctx, stmt2,
		addressNum,
		appKeyId,
	)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return int(orderId), nil
}

func (m *postgresDBRepo) UpdateOrderStatusAndCreateTransaction(orderId int, status string, txData models.Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := m.DB.Begin()

	if err != nil {
		return 0, err
	}

	stmt1 := `update orders set status = $1 where id = $2`

	_, err = m.DB.ExecContext(ctx, stmt1, status, orderId)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	stmt := `insert into transactions (tx_hash, sender, recipient, amount, payment_method_id, order_id, tx_type)
			values ($1, $2, $3, $4, $5, $6, $7)`

	row, err := m.DB.ExecContext(ctx, stmt,
		txData.TxHash,
		txData.Sender,
		txData.Recipient,
		txData.Amount,
		txData.PaymentMethodId,
		txData.OrderId,
		txData.TxType,
	)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	txId, err := row.LastInsertId()

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return int(txId), nil
}

func (m *postgresDBRepo) GetOrdersByAppKey(appKeyId int) ([]*models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, application_order_id, amount, received_address, status from orders where application_key_id = $1`

	rows, err := m.DB.QueryContext(ctx, query, appKeyId)

	var orders []*models.Order

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order models.Order
		rows.Scan(
			&order.Id,
			&order.ApplicationOrderId,
			&order.Amount,
			&order.ReceivedAddress,
			&order.Status,
		)

		orders = append(orders, &order)
	}

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// func (m *postgresDBRepo) GetOrderById(orderId int) (*models.Order, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	query := `select id, amount, application_key_id, application_order_id, path, received_address from orders where id = $1
// 	`

// 	row := m.DB.QueryRowContext(ctx, query, orderId)

// 	var order models.Order

// 	err := row.Scan(
// 		&order.Id,
// 		&order.Amount,
// 		&order.ApplicationKeyId,
// 		&order.ApplicationOrderId,
// 		&order.Path,
// 		&order.ReceivedAddress,
// 	)

// 	if err != nil {
// 		if err != sql.ErrNoRows {
// 			return nil, helpers.ErrOrderNotFound
// 		}
// 		return nil, err
// 	}

// 	return &order, nil
// }

// func (m *postgresDBRepo) GetOrderByAddress(address string) (*models.Order, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	query := `select id, amount, status from orders where received_address = $1
// 	`

// 	row := m.DB.QueryRowContext(ctx, query, address)

// 	var order models.Order

// 	err := row.Scan(
// 		&order.Id,
// 		&order.Amount,
// 		&order.Status,
// 	)

// 	if err != nil {
// 		if err != sql.ErrNoRows {
// 			return nil, helpers.ErrOrderNotFound
// 		}
// 		return nil, err
// 	}

// 	return &order, nil
// }

// func (m *postgresDBRepo) GetOrders(skip int) ([]*models.Order, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	query := `select id, application_order_id, amount, received_address, status, path, application_key_id from orders where id > $1 limit 10`

// 	rows, err := m.DB.QueryContext(ctx, query, skip)

// 	var orders []*models.Order

// 	if err != nil {
// 		return nil, err
// 	}

// 	for rows.Next() {
// 		var order models.Order
// 		rows.Scan(
// 			&order.Id,
// 			&order.ApplicationOrderId,
// 			&order.Amount,
// 			&order.ReceivedAddress,
// 			&order.Status,
// 			&order.Path,
// 			&order.ApplicationKeyId,
// 		)

// 		orders = append(orders, &order)
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	return orders, nil
// }

func (m *postgresDBRepo) GetOrdersToWithdraw(appKeyId int) ([]models.Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, amount, received_address, path from orders where application_key_id = $1 and status='confirm'`

	rows, err := m.DB.QueryContext(ctx, query, appKeyId)

	var orders []models.Order

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order models.Order
		rows.Scan(
			&order.Id,
			&order.Amount,
			&order.ReceivedAddress,
			&order.Path,
		)

		orders = append(orders, order)
	}

	if err != nil {
		return nil, err
	}

	return orders, nil
}

// func (m *postgresDBRepo) GetOrdersBalanceByAppKey(appKeyId int) (int, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	query := `select sum(amount) from orders where application_key_id = $1`

// 	row := m.DB.QueryRowContext(ctx, query, appKeyId)

// 	var balance int

// 	err := row.Scan(
// 		&balance,
// 	)

// 	if err != nil {
// 		return 0, err
// 	}

// 	return balance, nil
// }
