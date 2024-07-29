//go:build integration

package tests

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"homework-1/internal/model"
	"homework-1/internal/service"
	"homework-1/internal/storage"
	"os"
	"testing"
	"time"
)

type OrdersTestSuite struct {
	suite.Suite
	serv *service.Service
	st   *storage.Storage
	ctx  context.Context
}

func (suite *OrdersTestSuite) SetupSuite() {
	suite.serv, suite.st, suite.ctx = Init()
}

func (suite *OrdersTestSuite) SetupTest() {
	_ = suite.serv.Storage.WriteOrderWithUniqueId(suite.ctx, &model.Order{OrderID: 30, ClientID: 123, TakenAt: time.Now(), Pack: "box"})
	_ = suite.serv.Storage.WriteOrderWithUniqueId(suite.ctx, &model.Order{OrderID: 31, ClientID: 456, Pack: "box"})
}

func (suite *OrdersTestSuite) TearDownSuite() {
	suite.st.CloseStorage()
}

func (suite *OrdersTestSuite) TearDownTest() {
	DeleteOrder(suite.st, 30)
	DeleteOrder(suite.st, 31)
}

func (suite *OrdersTestSuite) TestListOfTakeBacks() {
	pageNumber := 1
	var buf bytes.Buffer
	originalStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := suite.serv.ListOfTakeBacks(suite.ctx, pageNumber)
	require.NoError(suite.T(), err)

	w.Close()
	os.Stdout = originalStdout

	buf.ReadFrom(r)
	actualOutput := buf.String()

	expectedOutput := fmt.Sprintf("Order id: %v\t taken back at: %v\n", 30, time.Time{})
	assert.Equal(suite.T(), expectedOutput, actualOutput)
}

func TestOrders(t *testing.T) {
	suite.Run(t, new(OrdersTestSuite))
}
