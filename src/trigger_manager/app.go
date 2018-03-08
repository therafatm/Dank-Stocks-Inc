package main

import (
	"log"
	"time"

	"common/logging"
	"transaction_service/queries/transdb"
	
	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

type Env struct {
	logger     logging.Logger
	tdb        transdb.TransactionDataStore
	quoteCache *redis.Client
}

func main() {
	logger := logging.NewLoggerConnection()
	tdb := transdb.NewTransactionDBConnection()
	quoteCache := transdb.NewQuoteCacheConnection()

	defer tdb.DB.Close()
	defer quoteCache.Close()

	env := &Env{quoteCache: quoteCache, logger: logger, tdb: tdb}

	log.Println("Trigger manager initiated:")
	const pollInterval = 2000
	for {
		time.Sleep(time.Millisecond * pollInterval)
		env.tdb.QueryAndExecuteCurrentTriggers(env.quoteCache, "1")
	}

}
