package redis

import (
	"testing"
)

func Test(t *testing.T) {
	redisCfg := &Config{Address: "localhost:6379", DB: 0}

	redisConn, err := New(redisCfg)

	if err != nil {
		t.Error(err)
	}

	redisKey := "testKey"
	redisValue := "testValue"

	redisConn.Set(redisKey, redisValue, 0)

	redisRead, err := redisConn.Get(redisKey)

	t.Logf("Redis test result: %v", redisRead == redisValue)

	if err != nil {
		t.Error(err)
	}
}
