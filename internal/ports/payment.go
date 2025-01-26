package ports

import "github.com/ioanzicu/microservices/order/internal/application/core/domain"

type PaymentPort interface {
	Charge(*domain.Order) error
}
