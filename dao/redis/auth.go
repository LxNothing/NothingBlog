package redis

import (
	"context"
	"time"
)

// 使用用户名缓存
func CatchVerifyCode(uid string, code string, expTime time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	key := getKeyWithPrefix(KeyVerifyCodeStrPrefix + uid)
	return rdb.Set(ctx, key, code, expTime).Err()
}
