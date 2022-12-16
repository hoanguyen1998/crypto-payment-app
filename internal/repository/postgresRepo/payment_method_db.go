package postgresRepo

import (
	"context"
	"time"

	"github.com/hoanguyen1998/crypto-payment-system/internal/models"
)

func (m *postgresDBRepo) GetPaymentMethods() ([]*models.PaymentMethod, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name from payment_methods`

	rows, err := m.DB.QueryContext(ctx, query)

	var paymentMethods []*models.PaymentMethod

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var paymentMethod models.PaymentMethod
		rows.Scan(
			&paymentMethod.Id,
			&paymentMethod.Name,
		)

		paymentMethods = append(paymentMethods, &paymentMethod)
	}

	if err != nil {
		return nil, err
	}

	return paymentMethods, nil
}
