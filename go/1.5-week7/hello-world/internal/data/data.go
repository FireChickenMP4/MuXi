package data

import (
	"sync"
	"sync/atomic"

	"hello-world/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo, NewUserRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	mu           sync.RWMutex
	users        map[int64]UserRecord
	accountIndex map[string]int64
	nextUserID   atomic.Int64
}

// NewData .
func NewData(c *conf.Data) (*Data, func(), error) {
	d := &Data{
		users:        make(map[int64]UserRecord),
		accountIndex: make(map[string]int64),
	}
	cleanup := func() {
		log.Info("closing the data resources")
	}
	return d, cleanup, nil
}

// UserRecord is the storage model used by the user repository.
type UserRecord struct {
	ID        int64
	Username  string
	Email     string
	Mobile    string
	Avatar    string
	Password  string
	CreatedAt int64
	UpdatedAt int64
}
