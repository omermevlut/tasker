package utils

import (
	"io/ioutil"
	"time"

	"github.com/go-redis/redis"
	"github.com/omermevlut/tasker/config"
)

// RedisUtilInterface ...
type RedisUtilInterface interface {
	SetDelayed(value interface{}, expiry float64)
	MoveExpiredItems(expiry int64)
	PopFromActiveQueue() string
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
func (r *RedisUtil) MoveExpiredItems(expiry int64) {
	file, _ := ioutil.ReadFile(config.LuaScripts.MoveExpired)

	r.Client.Eval(
		string(file),
		[]string{config.Queues.Delayed + r.queueSuffix, config.Queues.Default + r.queueSuffix},
		expiry,
	)
}

// PopFromActiveQueue ...
func (r *RedisUtil) PopFromActiveQueue() string {
	file, _ := ioutil.ReadFile(config.LuaScripts.PopActive)

	res := r.Client.Eval(
		string(file),
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
