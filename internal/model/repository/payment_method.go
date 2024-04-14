package repository

import (
	"context"
	mPaymentMethod "nub/internal/model/payment_method"
)

type PaymentMethodRepository interface {
	FindAll(ctx context.Context, query map[string][]string) ([]mPaymentMethod.PaymentMethod, error)
	FindById(ctx context.Context, id string) (mPaymentMethod.PaymentMethod, error)
	Create(ctx context.Context, payload mPaymentMethod.PaymentMethod) (mPaymentMethod.PaymentMethod, error)
	Update(ctx context.Context, payload mPaymentMethod.PaymentMethod) (mPaymentMethod.PaymentMethod, error)
}
