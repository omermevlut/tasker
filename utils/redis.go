package utils

import (
	"time"

	"github.com/go-redis/redis"
	"github.com/omermevlut/tasker/config"
)

// RedisUtilInterface ...
type RedisUtilInterface interface {
	SetDelayed(value interface{}, expiry float64)
	MoveExpiredItems(expiry int64, script string)
	PopFromActiveQueue(script string) string
}

// RedisUtil ...
type RedisUtil struct {
	Client      *redis.Client
	queueSuffix string
}

// NewRedisUtil ...
func NewRedisUtil(address, queueSuffix string) RedisUtilInterface {
	return &RedisUtil{
		Client:      redis.NewClient(&redis.Options{Addr: address, Password: "", DB: 0}),
		queueSuffix: queueSuffix,
	}
}

// SetDelayed task
func (r *RedisUtil) SetDelayed(value interface{}, expiry float64) {
	r.Client.ZAdd(config.Queues.Delayed+r.queueSuffix, redis.Z{Member: value, Score: expiry})
}

// MoveExpiredItems ...
func (r *RedisUtil) MoveExpiredItems(expiry int64, script string) {
	r.Client.Eval(
		script,
		[]string{config.Queues.Delayed + r.queueSuffix, config.Queues.Default + r.queueSuffix},
		expiry,
	)
}

// PopFromActiveQueue ...
func (r *RedisUtil) PopFromActiveQueue(script string) string {
	res := r.Client.Eval(
		script,
		[]string{
			config.Queues.Default + r.queueSuffix,
			config.Queues.Expired + r.queueSuffix,
			time.Now().Format(time.RFC3339),
		},
		time.Now().Unix(),
	)

	if res.Val() == nil {
		return ""
	}

	return res.Val().(string)
}
