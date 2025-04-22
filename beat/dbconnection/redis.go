package dbconnection

import(
	"context"
	"os"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

func CreateRedisClient(dbNo int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: os.Getenv("DBREDIS_ADDR"),
		Password: os.Getenv("DBREDIS_PASS"),
		DB: dbNo,
	})
	return rdb
}