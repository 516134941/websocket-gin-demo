package middlewares

import (
	"fmt"
	"test/websocket-gin-demo/message-chat/config"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

// Redis .
func Redis(cacheName string, tomlConfig *config.Config) gin.HandlerFunc {
	cacheConfig, ok := tomlConfig.RedisServerConf(cacheName)
	if !ok {
		panic(fmt.Sprintf("%v not set.", cacheName))
	}

	// 链接数据库
	pool := newPool(cacheConfig.Addr, cacheConfig.Password, cacheConfig.DB)

	return func(c *gin.Context) {
		c.Set(cacheName, pool)
		c.Next()
	}
}

// newPool New redis pool.
func newPool(server, password string, db int) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server,
				redis.DialPassword(password),
				redis.DialDatabase(db),
				redis.DialConnectTimeout(500*time.Millisecond),
				redis.DialReadTimeout(500*time.Millisecond),
				redis.DialWriteTimeout(500*time.Millisecond))
			if err != nil {
				return nil, err
			}

			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < 5*time.Second {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
