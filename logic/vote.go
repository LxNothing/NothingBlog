package logic

import (
	"NothingBlog/dao/redis"
	"NothingBlog/models"
	"strconv"
)

// 投票相关算法 - 参考阮一峰博客：https://www.ruanyifeng.com/blog/algorithm/
// 本文的方法基于redis，每投一个票分数加432 = 86400/200，含义是每天获得200票赞成，则总分数能够加一天（86400s） - 来自于《redis实战》

/*
	投票的逻辑分析：
		direction=1时 -> 表示投递赞成
			1. 用户已经投过赞成票
			2. 用户已经投过反对票
			3. 用户没有投过票
		direction=0时 -> 表示取消投票
			1. 用户没有投过票
			2. 投过赞成票或者反对票
		direction=1时 -> 表示投反对票
			1. 用户已经投过赞成票
			2. 用户已经投过反对票
			3. 用户没有投过票
	对投票的限制：这是对类似微博这种热点数据而言，避免用户对冷数据进行投票而影响到服务端的数据存储，因而限制
		文章的投票周期，比如自发布之日起多久内可以进行投票，
		但是对于博客的点赞功能来说，不应该设置这种限制
*/

// 实际执行投票功能

func VoteToBlog(uid int64, voteData *models.VoteDateParams) error {
	// 调用redis中的投票接口
	return redis.UpdateBlogVoteRecord(strconv.FormatInt(uid, 10),
		strconv.FormatInt(voteData.BlogId, 10),
		float64(voteData.Direction))
}
