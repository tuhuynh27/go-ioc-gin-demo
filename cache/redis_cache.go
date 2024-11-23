package cache

import (
	"context"

	"tuhuynh.com/go-ioc-gin-example/config"
)

type RedisCache struct {
	Component  struct{}
	Implements struct{}       `implements:"Cache"`
	Qualifier  struct{}       `value:"redis"`
	Config     *config.Config `autowired:"true"`
}

func (c *RedisCache) Get(ctx context.Context, key string) (interface{}, error) {
	client := c.Config.Redis
	if client == nil {
		panic("redis client is nil")
	}

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (c *RedisCache) Set(ctx context.Context, key string, value interface{}) error {
	client := c.Config.Redis
	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
