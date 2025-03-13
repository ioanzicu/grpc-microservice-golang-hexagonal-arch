package api

import (
	"context"
	"errors"
	"testing"

	"github.com/ioanzicu/microservices/order/internal/application/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type mockedPayment struct {
	mock.Mock
}

func (p *mockedPayment) Charge(ctx context.Context, order *domain.Order) error {
	args := p.Called(ctx, order)
	return args.Error(0)
}

type mockedDB struct {
	mock.Mock
}

func (d *mockedDB) Save(ctx context.Context, order *domain.Order) error {
	args := d.Called(ctx, order)
	return args.Error(0)
}

func (d *mockedDB) Get(ctx context.Context, id int64) (domain.Order, error) {
	args := d.Called(ctx, id)
	return args.Get(0).(domain.Order), args.Error(1)
}

func TestPlaceOrder(t *testing.T) {
	payment := new(mockedPayment)
	db := new(mockedDB)
	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)
	db.On("Save", mock.Anything, mock.Anything).Return(nil)

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(context.Background(), domain.Order{
		CustomerID: 123,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "Macbook",
				UnitPrice:   999.99,
				Quantity:    3,
			},
		},
		CreatedAt: 0,
	})
	assert.Nil(t, err)
}

func TestShouldErrorWhenDBPersistenceFail(t *testing.T) {
	payment := new(mockedPayment)
	db := new(mockedDB)
	payment.On("Charge", mock.Anything, mock.Anything).Return(nil)
	db.On("Save", mock.Anything, mock.Anything).Return(errors.New("connection error"))

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(context.Background(), domain.Order{
		CustomerID: 333,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "iphone",
				UnitPrice:   666.66,
				Quantity:    3,
			},
		},
		CreatedAt: 0,
	})
	assert.EqualError(t, err, "connection error")
}

func TestShouldReturnErrorWhenPaymetFail(t *testing.T) {
	payment := new(mockedPayment)
	db := new(mockedDB)
	payment.On("Charge", mock.Anything, mock.Anything).Return(errors.New("insufficient balance"))
	db.On("Save", mock.Anything, mock.Anything).Return(nil)

	application := NewApplication(db, payment)
	_, err := application.PlaceOrder(context.Background(), domain.Order{
		CustomerID: 333,
		OrderItems: []domain.OrderItem{
			{
				ProductCode: "bicycle",
				UnitPrice:   333.5,
				Quantity:    10,
			},
		},
		CreatedAt: 0,
	})
	st, _ := status.FromError(err)
	assert.Equal(t, st.Message(), "order creation failed")
	assert.Equal(t, st.Details()[0].(*errdetails.BadRequest).FieldViolations[0].Description, "insufficient balance")
	assert.Equal(t, st.Code(), codes.InvalidArgument)
}
