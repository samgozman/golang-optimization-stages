package stage5

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
	once     sync.Once
	response gin.H
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
	once.Do(func() {
		// 1. Read the content of the file dummy_data.json
		data := ReadFile("../object/dummy_data.json")

		// 2. Return a JSON response with a message
		users := ParseJSON(data)

		// 3. Calculate the total balances and transactions
		currents, pendings, transactionsSum, transactionsCount := GetTotals(users)

		response = gin.H{
			"current":            currents.String(),
			"pending":            pendings.String(),
			"transactions_sum":   transactionsSum.String(),
			"transactions_count": transactionsCount,
		}
	})

	c.JSON(200, response)
}

// ReadFile reads a file and returns its content as a string
func ReadFile(filePath string) []byte {
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to read file: %w", err))
	}

	return content
}

func ParseJSON(data []byte) []*object.User {
	var users []*object.User
	if err := json.Unmarshal(data, &users); err != nil {
		panic(fmt.Errorf("failed to unmarshal JSON: %w", err))
	}

	return users
}

// GetTotals calculates the total balances and transactions of the users.
func GetTotals(users []*object.User) (
	// Note: this can be returned as a struct, but it's not necessary for this example
	currents decimal.Decimal,
	pendings decimal.Decimal,
	transactionsSum decimal.Decimal,
	transactionsCount int,
) {
	for _, user := range users {
		current, _ := decimal.NewFromString(user.Balance.Current)
		currents = currents.Add(current)

		pending, _ := decimal.NewFromString(user.Balance.Pending)
		pendings = pendings.Add(pending)
		transactionsCount += len(user.Transactions)

		for _, transaction := range user.Transactions {
			amount, _ := decimal.NewFromString(transaction.Amount)
			transactionsSum = transactionsSum.Add(amount)
		}
	}

	return currents, pendings, transactionsSum, transactionsCount
}
