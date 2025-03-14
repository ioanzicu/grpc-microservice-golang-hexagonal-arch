package db

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/ioanzicu/microservices/order/internal/application/core/domain"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type OrderDatabaseTestSuite struct {
	suite.Suite
	DataSourceUrl string
}

const mysqlPassword string = "impossibletoguess"

func (o *OrderDatabaseTestSuite) SetUpSuite() {

	ctx := context.Background()
	port := "3306/tcp"
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("root:%s@tcp(localhost:%s)/orders?charset=utf8mb4&parseTime=True&loc=Local", mysqlPassword, port.Port())
	}

	req := testcontainers.ContainerRequest{
		Image:        "docker.io/mysql:8.8.30",
		ExposedPorts: []string{port},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": mysqlPassword,
			"MYSQL_DATABASE":      "orders",
		},
		WaitingFor: wait.ForSQL(nat.Port(port), "mysql", dbURL).WithStartupTimeout(30 * time.Second),
	}

	mysqlContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal("Failed to start Mysql.", err)
	}
	endpoint, _ := mysqlContainer.Endpoint(ctx, "")
	o.DataSourceUrl = fmt.Sprintf("root:%s@tcp(%s)/orders?charset=utf8mb4&parseTime=True&loc=Local", mysqlPassword, endpoint)
}

func (o *OrderDatabaseTestSuite) TestShouldSaveOrder() {
	adapter, err := NewAdapter(o.DataSourceUrl)
	o.Nil(err)
	saveErr := adapter.Save(context.Background(), &domain.Order{})
	o.Nil(saveErr)
}

func (o *OrderDatabaseTestSuite) TestShouldGetOrder() {
	adapter, _ := NewAdapter(o.DataSourceUrl)
	order := domain.NewOrder(2, []domain.OrderItem{
		{
			ProductCode: "OUS",
			Quantity:    3,
			UnitPrice:   3.33,
		},
	})
	adapter.Save(context.Background(), &order)
	ord, _ := adapter.Get(context.Background(), order.ID)
	o.Equal(int64(2), ord.CustomerID)
}

func TestOrderDatabaseTestSuite(t *testing.T) {
	suite.Run(t, new(OrderDatabaseTestSuite))
}
