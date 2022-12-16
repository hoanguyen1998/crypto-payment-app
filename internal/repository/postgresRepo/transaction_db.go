package postgresRepo

import (
	"context"
	"time"

	"github.com/hoanguyen1998/crypto-payment-system/internal/models"
)

func (m *postgresDBRepo) CreateTransaction(tx models.Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into transactions (tx_hash, sender, recipient, amount, payment_method_id, order_id, tx_type)
			values ($1, $2, $3, $4, $5, $6, $7)`

	row, err := m.DB.ExecContext(ctx, stmt,
		tx.TxHash,
		tx.Sender,
		tx.Recipient,
		tx.Amount,
		tx.PaymentMethodId,
		tx.OrderId,
		tx.TxType,
	)

	if err != nil {
		return 0, err
	}

	txId, err := row.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(txId), nil
}

func (m *postgresDBRepo) UpdateTransactionData(txId int, blockNumber int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `update transactions set block_number = $1 where id = $2`

	_, err := m.DB.ExecContext(ctx, query, blockNumber, txId)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) GetTransactionByHash(txHash string) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, sender, recipient, tx_type, order_id from transactions where tx_hash = $1
	`

	row := m.DB.QueryRowContext(ctx, query, txHash)

	var transaction models.Transaction

	err := row.Scan(
		&transaction.Id,
		&transaction.Sender,
		&transaction.Recipient,
		&transaction.TxType,
		&transaction.OrderId,
	)

	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (m *postgresDBRepo) GetLatestBlockNumber(paymentMethodId int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select block_number from blocks where payment_method_id=$1 order by id ASC limit 1`

	row := m.DB.QueryRowContext(ctx, query, paymentMethodId)

	var blockNumber int

	err := row.Scan(&blockNumber)

	if err != nil {
		return 0, err
	}

	return blockNumber, nil
}
