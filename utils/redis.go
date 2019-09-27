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
	Client *redis.Client
}

// NewRedisUtil ...
func NewRedisUtil(addres string) RedisUtilInterface {
	return &RedisUtil{
		Client: redis.NewClient(&redis.Options{Addr: addres, Password: "", DB: 0}),
	}
}

// SetDelayed task
func (r *RedisUtil) SetDelayed(value interface{}, expiry float64) {
	r.Client.ZAdd(config.Queues.Delayed, redis.Z{Member: value, Score: expiry})
}

// MoveExpiredItems ...
func (r *RedisUtil) MoveExpiredItems(expiry int64) {
	file, _ := ioutil.ReadFile(config.LuaScripts.MoveExpired)

	r.Client.Eval(
		string(file),
		[]string{config.Queues.Delayed, config.Queues.Default},
		expiry,
	)

	// log.Println("Got values from delayed...")
	// log.Println(res)
}

// PopFromActiveQueue ...
func (r *RedisUtil) PopFromActiveQueue() string {
	file, _ := ioutil.ReadFile(config.LuaScripts.PopActive)

	res := r.Client.Eval(
		string(file),
		[]string{config.Queues.Default, config.Queues.Expired, time.Now().Format(time.RFC3339)},
		time.Now().Unix(),
	)

	if res.Val() == nil {
		return ""
	}

	return res.Val().(string)
}
