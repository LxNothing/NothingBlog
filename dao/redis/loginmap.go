package redis

import (
	"context"
	"fmt"
	"time"
)

// 这里使用mysql作为映射感觉不太好，因为涉及到频繁的查询任务，性能不好，应该使用redis
// 将用户id和对应的token插入到login表
func InsertLoginInfo(user_id int64, token string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	err = rdb.Set(ctx, fmt.Sprintf("%d", user_id), token, time.Hour).Err()
	return err
}

// 从登录表中查询token
func QueryTokenByUserId(uid int64) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	return rdb.Get(ctx, fmt.Sprintf("%d", uid)).Result()
}
