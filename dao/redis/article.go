package redis

func CreateNewBlogZset(bid int64, score int64) error {
	/*ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 创建redis事务，保证下面两个zset同时成功
	pipeline := rdb.TxPipeline()
	now := float64(score)
	// 创建保存文章创建时间戳的zset
	pipeline.ZAdd(ctx, getKeyWithPrefix(KeyBlogTimeZset), redis.Z{
		Score:  now,
		Member: bid,
	})

	// 创建保存文章起始点赞分数的zset
	pipeline.ZAdd(ctx, getKeyWithPrefix(KeyBlogScoreZset), redis.Z{
		Score:  now,
		Member: bid,
	})

	_, err := pipeline.Exec(ctx)*/
	return nil
}

// 根据排序方式获取文章的ID列表
func GetBlogIdWithOrder(page int64, size int64, od string) ([]string, error) {
	/*redisKey := getKeyWithPrefix(KeyBlogTimeZset)
	if od == models.BlogOrderByScore {
		redisKey = getKeyWithPrefix(KeyBlogScoreZset)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	// ZRange 从zset中按照升序获取用户id
	return rdb.ZRevRange(ctx, redisKey, (page-1)*size, size-1).Result()*/
	return nil, nil
}
