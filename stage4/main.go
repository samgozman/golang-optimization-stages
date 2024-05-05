package stage4

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/shopspring/decimal"

	"github.com/samgozman/golang-optimization-stages/object"
)

var (
	once    sync.Once
	content []byte
)

// ServeApp just a simple server that runs one route
func ServeApp(ctx context.Context) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.GET("/json", GetJSONHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Errorf("listen: %s\n", err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		if err := srv.Shutdown(ctx); err != nil {
			panic(fmt.Errorf("Server Shutdown: %s\n", err))
		}
	case <-quit:
		if err := srv.Shutdown(ctx); err != nil {
			panic(fmt.Errorf("Server Shutdown: %s\n", err))
		}
	}
}

// GetJSONHandler is a simple handler that returns a JSON response with a message
func GetJSONHandler(c *gin.Context) {
	// 1. Read the content of the file dummy_data.json
	data := ReadFile("../object/dummy_data.json")

	// 2. Return a JSON response with a message
	users := ParseJSON(data)

	// 3. Calculate the total balances
	currents, pendings := BalancesTotals(users)

	// 4. Calculate the total transactions
	transactionsSum, transactionsCount := TransactionsTotals(users)

	c.JSON(200, gin.H{
		"current":            currents.String(),
		"pending":            pendings.String(),
		"transactions_sum":   transactionsSum.String(),
		"transactions_count": transactionsCount,
	})
}

// ReadFile reads a file and returns its content as a string
func ReadFile(filePath string) []byte {
	once.Do(func() {
		var err error
		content, err = os.ReadFile(filePath)
		if err != nil {
			panic(fmt.Errorf("failed to read file: %w", err))
		}
	})
	return content
}

func ParseJSON(data []byte) []object.User {
	var users []object.User
	if err := json.Unmarshal(data, &users); err != nil {
		panic(fmt.Errorf("failed to unmarshal JSON: %w", err))
	}

	return users
}

// BalancesTotals calculates the total balances of the users.
func BalancesTotals(users []object.User) (currents decimal.Decimal, pendings decimal.Decimal) {
	for _, user := range users {
		current, _ := decimal.NewFromString(user.Balance.Current)
		currents = currents.Add(current)

		pending, _ := decimal.NewFromString(user.Balance.Pending)
		pendings = pendings.Add(pending)
	}

	return
}

// TransactionsTotals calculates the total transactions of the users.
func TransactionsTotals(users []object.User) (sum decimal.Decimal, count int) {
	var transactionsSum decimal.Decimal
	var transactionsCount int

	for _, user := range users {
		for _, transaction := range user.Transactions {
			amount, _ := decimal.NewFromString(transaction.Amount)
			transactionsSum = transactionsSum.Add(amount)
			transactionsCount++
		}
	}

	return transactionsSum, transactionsCount
}
