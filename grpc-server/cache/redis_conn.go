package cache

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"server/grpc-server/config"
	"server/grpc-server/log"
	"time"
)

const (
	DEFAULT_REDIS_CONFIG = "redis.ini"
)

var RedisPool *redis.Pool
var NofindInCacheError = errors.New("no find in Cache Error")

func NewRedisPool(redisCfg *config.RedisConfig) (*redis.Pool, error) {
	if redisCfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	redisUrl := fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port)
	pool := &redis.Pool{
		MaxIdle:     redisCfg.MaxIdle,
		MaxActive:   redisCfg.MaxActive,
		IdleTimeout: time.Duration(redisCfg.IdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", redisUrl,
				redis.DialPassword(redisCfg.Password),
				redis.DialDatabase(redisCfg.DB),
				redis.DialConnectTimeout(time.Duration(redisCfg.Timeout)*time.Second),
				redis.DialReadTimeout(time.Duration(redisCfg.Timeout)*time.Second),
				redis.DialWriteTimeout(time.Duration(redisCfg.Timeout)*time.Second))

			if err != nil {
				return nil, err
			}

			return con, nil
		}, // Dial end
	}

	return pool, nil
}

func GetRedisClient() redis.Conn {
	if RedisPool == nil {
		return nil
	}

	return RedisPool.Get()
}

func init() {
	var err error
	cfg := config.NewRedisConfig()
	if err := cfg.LoadConfig("redis", DEFAULT_REDIS_CONFIG); err != nil {
		RedisPool = nil
		log.Log.Printf("Init NewRedisConfig LoadConfig fail with error: %+v", err)
		return
	}

	RedisPool, err = NewRedisPool(cfg)
	if err != nil {
		log.Log.Printf("Init NewRedisPool fail with error:%+v", err)
	}
}
