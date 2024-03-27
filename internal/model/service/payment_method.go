package service

import (
	"context"
	mPaymentMethod "findigitalservice/internal/model/payment_method"
)

type PaymentMethodService interface {
	FindAll(ctx context.Context, query map[string][]string) ([]mPaymentMethod.PaymentMethod, error)
	FindById(ctx context.Context, id string) (mPaymentMethod.PaymentMethod, error)
	Create(ctx context.Context, payload mPaymentMethod.PaymentMethod) (mPaymentMethod.PaymentMethod, error)
	Update(ctx context.Context, id string, payload mPaymentMethod.PaymentMethod) (mPaymentMethod.PaymentMethod, error)
}
