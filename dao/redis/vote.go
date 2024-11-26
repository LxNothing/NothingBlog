package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

// 实际执行 redis中投票数据更新操作

const (
	VoteLimitDuration = 7 * 24 * 3600
	ScorePreVoted     = 432 // 每次投票所增加或者减少的分数
)

var (
	ErrVoteTimeExpire = errors.New("超过投票时间")
)

/*
与投票相关会涉及到三个zset
1. wevision.blog.time 其分数为时间戳，member为文章id - 表示文章的创建时间，用于限制投票时间间隔
2. wevision.blog.score 分数为时间戳，member为文章id - 表示文章的当前分数，记录分数使用
3. wevision.blog.voted.文章id 其分数为-1，0，1，表示投票种类，member为用户ID，用于记录为某个文章投票的人
*/

func UpdateBlogVoteRecord(uid string, blogId string, newVoted float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 查询 文章的分数 - 分数是文章的发布时间（时间戳）
	createTime := rdb.ZScore(ctx, getKeyWithPrefix(KeyBlogTimeZset), blogId).Val()
	if float64(time.Now().Unix())-createTime > VoteLimitDuration {
		return ErrVoteTimeExpire
	}

	// 获取用户对当前文章的过去投票状态
	keyBlog := getKeyWithPrefix(KeyBlogVotedZsetPrefix + blogId)
	oldVoted := rdb.ZScore(ctx, keyBlog, uid).Val()
	/*
		newVoted = 1 时，表示投赞成票
			oldVoted = 0: newVoted - oldVoted=1 增加的分数为1 * 432
			oldVoted = 1 ： newVoted - oldVoted=0 增加的分数为0 * 432
			oldVoted = -1 ：newVoted - oldVoted=2 增加的分数为 2 * 432
		newVoted = 0时，表示取消投票
			oldVoted = 0: newVoted - oldVoted=0 增加的分数为0 * 432
			oldVoted = 1 ： newVoted - oldVoted=-1 增加的分数为-1 * 432
			oldVoted = -1 ：newVoted - oldVoted=1 增加的分数为 1 * 432
		newVoted = -1时，表示投反对票
			oldVoted = 0: newVoted - oldVoted=-1 增加的分数为-1 * 432
			oldVoted = 1 ： newVoted - oldVoted=-2 增加的分数为-2 * 432
			oldVoted = -1 ：newVoted - oldVoted=0 增加的分数为 0 * 432
	*/
	addScore := (newVoted - oldVoted) * ScorePreVoted
	// 更新 投票表中的分数 - 通过事务保证更新操作均完成
	pipeline := rdb.TxPipeline()
	pipeline.ZIncrBy(ctx, getKeyWithPrefix(KeyBlogScoreZset), addScore, blogId)
	// 更新用户的投票记录
	if newVoted == 0 { // 用户取消投票
		pipeline.ZRem(ctx, keyBlog, uid)
	} else {
		// 向对应的文章 zset中插入用户投票结果
		pipeline.ZAdd(ctx, keyBlog, redis.Z{
			Score:  newVoted,
			Member: uid,
		})
	}
	_, err := pipeline.Exec(ctx)
	return err
}
