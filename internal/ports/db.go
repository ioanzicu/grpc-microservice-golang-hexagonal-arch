package ports

import (
	"github.com/ioanzicu/microservices/order/internal/application/core/domain"
)

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(*domain.Order) error
}
